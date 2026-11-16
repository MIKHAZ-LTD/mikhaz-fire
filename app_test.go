package main

import (
	"fmt"
	"strings"
	"testing"
)

func findWorkingAccountForProject(t *testing.T, app *App, targetProjectID string) string {
	t.Helper()
	accounts := app.GetAccounts()
	for _, acc := range accounts {
		projs, err := app.ListProjects(acc.Email)
		if err != nil {
			continue
		}
		for _, p := range projs {
			if p.ProjectID == targetProjectID {
				return acc.Email
			}
		}
	}
	return ""
}

func TestProjects(t *testing.T) {
	app := NewApp()
	accounts := app.GetAccounts()
	t.Logf("Discovered %d accounts", len(accounts))

	allowed := []string{"minorgames", "textcup"}

	for _, acc := range accounts {
		t.Logf("Checking account: %s (source: %s, isSA: %v)", acc.Email, acc.Source, acc.IsServiceAcc)
		projs, err := app.ListProjects(acc.Email)
		if err != nil {
			t.Logf("  Error listing projects for %s: %v", acc.Email, err)
			continue
		}
		t.Logf("  Total projects for %s: %d", acc.Email, len(projs))
		for _, p := range projs {
			pid := strings.ToLower(p.ProjectID)
			isAllowed := false
			for _, target := range allowed {
				if strings.Contains(pid, target) || strings.Contains(strings.ToLower(p.DisplayName), target) {
					isAllowed = true
					break
				}
			}
			if isAllowed {
				t.Logf("  >>> FOUND TARGET PROJECT: ID=%s, DisplayName=%s", p.ProjectID, p.DisplayName)
				// Test listing root collections (read-only)
				cols, colErr := app.ListRootCollections(acc.Email, p.ProjectID, "")
				if colErr != nil {
					t.Logf("      Firestore cols error: %v", colErr)
				} else {
					t.Logf("      Firestore root collections (%d): %v", len(cols), cols)
				}
				// Test listing Auth users (limit 5, read-only)
				usersResp, authErr := app.ListAuthUsers(acc.Email, p.ProjectID, 5, "", "")
				if authErr != nil {
					t.Logf("      Auth users error: %v", authErr)
				} else if usersResp != nil {
					t.Logf("      Auth users loaded: %d", len(usersResp.Users))
				}
			}
		}
	}
}

func TestTextCupFiltersAndQueries(t *testing.T) {
	app := NewApp()
	const projectID = "textcup-55d52"
	email := findWorkingAccountForProject(t, app, projectID)
	if email == "" {
		t.Skipf("No active account found with access to %s", projectID)
	}
	t.Logf("Testing read-only filters and queries on %s via %s", projectID, email)

	// 1. Test ListRootCollections & ListDocuments
	cols, err := app.ListRootCollections(email, projectID, "")
	if err != nil {
		t.Fatalf("ListRootCollections failed: %v", err)
	}
	if len(cols) == 0 {
		t.Fatalf("expected non-empty collections in %s", projectID)
	}

	targetCol := "users"
	listResp, err := app.ListDocuments(email, projectID, "", targetCol, 5, "")
	if err != nil {
		t.Fatalf("ListDocuments(%s) failed: %v", targetCol, err)
	}
	rawDocs, _ := listResp["documents"].([]interface{})
	if len(rawDocs) == 0 {
		t.Fatalf("expected at least 1 document in %s", targetCol)
	}

	firstDoc, _ := rawDocs[0].(map[string]interface{})
	docName, _ := firstDoc["name"].(string)
	parts := strings.Split(docName, "/")
	firstDocID := parts[len(parts)-1]
	fields, _ := firstDoc["fields"].(map[string]interface{})

	// Pick a primitive field from the first document to test RunSimpleQuery & RunJSScript
	var sampleField string
	var sampleTypedVal map[string]interface{}
	var sampleJSVal interface{}
	for k, v := range fields {
		m, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if sv, ok := m["stringValue"].(string); ok {
			sampleField = k
			sampleTypedVal = map[string]interface{}{"stringValue": sv}
			sampleJSVal = sv
			break
		}
		if bv, ok := m["booleanValue"].(bool); ok {
			sampleField = k
			sampleTypedVal = map[string]interface{}{"booleanValue": bv}
			sampleJSVal = bv
			break
		}
	}

	if sampleField != "" {
		// 2a. Test RunSimpleQuery with EQUAL fieldFilter + limit
		res, err := app.RunSimpleQuery(email, projectID, "", map[string]interface{}{
			"structuredQuery": map[string]interface{}{
				"from": []map[string]interface{}{{"collectionId": targetCol}},
				"where": map[string]interface{}{
					"fieldFilter": map[string]interface{}{
						"field": map[string]interface{}{"fieldPath": sampleField},
						"op":    "EQUAL",
						"value": sampleTypedVal,
					},
				},
				"limit": 3,
			},
		})
		if err != nil {
			t.Fatalf("RunSimpleQuery EQUAL failed: %v", err)
		}
		resArr, _ := res.([]interface{})
		if len(resArr) == 0 {
			t.Fatalf("expected at least 1 doc for %s == %v", sampleField, sampleJSVal)
		}
		t.Logf("RunSimpleQuery(%s == %v) returned %d entry(s)", sampleField, sampleJSVal, len(resArr))

		// 2b. Test RunSimpleQuery with IS_NOT_NULL unaryFilter + orderBy
		resUnary, err := app.RunSimpleQuery(email, projectID, "", map[string]interface{}{
			"structuredQuery": map[string]interface{}{
				"from": []map[string]interface{}{{"collectionId": targetCol}},
				"where": map[string]interface{}{
					"unaryFilter": map[string]interface{}{
						"field": map[string]interface{}{"fieldPath": sampleField},
						"op":    "IS_NOT_NULL",
					},
				},
				"orderBy": []map[string]interface{}{
					{
						"field":     map[string]interface{}{"fieldPath": sampleField},
						"direction": "ASCENDING",
					},
				},
				"limit": 3,
			},
		})
		if err != nil {
			t.Fatalf("RunSimpleQuery IS_NOT_NULL + orderBy failed: %v", err)
		}
		resUnaryArr, _ := resUnary.([]interface{})
		t.Logf("RunSimpleQuery(%s IS_NOT_NULL) returned %d entry(s)", sampleField, len(resUnaryArr))

		// 3. Test RunJSScript with where, forEach, size, console.log, and db.doc
		jsScript := fmt.Sprintf(`
			async function run() {
				const snap = await db.collection(%q).where(%q, "==", %q).limit(3).get();
				console.log("Matched size:", snap.size);
				const ids = [];
				snap.forEach(d => {
					ids.push(d.id);
				});
				const single = await db.doc(%q).get();
				return { size: snap.size, empty: snap.empty, ids: ids, singleId: single.id, singleExists: single.exists };
			}
		`, targetCol, sampleField, fmt.Sprintf("%v", sampleJSVal), targetCol+"/"+firstDocID)

		// Only run string equality in %q when sampleJSVal is string; handle bool safely
		if bVal, isBool := sampleJSVal.(bool); isBool {
			jsScript = fmt.Sprintf(`
				async function run() {
					const snap = await db.collection(%q).where(%q, "==", %t).limit(3).get();
					console.log("Matched size:", snap.size);
					const ids = [];
					snap.forEach(d => {
						ids.push(d.id);
					});
					const single = await db.doc(%q).get();
					return { size: snap.size, empty: snap.empty, ids: ids, singleId: single.id, singleExists: single.exists };
				}
			`, targetCol, sampleField, bVal, targetCol+"/"+firstDocID)
		}

		jsRes, err := app.RunJSScript(email, projectID, "", jsScript)
		if err != nil {
			t.Fatalf("RunJSScript failed: %v", err)
		}
		jsMap, ok := jsRes.(map[string]interface{})
		if !ok {
			t.Fatalf("expected map from RunJSScript with console.log, got %T: %+v", jsRes, jsRes)
		}
		logs, _ := jsMap["logs"].([]string)
		resultMap, _ := jsMap["result"].(map[string]interface{})
		if len(logs) == 0 || resultMap == nil {
			t.Fatalf("expected logs and result from RunJSScript, got %+v", jsMap)
		}
		t.Logf("RunJSScript logs=%v result=%+v", logs, resultMap)
	}

	// 4. Test subcollection RunQuery if any sampled document in books/users has subcollections
	for _, colCandidate := range []string{"books", "users"} {
		colResp, err := app.ListDocuments(email, projectID, "", colCandidate, 5, "")
		if err != nil {
			continue
		}
		docs, _ := colResp["documents"].([]interface{})
		testedSub := false
		for _, rd := range docs {
			dm, ok := rd.(map[string]interface{})
			if !ok {
				continue
			}
			dn, _ := dm["name"].(string)
			idx := strings.Index(dn, "/documents/")
			if idx == -1 {
				continue
			}
			relDoc := dn[idx+len("/documents/"):]
			subs, err := app.ListDocumentSubcollections(email, projectID, "", relDoc)
			if err == nil && len(subs) > 0 {
				subPath := relDoc + "/" + subs[0]
				subRes, subErr := app.RunSimpleQuery(email, projectID, "", map[string]interface{}{
					"structuredQuery": map[string]interface{}{
						"from":  []map[string]interface{}{{"collectionId": subPath}},
						"limit": 3,
					},
				})
				if subErr != nil {
					t.Fatalf("RunSimpleQuery on subcollection %s failed: %v", subPath, subErr)
				}
				subArr, _ := subRes.([]interface{})
				t.Logf("Subcollection RunSimpleQuery(%s) returned %d entry(s)", subPath, len(subArr))
				testedSub = true
				break
			}
		}
		if testedSub {
			break
		}
	}

	// 5. Test Auth search on textcup-55d52
	usersAll, err := app.ListAuthUsers(email, projectID, 10, "", "")
	if err != nil {
		t.Fatalf("ListAuthUsers failed: %v", err)
	}
	if len(usersAll.Users) > 0 {
		needle := usersAll.Users[0].UID[:4]
		if usersAll.Users[0].Email != "" {
			parts := strings.Split(usersAll.Users[0].Email, "@")
			if len(parts[0]) >= 3 {
				needle = parts[0][:3]
			}
		}
		searchResp, err := app.ListAuthUsers(email, projectID, 50, "", needle)
		if err != nil {
			t.Fatalf("ListAuthUsers(search=%q) failed: %v", needle, err)
		}
		if len(searchResp.Users) == 0 {
			t.Fatalf("expected at least 1 Auth user matching %q", needle)
		}
		t.Logf("ListAuthUsers(search=%q) returned %d user(s)", needle, len(searchResp.Users))
	}
}

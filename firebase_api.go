package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Project struct {
	ProjectID     string `json:"projectId"`
	ProjectNumber string `json:"projectNumber"`
	DisplayName   string `json:"displayName"`
}

type AuthUser struct {
	UID              string                 `json:"uid"`
	Email            string                 `json:"email"`
	EmailVerified    bool                   `json:"emailVerified"`
	DisplayName      string                 `json:"displayName"`
	PhotoURL         string                 `json:"photoUrl"`
	Disabled         bool                   `json:"disabled"`
	CreatedAt        string                 `json:"createdAt"`
	LastLoginAt      string                 `json:"lastLoginAt"`
	CustomClaims     map[string]interface{} `json:"customClaims"`
	CustomAttributes string                 `json:"customAttributes,omitempty"`
	PasswordHash     string                 `json:"passwordHash,omitempty"`
	Salt             string                 `json:"salt,omitempty"`
	Providers        []string               `json:"providers,omitempty"`
	ProviderUserInfo []interface{}          `json:"providerUserInfo,omitempty"`
}

type AuthUsersResponse struct {
	Users         []AuthUser `json:"users"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}

// REST call helper
func doRequest(method, reqUrl, token string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Extract project ID to set the Google User Project quota header (required for ADC user credentials)
	if idx := strings.Index(reqUrl, "/projects/"); idx != -1 {
		rest := reqUrl[idx+10:]
		projectID := rest
		if endIdx := strings.Index(rest, "/"); endIdx != -1 {
			projectID = rest[:endIdx]
		}
		if colonIdx := strings.Index(projectID, ":"); colonIdx != -1 {
			projectID = projectID[:colonIdx]
		}
		if qIdx := strings.Index(projectID, "?"); qIdx != -1 {
			projectID = projectID[:qIdx]
		}
		if projectID != "" {
			req.Header.Set("x-goog-user-project", projectID)
		}
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// ListProjects lists Firebase projects.
func ListProjects(token string) ([]Project, error) {
	reqUrl := "https://firebase.googleapis.com/v1beta1/projects?pageSize=100"
	body, err := doRequest("GET", reqUrl, token, nil)
	if err != nil {
		// Fallback to Resource Manager if Firebase API fails or is not enabled
		reqUrl = "https://cloudresourcemanager.googleapis.com/v1/projects"
		body, err = doRequest("GET", reqUrl, token, nil)
		if err != nil {
			return nil, err
		}
		var res struct {
			Projects []struct {
				ProjectID string `json:"projectId"`
				Name      string `json:"name"`
			} `json:"projects"`
		}
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, err
		}
		var list []Project
		for _, p := range res.Projects {
			list = append(list, Project{
				ProjectID:   p.ProjectID,
				DisplayName: p.Name,
			})
		}
		return list, nil
	}

	var res struct {
		Results []Project `json:"results"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res.Results, nil
}

// ListAuthUsers retrieves users in Firebase Auth.
func ListAuthUsers(token string, projectId string, pageSize int, pageToken string) (*AuthUsersResponse, error) {
	if pageSize <= 0 {
		pageSize = 50
	}
	reqUrl := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/projects/%s/accounts:batchGet?maxResults=%d", projectId, pageSize)
	if pageToken != "" {
		reqUrl += "&nextPageToken=" + url.QueryEscape(pageToken)
	}

	body, err := doRequest("GET", reqUrl, token, nil)

	if err != nil {
		return nil, err
	}

	var raw struct {
		Users []struct {
			LocalID          string        `json:"localId"`
			Email            string        `json:"email"`
			EmailVerified    bool          `json:"emailVerified"`
			DisplayName      string        `json:"displayName"`
			PhotoURL         string        `json:"photoUrl"`
			Disabled         bool          `json:"disabled"`
			CreatedAt        string        `json:"createdAt"`
			LastLoginAt      string        `json:"lastLoginAt"`
			CustomAttributes string        `json:"customAttributes"` // stringified JSON
			PasswordHash     string        `json:"passwordHash"`
			Salt             string        `json:"salt"`
			ProviderUserInfo []interface{} `json:"providerUserInfo"`
		} `json:"users"`
		NextPageToken string `json:"nextPageToken"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	var users []AuthUser
	for _, u := range raw.Users {
		var claims map[string]interface{}
		if u.CustomAttributes != "" {
			_ = json.Unmarshal([]byte(u.CustomAttributes), &claims)
		}
		var providers []string
		for _, p := range u.ProviderUserInfo {
			if pMap, ok := p.(map[string]interface{}); ok {
				if pId, ok := pMap["providerId"].(string); ok {
					providers = append(providers, pId)
				}
			}
		}
		users = append(users, AuthUser{
			UID:              u.LocalID,
			Email:            u.Email,
			EmailVerified:    u.EmailVerified,
			DisplayName:      u.DisplayName,
			PhotoURL:         u.PhotoURL,
			Disabled:         u.Disabled,
			CreatedAt:        u.CreatedAt,
			LastLoginAt:      u.LastLoginAt,
			CustomClaims:     claims,
			CustomAttributes: u.CustomAttributes,
			PasswordHash:     u.PasswordHash,
			Salt:             u.Salt,
			Providers:        providers,
			ProviderUserInfo: u.ProviderUserInfo,
		})
	}

	return &AuthUsersResponse{
		Users:         users,
		NextPageToken: raw.NextPageToken,
	}, nil
}

// UpdateAuthUser updates a user in Firebase Auth.
func UpdateAuthUser(token string, projectId string, req map[string]interface{}) (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/projects/%s/accounts:update", projectId)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resBody, err := doRequest("POST", reqUrl, token, body)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(resBody, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteAuthUser deletes a user from Firebase Auth.
func DeleteAuthUser(token string, projectId string, uid string) error {
	reqUrl := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/projects/%s/accounts:delete", projectId)
	req := map[string]interface{}{
		"localId": uid,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = doRequest("POST", reqUrl, token, body)
	return err
}

// CreateAuthUser creates a user in Firebase Auth.
func CreateAuthUser(token string, projectId string, req map[string]interface{}) (map[string]interface{}, error) {
	reqUrl := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/projects/%s/accounts", projectId)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resBody, err := doRequest("POST", reqUrl, token, body)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(resBody, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListRootCollections lists root level collections.
func ListRootCollections(token string, projectId string, databaseId string) ([]string, error) {
	if databaseId == "" {
		databaseId = "(default)"
	}
	reqUrl := fmt.Sprintf("https://firestore.googleapis.com/v1/projects/%s/databases/%s/documents:listCollectionIds", projectId, databaseId)
	body, err := doRequest("POST", reqUrl, token, []byte("{}"))
	if err != nil {
		return nil, err
	}

	var res struct {
		CollectionIDs []string `json:"collectionIds"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res.CollectionIDs, nil
}

// ListDocuments lists documents in a Firestore collection.
func ListDocuments(token string, projectId string, databaseId string, collPath string, pageSize int, pageToken string) (map[string]interface{}, error) {
	if databaseId == "" {
		databaseId = "(default)"
	}
	// collPath might be e.g. "users" or "users/123/posts"
	reqUrl := fmt.Sprintf("https://firestore.googleapis.com/v1/projects/%s/databases/%s/documents/%s", projectId, databaseId, collPath)
	params := url.Values{}
	if pageSize > 0 {
		params.Set("pageSize", fmt.Sprintf("%d", pageSize))
	}
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	if len(params) > 0 {
		reqUrl += "?" + params.Encode()
	}

	body, err := doRequest("GET", reqUrl, token, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetDocument fetches a single document.
func GetDocument(token string, projectId string, databaseId string, docPath string) (map[string]interface{}, error) {
	if databaseId == "" {
		databaseId = "(default)"
	}
	reqUrl := fmt.Sprintf("https://firestore.googleapis.com/v1/projects/%s/databases/%s/documents/%s", projectId, databaseId, docPath)
	body, err := doRequest("GET", reqUrl, token, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// SaveDocument saves/overwrites/updates a document.
func SaveDocument(token string, projectId string, databaseId string, docPath string, fields map[string]interface{}, isNew bool) (map[string]interface{}, error) {
	if databaseId == "" {
		databaseId = "(default)"
	}
	
	// Format is projects/{projectId}/databases/{databaseId}/documents/{docPath}
	docName := fmt.Sprintf("projects/%s/databases/%s/documents/%s", projectId, databaseId, docPath)
	var doc map[string]interface{}
	if isNew {
		doc = map[string]interface{}{
			"fields": fields,
		}
	} else {
		doc = map[string]interface{}{
			"name":   docName,
			"fields": fields,
		}
	}
	
	body, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	var reqUrl string
	var method string

	if isNew {
		// Create: POST to parent collection with documentId parameter
		idx := strings.LastIndex(docPath, "/")
		var parent string
		var docId string
		if idx == -1 {
			parent = ""
			docId = docPath
		} else {
			parent = docPath[:idx]
			docId = docPath[idx+1:]
		}
		
		if parent != "" {
			reqUrl = fmt.Sprintf("https://firestore.googleapis.com/v1/projects/%s/databases/%s/documents/%s?documentId=%s", projectId, databaseId, parent, docId)
		} else {
			reqUrl = fmt.Sprintf("https://firestore.googleapis.com/v1/projects/%s/databases/%s/documents?documentId=%s", projectId, databaseId, docId)
		}
		method = "POST"
	} else {
		// Update: PATCH to the document url
		reqUrl = fmt.Sprintf("https://firestore.googleapis.com/v1/%s", docName)
		method = "PATCH"
	}

	resBody, err := doRequest(method, reqUrl, token, body)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(resBody, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteDocument deletes a Firestore document.
func DeleteDocument(token string, projectId string, databaseId string, docPath string) error {
	if databaseId == "" {
		databaseId = "(default)"
	}
	reqUrl := fmt.Sprintf("https://firestore.googleapis.com/v1/projects/%s/databases/%s/documents/%s", projectId, databaseId, docPath)
	_, err := doRequest("DELETE", reqUrl, token, nil)
	return err
}

// RunQuery executes a Firestore StructuredQuery.
func RunQuery(token string, projectId string, databaseId string, queryBody map[string]interface{}) (interface{}, error) {
	if databaseId == "" {
		databaseId = "(default)"
	}
	reqUrl := fmt.Sprintf("https://firestore.googleapis.com/v1/projects/%s/databases/%s/documents:runQuery", projectId, databaseId)
	body, err := json.Marshal(queryBody)
	if err != nil {
		return nil, err
	}

	resBody, err := doRequest("POST", reqUrl, token, body)
	if err != nil {
		return nil, err
	}

	var res interface{}
	if err := json.Unmarshal(resBody, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListDocumentSubcollections lists the subcollection IDs of a given document.
func ListDocumentSubcollections(token string, projectId string, databaseId string, docPath string) ([]string, error) {
	if databaseId == "" {
		databaseId = "(default)"
	}
	reqUrl := fmt.Sprintf("https://firestore.googleapis.com/v1/projects/%s/databases/%s/documents/%s:listCollectionIds", projectId, databaseId, docPath)
	body, err := doRequest("POST", reqUrl, token, []byte("{}"))
	if err != nil {
		if strings.Contains(err.Error(), "status 404") {
			return nil, nil // Not found usually means no subcollections exist
		}
		return nil, err
	}

	var res struct {
		CollectionIDs []string `json:"collectionIds"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res.CollectionIDs, nil
}

// DeleteDocumentRecursive recursively deletes a document and all its subcollections.
func DeleteDocumentRecursive(token string, projectId string, databaseId string, docPath string) error {
	// 1. List subcollections of the document
	subcolls, err := ListDocumentSubcollections(token, projectId, databaseId, docPath)
	if err != nil {
		return err
	}

	// 2. Recursively delete documents in each subcollection
	for _, colId := range subcolls {
		colPath := docPath + "/" + colId
		resp, err := ListDocuments(token, projectId, databaseId, colPath, 1000, "")
		if err != nil {
			return err
		}

		rawDocs, _ := resp["documents"].([]interface{})
		for _, rd := range rawDocs {
			docMap, ok := rd.(map[string]interface{})
			if !ok {
				continue
			}
			docName, _ := docMap["name"].(string)
			// Extract relative child document path
			idx := strings.Index(docName, "/documents/")
			if idx == -1 {
				continue
			}
			childDocPath := docName[idx+11:]

			// Delete child recursively
			err = DeleteDocumentRecursive(token, projectId, databaseId, childDocPath)
			if err != nil {
				return err
			}
		}
	}

	// 3. Delete the parent document
	return DeleteDocument(token, projectId, databaseId, docPath)
}

// CheckConflict checks if a document or collection exists in the destination.
func CheckConflict(token string, projectId string, databaseId string, path string, isCollection bool) (bool, string, error) {
	if databaseId == "" {
		databaseId = "(default)"
	}
	if isCollection {
		// Check if collection contains any documents
		res, err := ListDocuments(token, projectId, databaseId, path, 1, "")
		if err != nil {
			// If we get a 404 or resource not found, it means the parent document/path doesn't exist, hence no collection exists.
			if strings.Contains(err.Error(), "status 404") {
				return false, "No conflict (collection not found/empty)", nil
			}
			return false, "", err
		}
		docs, ok := res["documents"].([]interface{})
		if ok && len(docs) > 0 {
			return true, fmt.Sprintf("Conflict: Destination collection has %d or more document(s)", len(docs)), nil
		}
		return false, "No conflict (collection empty)", nil
	} else {
		// Check if document exists
		_, err := GetDocument(token, projectId, databaseId, path)
		if err != nil {
			if strings.Contains(err.Error(), "status 404") {
				return false, "No conflict (document not found)", nil
			}
			return false, "", err
		}
		return true, "Conflict: Destination document already exists", nil
	}
}

// CopyDocument copies a document from source to destination.
func CopyDocument(srcToken, srcProj, srcDb string, destToken, destProj, destDb string, docPath string, overwrite bool, recursive bool) (int, error) {
	// 1. Get the source document
	doc, err := GetDocument(srcToken, srcProj, srcDb, docPath)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch source document: %v", err)
	}

	fields, _ := doc["fields"].(map[string]interface{})

	// 2. Check conflict if not overwriting
	if !overwrite {
		exists, _, err := CheckConflict(destToken, destProj, destDb, docPath, false)
		if err != nil {
			return 0, fmt.Errorf("conflict check failed: %v", err)
		}
		if exists {
			return 0, fmt.Errorf("document already exists in destination and overwrite is false")
		}
	}

	// 3. Save to destination
	_, err = SaveDocument(destToken, destProj, destDb, docPath, fields, false) // false is fine to overwrite/create at exact path
	if err != nil {
		return 0, fmt.Errorf("failed to save document to destination: %v", err)
	}

	copiedCount := 1

	// 4. Handle subcollections if recursive is true
	if recursive {
		subcolls, err := ListDocumentSubcollections(srcToken, srcProj, srcDb, docPath)
		if err != nil {
			// Subcollections might not be allowed/supported, or just empty
			// We won't block if there is a 404 or list collections is not supported
			if !strings.Contains(err.Error(), "status 404") {
				return copiedCount, fmt.Errorf("copied document, but failed to list subcollections: %v", err)
			}
		}

		for _, subId := range subcolls {
			subColPath := docPath + "/" + subId
			n, err := CopyCollection(srcToken, srcProj, srcDb, destToken, destProj, destDb, subColPath, overwrite, recursive)
			copiedCount += n
			if err != nil {
				return copiedCount, err
			}
		}
	}

	return copiedCount, nil
}

// CopyCollection copies all documents in a collection.
func CopyCollection(srcToken, srcProj, srcDb string, destToken, destProj, destDb string, colPath string, overwrite bool, recursive bool) (int, error) {
	copiedCount := 0
	pageToken := ""

	for {
		resp, err := ListDocuments(srcToken, srcProj, srcDb, colPath, 100, pageToken)
		if err != nil {
			// If the collection doesn't exist/empty in source, returns error or empty. 404 is fine.
			if strings.Contains(err.Error(), "status 404") {
				return copiedCount, nil
			}
			return copiedCount, fmt.Errorf("failed to list source documents: %v", err)
		}

		rawDocs, _ := resp["documents"].([]interface{})
		for _, rd := range rawDocs {
			docMap, ok := rd.(map[string]interface{})
			if !ok {
				continue
			}
			docName, _ := docMap["name"].(string)
			// Extract relative child document path
			idx := strings.Index(docName, "/documents/")
			if idx == -1 {
				continue
			}
			childDocPath := docName[idx+11:]

			// Copy document
			n, err := CopyDocument(srcToken, srcProj, srcDb, destToken, destProj, destDb, childDocPath, overwrite, recursive)
			copiedCount += n
			if err != nil {
				return copiedCount, err
			}
		}

		nextPageToken, _ := resp["nextPageToken"].(string)
		if nextPageToken == "" {
			break
		}
		pageToken = nextPageToken
	}

	return copiedCount, nil
}

type ExportAuthResponse struct {
	Users          []AuthUser             `json:"users"`
	Count          int                    `json:"count"`
	HashParameters map[string]interface{} `json:"hashParameters"`
}

// ExportAuthUsers fetches all users across pages for a project along with default hash parameters.
func ExportAuthUsers(token string, projectId string) (*ExportAuthResponse, error) {
	var allUsers []AuthUser
	pageToken := ""

	for {
		resp, err := ListAuthUsers(token, projectId, 1000, pageToken)
		if err != nil {
			return nil, err
		}
		allUsers = append(allUsers, resp.Users...)
		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	return &ExportAuthResponse{
		Users: allUsers,
		Count: len(allUsers),
		HashParameters: map[string]interface{}{
			"hashAlgorithm": "SCRYPT",
			"signerKey":     "",
			"saltSeparator": "Bw==",
			"rounds":        8,
			"memoryCost":    14,
		},
	}, nil
}

type AuthImportRequest struct {
	Users         []map[string]interface{} `json:"users"`
	HashAlgorithm string                   `json:"hashAlgorithm"`
	SignerKey     string                   `json:"signerKey"`
	SaltSeparator string                   `json:"saltSeparator"`
	Rounds        int                      `json:"rounds"`
	MemoryCost    int                      `json:"memoryCost"`
}

type AuthImportResult struct {
	Success       bool                     `json:"success"`
	ImportedCount int                      `json:"importedCount"`
	FailedCount   int                      `json:"failedCount"`
	Errors        []map[string]interface{} `json:"errors,omitempty"`
}

func parseTimestampInt64(val interface{}) (string, bool) {
	if val == nil {
		return "", false
	}
	switch v := val.(type) {
	case float64:
		if v > 0 {
			return fmt.Sprintf("%.0f", v), true
		}
	case int64:
		if v > 0 {
			return fmt.Sprintf("%d", v), true
		}
	case int:
		if v > 0 {
			return fmt.Sprintf("%d", v), true
		}
	case string:
		vStr := strings.TrimSpace(v)
		if vStr == "" || vStr == "0" {
			return "", false
		}
		if _, err := strconv.ParseInt(vStr, 10, 64); err == nil {
			return vStr, true
		}
	}
	return "", false
}

func sanitizeUserForImport(u map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})

	// localId / uid / local_id
	if uid, ok := u["localId"].(string); ok && uid != "" {
		sanitized["localId"] = uid
	} else if uid, ok := u["uid"].(string); ok && uid != "" {
		sanitized["localId"] = uid
	} else if uid, ok := u["local_id"].(string); ok && uid != "" {
		sanitized["localId"] = uid
	}

	// email
	if email, ok := u["email"].(string); ok && email != "" {
		sanitized["email"] = email
	}

	// emailVerified / email_verified
	if emailVerified, ok := u["emailVerified"].(bool); ok {
		sanitized["emailVerified"] = emailVerified
	} else if emailVerified, ok := u["email_verified"].(bool); ok {
		sanitized["emailVerified"] = emailVerified
	}

	// displayName / display_name
	if displayName, ok := u["displayName"].(string); ok && displayName != "" {
		sanitized["displayName"] = displayName
	} else if displayName, ok := u["display_name"].(string); ok && displayName != "" {
		sanitized["displayName"] = displayName
	}

	// photoUrl / photoURL / photo_url
	if photoUrl, ok := u["photoUrl"].(string); ok && photoUrl != "" {
		sanitized["photoUrl"] = photoUrl
	} else if photoUrl, ok := u["photoURL"].(string); ok && photoUrl != "" {
		sanitized["photoUrl"] = photoUrl
	} else if photoUrl, ok := u["photo_url"].(string); ok && photoUrl != "" {
		sanitized["photoUrl"] = photoUrl
	}

	// disabled
	if disabled, ok := u["disabled"].(bool); ok {
		sanitized["disabled"] = disabled
	}

	// passwordHash / password_hash
	if pwdHash, ok := u["passwordHash"].(string); ok && pwdHash != "" {
		sanitized["passwordHash"] = pwdHash
	} else if pwdHash, ok := u["password_hash"].(string); ok && pwdHash != "" {
		sanitized["passwordHash"] = pwdHash
	}

	// salt
	if salt, ok := u["salt"].(string); ok && salt != "" {
		sanitized["salt"] = salt
	}

	// createdAt / created_at
	var createdVal interface{}
	if val, ok := u["createdAt"]; ok {
		createdVal = val
	} else if val, ok := u["created_at"]; ok {
		createdVal = val
	}
	if tsStr, ok := parseTimestampInt64(createdVal); ok {
		sanitized["createdAt"] = tsStr
	}

	// lastLoginAt / last_login_at
	var lastLoginVal interface{}
	if val, ok := u["lastLoginAt"]; ok {
		lastLoginVal = val
	} else if val, ok := u["last_login_at"]; ok {
		lastLoginVal = val
	}
	if tsStr, ok := parseTimestampInt64(lastLoginVal); ok {
		sanitized["lastLoginAt"] = tsStr
	}

	// customAttributes / customClaims / custom_attributes
	if customAttrs, ok := u["customAttributes"].(string); ok && customAttrs != "" {
		sanitized["customAttributes"] = customAttrs
	} else if customAttrs, ok := u["custom_attributes"].(string); ok && customAttrs != "" {
		sanitized["customAttributes"] = customAttrs
	} else if customClaims, ok := u["customClaims"].(map[string]interface{}); ok && len(customClaims) > 0 {
		if bytes, err := json.Marshal(customClaims); err == nil {
			sanitized["customAttributes"] = string(bytes)
		}
	}

	// providerUserInfo / provider_user_info
	if providerUserInfo, ok := u["providerUserInfo"].([]interface{}); ok {
		sanitized["providerUserInfo"] = providerUserInfo
	} else if providerUserInfo, ok := u["provider_user_info"].([]interface{}); ok {
		sanitized["providerUserInfo"] = providerUserInfo
	}

	return sanitized
}

// ImportAuthUsers imports users into Firebase Auth via Identity Toolkit accounts:upload API.
func ImportAuthUsers(token string, projectId string, req AuthImportRequest) (*AuthImportResult, error) {
	if len(req.Users) == 0 {
		return nil, fmt.Errorf("no users provided for import")
	}

	var sanitizedUsers []map[string]interface{}
	for _, u := range req.Users {
		sanitizedUsers = append(sanitizedUsers, sanitizeUserForImport(u))
	}

	var totalImported int
	var totalFailed int
	var allErrors []map[string]interface{}

	batchSize := 1000
	for i := 0; i < len(sanitizedUsers); i += batchSize {
		end := i + batchSize
		if end > len(sanitizedUsers) {
			end = len(sanitizedUsers)
		}
		batch := sanitizedUsers[i:end]

		payload := map[string]interface{}{
			"users": batch,
		}

		if req.HashAlgorithm != "" {
			payload["hashAlgorithm"] = req.HashAlgorithm
		}
		if req.SignerKey != "" {
			payload["signerKey"] = req.SignerKey
		}
		if req.SaltSeparator != "" {
			payload["saltSeparator"] = req.SaltSeparator
		} else if req.HashAlgorithm == "SCRYPT" && req.SignerKey != "" {
			payload["saltSeparator"] = "Bw=="
		}
		if req.Rounds > 0 {
			payload["rounds"] = req.Rounds
		}
		if req.MemoryCost > 0 {
			payload["memoryCost"] = req.MemoryCost
		}

		reqUrl := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/projects/%s/accounts:batchCreate", projectId)

		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal import payload: %v", err)
		}

		resBytes, err := doRequest("POST", reqUrl, token, bodyBytes)
		if err != nil && strings.Contains(err.Error(), "status 404") {
			// Fallback to accounts:upload with targetProjectId if batchCreate returns 404
			payload["targetProjectId"] = projectId
			bodyBytes, _ = json.Marshal(payload)
			fallbackUrl := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:upload?targetProjectId=%s", projectId)
			resBytes, err = doRequest("POST", fallbackUrl, token, bodyBytes)
		}

		if err != nil {
			return nil, fmt.Errorf("accounts import failed: %v", err)
		}

		var res struct {
			ErrorKind string `json:"kind"`
			Error     []struct {
				Index   int    `json:"index"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(resBytes, &res); err != nil {
			return nil, fmt.Errorf("failed to parse upload response: %v", err)
		}

		batchFailed := len(res.Error)
		batchSuccess := len(batch) - batchFailed
		totalImported += batchSuccess
		totalFailed += batchFailed

		for _, errDetail := range res.Error {
			userIndex := i + errDetail.Index
			userEmail := ""
			if userIndex < len(sanitizedUsers) {
				userEmail, _ = sanitizedUsers[userIndex]["email"].(string)
			}
			allErrors = append(allErrors, map[string]interface{}{
				"index":   userIndex,
				"email":   userEmail,
				"message": errDetail.Message,
			})
		}
	}

	return &AuthImportResult{
		Success:       totalFailed == 0,
		ImportedCount: totalImported,
		FailedCount:   totalFailed,
		Errors:        allErrors,
	}, nil
}




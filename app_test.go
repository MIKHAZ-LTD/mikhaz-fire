package main

import (
	"strings"
	"testing"
)

func TestProjects(t *testing.T) {
	app := NewApp()
	accounts := app.GetAccounts()
	t.Logf("Discovered %d accounts", len(accounts))

	allowed := []string{"minorgames", "textcup", "urearn"}

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
				// Test listing root collections
				cols, colErr := app.ListRootCollections(acc.Email, p.ProjectID, "")
				if colErr != nil {
					t.Logf("      Firestore cols error: %v", colErr)
				} else {
					t.Logf("      Firestore root collections (%d): %v", len(cols), cols)
				}
				// Test listing Auth users (limit 5)
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

package main

import (
	"context"
	"fmt"
	"strings"
)

// App struct
type App struct {
	ctx         context.Context
	credManager *CredentialManager
}

// NewApp creates a new App struct
func NewApp() *App {
	return &App{
		credManager: NewCredentialManager(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Scan credentials on startup to populate
	a.credManager.ScanLocalCredentials()
}

// GetAccounts retrieves discovered standard and service accounts.
func (a *App) GetAccounts() []Account {
	return a.credManager.ScanLocalCredentials()
}

// AddServiceAccount parses and adds a service account key in-memory.
func (a *App) AddServiceAccount(jsonContent string) (string, error) {
	return a.credManager.AddServiceAccount(jsonContent)
}

// ListProjects lists Google Cloud / Firebase projects.
func (a *App) ListProjects(email string) ([]Project, error) {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return nil, err
	}
	return ListProjects(token)
}

// ListAuthUsers retrieves users in Firebase Auth.
func (a *App) ListAuthUsers(email string, projectId string, pageSize int, pageToken string, search string) (*AuthUsersResponse, error) {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return nil, err
	}
	resp, err := ListAuthUsers(token, projectId, pageSize, pageToken)
	if err != nil {
		return nil, err
	}

	// Filter in-memory if a search query is provided
	if search != "" {
		searchQuery := strings.ToLower(search)
		var filteredUsers []AuthUser
		for _, u := range resp.Users {
			if strings.Contains(strings.ToLower(u.Email), searchQuery) ||
				strings.Contains(strings.ToLower(u.UID), searchQuery) ||
				strings.Contains(strings.ToLower(u.DisplayName), searchQuery) {
				filteredUsers = append(filteredUsers, u)
			}
		}
		resp.Users = filteredUsers
	}
	return resp, nil
}

// CreateAuthUser creates a new Firebase Auth user.
func (a *App) CreateAuthUser(email string, projectId string, userReq map[string]interface{}) (map[string]interface{}, error) {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return nil, err
	}
	return CreateAuthUser(token, projectId, userReq)
}

// UpdateAuthUser updates a Firebase Auth user.
func (a *App) UpdateAuthUser(email string, projectId string, userReq map[string]interface{}) (map[string]interface{}, error) {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return nil, err
	}
	return UpdateAuthUser(token, projectId, userReq)
}

// DeleteAuthUser deletes a Firebase Auth user.
func (a *App) DeleteAuthUser(email string, projectId string, uid string) error {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return err
	}
	return DeleteAuthUser(token, projectId, uid)
}

// ListRootCollections lists root collections in Firestore.
func (a *App) ListRootCollections(email string, projectId string, databaseId string) ([]string, error) {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return nil, err
	}
	return ListRootCollections(token, projectId, databaseId)
}

// ListDocuments retrieves documents inside a collection.
func (a *App) ListDocuments(email string, projectId string, databaseId string, collectionPath string, pageSize int, pageToken string) (map[string]interface{}, error) {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return nil, err
	}
	return ListDocuments(token, projectId, databaseId, collectionPath, pageSize, pageToken)
}

// GetDocument retrieves a single document.
func (a *App) GetDocument(email string, projectId string, databaseId string, documentPath string) (map[string]interface{}, error) {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return nil, err
	}
	return GetDocument(token, projectId, databaseId, documentPath)
}

// SaveDocument saves/overwrites/creates a document.
func (a *App) SaveDocument(email string, projectId string, databaseId string, documentPath string, fields map[string]interface{}, isNew bool) (map[string]interface{}, error) {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return nil, err
	}
	return SaveDocument(token, projectId, databaseId, documentPath, fields, isNew)
}

// DeleteDocument deletes a document (and optionally its subcollections recursively).
func (a *App) DeleteDocument(email string, projectId string, databaseId string, documentPath string, recursive bool) error {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return err
	}
	if recursive {
		return DeleteDocumentRecursive(token, projectId, databaseId, documentPath)
	}
	return DeleteDocument(token, projectId, databaseId, documentPath)
}

// RunJSScript executes a custom Javascript query script against Firestore.
func (a *App) RunJSScript(email string, projectId string, databaseId string, script string) (interface{}, error) {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return nil, err
	}
	return RunJSScript(token, projectId, databaseId, script)
}

// RunSimpleQuery runs a standard Firestore structured query.
func (a *App) RunSimpleQuery(email string, projectId string, databaseId string, structuredQuery map[string]interface{}) (interface{}, error) {
	_, token, err := a.getAuthToken(email)
	if err != nil {
		return nil, err
	}
	return RunQuery(token, projectId, databaseId, structuredQuery)
}

// CheckConflict checks if a document or collection exists in the destination before copy operations.
func (a *App) CheckConflict(destEmail string, projectId string, databaseId string, path string, nodeType string) (map[string]interface{}, error) {
	token, err := a.credManager.GetAccessToken(destEmail)
	if err != nil {
		return nil, err
	}
	isCollection := nodeType == "collection"
	conflict, message, err := CheckConflict(token, projectId, databaseId, path, isCollection)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"conflict": conflict,
		"message":  message,
	}, nil
}

// CopyNode copies a collection or document from a source Firebase project to a destination.
func (a *App) CopyNode(srcEmail string, srcProjectId string, srcDatabaseId string, destEmail string, destProjectId string, destDatabaseId string, path string, nodeType string, overwrite bool, recursive bool) (map[string]interface{}, error) {
	srcToken, err := a.credManager.GetAccessToken(srcEmail)
	if err != nil {
		return nil, fmt.Errorf("source auth failed: %v", err)
	}

	destToken, err := a.credManager.GetAccessToken(destEmail)
	if err != nil {
		return nil, fmt.Errorf("destination auth failed: %v", err)
	}

	var copiedCount int
	switch nodeType {
	case "collection":
		copiedCount, err = CopyCollection(srcToken, srcProjectId, srcDatabaseId, destToken, destProjectId, destDatabaseId, path, overwrite, recursive)
	case "document":
		copiedCount, err = CopyDocument(srcToken, srcProjectId, srcDatabaseId, destToken, destProjectId, destDatabaseId, path, overwrite, recursive)
	default:
		return nil, fmt.Errorf("invalid type: must be 'collection' or 'document'")
	}

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":     true,
		"copiedCount": copiedCount,
	}, nil
}

// Helper to get access token for a given account email.
func (a *App) getAuthToken(email string) (string, string, error) {
	if email == "" {
		return "", "", fmt.Errorf("missing email parameter")
	}
	token, err := a.credManager.GetAccessToken(email)
	if err != nil {
		return email, "", err
	}
	return email, token, nil
}

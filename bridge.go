package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RPCRequest struct {
	Method string        `json:"method"`
	Args   []interface{} `json:"args"`
}

type RPCResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func StartDevBridge(app *App, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(RPCResponse{Error: err.Error()})
			return
		}

		res, err := dispatchRPC(app, req)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(RPCResponse{Error: err.Error()})
		} else {
			json.NewEncoder(w).Encode(RPCResponse{Result: res})
		}
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return server.ListenAndServe()
}

func dispatchRPC(app *App, req RPCRequest) (interface{}, error) {
	getStr := func(idx int) string {
		if idx < len(req.Args) && req.Args[idx] != nil {
			if s, ok := req.Args[idx].(string); ok {
				return s
			}
			return fmt.Sprintf("%v", req.Args[idx])
		}
		return ""
	}
	getInt := func(idx int) int {
		if idx < len(req.Args) && req.Args[idx] != nil {
			if f, ok := req.Args[idx].(float64); ok {
				return int(f)
			}
			if i, ok := req.Args[idx].(int); ok {
				return i
			}
		}
		return 0
	}
	getBool := func(idx int) bool {
		if idx < len(req.Args) && req.Args[idx] != nil {
			if b, ok := req.Args[idx].(bool); ok {
				return b
			}
		}
		return false
	}
	getMap := func(idx int) map[string]interface{} {
		if idx < len(req.Args) && req.Args[idx] != nil {
			if m, ok := req.Args[idx].(map[string]interface{}); ok {
				return m
			}
		}
		return make(map[string]interface{})
	}

	switch req.Method {
	case "GetAccounts":
		return app.GetAccounts(), nil
	case "RemoveAccount":
		return nil, app.RemoveAccount(getStr(0))
	case "ListProjects":
		return app.ListProjects(getStr(0))
	case "ListAuthUsers":
		return app.ListAuthUsers(getStr(0), getStr(1), getInt(2), getStr(3), getStr(4))
	case "CreateAuthUser":
		return app.CreateAuthUser(getStr(0), getStr(1), getMap(2))
	case "UpdateAuthUser":
		return app.UpdateAuthUser(getStr(0), getStr(1), getMap(2))
	case "DeleteAuthUser":
		return nil, app.DeleteAuthUser(getStr(0), getStr(1), getStr(2))
	case "ExportAuthUsers":
		return app.ExportAuthUsers(getStr(0), getStr(1))
	case "ImportAuthUsers":
		var importReq AuthImportRequest
		if len(req.Args) > 2 && req.Args[2] != nil {
			bytes, _ := json.Marshal(req.Args[2])
			_ = json.Unmarshal(bytes, &importReq)
		}
		return app.ImportAuthUsers(getStr(0), getStr(1), importReq)
	case "ListRootCollections":
		return app.ListRootCollections(getStr(0), getStr(1), getStr(2))
	case "ListDocuments":
		return app.ListDocuments(getStr(0), getStr(1), getStr(2), getStr(3), getInt(4), getStr(5))
	case "GetDocument":
		return app.GetDocument(getStr(0), getStr(1), getStr(2), getStr(3))
	case "SaveDocument":
		return app.SaveDocument(getStr(0), getStr(1), getStr(2), getStr(3), getMap(4), getBool(5))
	case "DeleteDocument":
		return nil, app.DeleteDocument(getStr(0), getStr(1), getStr(2), getStr(3), getBool(4))
	case "ListDocumentSubcollections":
		return app.ListDocumentSubcollections(getStr(0), getStr(1), getStr(2), getStr(3))
	case "RunJSScript":
		return app.RunJSScript(getStr(0), getStr(1), getStr(2), getStr(3))
	case "RunSimpleQuery":
		return app.RunSimpleQuery(getStr(0), getStr(1), getStr(2), getMap(3))
	case "CheckConflict":
		return app.CheckConflict(getStr(0), getStr(1), getStr(2), getStr(3), getStr(4))
	case "CopyNode":
		return app.CopyNode(getStr(0), getStr(1), getStr(2), getStr(3), getStr(4), getStr(5), getStr(6), getStr(7), getBool(8), getBool(9))
	case "LoginGoogleAccount":
		return app.LoginGoogleAccount()
	case "AddServiceAccount":
		return app.AddServiceAccount(getStr(0))
	default:
		return nil, fmt.Errorf("unknown method: %s", req.Method)
	}
}

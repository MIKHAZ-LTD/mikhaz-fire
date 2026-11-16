package main

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

type JSQueryRunner struct {
	token      string
	projectId  string
	databaseId string
}

func NewJSQueryRunner(token, projectId, databaseId string) *JSQueryRunner {
	if databaseId == "" {
		databaseId = "(default)"
	}
	return &JSQueryRunner{
		token:      token,
		projectId:  projectId,
		databaseId: databaseId,
	}
}

// Convert native Go values to Firestore REST API typed values
func toFirestoreValue(val interface{}) map[string]interface{} {
	if val == nil {
		return map[string]interface{}{"nullValue": nil}
	}

	switch v := val.(type) {
	case bool:
		return map[string]interface{}{"booleanValue": v}
	case string:
		return map[string]interface{}{"stringValue": v}
	case int:
		return map[string]interface{}{"integerValue": fmt.Sprintf("%d", v)}
	case int32:
		return map[string]interface{}{"integerValue": fmt.Sprintf("%d", v)}
	case int64:
		return map[string]interface{}{"integerValue": fmt.Sprintf("%d", v)}
	case float32:
		if float32(int64(v)) == v {
			return map[string]interface{}{"integerValue": fmt.Sprintf("%d", int64(v))}
		}
		return map[string]interface{}{"doubleValue": float64(v)}
	case float64:
		if float64(int64(v)) == v {
			return map[string]interface{}{"integerValue": fmt.Sprintf("%d", int64(v))}
		}
		return map[string]interface{}{"doubleValue": v}
	case map[string]interface{}:
		if len(v) == 1 {
			if tStr, ok := v["__time__"].(string); ok {
				return map[string]interface{}{"timestampValue": tStr}
			}
			if gp, ok := v["__geopoint__"]; ok {
				return map[string]interface{}{"geoPointValue": gp}
			}
			if refStr, ok := v["__ref__"].(string); ok {
				return map[string]interface{}{"referenceValue": refStr}
			}
			if bStr, ok := v["__bytes__"].(string); ok {
				return map[string]interface{}{"bytesValue": bStr}
			}
		}
		fields := make(map[string]interface{})
		for k, valItem := range v {
			fields[k] = toFirestoreValue(valItem)
		}
		return map[string]interface{}{"mapValue": map[string]interface{}{"fields": fields}}
	case []interface{}:
		values := []interface{}{}
		for _, item := range v {
			values = append(values, toFirestoreValue(item))
		}
		return map[string]interface{}{"arrayValue": map[string]interface{}{"values": values}}
	default:
		return map[string]interface{}{"stringValue": fmt.Sprintf("%v", val)}
	}
}

// Convert Firestore REST API typed values to plain JS values
func fromFirestoreValue(fieldVal map[string]interface{}) interface{} {
	if fieldVal == nil {
		return nil
	}

	if val, ok := fieldVal["nullValue"]; ok {
		return val
	}
	if val, ok := fieldVal["booleanValue"]; ok {
		return val
	}
	if val, ok := fieldVal["stringValue"]; ok {
		return val
	}
	if val, ok := fieldVal["integerValue"]; ok {
		var i int64
		if str, ok := val.(string); ok {
			_, _ = fmt.Sscanf(str, "%d", &i)
		}
		return i
	}
	if val, ok := fieldVal["doubleValue"]; ok {
		return val
	}
	if val, ok := fieldVal["timestampValue"]; ok {
		return map[string]interface{}{
			"__time__": val,
		}
	}
	if val, ok := fieldVal["geoPointValue"]; ok {
		return map[string]interface{}{
			"__geopoint__": val,
		}
	}
	if val, ok := fieldVal["referenceValue"]; ok {
		return map[string]interface{}{
			"__ref__": val,
		}
	}
	if val, ok := fieldVal["bytesValue"]; ok {
		return map[string]interface{}{
			"__bytes__": val,
		}
	}
	if val, ok := fieldVal["mapValue"]; ok {
		m, ok := val.(map[string]interface{})
		if !ok {
			return nil
		}
		fields, ok := m["fields"].(map[string]interface{})
		if !ok {
			return make(map[string]interface{})
		}
		res := make(map[string]interface{})
		for k, v := range fields {
			if fv, ok := v.(map[string]interface{}); ok {
				res[k] = fromFirestoreValue(fv)
			}
		}
		return res
	}
	if val, ok := fieldVal["arrayValue"]; ok {
		arr, ok := val.(map[string]interface{})
		if !ok {
			return nil
		}
		values, ok := arr["values"].([]interface{})
		if !ok {
			return []interface{}{}
		}
		res := []interface{}{}
		for _, v := range values {
			if fv, ok := v.(map[string]interface{}); ok {
				res = append(res, fromFirestoreValue(fv))
			}
		}
		return res
	}

	return nil
}

// Document fields conversion helper
func fromFirestoreDocument(doc map[string]interface{}) map[string]interface{} {
	fields, ok := doc["fields"].(map[string]interface{})
	if !ok {
		return make(map[string]interface{})
	}
	res := make(map[string]interface{})
	for k, v := range fields {
		if fv, ok := v.(map[string]interface{}); ok {
			res[k] = fromFirestoreValue(fv)
		}
	}
	return res
}

// JS Objects mapping
type JSDocumentRef struct {
	vm     *goja.Runtime
	runner *JSQueryRunner
	Id     string `json:"id"`
	Path   string `json:"path"`
}

type JSDocumentSnapshot struct {
	Id     string                        `json:"id"`
	Path   string                        `json:"path"`
	Exists bool                          `json:"exists"`
	Ref    *JSDocumentRef                `json:"ref"`
	Data   func() map[string]interface{} `json:"data"`
}

type JSQuerySnapshot struct {
	vm    *goja.Runtime
	Docs  []*JSDocumentSnapshot `json:"docs"`
	Size  int                   `json:"size"`
	Empty bool                  `json:"empty"`
}

func (s *JSQuerySnapshot) ForEach(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		return goja.Undefined()
	}
	fn, ok := goja.AssertFunction(call.Arguments[0])
	if !ok {
		return goja.Undefined()
	}
	for _, doc := range s.Docs {
		_, err := fn(goja.Undefined(), s.vm.ToValue(doc))
		if err != nil {
			panic(err)
		}
	}
	return goja.Undefined()
}

type JSCollectionRef struct {
	vm     *goja.Runtime
	runner *JSQueryRunner
	Path   string
}

type JSQueryRef struct {
	vm             *goja.Runtime
	runner         *JSQueryRunner
	collPath       string
	allDescendants bool
	limit          int
	filters        []map[string]interface{}
	sorts          []map[string]interface{}
}

func (c *JSCollectionRef) Doc(id string) *JSDocumentRef {
	docPath := strings.TrimRight(c.Path, "/") + "/" + strings.TrimLeft(id, "/")
	return &JSDocumentRef{
		vm:     c.vm,
		runner: c.runner,
		Id:     id,
		Path:   docPath,
	}
}

func (c *JSCollectionRef) Limit(n int) *JSQueryRef {
	q := &JSQueryRef{vm: c.vm, runner: c.runner, collPath: c.Path}
	return q.Limit(n)
}

func (c *JSCollectionRef) Where(field, op string, val interface{}) *JSQueryRef {
	q := &JSQueryRef{vm: c.vm, runner: c.runner, collPath: c.Path}
	return q.Where(field, op, val)
}

func (c *JSCollectionRef) OrderBy(field, dir string) *JSQueryRef {
	q := &JSQueryRef{vm: c.vm, runner: c.runner, collPath: c.Path}
	return q.OrderBy(field, dir)
}

func (c *JSCollectionRef) Get() goja.Value {
	q := &JSQueryRef{vm: c.vm, runner: c.runner, collPath: c.Path}
	return q.Get()
}

func (q *JSQueryRef) clone() *JSQueryRef {
	filtersCopy := make([]map[string]interface{}, len(q.filters))
	copy(filtersCopy, q.filters)
	sortsCopy := make([]map[string]interface{}, len(q.sorts))
	copy(sortsCopy, q.sorts)
	return &JSQueryRef{
		vm:             q.vm,
		runner:         q.runner,
		collPath:       q.collPath,
		allDescendants: q.allDescendants,
		limit:          q.limit,
		filters:        filtersCopy,
		sorts:          sortsCopy,
	}
}

func (q *JSQueryRef) Limit(n int) *JSQueryRef {
	next := q.clone()
	next.limit = n
	return next
}

func (q *JSQueryRef) Where(field, op string, val interface{}) *JSQueryRef {
	opMap := map[string]string{
		"==":                 "EQUAL",
		"=":                  "EQUAL",
		"<":                  "LESS_THAN",
		"<=":                 "LESS_THAN_OR_EQUAL",
		">":                  "GREATER_THAN",
		">=":                 "GREATER_THAN_OR_EQUAL",
		"array-contains":     "ARRAY_CONTAINS",
		"in":                 "IN",
		"array-contains-any": "ARRAY_CONTAINS_ANY",
		"not-in":             "NOT_IN",
		"!=":                 "NOT_EQUAL",
	}

	mappedOp, ok := opMap[op]
	if !ok {
		mappedOp = op
	}

	var filter map[string]interface{}
	if val == nil && (mappedOp == "EQUAL" || mappedOp == "NOT_EQUAL") {
		unaryOp := "IS_NULL"
		if mappedOp == "NOT_EQUAL" {
			unaryOp = "IS_NOT_NULL"
		}
		filter = map[string]interface{}{
			"unaryFilter": map[string]interface{}{
				"field": map[string]interface{}{
					"fieldPath": field,
				},
				"op": unaryOp,
			},
		}
	} else {
		filter = map[string]interface{}{
			"fieldFilter": map[string]interface{}{
				"field": map[string]interface{}{
					"fieldPath": field,
				},
				"op":    mappedOp,
				"value": toFirestoreValue(val),
			},
		}
	}

	next := q.clone()
	next.filters = append(next.filters, filter)
	return next
}

func (q *JSQueryRef) OrderBy(field, dir string) *JSQueryRef {
	direction := "ASCENDING"
	upperDir := strings.ToUpper(strings.TrimSpace(dir))
	if upperDir == "DESC" || upperDir == "DESCENDING" {
		direction = "DESCENDING"
	}
	sort := map[string]interface{}{
		"field": map[string]interface{}{
			"fieldPath": field,
		},
		"direction": direction,
	}
	next := q.clone()
	next.sorts = append(next.sorts, sort)
	return next
}

func (q *JSQueryRef) Get() goja.Value {
	promise, resolve, reject := q.vm.NewPromise()

	// Execute synchronously in the same thread
	docs, err := q.execute()
	if err != nil {
		_ = reject(q.vm.ToValue(err.Error()))
	} else {
		_ = resolve(q.vm.ToValue(docs))
	}

	return q.vm.ToValue(promise)
}

func (q *JSQueryRef) buildSnapshotDoc(docMap map[string]interface{}) *JSDocumentSnapshot {
	docName, _ := docMap["name"].(string)
	parts := strings.Split(docName, "/")
	docId := ""
	if len(parts) > 0 {
		docId = parts[len(parts)-1]
	}

	docRelPath := q.collPath + "/" + docId
	if idx := strings.Index(docName, "/documents/"); idx != -1 {
		docRelPath = docName[idx+len("/documents/"):]
	}

	ref := &JSDocumentRef{
		vm:     q.vm,
		runner: q.runner,
		Id:     docId,
		Path:   docRelPath,
	}

	fieldsData := fromFirestoreDocument(docMap)
	return &JSDocumentSnapshot{
		Id:     docId,
		Path:   docRelPath,
		Exists: true,
		Ref:    ref,
		Data: func() map[string]interface{} {
			return fieldsData
		},
	}
}

func (q *JSQueryRef) execute() (*JSQuerySnapshot, error) {
	if !q.allDescendants && len(q.filters) == 0 && len(q.sorts) == 0 {
		resp, err := ListDocuments(q.runner.token, q.runner.projectId, q.runner.databaseId, q.collPath, q.limit, "")
		if err != nil {
			return nil, err
		}

		rawDocs, _ := resp["documents"].([]interface{})
		snapDocs := []*JSDocumentSnapshot{}
		for _, rd := range rawDocs {
			docMap, ok := rd.(map[string]interface{})
			if !ok {
				continue
			}
			snapDocs = append(snapDocs, q.buildSnapshotDoc(docMap))
		}
		return &JSQuerySnapshot{
			vm:    q.vm,
			Docs:  snapDocs,
			Size:  len(snapDocs),
			Empty: len(snapDocs) == 0,
		}, nil
	}

	var filter map[string]interface{}
	if len(q.filters) == 1 {
		filter = q.filters[0]
	} else if len(q.filters) > 1 {
		filter = map[string]interface{}{
			"compositeFilter": map[string]interface{}{
				"op":      "AND",
				"filters": q.filters,
			},
		}
	}

	fromEntry := map[string]interface{}{
		"collectionId": q.collPath,
	}
	if q.allDescendants {
		fromEntry["allDescendants"] = true
	}

	structuredQuery := map[string]interface{}{
		"from": []map[string]interface{}{fromEntry},
	}

	if filter != nil {
		structuredQuery["where"] = filter
	}
	if len(q.sorts) > 0 {
		structuredQuery["orderBy"] = q.sorts
	}
	if q.limit > 0 {
		structuredQuery["limit"] = q.limit
	}

	body := map[string]interface{}{
		"structuredQuery": structuredQuery,
	}

	res, err := RunQuery(q.runner.token, q.runner.projectId, q.runner.databaseId, body)
	if err != nil {
		return nil, err
	}

	results, ok := res.([]interface{})
	if !ok {
		return &JSQuerySnapshot{
			vm:    q.vm,
			Docs:  []*JSDocumentSnapshot{},
			Size:  0,
			Empty: true,
		}, nil
	}

	snapDocs := []*JSDocumentSnapshot{}
	for _, r := range results {
		resMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		docMap, ok := resMap["document"].(map[string]interface{})
		if !ok {
			continue
		}
		snapDocs = append(snapDocs, q.buildSnapshotDoc(docMap))
	}

	return &JSQuerySnapshot{
		vm:    q.vm,
		Docs:  snapDocs,
		Size:  len(snapDocs),
		Empty: len(snapDocs) == 0,
	}, nil
}

func (d *JSDocumentRef) Collection(subCol string) *JSCollectionRef {
	return &JSCollectionRef{
		vm:     d.vm,
		runner: d.runner,
		Path:   strings.TrimRight(d.Path, "/") + "/" + strings.TrimLeft(subCol, "/"),
	}
}

func (d *JSDocumentRef) Get() goja.Value {
	promise, resolve, reject := d.vm.NewPromise()

	docMap, err := GetDocument(d.runner.token, d.runner.projectId, d.runner.databaseId, d.Path)
	if err != nil {
		_ = reject(d.vm.ToValue(err.Error()))
	} else {
		fieldsData := fromFirestoreDocument(docMap)
		snap := &JSDocumentSnapshot{
			Id:     d.Id,
			Path:   d.Path,
			Exists: true,
			Ref:    d,
			Data: func() map[string]interface{} {
				return fieldsData
			},
		}
		_ = resolve(d.vm.ToValue(snap))
	}

	return d.vm.ToValue(promise)
}

func (d *JSDocumentRef) Set(data map[string]interface{}) goja.Value {
	promise, resolve, reject := d.vm.NewPromise()

	fields := make(map[string]interface{})
	for k, v := range data {
		fields[k] = toFirestoreValue(v)
	}
	_, err := SaveDocument(d.runner.token, d.runner.projectId, d.runner.databaseId, d.Path, fields, false)
	if err != nil {
		_ = reject(d.vm.ToValue(err.Error()))
	} else {
		_ = resolve(goja.Undefined())
	}

	return d.vm.ToValue(promise)
}

func (d *JSDocumentRef) Update(data map[string]interface{}) goja.Value {
	return d.Set(data)
}

func (d *JSDocumentRef) Delete() goja.Value {
	promise, resolve, reject := d.vm.NewPromise()

	err := DeleteDocument(d.runner.token, d.runner.projectId, d.runner.databaseId, d.Path)
	if err != nil {
		_ = reject(d.vm.ToValue(err.Error()))
	} else {
		_ = resolve(goja.Undefined())
	}

	return d.vm.ToValue(promise)
}

func formatConsoleArgs(args []goja.Value) string {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == nil || goja.IsUndefined(arg) {
			parts = append(parts, "undefined")
			continue
		}
		if goja.IsNull(arg) {
			parts = append(parts, "null")
			continue
		}
		exported := arg.Export()
		if _, ok := exported.(string); ok {
			parts = append(parts, fmt.Sprintf("%v", exported))
		} else {
			parts = append(parts, fmt.Sprintf("%+v", cleanGojaExport(exported)))
		}
	}
	return strings.Join(parts, " ")
}

// Main execution function
func RunJSScript(token, projectId, databaseId, script string) (interface{}, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	runner := NewJSQueryRunner(token, projectId, databaseId)
	var logs []string

	// Inject console object
	consoleObj := vm.NewObject()
	appendLog := func(prefix string) func(goja.FunctionCall) goja.Value {
		return func(call goja.FunctionCall) goja.Value {
			msg := formatConsoleArgs(call.Arguments)
			if prefix != "" {
				msg = prefix + " " + msg
			}
			logs = append(logs, msg)
			return goja.Undefined()
		}
	}
	_ = consoleObj.Set("log", appendLog(""))
	_ = consoleObj.Set("info", appendLog("[INFO]"))
	_ = consoleObj.Set("warn", appendLog("[WARN]"))
	_ = consoleObj.Set("error", appendLog("[ERROR]"))
	_ = vm.Set("console", consoleObj)

	// Inject global db object
	dbObj := vm.NewObject()
	_ = dbObj.Set("collection", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("collection() requires a collection path"))
		}
		path := strings.Trim(call.Arguments[0].String(), "/")
		return vm.ToValue(&JSCollectionRef{
			vm:     vm,
			runner: runner,
			Path:   path,
		})
	})
	_ = dbObj.Set("doc", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("doc() requires a document path"))
		}
		docPath := strings.Trim(call.Arguments[0].String(), "/")
		parts := strings.Split(docPath, "/")
		docId := ""
		if len(parts) > 0 {
			docId = parts[len(parts)-1]
		}
		return vm.ToValue(&JSDocumentRef{
			vm:     vm,
			runner: runner,
			Id:     docId,
			Path:   docPath,
		})
	})
	_ = dbObj.Set("collectionGroup", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("collectionGroup() requires a collection ID"))
		}
		colId := strings.Trim(call.Arguments[0].String(), "/")
		return vm.ToValue(&JSQueryRef{
			vm:             vm,
			runner:         runner,
			collPath:       colId,
			allDescendants: true,
		})
	})
	_ = vm.Set("db", dbObj)

	// Evaluate script
	_, err := vm.RunString(script)
	if err != nil {
		return nil, fmt.Errorf("JS evaluation error: %v", err)
	}

	runFn, ok := goja.AssertFunction(vm.Get("run"))
	if !ok {
		return nil, fmt.Errorf("run() function is not defined. It must be: async function run() { ... }")
	}

	resVal, err := runFn(goja.Undefined())
	if err != nil {
		return nil, fmt.Errorf("JS execution error: %v", err)
	}

	var exported interface{}
	// If the return value is a Promise, extract its result directly
	if promise, ok := resVal.Export().(*goja.Promise); ok {
		if promise.State() == goja.PromiseStateFulfilled {
			exported = cleanGojaExport(promise.Result().Export())
		} else if promise.State() == goja.PromiseStateRejected {
			return nil, fmt.Errorf("Promise rejected: %v", promise.Result().Export())
		} else {
			return nil, fmt.Errorf("Promise is pending (synchronous REST execution failed to resolve)")
		}
	} else {
		exported = cleanGojaExport(resVal.Export())
	}

	if len(logs) > 0 {
		if m, ok := exported.(map[string]interface{}); ok {
			if _, hasLogs := m["logs"]; hasLogs {
				return m, nil
			}
		}
		return map[string]interface{}{
			"logs":   logs,
			"result": exported,
		}, nil
	}

	return exported, nil
}

// Clean Goja export types to avoid structures that json.Marshal cannot handle (like functions)
func cleanGojaExport(val interface{}) interface{} {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case *JSQuerySnapshot:
		docs := []map[string]interface{}{}
		for _, d := range v.Docs {
			docs = append(docs, map[string]interface{}{
				"id":   d.Id,
				"path": d.Ref.Path,
				"data": d.Data(),
			})
		}
		return docs
	case []*JSDocumentSnapshot:
		docs := []map[string]interface{}{}
		for _, d := range v {
			docs = append(docs, map[string]interface{}{
				"id":   d.Id,
				"path": d.Ref.Path,
				"data": d.Data(),
			})
		}
		return docs
	case *JSDocumentSnapshot:
		return map[string]interface{}{
			"id":   v.Id,
			"path": v.Ref.Path,
			"data": v.Data(),
		}
	case *JSDocumentRef:
		return map[string]interface{}{
			"id":   v.Id,
			"path": v.Path,
		}
	case map[string]interface{}:
		res := make(map[string]interface{})
		for k, valItem := range v {
			res[k] = cleanGojaExport(valItem)
		}
		return res
	case []interface{}:
		res := []interface{}{}
		for _, item := range v {
			res = append(res, cleanGojaExport(item))
		}
		return res
	default:
		return val
	}
}

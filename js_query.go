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
	case int64:
		return map[string]interface{}{"integerValue": fmt.Sprintf("%d", v)}
	case float64:
		if float64(int64(v)) == v {
			return map[string]interface{}{"integerValue": fmt.Sprintf("%d", int64(v))}
		}
		return map[string]interface{}{"doubleValue": v}
	case map[string]interface{}:
		if tStr, ok := v["__time__"].(string); ok && len(v) == 1 {
			return map[string]interface{}{"timestampValue": tStr}
		}
		fields := make(map[string]interface{})
		for k, valItem := range v {
			fields[k] = toFirestoreValue(valItem)
		}
		return map[string]interface{}{"mapValue": map[string]interface{}{"fields": fields}}
	case []interface{}:
		var values []interface{}
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
		_, _ = fmt.Sscanf(val.(string), "%d", &i)
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
	if val, ok := fieldVal["referenceValue"]; ok {
		return val
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
		var res []interface{}
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
	Id   string                 `json:"id"`
	Ref  *JSDocumentRef         `json:"ref"`
	Data func() map[string]interface{} `json:"data"`
}

type JSQuerySnapshot struct {
	Docs []*JSDocumentSnapshot `json:"docs"`
}

type JSCollectionRef struct {
	vm     *goja.Runtime
	runner *JSQueryRunner
	Path   string
}

type JSQueryRef struct {
	vm       *goja.Runtime
	runner   *JSQueryRunner
	collPath string
	limit    int
	filters  []map[string]interface{}
	sorts    []map[string]interface{}
}

func (c *JSCollectionRef) Doc(id string) *JSDocumentRef {
	docPath := c.Path + "/" + id
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

func (q *JSQueryRef) Limit(n int) *JSQueryRef {
	q.limit = n
	return q
}

func (q *JSQueryRef) Where(field, op string, val interface{}) *JSQueryRef {
	opMap := map[string]string{
		"==":             "EQUAL",
		"=":              "EQUAL",
		"<":              "LESS_THAN",
		"<=":             "LESS_THAN_OR_EQUAL",
		">":              "GREATER_THAN",
		">=":             "GREATER_THAN_OR_EQUAL",
		"array-contains": "ARRAY_CONTAINS",
		"in":             "IN",
		"array-contains-any": "ARRAY_CONTAINS_ANY",
		"not-in":         "NOT_IN",
		"!=":             "NOT_EQUAL",
	}

	mappedOp, ok := opMap[op]
	if !ok {
		mappedOp = op
	}

	filter := map[string]interface{}{
		"fieldFilter": map[string]interface{}{
			"field": map[string]interface{}{
				"fieldPath": field,
			},
			"op":    mappedOp,
			"value": toFirestoreValue(val),
		},
	}
	q.filters = append(q.filters, filter)
	return q
}

func (q *JSQueryRef) OrderBy(field, dir string) *JSQueryRef {
	direction := "ASCENDING"
	if strings.ToUpper(dir) == "DESC" || strings.ToUpper(dir) == "DESCENDING" {
		direction = "DESCENDING"
	}
	sort := map[string]interface{}{
		"field": map[string]interface{}{
			"fieldPath": field,
		},
		"direction": direction,
	}
	q.sorts = append(q.sorts, sort)
	return q
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

func (q *JSQueryRef) execute() (*JSQuerySnapshot, error) {
	if len(q.filters) == 0 && len(q.sorts) == 0 {
		resp, err := ListDocuments(q.runner.token, q.runner.projectId, q.runner.databaseId, q.collPath, q.limit, "")
		if err != nil {
			return nil, err
		}

		rawDocs, _ := resp["documents"].([]interface{})
		var snapDocs []*JSDocumentSnapshot
		for _, rd := range rawDocs {
			docMap, ok := rd.(map[string]interface{})
			if !ok {
				continue
			}
			docName, _ := docMap["name"].(string)
			parts := strings.Split(docName, "/")
			docId := parts[len(parts)-1]

			ref := &JSDocumentRef{
				vm:     q.vm,
				runner: q.runner,
				Id:     docId,
				Path:   q.collPath + "/" + docId,
			}
			
			fieldsData := fromFirestoreDocument(docMap)
			snapDocs = append(snapDocs, &JSDocumentSnapshot{
				Id:  docId,
				Ref: ref,
				Data: func() map[string]interface{} {
					return fieldsData
				},
			})
		}
		return &JSQuerySnapshot{Docs: snapDocs}, nil
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

	structuredQuery := map[string]interface{}{
		"from": []map[string]interface{}{
			{
				"collectionId": q.collPath,
			},
		},
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
		return &JSQuerySnapshot{Docs: []*JSDocumentSnapshot{}}, nil
	}

	var snapDocs []*JSDocumentSnapshot
	for _, r := range results {
		resMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		docMap, ok := resMap["document"].(map[string]interface{})
		if !ok {
			continue
		}

		docName, _ := docMap["name"].(string)
		parts := strings.Split(docName, "/")
		docId := parts[len(parts)-1]

		ref := &JSDocumentRef{
			vm:     q.vm,
			runner: q.runner,
			Id:     docId,
			Path:   q.collPath + "/" + docId,
		}

		fieldsData := fromFirestoreDocument(docMap)
		snapDocs = append(snapDocs, &JSDocumentSnapshot{
			Id:  docId,
			Ref: ref,
			Data: func() map[string]interface{} {
				return fieldsData
			},
		})
	}

	return &JSQuerySnapshot{Docs: snapDocs}, nil
}

func (d *JSDocumentRef) Get() goja.Value {
	promise, resolve, reject := d.vm.NewPromise()
	
	docMap, err := GetDocument(d.runner.token, d.runner.projectId, d.runner.databaseId, d.Path)
	if err != nil {
		_ = reject(d.vm.ToValue(err.Error()))
	} else {
		fieldsData := fromFirestoreDocument(docMap)
		snap := &JSDocumentSnapshot{
			Id:  d.Id,
			Ref: d,
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

// Main execution function
func RunJSScript(token, projectId, databaseId, script string) (interface{}, error) {
	vm := goja.New()

	runner := NewJSQueryRunner(token, projectId, databaseId)

	// Inject global db object
	dbObj := vm.NewObject()
	_ = dbObj.Set("collection", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("collection() requires a collection path"))
		}
		path := call.Arguments[0].String()
		return vm.ToValue(&JSCollectionRef{
			vm:     vm,
			runner: runner,
			Path:   path,
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

	// If the return value is a Promise, extract its result directly
	if promise, ok := resVal.Export().(*goja.Promise); ok {
		if promise.State() == goja.PromiseStateFulfilled {
			return cleanGojaExport(promise.Result().Export()), nil
		} else if promise.State() == goja.PromiseStateRejected {
			return nil, fmt.Errorf("Promise rejected: %v", promise.Result().Export())
		} else {
			return nil, fmt.Errorf("Promise is pending (synchronous REST execution failed to resolve)")
		}
	}

	return cleanGojaExport(resVal.Export()), nil
}

// Clean Goja export types to avoid structures that json.Marshal cannot handle (like functions)
func cleanGojaExport(val interface{}) interface{} {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case *JSQuerySnapshot:
		var docs []map[string]interface{}
		for _, d := range v.Docs {
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
	case map[string]interface{}:
		res := make(map[string]interface{})
		for k, valItem := range v {
			res[k] = cleanGojaExport(valItem)
		}
		return res
	case []interface{}:
		var res []interface{}
		for _, item := range v {
			res = append(res, cleanGojaExport(item))
		}
		return res
	default:
		return val
	}
}

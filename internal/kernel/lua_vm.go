/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package kernel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/vayload/vayload/internal/shared/ds"
	httpi "github.com/vayload/vayload/pkg/http"
	lua "github.com/yuin/gopher-lua"
)

type ScriptVMEngine interface {
	PrepareSandbox()
	Eval(script string) error
	LoadFile(filename string) error

	// For binding more objects or function to cm
	Bindings(fn func(*lua.LState))
	Close()
}

type ScriptModule interface {
	Load(lua *lua.LState) int
}

type scriptVMEngine struct {
	state *lua.LState
	mu    sync.RWMutex

	namespaces *ds.HashMap[string, *lua.LState]
}

func NewScriptingEngine() *scriptVMEngine {
	L := lua.NewState(lua.Options{
		SkipOpenLibs:        true,
		IncludeGoStackTrace: true,
	})

	// Open ONLY safe base libraries
	// 'base', 'table', 'string', 'math' are usually safe.
	// 'io', 'os', 'package', 'debug' are UNSAFE.
	safeLibs := map[string]lua.LGFunction{
		lua.LoadLibName:   lua.OpenPackage, // Needed for require, but we must restrict it
		lua.BaseLibName:   lua.OpenBase,
		lua.TabLibName:    lua.OpenTable,
		lua.StringLibName: lua.OpenString,
		lua.MathLibName:   lua.OpenMath,
	}

	for name, lib := range safeLibs {
		L.Push(L.NewFunction(lib))
		L.Push(lua.LString(name))
		L.Call(1, 0)
	}

	modules := map[string]ScriptModule{
		"vayload:http": newHttpModule(),
	}

	for name, module := range modules {
		L.PreloadModule(name, module.Load)
	}

	injectJSONFunctions(L)

	return &scriptVMEngine{
		state:      L,
		namespaces: ds.NewHashMap[string, *lua.LState](),
	}
}

func (e *scriptVMEngine) Close() {
	e.state.Close()
}

// PrepareSandbox removes any restricted globals if open libs leaked them,
// and ensures only whitelisted APIs are available.
func (e *scriptVMEngine) PrepareSandbox() {
	e.state.SetGlobal("os", lua.LNil) // Ensure OS is not reachable
	e.state.SetGlobal("io", lua.LNil) // Ensure IO is not reachable
}

func (e *scriptVMEngine) Bindings(fn func(*lua.LState)) {
	fn(e.state)
}

// Eval runs a script string
func (e *scriptVMEngine) Eval(script string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.state.DoString(script)
}

// LoadFile loads a script from file (used by internal loader)
func (e *scriptVMEngine) LoadFile(path string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.state.DoFile(path)
}

func (e *scriptVMEngine) GetState() *lua.LState {
	return e.state
}

type httpModule struct {
	client *httpi.HttpClient
}

func newHttpModule() *httpModule {
	return &httpModule{
		client: httpi.NewHttpClient(),
	}
}

func (m *httpModule) Load(L *lua.LState) int {
	exports := map[string]lua.LGFunction{
		"get": func(L *lua.LState) int {
			url := L.CheckString(1)
			headers := m.extractHeaders(L, 2)

			resp, err := m.makeHTTPRequest("GET", url, "", headers)
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}

			responseTable := m.createHTTPResponseTable(L, resp)
			L.Push(responseTable)
			L.Push(lua.LNil)
			return 2
		},
		"post": func(L *lua.LState) int {
			url := L.CheckString(1)
			body := L.CheckString(2)

			headers := m.extractHeaders(L, 4)

			resp, err := m.makeHTTPRequest("POST", url, body, headers)
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}

			responseTable := m.createHTTPResponseTable(L, resp)
			L.Push(responseTable)
			L.Push(lua.LNil)
			return 2
		},
		"put": func(L *lua.LState) int {
			url := L.CheckString(1)
			body := L.CheckString(2)

			headers := m.extractHeaders(L, 4)

			resp, err := m.makeHTTPRequest("PUT", url, body, headers)
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}

			responseTable := m.createHTTPResponseTable(L, resp)
			L.Push(responseTable)
			L.Push(lua.LNil)
			return 2
		},
		"delete": func(L *lua.LState) int {
			url := L.CheckString(1)
			headers := m.extractHeaders(L, 2)

			resp, err := m.makeHTTPRequest("DELETE", url, "", headers)
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}

			responseTable := m.createHTTPResponseTable(L, resp)
			L.Push(responseTable)
			L.Push(lua.LNil)
			return 2
		},
		"patch": func(L *lua.LState) int {
			url := L.CheckString(1)
			body := L.CheckString(2)

			headers := m.extractHeaders(L, 4)

			resp, err := m.makeHTTPRequest("PATCH", url, body, headers)
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}

			responseTable := m.createHTTPResponseTable(L, resp)
			L.Push(responseTable)
			L.Push(lua.LNil)
			return 2
		},
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	L.SetField(mod, "name", lua.LString("value"))

	L.Push(mod)
	return 1
}

func (mod *httpModule) makeHTTPRequest(method, url, body string, headers map[string]string) (*http.Response, error) {
	ctx := context.Background()
	switch method {
	case "GET":
		return mod.client.Get(ctx, url, httpi.HttpClientConfig{
			Headers: headers,
		})
	case "POST":
		return mod.client.Post(ctx, url, []byte(body), httpi.HttpClientConfig{
			Headers: headers,
		})
	case "PUT":
		return mod.client.Put(ctx, url, []byte(body), httpi.HttpClientConfig{
			Headers: headers,
		})
	case "DELETE":
		return mod.client.Delete(ctx, url, httpi.HttpClientConfig{
			Headers: headers,
		})
	case "PATCH":
		return mod.client.Patch(ctx, url, []byte(body), httpi.HttpClientConfig{
			Headers: headers,
		})
	default:
		return nil, fmt.Errorf("método HTTP no soportado: %s", method)
	}
}

func (mod *httpModule) extractHeaders(L *lua.LState, position int) map[string]string {
	if L.GetTop() < position || L.Get(position).Type() != lua.LTTable {
		return nil
	}

	table := L.CheckTable(position)
	headers := make(map[string]string)

	table.ForEach(func(key, value lua.LValue) {
		headers[key.String()] = value.String()
	})

	return headers
}

func (mod *httpModule) createHTTPResponseTable(L *lua.LState, resp *http.Response) *lua.LTable {
	responseTable := L.NewTable()
	L.SetField(responseTable, "status", lua.LNumber(resp.StatusCode))
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return responseTable
	}
	L.SetField(responseTable, "body", lua.LString(string(body)))

	headersTable := L.NewTable()
	for key, values := range resp.Header {
		if len(values) > 0 {
			L.SetField(headersTable, key, lua.LString(values[0]))
		}
	}
	L.SetField(responseTable, "headers", headersTable)

	return responseTable
}

// Global functions
func injectJSONFunctions(L *lua.LState) {
	jsonTable := L.NewTable()

	// Json.encode(value, space?, replacer?)
	L.SetField(jsonTable, "encode", L.NewFunction(func(L *lua.LState) int {
		val := L.CheckAny(1)

		var replacer *lua.LFunction
		var space string

		if L.GetTop() >= 2 && L.Get(2) != lua.LNil {
			switch s := L.Get(2).(type) {
			case lua.LNumber:
				space = strings.Repeat(" ", int(s))
			case lua.LString:
				space = string(s)
			}
		}

		if L.GetTop() >= 3 && L.Get(3) != lua.LNil {
			if fn, ok := L.Get(3).(*lua.LFunction); ok {
				replacer = fn
			}
		}

		data := luaValueToInterface(val)

		var jsonBytes []byte
		var err error
		if space != "" {
			jsonBytes, err = json.MarshalIndent(data, "", space)
		} else {
			jsonBytes, err = json.Marshal(data)
		}

		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		jsonStr := string(jsonBytes)

		if replacer != nil {
			// Todo: Replace implementation
		}

		L.Push(lua.LString(jsonStr))
		return 1
	}))

	// Json.decode(text, reviver?)
	L.SetField(jsonTable, "decode", L.NewFunction(func(L *lua.LState) int {
		jsonStr := L.CheckString(1)

		var reviver *lua.LFunction
		if L.GetTop() >= 2 && L.Get(2) != lua.LNil {
			if fn, ok := L.Get(2).(*lua.LFunction); ok {
				reviver = fn
			}
		}

		var data any
		err := json.Unmarshal([]byte(jsonStr), &data)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		result := interfaceToLuaValue(L, data)

		if reviver != nil {
			L.Push(reviver)
			L.Push(lua.LNil) // key
			L.Push(result)   // value
			if err := L.PCall(2, 1, nil); err == nil {
				result = L.Get(-1)
			}
		}

		L.Push(result)
		return 1
	}))

	L.SetGlobal("Json", jsonTable)
}

func mapToLuaTable(L *lua.LState, m map[string]any) *lua.LTable {
	table := L.NewTable()
	for key, value := range m {
		luaValue := interfaceToLuaValue(L, value)
		L.SetField(table, key, luaValue)
	}

	return table
}

func interfaceToLuaValue(L *lua.LState, value any) lua.LValue {
	if value == nil {
		return lua.LNil
	}

	switch v := value.(type) {
	case bool:
		return lua.LBool(v)
	case int:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []byte:
		return lua.LString(string(v))
	case time.Time:
		return lua.LString(v.Format(time.RFC3339))
	case map[string]any:
		return mapToLuaTable(L, v)
	case []any:
		table := L.NewTable()
		for i, item := range v {
			L.RawSetInt(table, i+1, interfaceToLuaValue(L, item))
		}
		return table
	default:
		return lua.LString(fmt.Sprintf("%v", v))
	}
}

func tableToInterface(tbl *lua.LTable) any {
	isArray := true
	maxIndex := 0
	count := 0
	arrayElems := make([]any, 0)

	tbl.ForEach(func(key, value lua.LValue) {
		count++
		if key.Type() == lua.LTNumber {
			i := int(key.(lua.LNumber))
			if i > maxIndex {
				maxIndex = i
			}
			arrayElems = append(arrayElems, luaValueToInterface(value))
		} else {
			isArray = false
		}
	})

	// Si es array "limpio" (1..N)
	if isArray && count == maxIndex {
		return arrayElems
	}

	// Caso objeto (incluye mezcla)
	obj := make(map[string]any)
	tbl.ForEach(func(key, value lua.LValue) {
		obj[key.String()] = luaValueToInterface(value)
	})
	return obj
}

func luaValueToInterface(value lua.LValue) any {
	switch v := value.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return string(v)
	case *lua.LTable:
		return tableToInterface(v)
	default:
		return v.String()
	}
}

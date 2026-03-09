package performance

import (
	"fmt"
	"testing"

	"github.com/Shopify/go-lua"
	"github.com/d5/tengo/v2"
	"github.com/dop251/goja"
	glua "github.com/yuin/gopher-lua"
)

// Datos que el sistema Go le pasa al script
type User struct {
	Name  string
	Level int
}

// Función de Go que el script debe invocar
func GoCheckPermission(level int) bool {
	return level > 10
}

// ─────────────────────────────────────────────────────────────────────────────
// BENCHMARK: INTEROPERABILIDAD REAL
// ─────────────────────────────────────────────────────────────────────────────

func BenchmarkPlugin_GopherLua(b *testing.B) {
	script := `if checkPermission(user.Level) then result = "ok" else result = "denied" end`
	for i := 0; i < b.N; i++ {
		L := glua.NewState()
		// Inyectar Datos
		u := L.NewTable()
		L.SetField(u, "Level", glua.LNumber(15))
		L.SetGlobal("user", u)
		// Inyectar Función
		L.SetGlobal("checkPermission", L.NewFunction(func(L *glua.LState) int {
			lvl := L.CheckInt(1)
			L.Push(glua.LBool(GoCheckPermission(lvl)))
			return 1
		}))
		L.DoString(script)
		_ = L.GetGlobal("result")
		L.Close()
	}
}

func BenchmarkPlugin_Goja(b *testing.B) {
	script := `if (checkPermission(user.Level)) { result = "ok"; } else { result = "denied"; }`
	for i := 0; i < b.N; i++ {
		vm := goja.New()
		// Goja usa reflexión, es más "automático"
		vm.Set("user", User{Level: 15})
		vm.Set("checkPermission", GoCheckPermission)
		vm.RunString(script)
		_ = vm.Get("result")
	}
}

func objectToInt(obj tengo.Object) (int, error) {
	switch o := obj.(type) {
	case *tengo.Int:
		return int(o.Value), nil
	case *tengo.Float:
		return int(o.Value), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", obj)
	}
}

func BenchmarkPlugin_Tengo(b *testing.B) {
	script := `if checkPermission(user.Level) { result = "ok" } else { result = "denied" }`
	for i := 0; i < b.N; i++ {
		s := tengo.NewScript([]byte(script))
		s.Add("user", map[string]interface{}{"Level": 15})
		s.Add("checkPermission", func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) == 0 {
				return nil, fmt.Errorf("expected at least one argument")
			}

			lvl, _ := objectToInt(args[0])
			res := GoCheckPermission(int(lvl))
			if res {
				return tengo.TrueValue, nil
			}
			return tengo.FalseValue, nil
		})
		c, _ := s.Run()
		_ = c.Get("result")
	}
}

func BenchmarkPlugin_GoLua_Shopify(b *testing.B) {
	script := `if checkPermission(user.Level) then result = "ok" else result = "denied" end`
	for i := 0; i < b.N; i++ {
		L := lua.NewState()
		lua.BaseOpen(L)
		// En Go-Lua la manipulación de la pila es más manual (estilo C)
		L.NewTable()
		L.PushInteger(15)
		L.SetField(-2, "Level")
		L.SetGlobal("user")

		L.PushGoFunction(func(L *lua.State) int {
			lvl, _ := L.ToInteger(1)
			L.PushBoolean(GoCheckPermission(lvl))
			return 1
		})
		L.SetGlobal("checkPermission")

		lua.DoString(L, script)
		L.Global("result")
	}
}

package json5_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/vayload/vayload/pkg/encoding/json5"
)

// =========================
// Casos combinados JSON5
// =========================
var json5TestCases = []struct {
	name   string
	json   string
	verify func(map[string]any) error
}{
	{
		name: "all basics",
		json: `
		{
			// Strings
			simple: 'hello',
			doubleQuoted: "world",
			multiline: "line1 \
line2",
			escaped: "tab:\t newline:\n quote:\" backslash:\\",
			unicode: "\u0041\u00DF\u6771",

			// Numbers
			int: 42,
			float: 3.14,
			neg: -7,
			exp: 1e3,
			hex: 0xDEAD,
			posInf: Infinity,
			negInf: -Infinity,
			nanVal: NaN,

			// Booleans / null
			t: true,
			f: false,
			n: null,

			// Arrays & Objects
			arr: [1,2,3,],
			nestedObj: {a:1,b:{c:2,},},

			// Comentarios
			/* multi-line
			   comment */
		}
		`,
		verify: func(m map[string]any) error {
			// Strings
			if m["simple"] != "hello" {
				return fmt.Errorf("simple string mismatch")
			}
			if m["doubleQuoted"] != "world" {
				return fmt.Errorf("doubleQuoted mismatch")
			}
			if m["multiline"] != "line1 line2" {
				return fmt.Errorf("multiline mismatch")
			}
			if m["escaped"] != "tab:\t newline:\n quote:\" backslash:\\" {
				return fmt.Errorf("escaped mismatch")
			}
			if m["unicode"] != "Aß東" {
				return fmt.Errorf("unicode mismatch")
			}

			// Numbers
			if m["int"].(float64) != 42 {
				return fmt.Errorf("int mismatch")
			}
			if m["float"].(float64) != 3.14 {
				return fmt.Errorf("float mismatch")
			}
			if m["neg"].(float64) != -7 {
				return fmt.Errorf("neg mismatch")
			}
			if m["exp"].(float64) != 1e3 {
				return fmt.Errorf("exp mismatch")
			}
			if m["hex"].(float64) != 0xDEAD {
				return fmt.Errorf("hex mismatch")
			}
			if m["posInf"].(float64) != math.Inf(1) {
				return fmt.Errorf("posInf mismatch")
			}
			if m["negInf"].(float64) != math.Inf(-1) {
				return fmt.Errorf("negInf mismatch")
			}
			if !math.IsNaN(m["nanVal"].(float64)) {
				return fmt.Errorf("nanVal mismatch")
			}

			// Booleans / null
			if m["t"] != true || m["f"] != false || m["n"] != nil {
				return fmt.Errorf("boolean/null mismatch")
			}

			// Arrays & nested objects
			arr := m["arr"].([]any)
			if len(arr) != 3 || arr[0].(float64) != 1 {
				return fmt.Errorf("array mismatch")
			}
			nested := m["nestedObj"].(map[string]any)
			if nested["a"].(float64) != 1 {
				return fmt.Errorf("nestedObj.a mismatch")
			}
			if nested["b"].(map[string]any)["c"].(float64) != 2 {
				return fmt.Errorf("nestedObj.b.c mismatch")
			}

			return nil
		},
	},
	// Aquí puedes agregar más combinaciones automáticas de tests complejos
}

// =========================
// Test runner automático
// =========================
func TestJSON5Spec(t *testing.T) {
	for _, tc := range json5TestCases {
		t.Run(tc.name, func(t *testing.T) {
			var m map[string]any
			if err := json5.Unmarshal([]byte(tc.json), &m); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if err := tc.verify(m); err != nil {
				t.Fatalf("Verification failed: %v", err)
			}

			// Marshal roundtrip
			out, err := json5.Marshal(m)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			var round map[string]any
			if err := json5.Unmarshal(out, &round); err != nil {
				t.Fatalf("Roundtrip unmarshal failed: %v", err)
			}

			if err := tc.verify(round); err != nil {
				t.Fatalf("Roundtrip verification failed: %v", err)
			}
		})
	}
}

// // =========================
// // Streaming tests (Decoder / Encoder)
// // =========================
// func TestJSON5Streaming(t *testing.T) {
// 	type Data struct {
// 		Name string  `json:"name"`
// 		Val  float64 `json:"val"`
// 		Arr  []int   `json:"arr"`
// 		Obj  map[string]string `json:"obj"`
// 	}

// 	original := Data{
// 		Name: "Test",
// 		Val:  123.45,
// 		Arr:  []int{1,2,3},
// 		Obj:  map[string]string{"a":"A","b":"B"},
// 	}

// 	var buf bytes.Buffer
// 	enc := json5.NewEncoder(&buf)
// 	if err := enc.Encode(original); err != nil {
// 		t.Fatalf("Encode failed: %v", err)
// 	}

// 	var decoded Data
// 	dec := json5.NewDecoder(&buf)
// 	if err := dec.Decode(&decoded); err != nil {
// 		t.Fatalf("Decode failed: %v", err)
// 	}

// 	if decoded.Name != original.Name || decoded.Val != original.Val || len(decoded.Arr) != len(original.Arr) {
// 		t.Errorf("Streaming mismatch: %+v", decoded)
// 	}
// }

// =========================
// Benchmarks automáticos
// =========================
func BenchmarkJSON5MarshalFull(b *testing.B) {
	data := map[string]any{
		"string": "value", "int": 42, "float": 3.14, "bool": true,
		"arr": []any{1, 2, 3, 4, 5},
		"obj": map[string]any{"nested": "yes", "arr": []int{1, 2}},
	}

	for b.Loop() {
		_, _ = json5.Marshal(data)
	}
}

func BenchmarkJSON5UnmarshalFull(b *testing.B) {
	data := []byte(`{
		"string":"value","int":42,"float":3.14,"bool":true,
		"arr":[1,2,3,4,5],"obj":{"nested":"yes","arr":[1,2]}
	}`)

	for b.Loop() {
		var m map[string]any
		_ = json5.Unmarshal(data, &m)
	}
}

package performance

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/go-lua"
	"github.com/d5/tengo/v2"
	"github.com/dop251/goja"
	glua "github.com/yuin/gopher-lua"
)

// ─────────────────────────────────────────────────────────────────────────────
// Scripts
// ─────────────────────────────────────────────────────────────────────────────

const fibCodeLua = `
function fib(n)
    if n < 2 then return n end
    return fib(n - 1) + fib(n - 2)
end
fib(20)
`

const fibCodeJS = `
function fib(n) {
    if (n < 2) return n;
    return fib(n - 1) + fib(n - 2);
}
fib(20);
`

const fibCodeTengo = `
fib := func(n) {
    if n < 2 { return n }
    return fib(n - 1) + fib(n - 2)
}
fib(20)
`

const loopCodeLua = `
local sum = 0
for i = 1, 1000 do
    sum = sum + i
end
`

const loopCodeJS = `
let sum = 0;
for (let i = 1; i <= 1000; i++) {
    sum += i;
}
`

const loopCodeTengo = `
sum := 0
for i := 1; i <= 1000; i++ {
    sum += i
}
`

// Heavier workload: string manipulation
const stringCodeLua = `
local s = ""
for i = 1, 200 do
    s = s .. tostring(i)
end
`

const stringCodeJS = `
let s = "";
for (let i = 1; i <= 200; i++) {
    s += i.toString();
}
`

const stringCodeTengo = `
s := ""
for i := 1; i <= 200; i++ {
    s += string(i)
}
`

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func gcStats() (uint64, uint32) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.TotalAlloc, ms.NumGC
}

// measureColdStart runs fn N times each in a brand-new VM (real cold start).
// Returns slice of individual durations.
func measureColdStart(n int, fn func() error) ([]time.Duration, error) {
	durations := make([]time.Duration, n)
	for i := 0; i < n; i++ {
		t0 := time.Now()
		if err := fn(); err != nil {
			return nil, err
		}
		durations[i] = time.Since(t0)
	}
	return durations, nil
}

func avgDuration(ds []time.Duration) time.Duration {
	var total time.Duration
	for _, d := range ds {
		total += d
	}
	return total / time.Duration(len(ds))
}

func minMax(ds []time.Duration) (time.Duration, time.Duration) {
	mn, mx := ds[0], ds[0]
	for _, d := range ds[1:] {
		if d < mn {
			mn = d
		}
		if d > mx {
			mx = d
		}
	}
	return mn, mx
}

// ─────────────────────────────────────────────────────────────────────────────
// COLD START – each iteration spins up a brand-new VM + compiles + runs
// ─────────────────────────────────────────────────────────────────────────────

func BenchmarkColdStart_GopherLua(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L := glua.NewState()
		fn, err := L.LoadString(fibCodeLua)
		if err != nil {
			b.Fatal(err)
		}
		L.Push(fn)
		if err := L.PCall(0, 0, nil); err != nil {
			b.Fatal(err)
		}
		L.Close()
	}
}

func BenchmarkColdStart_GoLua(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L := lua.NewState()
		if err := lua.LoadString(L, fibCodeLua); err != nil {
			b.Fatal(err)
		}
		if err := L.ProtectedCall(0, 0, 0); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkColdStart_Goja(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm := goja.New()
		prog, err := goja.Compile("fib.js", fibCodeJS, false)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := vm.RunProgram(prog); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkColdStart_Tengo(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		script := tengo.NewScript([]byte(fibCodeTengo))
		compiled, err := script.Compile()
		if err != nil {
			b.Fatal(err)
		}
		if err := compiled.RunContext(context.Background()); err != nil {
			b.Fatal(err)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// WARM EXECUTION – VM created once, script pre-compiled, only Run measured
// ─────────────────────────────────────────────────────────────────────────────

func BenchmarkWarm_GopherLua_Fib(b *testing.B) {
	L := glua.NewState()
	defer L.Close()
	fn, err := L.LoadString(fibCodeLua)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L.Push(fn)
		if err := L.PCall(0, 0, nil); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWarm_Goja_Fib(b *testing.B) {
	vm := goja.New()
	prog, err := goja.Compile("fib.js", fibCodeJS, false)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := vm.RunProgram(prog); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWarm_Tengo_Fib(b *testing.B) {
	script := tengo.NewScript([]byte(fibCodeTengo))
	compiled, err := script.Compile()
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := compiled.RunContext(context.Background()); err != nil {
			b.Fatal(err)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// LOOP benchmarks (warm)
// ─────────────────────────────────────────────────────────────────────────────

func BenchmarkWarm_GopherLua_Loop(b *testing.B) {
	L := glua.NewState()
	defer L.Close()
	fn, _ := L.LoadString(loopCodeLua)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L.Push(fn)
		L.PCall(0, 0, nil)
	}
}

func BenchmarkWarm_Goja_Loop(b *testing.B) {
	vm := goja.New()
	prog, _ := goja.Compile("loop.js", loopCodeJS, false)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.RunProgram(prog)
	}
}

func BenchmarkWarm_Tengo_Loop(b *testing.B) {
	script := tengo.NewScript([]byte(loopCodeTengo))
	compiled, _ := script.Compile()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compiled.Run()
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// STRING benchmarks (warm)
// ─────────────────────────────────────────────────────────────────────────────

func BenchmarkWarm_GopherLua_String(b *testing.B) {
	L := glua.NewState()
	defer L.Close()
	fn, _ := L.LoadString(stringCodeLua)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L.Push(fn)
		L.PCall(0, 0, nil)
	}
}

func BenchmarkWarm_Goja_String(b *testing.B) {
	vm := goja.New()
	prog, _ := goja.Compile("string.js", stringCodeJS, false)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.RunProgram(prog)
	}
}

func BenchmarkWarm_Tengo_String(b *testing.B) {
	script := tengo.NewScript([]byte(stringCodeTengo))
	compiled, _ := script.Compile()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compiled.Run()
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// CONCURRENT execution – multiple goroutines sharing/isolating VMs
// ─────────────────────────────────────────────────────────────────────────────

func BenchmarkConcurrent_GopherLua(b *testing.B) {
	// gopher-lua states are NOT goroutine-safe → one per goroutine
	b.ReportAllocs()
	b.SetParallelism(8)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		L := glua.NewState()
		defer L.Close()
		fn, err := L.LoadString(fibCodeLua)
		if err != nil {
			b.Fatal(err)
		}
		for pb.Next() {
			L.Push(fn)
			L.PCall(0, 0, nil)
		}
	})
}

func BenchmarkConcurrent_Goja(b *testing.B) {
	// goja is NOT goroutine-safe → one VM per goroutine
	prog, err := goja.Compile("fib.js", fibCodeJS, false)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.SetParallelism(8)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		vm := goja.New()
		for pb.Next() {
			vm.RunProgram(prog)
		}
	})
}

func BenchmarkConcurrent_Tengo(b *testing.B) {
	// Tengo compiled scripts are safe to Clone per goroutine
	script := tengo.NewScript([]byte(fibCodeTengo))
	compiled, err := script.Compile()
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.SetParallelism(8)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		c := compiled.Clone()
		for pb.Next() {
			c.RunContext(context.Background())
		}
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// MEMORY PRESSURE – force GC and measure heap impact
// ─────────────────────────────────────────────────────────────────────────────

func BenchmarkMemory_GopherLua(b *testing.B) {
	b.ReportAllocs()
	runtime.GC()
	var msBefore, msAfter runtime.MemStats
	runtime.ReadMemStats(&msBefore)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L := glua.NewState()
		fn, _ := L.LoadString(fibCodeLua)
		L.Push(fn)
		L.PCall(0, 0, nil)
		L.Close()
	}
	b.StopTimer()

	runtime.GC()
	runtime.ReadMemStats(&msAfter)
	b.ReportMetric(float64(msAfter.TotalAlloc-msBefore.TotalAlloc)/float64(b.N), "B/op-total")
	b.ReportMetric(float64(msAfter.NumGC-msBefore.NumGC), "gc-cycles")
}

func BenchmarkMemory_Goja(b *testing.B) {
	b.ReportAllocs()
	runtime.GC()
	var msBefore, msAfter runtime.MemStats
	runtime.ReadMemStats(&msBefore)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm := goja.New()
		prog, _ := goja.Compile("fib.js", fibCodeJS, false)
		vm.RunProgram(prog)
	}
	b.StopTimer()

	runtime.GC()
	runtime.ReadMemStats(&msAfter)
	b.ReportMetric(float64(msAfter.TotalAlloc-msBefore.TotalAlloc)/float64(b.N), "B/op-total")
	b.ReportMetric(float64(msAfter.NumGC-msBefore.NumGC), "gc-cycles")
}

func BenchmarkMemory_Tengo(b *testing.B) {
	b.ReportAllocs()
	runtime.GC()
	var msBefore, msAfter runtime.MemStats
	runtime.ReadMemStats(&msBefore)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := tengo.NewScript([]byte(fibCodeTengo))
		c, _ := s.Compile()
		c.RunContext(context.Background())
	}
	b.StopTimer()

	runtime.GC()
	runtime.ReadMemStats(&msAfter)
	b.ReportMetric(float64(msAfter.TotalAlloc-msBefore.TotalAlloc)/float64(b.N), "B/op-total")
	b.ReportMetric(float64(msAfter.NumGC-msBefore.NumGC), "gc-cycles")
}

// ─────────────────────────────────────────────────────────────────────────────
// DETAILED COLD START TEST (not a benchmark – prints human-readable report)
// ─────────────────────────────────────────────────────────────────────────────

func TestColdStartDetailed(t *testing.T) {
	const runs = 50
	runtime.GC()

	type result struct {
		name     string
		avg      time.Duration
		min      time.Duration
		max      time.Duration
		alloc    uint64
		gcCycles uint32
	}

	var results []result

	measure := func(name string, fn func() error) {
		allocBefore, gcBefore := gcStats()
		ds, err := measureColdStart(runs, fn)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		runtime.GC()
		allocAfter, gcAfter := gcStats()
		mn, mx := minMax(ds)
		results = append(results, result{
			name:     name,
			avg:      avgDuration(ds),
			min:      mn,
			max:      mx,
			alloc:    allocAfter - allocBefore,
			gcCycles: gcAfter - gcBefore,
		})
	}

	measure("GopherLua", func() error {
		L := glua.NewState()
		defer L.Close()
		fn, err := L.LoadString(fibCodeLua)
		if err != nil {
			return err
		}
		L.Push(fn)
		return L.PCall(0, 0, nil)
	})

	measure("GoLua", func() error {
		L := lua.NewState()
		if err := lua.LoadString(L, fibCodeLua); err != nil {
			return err
		}
		return L.ProtectedCall(0, 0, 0)
	})

	measure("Goja", func() error {
		vm := goja.New()
		prog, err := goja.Compile("fib.js", fibCodeJS, false)
		if err != nil {
			return err
		}
		_, err = vm.RunProgram(prog)
		return err
	})

	measure("Tengo", func() error {
		s := tengo.NewScript([]byte(fibCodeTengo))
		c, err := s.Compile()
		if err != nil {
			return err
		}
		return c.RunContext(context.Background())
	})

	fmt.Printf("\n%-14s %12s %12s %12s %14s %10s\n",
		"Engine", "Avg", "Min", "Max", "TotalAlloc", "GC cycles")
	fmt.Printf("%s\n", "────────────────────────────────────────────────────────────────────────────")
	for _, r := range results {
		fmt.Printf("%-14s %12s %12s %12s %12.1f KB %10d\n",
			r.name, r.avg, r.min, r.max,
			float64(r.alloc)/1024, r.gcCycles)
	}
	fmt.Println()
}

// ─────────────────────────────────────────────────────────────────────────────
// POOL / REUSE TEST – simulate real-world VM pool pattern
// ─────────────────────────────────────────────────────────────────────────────

func BenchmarkPool_GopherLua(b *testing.B) {
	pool := sync.Pool{
		New: func() interface{} {
			return glua.NewState()
		},
	}
	fn_cache := make(map[*glua.LState]*glua.LFunction)
	var mu sync.Mutex

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L := pool.Get().(*glua.LState)
		mu.Lock()
		fn, ok := fn_cache[L]
		if !ok {
			fn, _ = L.LoadString(fibCodeLua)
			fn_cache[L] = fn
		}
		mu.Unlock()
		L.Push(fn)
		L.PCall(0, 0, nil)
		pool.Put(L)
	}
}

func BenchmarkPool_Goja(b *testing.B) {
	prog, _ := goja.Compile("fib.js", fibCodeJS, false)
	pool := sync.Pool{
		New: func() interface{} {
			return goja.New()
		},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm := pool.Get().(*goja.Runtime)
		vm.RunProgram(prog)
		pool.Put(vm)
	}
}

func BenchmarkPool_Tengo(b *testing.B) {
	script := tengo.NewScript([]byte(fibCodeTengo))
	compiled, _ := script.Compile()

	pool := sync.Pool{
		New: func() interface{} {
			return compiled.Clone()
		},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := pool.Get().(*tengo.Compiled)
		c.RunContext(context.Background())
		pool.Put(c)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// SUSTAINED LOAD – many iterations to expose GC pressure over time
// ─────────────────────────────────────────────────────────────────────────────

func TestSustainedLoad(t *testing.T) {
	const duration = 3 * time.Second

	type result struct {
		name      string
		ops       int
		gcCycles  uint32
		heapAlloc uint64
	}

	run := func(name string, fn func()) result {
		runtime.GC()
		var msBefore runtime.MemStats
		runtime.ReadMemStats(&msBefore)

		ops := 0
		deadline := time.Now().Add(duration)
		for time.Now().Before(deadline) {
			fn()
			ops++
		}

		runtime.GC()
		var msAfter runtime.MemStats
		runtime.ReadMemStats(&msAfter)

		return result{
			name:      name,
			ops:       ops,
			gcCycles:  msAfter.NumGC - msBefore.NumGC,
			heapAlloc: msAfter.TotalAlloc - msBefore.TotalAlloc,
		}
	}

	// GopherLua warm
	glL := glua.NewState()
	glFn, _ := glL.LoadString(fibCodeLua)
	rGopherLua := run("GopherLua(warm)", func() {
		glL.Push(glFn)
		glL.PCall(0, 0, nil)
	})
	glL.Close()

	// Goja warm
	gojaVM := goja.New()
	gojaProg, _ := goja.Compile("fib.js", fibCodeJS, false)
	rGoja := run("Goja(warm)", func() {
		gojaVM.RunProgram(gojaProg)
	})

	// Tengo warm
	tengoScript := tengo.NewScript([]byte(fibCodeTengo))
	tengoCompiled, _ := tengoScript.Compile()
	rTengo := run("Tengo(warm)", func() {
		tengoCompiled.RunContext(context.Background())
	})

	// GopherLua cold
	rGopherLuaCold := run("GopherLua(cold)", func() {
		L := glua.NewState()
		fn, _ := L.LoadString(fibCodeLua)
		L.Push(fn)
		L.PCall(0, 0, nil)
		L.Close()
	})

	// Goja cold
	rGojaCold := run("Goja(cold)", func() {
		vm := goja.New()
		prog, _ := goja.Compile("fib.js", fibCodeJS, false)
		vm.RunProgram(prog)
	})

	// Tengo cold
	rTengoCold := run("Tengo(cold)", func() {
		s := tengo.NewScript([]byte(fibCodeTengo))
		c, _ := s.Compile()
		c.RunContext(context.Background())
	})

	allResults := []result{rGopherLua, rGoja, rTengo, rGopherLuaCold, rGojaCold, rTengoCold}

	fmt.Printf("\n=== Sustained Load (%v) ===\n", duration)
	fmt.Printf("%-20s %10s %12s %14s %12s\n",
		"Engine", "Ops", "Ops/sec", "TotalAlloc", "GC cycles")
	fmt.Printf("%s\n", "────────────────────────────────────────────────────────────────────────────")
	for _, r := range allResults {
		opsPerSec := float64(r.ops) / duration.Seconds()
		fmt.Printf("%-20s %10d %12.0f %12.1f MB %12d\n",
			r.name, r.ops, opsPerSec,
			float64(r.heapAlloc)/1024/1024, r.gcCycles)
	}
	fmt.Println()
}

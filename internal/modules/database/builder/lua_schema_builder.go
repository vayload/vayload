package builder

import (
	"github.com/vayload/vayload/internal/modules/database/migrator"
	lua "github.com/yuin/gopher-lua"
)

const (
	luaBlueprintType = "schema_blueprint"
	luaColumnType    = "schema_column"
	luaForeignType   = "schema_foreign"
)

type schemaContext struct {
	builder *SchemaBuilder
}

func RegisterSchemaBuilderAPI(L *lua.LState, sb *SchemaBuilder) {
	registerTypes(L)
	ctx := &schemaContext{builder: sb}
	schema := L.NewTable()

	L.SetFuncs(schema, map[string]lua.LGFunction{
		"create":       ctx.schemaCreate,
		"table":        ctx.schemaTable,
		"drop":         ctx.schemaDrop,
		"dropIfExists": ctx.schemaDropIfExists,
		"index":        ctx.schemaIndex,
		"dropIndex":    ctx.schemaDropIndex,
		"unique":       ctx.schemaUnique,
		"dropUnique":   ctx.schemaDropUnique,
		"addColumn":    ctx.schemaAddColumn,
		"dropColumn":   ctx.schemaDropColumn,
		"renameColumn": ctx.schemaRenameColumn,
		"hasTable":     ctx.schemaHasTable,
		"hasColumn":    ctx.schemaHasColumn,
	})

	L.SetGlobal("Schema", schema)
}

func registerTypes(L *lua.LState) {
	bp := L.NewTypeMetatable(luaBlueprintType)
	L.SetField(bp, "__index", L.SetFuncs(L.NewTable(), blueprintMethods))

	col := L.NewTypeMetatable(luaColumnType)
	L.SetField(col, "__index", L.SetFuncs(L.NewTable(), columnMethods))

	fk := L.NewTypeMetatable(luaForeignType)
	L.SetField(fk, "__index", L.SetFuncs(L.NewTable(), foreignMethods))
}

func newBlueprint(L *lua.LState, bp *migrator.Blueprint) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = bp

	L.SetMetatable(ud, L.GetTypeMetatable(luaBlueprintType))

	return ud
}

func newColumn(L *lua.LState, col *migrator.Column) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = col

	L.SetMetatable(ud, L.GetTypeMetatable(luaColumnType))

	return ud
}

func newForeign(L *lua.LState, fk *migrator.ForeignColumn) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = fk
	L.SetMetatable(ud, L.GetTypeMetatable(luaForeignType))

	return ud
}

func checkBlueprint(L *lua.LState, idx int) *migrator.Blueprint {
	ud := L.CheckUserData(idx)

	if v, ok := ud.Value.(*migrator.Blueprint); ok {
		return v
	}

	L.ArgError(idx, "blueprint expected")
	return nil
}

func checkColumn(L *lua.LState, idx int) *migrator.Column {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*migrator.Column); ok {
		return v
	}

	L.ArgError(idx, "column expected")
	return nil
}

func checkForeign(L *lua.LState, idx int) *migrator.ForeignColumn {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*migrator.ForeignColumn); ok {
		return v
	}

	L.ArgError(idx, "foreign expected")
	return nil
}

func tableToStrings(tbl *lua.LTable) []string {
	n := tbl.Len()
	out := make([]string, n)
	for i := 1; i <= n; i++ {
		out[i-1] = tbl.RawGetInt(i).String()
	}

	return out
}

///////////////////////////
//// Schema operations
///////////////////////////

func (ctx *schemaContext) schemaCreate(L *lua.LState) int {
	table := L.CheckString(1)
	cb := L.CheckFunction(2)
	bp := migrator.NewBlueprint(table)

	L.Push(cb)
	L.Push(newBlueprint(L, bp))

	if err := L.PCall(1, 0, nil); err != nil {

		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// if err := ctx.builder.CreateBlueprint(bp); err != nil {

	// 	L.Push(lua.LNil)
	// 	L.Push(lua.LString(err.Error()))
	// 	return 2
	// }

	L.Push(lua.LTrue)
	return 1
}

func (ctx *schemaContext) schemaTable(L *lua.LState) int {
	table := L.CheckString(1)
	cb := L.CheckFunction(2)
	bp := migrator.NewBlueprint(table)
	bp.IsAltering = true

	L.Push(cb)
	L.Push(newBlueprint(L, bp))

	if err := L.PCall(1, 0, nil); err != nil {

		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// if err := ctx.builder.AlterBlueprint(bp); err != nil {

	// 	L.Push(lua.LNil)
	// 	L.Push(lua.LString(err.Error()))
	// 	return 2
	// }

	L.Push(lua.LTrue)
	return 1
}

func (ctx *schemaContext) schemaDrop(L *lua.LState) int {
	table := L.CheckString(1)
	err := ctx.builder.Drop(table)

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

func (ctx *schemaContext) schemaDropIfExists(L *lua.LState) int {
	table := L.CheckString(1)
	err := ctx.builder.DropIfExists(table)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

func (ctx *schemaContext) schemaIndex(L *lua.LState) int {
	table := L.CheckString(1)
	cols := tableToStrings(L.CheckTable(2))
	name := ""
	if L.GetTop() >= 3 {
		name = L.CheckString(3)
	}

	err := ctx.builder.AddIndex(table, cols, name, false)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

func (ctx *schemaContext) schemaUnique(L *lua.LState) int {
	table := L.CheckString(1)
	cols := tableToStrings(L.CheckTable(2))

	name := ""
	if L.GetTop() >= 3 {
		name = L.CheckString(3)
	}

	err := ctx.builder.AddIndex(table, cols, name, true)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

func (ctx *schemaContext) schemaDropIndex(L *lua.LState) int {
	table := L.CheckString(1)
	name := L.CheckString(2)

	err := ctx.builder.DropIndex(table, name)
	if err != nil {

		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

func (ctx *schemaContext) schemaDropUnique(L *lua.LState) int {
	return ctx.schemaDropIndex(L)
}

func (ctx *schemaContext) schemaAddColumn(L *lua.LState) int {
	table := L.CheckString(1)
	col := L.CheckString(2)
	typ := L.CheckString(3)

	err := ctx.builder.AddColumn(table, col, typ)
	if err != nil {

		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

func (ctx *schemaContext) schemaDropColumn(L *lua.LState) int {
	table := L.CheckString(1)
	col := L.CheckString(2)

	err := ctx.builder.DropColumn(table, col)
	if err != nil {

		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

func (ctx *schemaContext) schemaRenameColumn(L *lua.LState) int {
	table := L.CheckString(1)
	from := L.CheckString(2)
	to := L.CheckString(3)

	err := ctx.builder.RenameColumn(table, from, to)
	if err != nil {

		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

func (ctx *schemaContext) schemaHasTable(L *lua.LState) int {
	table := L.CheckString(1)

	ok, err := ctx.builder.HasTable(table)
	if err != nil {

		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(ok))
	return 1
}

func (ctx *schemaContext) schemaHasColumn(L *lua.LState) int {
	table := L.CheckString(1)
	col := L.CheckString(2)

	ok, err := ctx.builder.HasColumn(table, col)
	if err != nil {

		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(ok))
	return 1
}

///////////////////////////
//// Blueprint methods
///////////////////////////

var blueprintMethods = map[string]lua.LGFunction{
	"id":           bpID,
	"string":       bpString,
	"text":         bpText,
	"integer":      bpInteger,
	"bigInteger":   bpBigInteger,
	"boolean":      bpBoolean,
	"timestamp":    bpTimestamp,
	"decimal":      bpDecimal,
	"float":        bpFloat,
	"json":         bpJSON,
	"timestamps":   bpTimestamps,
	"softDeletes":  bpSoftDeletes,
	"foreign":      bpForeign,
	"index":        bpIndex,
	"unique":       bpUnique,
	"dropColumn":   bpDropColumn,
	"renameColumn": bpRenameColumn,
	"dropIndex":    bpDropIndex,
}

func bpID(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	bp.ID()
	return 0
}

func bpString(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	name := L.CheckString(2)

	length := 255
	if L.GetTop() >= 3 {
		length = L.CheckInt(3)
	}

	L.Push(newColumn(L, bp.String(name, length)))
	return 1
}

func bpText(L *lua.LState) int {

	bp := checkBlueprint(L, 1)
	name := L.CheckString(2)

	L.Push(newColumn(L, bp.Text(name)))
	return 1
}

func bpInteger(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	name := L.CheckString(2)

	L.Push(newColumn(L, bp.Integer(name)))
	return 1
}

func bpBigInteger(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	name := L.CheckString(2)

	L.Push(newColumn(L, bp.BigInteger(name)))
	return 1
}

func bpBoolean(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	name := L.CheckString(2)

	L.Push(newColumn(L, bp.Boolean(name)))
	return 1
}

func bpTimestamp(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	name := L.CheckString(2)

	L.Push(newColumn(L, bp.Timestamp(name)))
	return 1
}

func bpDecimal(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	name := L.CheckString(2)

	p := 8
	s := 2

	if L.GetTop() >= 3 {
		p = L.CheckInt(3)
	}

	if L.GetTop() >= 4 {
		s = L.CheckInt(4)
	}

	L.Push(newColumn(L, bp.Decimal(name, p, s)))
	return 1
}

func bpFloat(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	name := L.CheckString(2)

	L.Push(newColumn(L, bp.Float(name)))
	return 1
}

func bpJSON(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	name := L.CheckString(2)

	L.Push(newColumn(L, bp.JSON(name)))
	return 1
}

func bpTimestamps(L *lua.LState) int {
	checkBlueprint(L, 1).Timestamps()
	return 0
}

func bpSoftDeletes(L *lua.LState) int {
	checkBlueprint(L, 1).SoftDeletes()
	return 0
}

func bpForeign(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	name := L.CheckString(2)

	L.Push(newForeign(L, bp.Foreign(name)))
	return 1
}

func bpIndex(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	cols := tableToStrings(L.CheckTable(2))
	name := ""
	if L.GetTop() >= 3 {
		name = L.CheckString(3)
	}

	bp.Index(cols, name)
	return 0
}

func bpUnique(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	cols := tableToStrings(L.CheckTable(2))
	name := ""
	if L.GetTop() >= 3 {
		name = L.CheckString(3)
	}

	bp.UniqueIndex(cols, name)
	return 0
}

func bpDropColumn(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	col := L.CheckString(2)
	bp.DropColumn(col)
	return 0
}

func bpRenameColumn(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	from := L.CheckString(2)
	to := L.CheckString(3)
	bp.RenameColumn(from, to)
	return 0
}

func bpDropIndex(L *lua.LState) int {
	bp := checkBlueprint(L, 1)
	name := L.CheckString(2)
	bp.DropIndex(name)
	return 0
}

///////////////////////////
//// Column modifiers
///////////////////////////

var columnMethods = map[string]lua.LGFunction{
	"nullable":      columnNullable,
	"default":       columnDefault,
	"unique":        columnUnique,
	"unsigned":      columnUnsigned,
	"autoIncrement": columnAutoIncrement,
	"primary":       columnPrimary,
}

func columnNullable(L *lua.LState) int {
	checkColumn(L, 1).Nullable()
	L.Push(L.Get(1))
	return 1
}

func columnDefault(L *lua.LState) int {
	col := checkColumn(L, 1)
	val := L.Get(2)
	switch v := val.(type) {

	case lua.LString:
		col.Default(string(v))

	case lua.LNumber:
		col.Default(float64(v))

	case lua.LBool:
		col.Default(bool(v))

	default:
		col.Default(nil)
	}

	L.Push(L.Get(1))
	return 1
}

func columnUnique(L *lua.LState) int {
	checkColumn(L, 1).Unique()
	L.Push(L.Get(1))
	return 1
}

func columnUnsigned(L *lua.LState) int {
	checkColumn(L, 1).Unsigned()
	L.Push(L.Get(1))
	return 1
}

func columnAutoIncrement(L *lua.LState) int {
	checkColumn(L, 1).AutoIncrement()
	L.Push(L.Get(1))
	return 1
}

func columnPrimary(L *lua.LState) int {
	checkColumn(L, 1).Primary()
	L.Push(L.Get(1))
	return 1
}

///////////////////////////
//// Foreign key
///////////////////////////

var foreignMethods = map[string]lua.LGFunction{
	"references": foreignReferences,
	"on":         foreignOn,
	"onDelete":   foreignOnDelete,
	"nullable":   foreignNullable,
}

func foreignReferences(L *lua.LState) int {
	fk := checkForeign(L, 1)
	col := L.CheckString(2)
	fk.References(col)
	L.Push(L.Get(1))
	return 1
}

func foreignOn(L *lua.LState) int {
	fk := checkForeign(L, 1)
	table := L.CheckString(2)
	fk.On(table)
	L.Push(L.Get(1))
	return 1
}

func foreignOnDelete(L *lua.LState) int {
	fk := checkForeign(L, 1)
	action := L.CheckString(2)
	fk.OnDelete(action)
	L.Push(L.Get(1))
	return 1
}

func foreignNullable(L *lua.LState) int {
	checkForeign(L, 1).Nullable()
	L.Push(L.Get(1))
	return 1
}

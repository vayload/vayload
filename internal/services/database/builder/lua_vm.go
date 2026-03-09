package builder

import (
	"github.com/vayload/vayload/internal/services/database/migrator"
	lua "github.com/yuin/gopher-lua"
)

// RegisterSchemaBuilderAPI registra el SchemaBuilder en Lua
func RegisterSchemaBuilderAPI(L *lua.LState, sb *SchemaBuilder) {
	schemaTable := L.NewTable()

	// Métodos de tabla
	L.SetField(schemaTable, "create", L.NewFunction(luaSchemaCreate(sb)))
	L.SetField(schemaTable, "drop", L.NewFunction(luaSchemaDrop(sb)))
	L.SetField(schemaTable, "dropIfExists", L.NewFunction(luaSchemaDropIfExists(sb)))
	L.SetField(schemaTable, "table", L.NewFunction(luaSchemaTable(sb)))

	// Métodos de índices
	L.SetField(schemaTable, "index", L.NewFunction(luaSchemaIndex(sb)))
	L.SetField(schemaTable, "dropIndex", L.NewFunction(luaSchemaDropIndex(sb)))
	L.SetField(schemaTable, "unique", L.NewFunction(luaSchemaUnique(sb)))
	L.SetField(schemaTable, "dropUnique", L.NewFunction(luaSchemaDropUnique(sb)))

	// Métodos de columnas
	L.SetField(schemaTable, "addColumn", L.NewFunction(luaSchemaAddColumn(sb)))
	L.SetField(schemaTable, "dropColumn", L.NewFunction(luaSchemaDropColumn(sb)))
	L.SetField(schemaTable, "renameColumn", L.NewFunction(luaSchemaRenameColumn(sb)))

	// Métodos de utilidad
	L.SetField(schemaTable, "hasTable", L.NewFunction(luaSchemaHasTable(sb)))
	L.SetField(schemaTable, "hasColumn", L.NewFunction(luaSchemaHasColumn(sb)))

	L.SetGlobal("Schema", schemaTable)
}

// luaSchemaCreate maneja Schema.create(table_name, callback)
func luaSchemaCreate(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)
		callback := L.CheckFunction(2)

		bp := migrator.NewBlueprint(tableName)
		bpTable := createBlueprintTable(L, bp)

		L.Push(callback)
		L.Push(bpTable)
		if err := L.PCall(1, 0, nil); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		if err := sb.Create(tableName, func(b *migrator.Blueprint) {
			*b = *bp
		}); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LTrue)
		return 1
	}
}

// luaSchemaTable maneja Schema.table(table_name, callback) para modificar tablas
func luaSchemaTable(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)
		callback := L.CheckFunction(2)

		bp := migrator.NewBlueprint(tableName)
		bp.IsAltering = true
		bpTable := createBlueprintTable(L, bp)

		L.Push(callback)
		L.Push(bpTable)
		if err := L.PCall(1, 0, nil); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		if err := sb.Table(tableName, func(b *migrator.Blueprint) {
			*b = *bp
		}); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LTrue)
		return 1
	}
}

// luaSchemaDrop maneja Schema.drop(table_name)
func luaSchemaDrop(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)

		if err := sb.Drop(tableName); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LTrue)
		return 1
	}
}

// luaSchemaDropIfExists maneja Schema.dropIfExists(table_name)
func luaSchemaDropIfExists(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)

		if err := sb.DropIfExists(tableName); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LTrue)
		return 1
	}
}

// luaSchemaIndex maneja Schema.index(table, columns, name?)
func luaSchemaIndex(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)
		columnsTable := L.CheckTable(2)
		indexName := ""
		if L.GetTop() >= 3 {
			indexName = L.CheckString(3)
		}

		var columns []string
		columnsTable.ForEach(func(_, v lua.LValue) {
			columns = append(columns, lua.LVAsString(v))
		})

		if err := sb.AddIndex(tableName, columns, indexName, false); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LTrue)
		return 1
	}
}

// luaSchemaDropIndex maneja Schema.dropIndex(table, name)
func luaSchemaDropIndex(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)
		indexName := L.CheckString(2)

		if err := sb.DropIndex(tableName, indexName); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LTrue)
		return 1
	}
}

// luaSchemaUnique maneja Schema.unique(table, columns, name?)
func luaSchemaUnique(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)
		columnsTable := L.CheckTable(2)
		indexName := ""
		if L.GetTop() >= 3 {
			indexName = L.CheckString(3)
		}

		var columns []string
		columnsTable.ForEach(func(_, v lua.LValue) {
			columns = append(columns, lua.LVAsString(v))
		})

		if err := sb.AddIndex(tableName, columns, indexName, true); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LTrue)
		return 1
	}
}

// luaSchemaDropUnique maneja Schema.dropUnique(table, name)
func luaSchemaDropUnique(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)
		indexName := L.CheckString(2)

		if err := sb.DropIndex(tableName, indexName); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LTrue)
		return 1
	}
}

// luaSchemaAddColumn maneja Schema.addColumn(table, name, type)
func luaSchemaAddColumn(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)
		columnName := L.CheckString(2)
		columnType := L.CheckString(3)

		if err := sb.AddColumn(tableName, columnName, columnType); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LTrue)
		return 1
	}
}

// luaSchemaDropColumn maneja Schema.dropColumn(table, name)
func luaSchemaDropColumn(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)
		columnName := L.CheckString(2)

		if err := sb.DropColumn(tableName, columnName); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LTrue)
		return 1
	}
}

// luaSchemaRenameColumn maneja Schema.renameColumn(table, from, to)
func luaSchemaRenameColumn(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)
		fromName := L.CheckString(2)
		toName := L.CheckString(3)

		if err := sb.RenameColumn(tableName, fromName, toName); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LTrue)
		return 1
	}
}

// luaSchemaHasTable maneja Schema.hasTable(table)
func luaSchemaHasTable(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)

		exists, err := sb.HasTable(tableName)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LBool(exists))
		return 1
	}
}

// luaSchemaHasColumn maneja Schema.hasColumn(table, column)
func luaSchemaHasColumn(sb *SchemaBuilder) lua.LGFunction {
	return func(L *lua.LState) int {
		tableName := L.CheckString(1)
		columnName := L.CheckString(2)

		exists, err := sb.HasColumn(tableName, columnName)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LBool(exists))
		return 1
	}
}

// createBlueprintTable crea una tabla Lua con la API de Blueprint
func createBlueprintTable(L *lua.LState, bp *migrator.Blueprint) *lua.LTable {
	table := L.NewTable()

	// Métodos de columnas básicas
	L.SetField(table, "id", L.NewFunction(luaBlueprintID(bp)))
	L.SetField(table, "string", L.NewFunction(luaBlueprintString(bp)))
	L.SetField(table, "text", L.NewFunction(luaBlueprintText(bp)))
	L.SetField(table, "integer", L.NewFunction(luaBlueprintInteger(bp)))
	L.SetField(table, "bigInteger", L.NewFunction(luaBlueprintBigInteger(bp)))
	L.SetField(table, "boolean", L.NewFunction(luaBlueprintBoolean(bp)))
	L.SetField(table, "timestamp", L.NewFunction(luaBlueprintTimestamp(bp)))
	L.SetField(table, "decimal", L.NewFunction(luaBlueprintDecimal(bp)))
	L.SetField(table, "float", L.NewFunction(luaBlueprintFloat(bp)))
	L.SetField(table, "json", L.NewFunction(luaBlueprintJSON(bp)))

	// Helpers
	L.SetField(table, "timestamps", L.NewFunction(luaBlueprintTimestamps(bp)))
	L.SetField(table, "softDeletes", L.NewFunction(luaBlueprintSoftDeletes(bp)))

	// Foreign keys
	L.SetField(table, "foreign", L.NewFunction(luaBlueprintForeign(bp)))

	// Índices dentro del blueprint
	L.SetField(table, "index", L.NewFunction(luaBlueprintIndex(bp)))
	L.SetField(table, "unique", L.NewFunction(luaBlueprintUniqueIndex(bp)))

	// Modificación de columnas (para Schema.table)
	L.SetField(table, "dropColumn", L.NewFunction(luaBlueprintDropColumn(bp)))
	L.SetField(table, "renameColumn", L.NewFunction(luaBlueprintRenameColumn(bp)))
	L.SetField(table, "dropIndex", L.NewFunction(luaBlueprintDropIndex(bp)))

	return table
}

// Métodos de columnas básicas
func luaBlueprintID(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		bp.ID()
		return 0
	}
}

func luaBlueprintString(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.CheckString(1)
		length := 255
		if L.GetTop() >= 2 {
			length = L.CheckInt(2)
		}
		col := bp.String(name, length)
		L.Push(createColumnTable(L, col))
		return 1
	}
}

func luaBlueprintText(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.CheckString(1)
		col := bp.Text(name)
		L.Push(createColumnTable(L, col))
		return 1
	}
}

func luaBlueprintInteger(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.CheckString(1)
		col := bp.Integer(name)
		L.Push(createColumnTable(L, col))
		return 1
	}
}

func luaBlueprintBigInteger(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.CheckString(1)
		col := bp.BigInteger(name)
		L.Push(createColumnTable(L, col))
		return 1
	}
}

func luaBlueprintBoolean(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.CheckString(1)
		col := bp.Boolean(name)
		L.Push(createColumnTable(L, col))
		return 1
	}
}

func luaBlueprintTimestamp(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.CheckString(1)
		col := bp.Timestamp(name)
		L.Push(createColumnTable(L, col))
		return 1
	}
}

func luaBlueprintDecimal(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.CheckString(1)
		precision := 8
		scale := 2
		if L.GetTop() >= 2 {
			precision = L.CheckInt(2)
		}
		if L.GetTop() >= 3 {
			scale = L.CheckInt(3)
		}
		col := bp.Decimal(name, precision, scale)
		L.Push(createColumnTable(L, col))
		return 1
	}
}

func luaBlueprintFloat(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.CheckString(1)
		col := bp.Float(name)
		L.Push(createColumnTable(L, col))
		return 1
	}
}

func luaBlueprintJSON(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.CheckString(1)
		col := bp.JSON(name)
		L.Push(createColumnTable(L, col))
		return 1
	}
}

func luaBlueprintTimestamps(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		bp.Timestamps()
		return 0
	}
}

func luaBlueprintSoftDeletes(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		bp.SoftDeletes()
		return 0
	}
}

func luaBlueprintForeign(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.CheckString(1)
		col := bp.Foreign(name)
		L.Push(createForeignColumnTable(L, col))
		return 1
	}
}

// Métodos de índices
func luaBlueprintIndex(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		columnsTable := L.CheckTable(1)
		indexName := ""
		if L.GetTop() >= 2 {
			indexName = L.CheckString(2)
		}

		var columns []string
		columnsTable.ForEach(func(_, v lua.LValue) {
			columns = append(columns, lua.LVAsString(v))
		})

		bp.Index(columns, indexName)
		return 0
	}
}

func luaBlueprintUniqueIndex(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		columnsTable := L.CheckTable(1)
		indexName := ""
		if L.GetTop() >= 2 {
			indexName = L.CheckString(2)
		}

		var columns []string
		columnsTable.ForEach(func(_, v lua.LValue) {
			columns = append(columns, lua.LVAsString(v))
		})

		bp.UniqueIndex(columns, indexName)
		return 0
	}
}

func luaBlueprintDropColumn(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		columnName := L.CheckString(1)
		bp.DropColumn(columnName)
		return 0
	}
}

func luaBlueprintRenameColumn(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		fromName := L.CheckString(1)
		toName := L.CheckString(2)
		bp.RenameColumn(fromName, toName)
		return 0
	}
}

func luaBlueprintDropIndex(bp *migrator.Blueprint) lua.LGFunction {
	return func(L *lua.LState) int {
		indexName := L.CheckString(1)
		bp.DropIndex(indexName)
		return 0
	}
}

// createColumnTable crea métodos encadenables para Column
func createColumnTable(L *lua.LState, col *migrator.Column) *lua.LTable {
	table := L.NewTable()

	L.SetField(table, "nullable", L.NewFunction(func(L *lua.LState) int {
		col.Nullable()
		L.Push(table)
		return 1
	}))

	L.SetField(table, "default", L.NewFunction(func(L *lua.LState) int {
		val := L.Get(1)
		col.Default(luaValueToGo(val))
		L.Push(table)
		return 1
	}))

	L.SetField(table, "unique", L.NewFunction(func(L *lua.LState) int {
		col.Unique()
		L.Push(table)
		return 1
	}))

	L.SetField(table, "unsigned", L.NewFunction(func(L *lua.LState) int {
		col.Unsigned()
		L.Push(table)
		return 1
	}))

	L.SetField(table, "autoIncrement", L.NewFunction(func(L *lua.LState) int {
		col.AutoIncrement()
		L.Push(table)
		return 1
	}))

	L.SetField(table, "primary", L.NewFunction(func(L *lua.LState) int {
		col.Primary()
		L.Push(table)
		return 1
	}))

	return table
}

// createForeignColumnTable crea métodos para ForeignColumn
func createForeignColumnTable(L *lua.LState, col *migrator.ForeignColumn) *lua.LTable {
	table := L.NewTable()

	L.SetField(table, "references", L.NewFunction(func(L *lua.LState) int {
		column := L.CheckString(1)
		col.References(column)
		L.Push(table)
		return 1
	}))

	L.SetField(table, "on", L.NewFunction(func(L *lua.LState) int {
		tableName := L.CheckString(1)
		col.On(tableName)
		L.Push(table)
		return 1
	}))

	L.SetField(table, "onDelete", L.NewFunction(func(L *lua.LState) int {
		action := L.CheckString(1)
		col.OnDelete(action)
		L.Push(table)
		return 1
	}))

	L.SetField(table, "nullable", L.NewFunction(func(L *lua.LState) int {
		col.Nullable()
		L.Push(table)
		return 1
	}))

	return table
}

func luaValueToGo(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case lua.LBool:
		return bool(v)
	case *lua.LNilType:
		return nil
	default:
		return nil
	}
}

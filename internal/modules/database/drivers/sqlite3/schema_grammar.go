package sqlite3

import (
	"bytes"
	"strconv"
	"sync"

	"github.com/vayload/vayload/internal/modules/database/grammar"
	"github.com/vayload/vayload/internal/modules/database/migrator"
	database_types "github.com/vayload/vayload/internal/modules/database/types"
)

type SQLiteSchemaGrammar struct{}

func NewSchemaGrammar() *SQLiteSchemaGrammar {
	return &SQLiteSchemaGrammar{}
}

var builderPool = sync.Pool{
	New: func() any {
		b := new(bytes.Buffer)
		b.Grow(256)
		return b
	},
}

func getBuilder() *bytes.Buffer {
	b := builderPool.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

func putBuilder(b *bytes.Buffer) {
	builderPool.Put(b)
}

func (g *SQLiteSchemaGrammar) CreateTable(bp migrator.Blueprint) (string, error) {
	buf := getBuilder()
	defer putBuilder(buf)

	buf.WriteString("CREATE TABLE ")
	buf.WriteString(bp.TableName)
	buf.WriteString(" (")

	first := true
	for _, c := range bp.Columns {
		if !first {
			buf.WriteString(", ")
		}
		first = false

		buf.WriteString(c.Name)
		buf.WriteByte(' ')
		CompileColType(buf, c.Name, database_types.DataType{
			Kind:      c.Type,
			Length:    c.Length,
			Precision: c.Precision,
			Scale:     c.Scale,
			Nullable:  c.IsNullable,
		})

		if c.IsPrimary {
			buf.WriteString(" PRIMARY KEY")
		}

		if c.IsAutoInc && c.Type == database_types.TypeBigInt {
			buf.WriteString(" AUTOINCREMENT")
		}

		if !c.IsNullable {
			buf.WriteString(" NOT NULL")
		}

		if c.IsUnique {
			buf.WriteString(" UNIQUE")
		}

		if c.DefaultVal != nil {
			buf.WriteString(" DEFAULT ")
			writeDefault(buf, c.DefaultVal)
		}
	}

	for _, fc := range bp.ForeignColumns {
		if !first {
			buf.WriteString(", ")
		}
		first = false

		buf.WriteString("FOREIGN KEY (")
		buf.WriteString(fc.Name)
		buf.WriteString(") REFERENCES ")
		buf.WriteString(fc.ForeignReference.Table)
		buf.WriteByte('(')
		buf.WriteString(fc.ForeignReference.Column)
		buf.WriteByte(')')

		if fc.OnDeleteAction != "" {
			buf.WriteString(" ON DELETE ")
			switch fc.OnDeleteAction {
			case "cascade":
				buf.WriteString("CASCADE")
			case "restrict":
				buf.WriteString("RESTRICT")
			case "set null":
				buf.WriteString("SET NULL")
			case "set default":
				buf.WriteString("SET DEFAULT")
			}
		}
	}

	buf.WriteString(");")

	return bytesToString(buf.Bytes()), nil
}

func (g *SQLiteSchemaGrammar) DropTable(name string, ifExists bool) (string, error) {
	buf := getBuilder()
	defer putBuilder(buf)

	buf.WriteString("DROP TABLE ")
	if ifExists {
		buf.WriteString("IF EXISTS ")
	}
	buf.WriteString(name)
	buf.WriteString(";")

	return bytesToString(buf.Bytes()), nil
}

func (g *SQLiteSchemaGrammar) AddColumn(table string, column *migrator.Column) (string, error) {
	buf := getBuilder()
	defer putBuilder(buf)

	buf.WriteString("ALTER TABLE ")
	buf.WriteString(table)
	buf.WriteString(" ADD COLUMN ")

	buf.WriteString(column.Name)
	buf.WriteByte(' ')
	CompileColType(buf, column.Name, database_types.DataType{
		Kind:      column.Type,
		Length:    column.Length,
		Precision: column.Precision,
		Scale:     column.Scale,
		Nullable:  column.IsNullable,
	})

	if !column.IsNullable {
		buf.WriteString(" NOT NULL")
	}

	if column.DefaultVal != nil {
		buf.WriteString(" DEFAULT ")
		writeDefault(buf, column.DefaultVal)
	}

	buf.WriteString(";")

	return bytesToString(buf.Bytes()), nil
}

func (g *SQLiteSchemaGrammar) DropColumn(table string, columnName string) (string, error) {
	buf := getBuilder()
	defer putBuilder(buf)

	buf.WriteString("ALTER TABLE ")
	buf.WriteString(table)
	buf.WriteString(" DROP COLUMN ")
	buf.WriteString(columnName)
	buf.WriteString(";")

	return bytesToString(buf.Bytes()), nil
}

func (g *SQLiteSchemaGrammar) RenameColumn(table string, from string, to string) (string, error) {
	buf := getBuilder()
	defer putBuilder(buf)

	buf.WriteString("ALTER TABLE ")
	buf.WriteString(table)
	buf.WriteString(" RENAME COLUMN ")
	buf.WriteString(from)
	buf.WriteString(" TO ")
	buf.WriteString(to)
	buf.WriteString(";")

	return bytesToString(buf.Bytes()), nil
}

func (g *SQLiteSchemaGrammar) DropIndex(_ string, name string) (string, error) {
	buf := getBuilder()
	defer putBuilder(buf)

	buf.WriteString("DROP INDEX ")
	buf.WriteString(name)
	buf.WriteString(";")

	return bytesToString(buf.Bytes()), nil
}

func (g *SQLiteSchemaGrammar) HasTable(table string) (string, error) {
	buf := getBuilder()
	defer putBuilder(buf)

	buf.WriteString("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='")
	buf.WriteString(table)
	buf.WriteString("';")

	return bytesToString(buf.Bytes()), nil
}

func (g *SQLiteSchemaGrammar) HasColumn(table string, column string) (string, error) {
	buf := getBuilder()
	defer putBuilder(buf)

	buf.WriteString("SELECT COUNT(*) FROM pragma_table_info('")
	buf.WriteString(table)
	buf.WriteString("') WHERE name='")
	buf.WriteString(column)
	buf.WriteString("';")

	return bytesToString(buf.Bytes()), nil
}

func (g *SQLiteSchemaGrammar) AddIndex(table string, columns []string, name string, unique bool) (string, error) {
	buf := getBuilder()
	defer putBuilder(buf)

	buf.WriteString("CREATE ")

	if unique {
		buf.WriteString("UNIQUE ")
	}

	buf.WriteString("INDEX ")
	buf.WriteString(name)
	buf.WriteString(" ON ")
	buf.WriteString(table)
	buf.WriteString(" (")

	for i, c := range columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(c)
	}

	buf.WriteString(");")

	return bytesToString(buf.Bytes()), nil
}

func writeDefault(b *bytes.Buffer, val any) {
	switch v := val.(type) {

	case string:
		if v == "now()" {
			b.WriteString("CURRENT_TIMESTAMP")
			return
		}

		b.WriteByte('\'')
		b.WriteString(v)
		b.WriteByte('\'')

	case int:
		var tmp [20]byte
		b.Write(strconv.AppendInt(tmp[:0], int64(v), 10))

	case int64:
		var tmp [20]byte
		b.Write(strconv.AppendInt(tmp[:0], v, 10))

	case float64:
		var tmp [20]byte
		b.Write(strconv.AppendFloat(tmp[:0], v, 'f', -1, 64))

	case bool:
		if v {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
	}
}

var _ grammar.SchemaGrammar = (*SQLiteSchemaGrammar)(nil)

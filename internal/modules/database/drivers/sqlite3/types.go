package sqlite3

import (
	"bytes"
	"strconv"

	database_types "github.com/vayload/vayload/internal/modules/database/types"
)

var TypeMap map[database_types.DataTypeKind]string
var LookupMap map[string]database_types.DataTypeKind

func init() {
	// SQLite does not support all data types, so we map them to the closest equivalent
	// See: https://www.sqlite.org/datatype3.html
	TypeMap = map[database_types.DataTypeKind]string{
		// ===== NUMERIC =====
		database_types.TypeTinyInt:   "INTEGER",
		database_types.TypeSmallInt:  "INTEGER",
		database_types.TypeInt:       "INTEGER",
		database_types.TypeBigInt:    "INTEGER",
		database_types.TypeSerial:    "INTEGER",
		database_types.TypeBigSerial: "INTEGER",

		database_types.TypeFloat:   "REAL",
		database_types.TypeDouble:  "REAL",
		database_types.TypeReal:    "REAL",
		database_types.TypeDecimal: "NUMERIC",
		database_types.TypeNumeric: "NUMERIC",

		// ===== BOOLEAN =====
		database_types.TypeBool: "INTEGER",

		// ===== TEXT =====
		database_types.TypeChar:    "TEXT",
		database_types.TypeVarchar: "TEXT",
		database_types.TypeText:    "TEXT",
		database_types.TypeCitext:  "TEXT",

		// ===== DATE / TIME =====
		database_types.TypeDate:        "TEXT",
		database_types.TypeTime:        "TEXT",
		database_types.TypeDateTime:    "TEXT",
		database_types.TypeTimestamp:   "TEXT",
		database_types.TypeTimestampTZ: "TEXT",
		database_types.TypeInterval:    "TEXT",

		// ===== BINARY =====
		database_types.TypeBinary:    "BLOB",
		database_types.TypeVarBinary: "BLOB",
		database_types.TypeBlob:      "BLOB",
		database_types.TypeBytea:     "BLOB",

		// ===== STRUCTURED =====
		database_types.TypeJSON: "TEXT",
		database_types.TypeUUID: "TEXT",
		database_types.TypeEnum: "TEXT",

		// ===== NETWORK / SPECIAL =====
		database_types.TypeInet:    "TEXT",
		database_types.TypeCIDR:    "TEXT",
		database_types.TypeMacAddr: "TEXT",
		database_types.TypeXML:     "TEXT",
	}

	LookupMap = make(map[string]database_types.DataTypeKind)
	for k, v := range TypeMap {
		LookupMap[v] = k
	}
}

func GetType(kind database_types.DataTypeKind) string {
	return TypeMap[kind]
}

func GetKind(typ string) database_types.DataTypeKind {
	return LookupMap[typ]
}

func CompileColType(buf *bytes.Buffer, col string, dt database_types.DataType) {
	base, ok := TypeMap[dt.Kind]
	if !ok {
		return
	}

	switch dt.Kind {
	case database_types.TypeEnum:
		buf.WriteString("TEXT CHECK (")
		buf.WriteString(col)
		buf.WriteString(" IN (")
		for i, v := range dt.EnumValues {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteByte('\'')
			buf.WriteString(v)
			buf.WriteByte('\'')
		}
		buf.WriteString(")")
		return

	// ===== CASE-INSENSITIVE TEXT =====
	case database_types.TypeDecimal, database_types.TypeNumeric:
		buf.WriteString("DECIMAL(")
		var tmp [20]byte
		buf.Write(strconv.AppendInt(tmp[:0], int64(dt.Precision), 10))
		buf.WriteString(",")
		buf.Write(strconv.AppendInt(tmp[:0], int64(dt.Scale), 10))
		buf.WriteString(")")
		return

	// ===== UUID =====
	case database_types.TypeUUID:
		buf.WriteString("TEXT CHECK (length(")
		buf.WriteString(col)
		buf.WriteString(") = 36)")
		return
	}

	buf.WriteString(base)
}

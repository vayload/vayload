package postgres

import (
	"fmt"
	"strings"

	database_types "github.com/vayload/vayload/internal/modules/database/types"
)

var TypeMap map[database_types.DataTypeKind]string

func init() {
	// See: https://www.postgresql.org/docs/18/datatype.html
	TypeMap = map[database_types.DataTypeKind]string{
		// ===== NUMERIC =====
		database_types.TypeTinyInt:   "SMALLINT",
		database_types.TypeSmallInt:  "SMALLINT",
		database_types.TypeInt:       "INTEGER",
		database_types.TypeBigInt:    "BIGINT",
		database_types.TypeSerial:    "SERIAL",
		database_types.TypeBigSerial: "BIGSERIAL",

		database_types.TypeFloat:   "REAL",
		database_types.TypeDouble:  "DOUBLE PRECISION",
		database_types.TypeReal:    "REAL",
		database_types.TypeDecimal: "DECIMAL",
		database_types.TypeNumeric: "NUMERIC",

		// ===== BOOLEAN =====
		database_types.TypeBool: "BOOLEAN",

		// ===== TEXT =====
		database_types.TypeChar:    "CHAR",
		database_types.TypeVarchar: "VARCHAR",
		database_types.TypeText:    "TEXT",
		database_types.TypeCitext:  "CITEXT",

		// ===== DATE / TIME =====
		database_types.TypeDate:        "DATE",
		database_types.TypeTime:        "TIME",
		database_types.TypeDateTime:    "TIMESTAMP",
		database_types.TypeTimestamp:   "TIMESTAMP",
		database_types.TypeTimestampTZ: "TIMESTAMPTZ",
		database_types.TypeInterval:    "INTERVAL",

		// ===== BINARY =====
		database_types.TypeBinary:    "BYTEA",
		database_types.TypeVarBinary: "BYTEA",
		database_types.TypeBlob:      "BYTEA",
		database_types.TypeBytea:     "BYTEA",

		// ===== STRUCTURED =====
		database_types.TypeJSON: "JSONB",
		database_types.TypeUUID: "UUID",
		database_types.TypeEnum: "ENUM",

		// ===== NETWORK / SPECIAL =====
		database_types.TypeInet:    "INET",
		database_types.TypeCIDR:    "CIDR",
		database_types.TypeMacAddr: "MACADDR",
		database_types.TypeXML:     "XML",
	}
}

func GetType(kind database_types.DataTypeKind) string {
	return TypeMap[kind]
}

func CompileColType(col string, dt database_types.DataType) string {
	base, ok := TypeMap[dt.Kind]
	if !ok {
		return "TEXT"
	}

	switch dt.Kind {

	case database_types.TypeEnum:
		// TODO: Implement ENUM type, use buffer to avoid allocs
		return fmt.Sprintf(
			"TEXT CHECK (%s IN (%s))",
			col,
			strings.Join(database_types.QuoteAll(dt.EnumValues), ", "),
		)

	// ===== CASE-INSENSITIVE TEXT =====
	case database_types.TypeCitext:
		return "TEXT COLLATE NOCASE"

	// ===== UUID =====
	case database_types.TypeChar, database_types.TypeVarchar:
		if dt.Length > 0 {
			return fmt.Sprintf("%s(%d)", base, dt.Length)
		}
		return base
	}

	return base
}

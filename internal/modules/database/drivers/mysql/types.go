package mysql

import (
	"fmt"
	"strings"

	database_types "github.com/vayload/vayload/internal/modules/database/types"
)

var TypeMap map[database_types.DataTypeKind]string

func init() {
	// See: https://dev.mysql.com/doc/refman/8.4/en/data-types.html
	TypeMap = map[database_types.DataTypeKind]string{
		// ===== NUMERIC =====
		database_types.TypeTinyInt:   "TINYINT",
		database_types.TypeSmallInt:  "SMALLINT",
		database_types.TypeInt:       "INT",
		database_types.TypeBigInt:    "BIGINT",
		database_types.TypeSerial:    "BIGINT AUTO_INCREMENT",
		database_types.TypeBigSerial: "BIGINT AUTO_INCREMENT",

		database_types.TypeFloat:   "FLOAT",
		database_types.TypeDouble:  "DOUBLE",
		database_types.TypeReal:    "REAL",
		database_types.TypeDecimal: "DECIMAL",
		database_types.TypeNumeric: "NUMERIC",

		// ===== BOOLEAN =====
		database_types.TypeBool: "TINYINT(1)",

		// ===== TEXT =====
		database_types.TypeChar:    "CHAR",
		database_types.TypeVarchar: "VARCHAR",
		database_types.TypeText:    "TEXT",
		database_types.TypeCitext:  "TEXT",

		// ===== DATE / TIME =====
		database_types.TypeDate:        "DATE",
		database_types.TypeTime:        "TIME",
		database_types.TypeDateTime:    "DATETIME",
		database_types.TypeTimestamp:   "TIMESTAMP",
		database_types.TypeTimestampTZ: "TIMESTAMP",
		database_types.TypeInterval:    "VARCHAR(255)",

		// ===== BINARY =====
		database_types.TypeBinary:    "BINARY",
		database_types.TypeVarBinary: "VARBINARY",
		database_types.TypeBlob:      "BLOB",
		database_types.TypeBytea:     "BLOB",

		// ===== STRUCTURED =====
		database_types.TypeJSON: "JSON",
		database_types.TypeUUID: "CHAR(36)",
		database_types.TypeEnum: "ENUM",

		// ===== NETWORK / SPECIAL =====
		database_types.TypeInet:    "VARCHAR(45)",
		database_types.TypeCIDR:    "VARCHAR(45)",
		database_types.TypeMacAddr: "VARCHAR(17)",
		database_types.TypeXML:     "TEXT",
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
		return fmt.Sprintf("ENUM (%s)", strings.Join(database_types.QuoteAll(dt.EnumValues), ", "))
	// ===== CASE-INSENSITIVE TEXT =====
	case database_types.TypeCitext:
		return "TEXT COLLATE NOCASE"

	// ===== UUID =====
	case database_types.TypeDecimal, database_types.TypeNumeric:
		if dt.Precision > 0 {
			return fmt.Sprintf("NUMERIC(%d,%d)", dt.Precision, dt.Scale)
		}
		return "NUMERIC"

	case database_types.TypeChar, database_types.TypeVarchar:
		if dt.Length > 0 {
			return fmt.Sprintf("%s(%d)", base, dt.Length)
		}
		return base

	case database_types.TypeBinary, database_types.TypeVarBinary:
		if dt.Length > 0 {
			return fmt.Sprintf("%s(%d)", base, dt.Length)
		}
		return base
	}

	return base
}

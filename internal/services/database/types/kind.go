package database_types

import (
	"fmt"
	"strings"
)

type DataTypeKind int

const (
	TypeUnknown DataTypeKind = iota

	// ===== NUMERIC =====
	TypeTinyInt
	TypeSmallInt
	TypeInt
	TypeBigInt
	TypeSerial
	TypeBigSerial
	TypeFloat
	TypeDouble
	TypeReal
	TypeDecimal
	TypeNumeric

	// ===== BOOLEAN =====
	TypeBool

	// ===== TEXT =====
	TypeChar
	TypeVarchar
	TypeText
	TypeCitext

	// ===== DATE / TIME =====
	TypeDate
	TypeTime
	TypeDateTime
	TypeTimestamp
	TypeTimestampTZ
	TypeInterval

	// ===== BINARY =====
	TypeBinary
	TypeVarBinary
	TypeBlob
	TypeBytea

	// ===== STRUCTURED =====
	TypeJSON
	TypeUUID
	TypeEnum

	// ===== NETWORK / SPECIAL =====
	TypeInet
	TypeCIDR
	TypeMacAddr
	TypeXML
)

type DataType struct {
	Kind       DataTypeKind
	Length     int
	Precision  int
	Scale      int
	EnumValues []string
	Nullable   bool
}

type PostgresDialect struct{}

func (PostgresDialect) ColumnTypeSQL(_ string, t DataType) string {
	switch t.Kind {

	case TypeSmallInt:
		return "SMALLINT"
	case TypeInt:
		return "INTEGER"
	case TypeBigInt:
		return "BIGINT"
	case TypeSerial:
		return "SERIAL"
	case TypeBigSerial:
		return "BIGSERIAL"

	case TypeDecimal, TypeNumeric:
		if t.Precision > 0 {
			return fmt.Sprintf("NUMERIC(%d,%d)", t.Precision, t.Scale)
		}
		return "NUMERIC"

	case TypeReal:
		return "REAL"
	case TypeDouble:
		return "DOUBLE PRECISION"

	case TypeBool:
		return "BOOLEAN"

	case TypeChar:
		return fmt.Sprintf("CHAR(%d)", t.Length)
	case TypeVarchar:
		return fmt.Sprintf("VARCHAR(%d)", t.Length)
	case TypeText:
		return "TEXT"
	case TypeCitext:
		return "CITEXT"

	case TypeDate:
		return "DATE"
	case TypeTime:
		return "TIME"
	case TypeTimestamp:
		return "TIMESTAMP"
	case TypeTimestampTZ:
		return "TIMESTAMPTZ"
	case TypeInterval:
		return "INTERVAL"

	case TypeBytea:
		return "BYTEA"

	case TypeJSON:
		return "JSONB"
	case TypeUUID:
		return "UUID"
	case TypeEnum:
		return fmt.Sprintf("ENUM (%s)", strings.Join(quoteAll(t.EnumValues), ", "))

	case TypeInet:
		return "INET"
	case TypeCIDR:
		return "CIDR"
	case TypeMacAddr:
		return "MACADDR"
	case TypeXML:
		return "XML"

	default:
		return "TEXT"
	}
}

type MySQLDialect struct{}

func (MySQLDialect) ColumnTypeSQL(_ string, t DataType) string {
	switch t.Kind {

	case TypeTinyInt:
		return "TINYINT"
	case TypeSmallInt:
		return "SMALLINT"
	case TypeInt:
		return "INT"
	case TypeBigInt:
		return "BIGINT"

	case TypeDecimal, TypeNumeric:
		return fmt.Sprintf("DECIMAL(%d,%d)", t.Precision, t.Scale)

	case TypeFloat:
		return "FLOAT"
	case TypeDouble:
		return "DOUBLE"

	case TypeBool:
		return "BOOLEAN"

	case TypeChar:
		return fmt.Sprintf("CHAR(%d)", t.Length)
	case TypeVarchar:
		return fmt.Sprintf("VARCHAR(%d)", t.Length)
	case TypeText:
		return "TEXT"

	case TypeDate:
		return "DATE"
	case TypeTime:
		return "TIME"
	case TypeDateTime:
		return "DATETIME"
	case TypeTimestamp:
		return "TIMESTAMP"

	case TypeBinary:
		return fmt.Sprintf("BINARY(%d)", t.Length)
	case TypeVarBinary:
		return fmt.Sprintf("VARBINARY(%d)", t.Length)
	case TypeBlob:
		return "BLOB"

	case TypeJSON:
		return "JSON"
	case TypeUUID:
		return "CHAR(36)"
	case TypeEnum:
		return fmt.Sprintf("ENUM(%s)", strings.Join(quoteAll(t.EnumValues), ", "))

	default:
		return "TEXT"
	}
}

type SQLiteDialect struct{}

func (SQLiteDialect) ColumnTypeSQL(col string, t DataType) string {
	switch t.Kind {

	// ===== ENUM =====
	case TypeEnum:
		return fmt.Sprintf(
			"TEXT CHECK (%s IN (%s))",
			col,
			strings.Join(quoteAll(t.EnumValues), ", "),
		)

	// ===== CASE-INSENSITIVE TEXT =====
	case TypeCitext:
		return "TEXT COLLATE NOCASE"

	// ===== NETWORK TYPES =====
	case TypeInet, TypeCIDR, TypeMacAddr:
		return "TEXT"

	// ===== UUID =====
	case TypeUUID:
		return "TEXT CHECK (length(" + col + ") = 36)"

	// ===== JSON =====
	case TypeJSON:
		return "TEXT"

	// ===== NUMERIC =====
	case TypeBool, TypeTinyInt, TypeSmallInt, TypeInt, TypeBigInt:
		return "INTEGER"

	case TypeDecimal, TypeNumeric, TypeFloat, TypeDouble, TypeReal:
		return "REAL"

	// ===== BINARY =====
	case TypeBinary, TypeVarBinary, TypeBlob:
		return "BLOB"

	// ===== DATE / TIME =====
	case TypeDate, TypeTime, TypeDateTime, TypeTimestamp, TypeTimestampTZ:
		return "TEXT" // ISO8601

	default:
		return "TEXT"
	}
}

func quoteAll(vals []string) []string {
	out := make([]string, len(vals))
	for i, v := range vals {
		out[i] = fmt.Sprintf("'%s'", v)
	}
	return out
}

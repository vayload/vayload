package database_types

import "unsafe"

type DataTypeKind uint8

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

func QuoteAll(vals []string) []string {
	n := len(vals)
	out := make([]string, n)
	for i, v := range vals {
		b := make([]byte, len(v)+2)
		b[0] = '\''
		copy(b[1:], v)
		b[len(b)-1] = '\''
		out[i] = unsafeString(b)
	}

	return out
}

func unsafeString(b []byte) string {
	return unsafe.String(&b[0], len(b))
}

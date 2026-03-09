package json5

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

func Marshal(v any) ([]byte, error) {
	e := newEncoder()
	defer freeEncoder(e)
	err := e.marshal(v)
	if err != nil {
		return nil, err
	}

	result := make([]byte, e.buf.Len())
	copy(result, e.buf.Bytes())
	return result, nil
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	e := &Encoder{
		buf:    bytes.NewBuffer(make([]byte, 0, 512)),
		indent: indent,
	}
	err := e.marshal(v)
	if err != nil {
		return nil, err
	}

	return e.buf.Bytes(), nil
}

func Unmarshal(data []byte, v any) error {
	d := newDecoder(data)
	return d.decode(v)
}

var encPool = sync.Pool{
	New: func() any {
		return &Encoder{buf: bytes.NewBuffer(make([]byte, 0, 512))}
	},
}

type Encoder struct {
	buf    *bytes.Buffer
	indent string
	depth  int
}

func newEncoder() *Encoder {
	return encPool.Get().(*Encoder)
}

func freeEncoder(e *Encoder) {
	e.buf.Reset()
	encPool.Put(e)
}

// inline for fast operations
func unsafeString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func (e *Encoder) marshal(v any) error {
	rv := reflect.ValueOf(v)
	return e.reflectValue(rv)
}

func (e *Encoder) reflectValue(v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		e.writeString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.buf.WriteString(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		e.buf.WriteString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		if math.IsInf(f, 1) {
			e.buf.WriteString("Infinity")
		} else if math.IsInf(f, -1) {
			e.buf.WriteString("-Infinity")
		} else if math.IsNaN(f) {
			e.buf.WriteString("NaN")
		} else {
			b := strconv.AppendFloat(nil, f, 'g', -1, 64)
			e.buf.Write(b)
		}
	case reflect.Bool:
		if v.Bool() {
			e.buf.WriteString("true")
		} else {
			e.buf.WriteString("false")
		}
	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice && v.IsNil() {
			e.buf.WriteString("null")
			return nil
		}
		e.buf.WriteByte('[')
		count := v.Len()
		e.depth++
		for i := range count {
			if i > 0 {
				e.buf.WriteByte(',')
			}
			if e.indent != "" {
				e.buf.WriteByte('\n')
				e.buf.WriteString(strings.Repeat(e.indent, e.depth))
			}
			if err := e.reflectValue(v.Index(i)); err != nil {
				return err
			}
		}
		e.depth--
		if count > 0 && e.indent != "" {
			e.buf.WriteByte('\n')
			e.buf.WriteString(strings.Repeat(e.indent, e.depth))
		}
		e.buf.WriteByte(']')
	case reflect.Map:
		if v.IsNil() {
			e.buf.WriteString("null")
			return nil
		}
		e.buf.WriteByte('{')
		iter := v.MapRange()
		first := true
		e.depth++
		for iter.Next() {
			if !first {
				e.buf.WriteByte(',')
			}
			if e.indent != "" {
				e.buf.WriteByte('\n')
				e.buf.WriteString(strings.Repeat(e.indent, e.depth))
			}
			first = false
			k := iter.Key()
			if k.Kind() != reflect.String {
				return errors.New("JSON5 maps must have string keys")
			}
			keyStr := k.String()
			if isSimpleIdentifier(keyStr) {
				e.buf.WriteString(keyStr)
			} else {
				e.writeString(keyStr)
			}
			e.buf.WriteByte(':')
			if e.indent != "" {
				e.buf.WriteByte(' ')
			}
			if err := e.reflectValue(iter.Value()); err != nil {
				return err
			}
		}
		e.depth--
		if !first && e.indent != "" {
			e.buf.WriteByte('\n')
			e.buf.WriteString(strings.Repeat(e.indent, e.depth))
		}
		e.buf.WriteByte('}')
	case reflect.Struct:
		e.buf.WriteByte('{')
		t := v.Type()
		first := true
		e.depth++
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" { // Skip unexported
				continue
			}
			tag := field.Tag.Get("json")
			if tag == "-" {
				continue
			}
			name, opts, _ := strings.Cut(tag, ",")
			if name == "" {
				name = field.Name
			}

			val := v.Field(i)
			if strings.Contains(opts, "omitempty") && isEmptyValue(val) {
				continue
			}

			if !first {
				e.buf.WriteByte(',')
			}
			if e.indent != "" {
				e.buf.WriteByte('\n')
				e.buf.WriteString(strings.Repeat(e.indent, e.depth))
			}
			first = false

			if isSimpleIdentifier(name) {
				e.buf.WriteString(name)
			} else {
				e.writeString(name)
			}
			e.buf.WriteByte(':')
			if e.indent != "" {
				e.buf.WriteByte(' ')
			}
			if err := e.reflectValue(val); err != nil {
				return err
			}
		}
		e.depth--
		if !first && e.indent != "" {
			e.buf.WriteByte('\n')
			e.buf.WriteString(strings.Repeat(e.indent, e.depth))
		}
		e.buf.WriteByte('}')
	case reflect.Interface, reflect.Pointer:
		if v.IsNil() {
			e.buf.WriteString("null")
			return nil
		}
		return e.reflectValue(v.Elem())
	default:
		return fmt.Errorf("unsupported type: %v", v.Kind())
	}
	return nil
}

func (e *Encoder) writeString(s string) {
	e.buf.WriteByte('"')
	for i := range s {
		b := s[i]
		if b == '\\' || b == '"' || b < 0x20 {
			e.writeStringSlow(s[i:])
			return
		}
		e.buf.WriteByte(b)
	}
	e.buf.WriteByte('"')
}

func (e *Encoder) writeStringSlow(s string) {
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '\\':
			e.buf.WriteString(`\\`)
		case '"':
			e.buf.WriteString(`\"`)
		case '\n':
			e.buf.WriteString(`\n`)
		case '\r':
			e.buf.WriteString(`\r`)
		case '\t':
			e.buf.WriteString(`\t`)
		default:
			if c < 0x20 {
				e.buf.WriteString(fmt.Sprintf("\\u%04x", c))
			} else {
				e.buf.WriteByte(c)
			}
		}
	}
	e.buf.WriteByte('"')
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	}
	return false
}

func isSimpleIdentifier(key string) bool {
	if len(key) == 0 {
		return false
	}
	for i, ch := range key {
		if ch == '$' || ch == '_' {
			continue
		}
		if ch >= 'a' && ch <= 'z' {
			continue
		}
		if ch >= 'A' && ch <= 'Z' {
			continue
		}
		if i > 0 && ch >= '0' && ch <= '9' {
			continue
		}
		return false
	}
	return true
}

type Decoder struct {
	data []byte
	pos  int
	len  int
}

func newDecoder(data []byte) *Decoder {
	return &Decoder{data: data, pos: 0, len: len(data)}
}

func (d *Decoder) decode(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("json5: Unmarshal(non-pointer or nil)")
	}

	d.skipSpace()
	return d.parseValue(rv.Elem())
}

func (d *Decoder) parseValue(v reflect.Value) error {
	if d.pos >= d.len {
		return errors.New("unexpected end of data")
	}

	c := d.data[d.pos]

	switch c {
	case '{':
		return d.parseObject(v)
	case '[':
		return d.parseArray(v)
	case '"', '\'':
		s, err := d.parseString()
		if err != nil {
			return err
		}
		return d.assignString(v, s)
	case 't': // true
		if d.match("true") {
			return d.assignBool(v, true)
		}
	case 'f': // false
		if d.match("false") {
			return d.assignBool(v, false)
		}
	case 'n': // null
		if d.match("null") {
			return d.assignNull(v)
		}
	case 'N': // NaN
		if d.match("NaN") {
			return d.assignFloat(v, math.NaN())
		}
	case 'I': // Infinity
		if d.match("Infinity") {
			return d.assignFloat(v, math.Inf(1))
		}
	case '+': // +Infinity or number
		if d.pos+8 < d.len && unsafeString(d.data[d.pos:d.pos+9]) == "+Infinity" {
			d.pos += 9
			return d.assignFloat(v, math.Inf(1))
		}
		fallthrough
	case '-', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// Handle -Infinity
		if c == '-' && d.pos+8 < d.len && unsafeString(d.data[d.pos:d.pos+9]) == "-Infinity" {
			d.pos += 9
			return d.assignFloat(v, math.Inf(-1))
		}
		return d.parseNumber(v)
	}

	return fmt.Errorf("unexpected char '%c' at %d", c, d.pos)
}

func (d *Decoder) assignString(v reflect.Value, s string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Interface:
		v.Set(reflect.ValueOf(s))
	default:
		return fmt.Errorf("cannot unmarshal string into %v", v.Type())
	}
	return nil
}

func (d *Decoder) assignBool(v reflect.Value, b bool) error {
	switch v.Kind() {
	case reflect.Bool:
		v.SetBool(b)
	case reflect.Interface:
		v.Set(reflect.ValueOf(b))
	default:
		return fmt.Errorf("cannot unmarshal bool into %v", v.Type())
	}
	return nil
}

func (d *Decoder) assignNull(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Interface, reflect.Pointer, reflect.Map, reflect.Slice:
		v.SetZero()
	}
	return nil
}

func (d *Decoder) assignFloat(v reflect.Value, f float64) error {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		v.SetFloat(f)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(f))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(f))
	case reflect.Interface:
		v.Set(reflect.ValueOf(f))
	default:
		return fmt.Errorf("cannot unmarshal number into %v", v.Type())
	}
	return nil
}

func (d *Decoder) parseObject(v reflect.Value) error {
	d.pos++
	d.skipSpace()

	if v.Kind() == reflect.Interface && v.IsNil() {
		m := make(map[string]any)
		v.Set(reflect.ValueOf(m))
		v = v.Elem()
	} else if v.Kind() == reflect.Map && v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}

	for d.pos < d.len {
		if d.data[d.pos] == '}' {
			d.pos++
			return nil
		}

		// Parse Key
		key, err := d.parseKey()
		if err != nil {
			return err
		}

		d.skipSpace()
		if d.pos >= d.len || d.data[d.pos] != ':' {
			return errors.New("expected ':'")
		}
		d.pos++ // consume ':'
		d.skipSpace()

		// Set Value
		if v.Kind() == reflect.Map {
			elemType := v.Type().Elem()
			newVal := reflect.New(elemType).Elem()
			if err := d.parseValue(newVal); err != nil {
				return err
			}
			v.SetMapIndex(reflect.ValueOf(key), newVal)
		} else if v.Kind() == reflect.Struct {
			// Fast path for field lookup
			field := d.findStructField(v, key)
			if field.IsValid() && field.CanSet() {
				if err := d.parseValue(field); err != nil {
					return err
				}
			} else {
				// Ignorar campo desconocido, pero debemos consumir el valor
				var dummy any
				if err := d.parseValue(reflect.ValueOf(&dummy).Elem()); err != nil {
					return err
				}
			}
		} else {
			// Interface or unsupported
			var val any
			if err := d.parseValue(reflect.ValueOf(&val).Elem()); err != nil {
				return err
			}
		}

		d.skipSpace()
		if d.pos < d.len && d.data[d.pos] == ',' {
			d.pos++
			d.skipSpace()
		}
	}
	return errors.New("expected '}'")
}

func (d *Decoder) findStructField(v reflect.Value, key string) reflect.Value {
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		tagName, _, _ := strings.Cut(tag, ",")
		if tagName == key {
			return v.Field(i)
		}
	}
	if f := v.FieldByName(key); f.IsValid() {
		return f
	}
	return v.FieldByNameFunc(func(n string) bool {
		return strings.EqualFold(n, key)
	})
}

func (d *Decoder) parseArray(v reflect.Value) error {
	d.pos++
	d.skipSpace()

	if v.Kind() == reflect.Interface && v.IsNil() {
		arr := make([]any, 0)
		v.Set(reflect.ValueOf(arr))
		v = v.Elem()
	}

	i := 0
	for d.pos < d.len {
		if d.data[d.pos] == ']' {
			d.pos++
			return nil
		}

		// Preparar destino
		var currentVal reflect.Value
		if v.Kind() == reflect.Slice {
			if i >= v.Cap() {
				newCap := v.Cap() + v.Cap()/2
				if newCap < 4 {
					newCap = 4
				}
				newSlice := reflect.MakeSlice(v.Type(), v.Len(), newCap)
				reflect.Copy(newSlice, v)
				v.Set(newSlice)
			}
			if i >= v.Len() {
				v.SetLen(i + 1)
			}
			currentVal = v.Index(i)
		} else if v.Kind() == reflect.Array {
			if i < v.Len() {
				currentVal = v.Index(i)
			} else {
				// Array lleno, necesitamos consumir y descartar
				var dummy any
				currentVal = reflect.ValueOf(&dummy).Elem()
			}
		} else {
			// Interface genérica dentro de array
			var dummy any
			currentVal = reflect.ValueOf(&dummy).Elem()
		}

		if err := d.parseValue(currentVal); err != nil {
			return err
		}

		// Si estábamos trabajando con []any (interface), debemos hacer append manual si v era interface{} original
		// (Nota: la lógica arriba maneja reflect.Slice directamente sobre v, que es un puntero al valor original si se pasó correctamente)

		i++
		d.skipSpace()
		if d.pos < d.len && d.data[d.pos] == ',' {
			d.pos++
			d.skipSpace()
		}
	}
	return errors.New("expected ']'")
}

func (d *Decoder) parseKey() (string, error) {
	c := d.data[d.pos]
	if c == '"' || c == '\'' {
		return d.parseString()
	}
	start := d.pos
	for d.pos < d.len {
		c := d.data[d.pos]
		if isIdentifierChar(c) {
			d.pos++
		} else {
			break
		}
	}
	if d.pos == start {
		return "", errors.New("invalid key")
	}

	return string(d.data[start:d.pos]), nil
}

func isIdentifierChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '$' || c > 127
}

func (d *Decoder) parseString() (string, error) {
	quote := d.data[d.pos]
	d.pos++
	start := d.pos

	for d.pos < d.len {
		c := d.data[d.pos]
		if c == quote {
			s := d.data[start:d.pos]
			d.pos++
			return string(s), nil
		}
		if c == '\\' || c == '\n' || c == '\r' {
			break
		}
		d.pos++
	}

	var sb strings.Builder
	sb.Write(d.data[start:d.pos])

	for d.pos < d.len {
		c := d.data[d.pos]
		if c == quote {
			d.pos++
			return sb.String(), nil
		}
		if c == '\n' || c == '\r' {
			if c == '\r' && d.pos+1 < d.len && d.data[d.pos+1] == '\n' {
				d.pos++
			}
			return "", errors.New("unescaped newline in string")
		}

		if c == '\\' {
			d.pos++
			if d.pos >= d.len {
				return "", errors.New("unexpected end of string")
			}
			esc := d.data[d.pos]
			switch esc {
			case 'b':
				sb.WriteByte('\b')
			case 'f':
				sb.WriteByte('\f')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case 'v':
				sb.WriteByte('\v')
			case '0':
				sb.WriteByte(0)
			case '\n':
			case '\r':
				if d.pos+1 < d.len && d.data[d.pos+1] == '\n' {
					d.pos++
				}
			case 'x':
				d.pos++
				if d.pos+2 > d.len {
					return "", errors.New("invalid hex escape")
				}
				val, err := strconv.ParseUint(unsafeString(d.data[d.pos:d.pos+2]), 16, 8)
				if err != nil {
					return "", err
				}
				sb.WriteByte(byte(val))
				d.pos++
			case 'u':
				d.pos++
				if d.pos+4 > d.len {
					return "", errors.New("invalid unicode escape")
				}
				val, err := strconv.ParseUint(unsafeString(d.data[d.pos:d.pos+4]), 16, 32)
				if err != nil {
					return "", err
				}
				sb.WriteRune(rune(val))
				d.pos += 3
			default:
				sb.WriteByte(esc)
			}
		} else {
			sb.WriteByte(c)
		}
		d.pos++
	}
	return "", errors.New("unexpected end of string")
}

func (d *Decoder) parseNumber(v reflect.Value) error {
	start := d.pos

	isHex := false
	if d.data[d.pos] == '0' && d.pos+1 < d.len && (d.data[d.pos+1] == 'x' || d.data[d.pos+1] == 'X') {
		isHex = true
		d.pos += 2
		for d.pos < d.len {
			c := d.data[d.pos]
			if isHexDigit(c) {
				d.pos++
			} else {
				break
			}
		}
	} else {
		if d.data[d.pos] == '-' || d.data[d.pos] == '+' {
			d.pos++
		}
		for d.pos < d.len {
			c := d.data[d.pos]
			if c >= '0' && c <= '9' {
				d.pos++
			} else if c == '.' || c == 'e' || c == 'E' || c == '-' || c == '+' {
				d.pos++
			} else {
				break
			}
		}
	}

	numStr := unsafeString(d.data[start:d.pos])

	if isHex {
		i, err := strconv.ParseInt(numStr[2:], 16, 64)
		if err != nil {
			return err
		}
		return d.assignFloat(v, float64(i))
	}

	f, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return err
	}
	return d.assignFloat(v, f)
}

func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func (d *Decoder) skipSpace() {
	for d.pos < d.len {
		c := d.data[d.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f' || c == 0xA0 {
			d.pos++
			continue
		}
		if c == '/' && d.pos+1 < d.len {
			next := d.data[d.pos+1]
			if next == '/' {
				d.pos += 2
				for d.pos < d.len && d.data[d.pos] != '\n' && d.data[d.pos] != '\r' {
					d.pos++
				}
				continue
			} else if next == '*' {
				d.pos += 2
				for d.pos+1 < d.len {
					if d.data[d.pos] == '*' && d.data[d.pos+1] == '/' {
						d.pos += 2
						break
					}
					d.pos++
				}
				continue
			}
		}
		break
	}
}

func (d *Decoder) match(s string) bool {
	if d.pos+len(s) <= d.len && unsafeString(d.data[d.pos:d.pos+len(s)]) == s {
		d.pos += len(s)
		return true
	}
	return false
}

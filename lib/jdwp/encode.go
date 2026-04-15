package jdwp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"
)

func Unmashal(bf []byte, v any) error {
	r := bytes.NewReader(bf)
	return decode(r, reflect.ValueOf(v))
}

func Decode(r io.Reader, v any) error {
	return decode(r, reflect.ValueOf(v))
}

func decode(r io.Reader, v reflect.Value) error {
	switch v.Kind() {

	// Pointer
	case reflect.Pointer:
		return decode(r, v.Elem())

	// Numbers
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if err := binary.Read(r, binary.BigEndian, v.Addr().Interface()); err != nil {
			return fmt.Errorf("jdwp slice read number")
		}
		return nil

	// Struct
	case reflect.Struct:
		t := v.Type()
		for i := range v.NumField() {
			if SkipField(t.Field(i)) {
				continue
			}
			if err := decode(r, v.Field(i)); err != nil {
				return fmt.Errorf("jdwp unmashal field[%d]: %w", i, err)
			}
		}
		return nil

	// Slice
	case reflect.Slice:
		t := v.Type()
		var length uint32
		if err := binary.Read(r, binary.BigEndian, &length); err != nil {
			return fmt.Errorf("jdwp slice read length")
		}
		for i := range length {
			el := reflect.New(t.Elem()).Elem()
			if err := decode(r, el); err != nil {
				return fmt.Errorf("jdwp slice read (%d/%d): %w", i, length, err)
			}
			v.Set(reflect.Append(v, el))
		}
		return nil

	// String
	case reflect.String:
		var length uint32
		if err := binary.Read(r, binary.BigEndian, &length); err != nil {
			return fmt.Errorf("jdwp string read length")
		}
		str := make([]byte, length)
		if n, err := r.Read(str); n != int(length) || err != nil {
			return fmt.Errorf("jdwp string read bytes (%d/%d): %w", n, length, err)
		}
		v.SetString(string(str))
		return nil

	// Unsupported
	default:
		return fmt.Errorf("jdwp unmashal: unsupported type %s", v.Kind())
	}
}

func SkipField(s reflect.StructField) bool {
	tags := strings.Split(s.Tag.Get("jdwp"), ",")
	return slices.Contains(tags, "ignore")
}

func Encode(w io.Writer, command any) error {
	// TODO
	panic("not implemented")
}

func Marshal(any) ([]byte, error) {
	// TODO
	panic("not implemented")
}

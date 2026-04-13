package jdwp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func Unmashal(bf []byte, v any) error {
	r := bytes.NewReader(bf)
	return unmashal(r, reflect.ValueOf(v))
}

func unmashal(r io.Reader, v reflect.Value) error {
	switch v.Kind() {

	// Pointer
	case reflect.Pointer:
		return unmashal(r, v.Elem())

	// Numbers
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if err := binary.Read(r, binary.BigEndian, v.Addr().Interface()); err != nil {
			return fmt.Errorf("jdwp slice read number")
		}
		return nil

	// Struct
	case reflect.Struct:
		for i := range v.NumField() {
			if err := unmashal(r, v.Field(i)); err != nil {
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
			if err := unmashal(r, el); err != nil {
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

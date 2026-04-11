package jdwp

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf8"
)

func binEnc(buf []byte, data any) (int, error) {
	return binary.Encode(buf, binary.BigEndian, data)
}

func binDec(buf []byte, data any) (int, error) {
	return binary.Decode(buf, binary.BigEndian, data)
}

func binRead(r io.Reader, data any) error {
	return binary.Read(r, binary.BigEndian, data)
}

func binWrite(w io.Writer, data any) error {
	return binary.Write(w, binary.BigEndian, data)
}

func binWriteBytes(w io.Writer, data []byte) error {
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("jwdp binWrite: %w", err)
	}
	return nil
}

func encodeString(v string) []byte {
	size := int32(utf8.RuneCountInString(v))
	buf := make([]byte, 0, size+4)
	if _, err := binEnc(buf[0:4], size); err != nil {
		panic(fmt.Errorf("jdwp encode string size: %w", err))
	}
	if _, err := binEnc(buf[4:], v); err != nil {
		panic(fmt.Errorf("jdwp encode string value: %w", err))
	}
	return buf
}

func encodeInt(v int32) []byte {
	buf := make([]byte, 4)
	if _, err := binEnc(buf, v); err != nil {
		panic(fmt.Errorf("jdwp encode int: %w", err))
	}
	return buf
}

func parseString(r io.Reader) (string, error) {
	var size int32
	if err := binRead(r, &size); err != nil {
		return "", fmt.Errorf("jdwp decode string size: %w", err)
	}
	buf := make([]byte, size)
	if err := binRead(r, buf); err != nil {
		return "", fmt.Errorf("jdwp decode string body: %w", err)
	}
	return string(buf), nil
}

func parseInt(r io.Reader) (int32, error) {
	var value int32
	if err := binRead(r, &value); err != nil {
		return 0, fmt.Errorf("jdwp parse int: %w", err)
	}
	return value, nil
}

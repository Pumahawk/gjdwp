package jdwp

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type Mytype struct {
	Id int32
}

type Group struct {
	N1 Mytype
	N2 Mytype
}

func TestUnmarshalNumber(t *testing.T) {
	bf := prepareBuffer(int32(10))
	var v int32
	err := Unmashal(bf, &v)
	if err != nil {
		t.Error(err)
	}
	if v != 10 {
		t.Errorf("v expected 10, got %d", v)
	}
}

func TestUnmarshalStruct(t *testing.T) {
	// Basic field
	{
		bf := prepareBuffer(int32(11), int32(22))
		type str struct {
			A int32
			B int32
		}
		var v str
		if err := Unmashal(bf, &v); err != nil {
			t.Error(err)
		}
		if v.A != 11 {
			t.Errorf("v.A expected 11, got %d", v.A)
		}
		if v.B != 22 {
			t.Errorf("v.B expected 22, got %d", v.B)
		}
	}

	// Ignore field
	{
		bf := prepareBuffer(int32(11), int32(22))
		type str struct {
			A int32 `jdwp:"ignore"`
			B int32
			C int32 `jdwp:"mock1,ignore,mock2"`
			D int32
		}
		var v str
		if err := Unmashal(bf, &v); err != nil {
			t.Error(err)
		}
		if v.A != 0 {
			t.Errorf("v.A expected 0, got %d", v.A)
		}
		if v.B != 11 {
			t.Errorf("v.B expected 11, got %d", v.B)
		}
		if v.C != 0 {
			t.Errorf("v.C expected 0, got %d", v.B)
		}
		if v.D != 22 {
			t.Errorf("v.D expected 22, got %d", v.B)
		}
	}
}

func TestUnmarshalSlice(t *testing.T) {
	bf := prepareBuffer(int32(2), int32(13), int32(123))
	var s []int32
	if err := Unmashal(bf, &s); err != nil {
		t.Error(err)
		return
	}
	if len(s) != 2 {
		t.Errorf("len expected 2, got %d", len(s))
		return
	}
	if s[0] != 13 {
		t.Errorf("s[0] expected 13, got %d", s[0])
	}
	if s[1] != 123 {
		t.Errorf("s[0] expected 123, got %d", s[1])
	}
}

func TestUnmarshalString(t *testing.T) {
	bf := &bytes.Buffer{}
	binary.Write(bf, binary.BigEndian, int32(13))
	if n, err := bf.Write([]byte("Hello, World!")); n != 13 || err != nil {
		t.Errorf("read %d: %s", n, err)
	}
	var s string
	if err := Unmashal(bf.Bytes(), &s); err != nil {
		t.Error(err)
		return
	}
	if s != "Hello, World!" {
		t.Errorf("string expected %q, got %q", "Hello, World!", s)
		return
	}
}

func prepareBuffer(v ...any) []byte {
	bf := &bytes.Buffer{}
	for _, v := range v {
		if err := binary.Write(bf, binary.BigEndian, v); err != nil {
			panic(err)
		}
	}
	return bf.Bytes()
}

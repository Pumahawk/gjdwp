package jdwp

import (
	"bytes"
	"slices"
	"testing"
)

func TestReadPack(t *testing.T) {
	data := []byte("Hello, World!")
	length := int32(13)
	id := int32(17)
	flags := ReplyFlag
	errorcode := int16(130)
	bf := bytes.Buffer{}
	must(binWrite(&bf, length))
	must(binWrite(&bf, id))
	must(binWrite(&bf, flags))
	must(binWrite(&bf, errorcode))
	must(binWriteBytes(&bf, data))

	if pack, err := readPack(&bf); err != nil {
		t.Errorf("read pack with error: %v", err)
	} else {
		if id := pack.baseheader.Id; id != 17 {
			t.Errorf("pack.baseheader.Id == 17: id=%d", id)
		}
		if length := pack.baseheader.Length; length != 13 {
			t.Errorf("pack.baseheader.Length == 13: length=%d", id)
		}
		if length := pack.baseheader.Length; length != 13 {
			t.Errorf("pack.baseheader.Length == 13: length=%d", id)
		}
		if dataresp := pack.data; !slices.Equal(dataresp, data) {
			t.Errorf("slices.Equal(dataresp, data)")
		}
	}

}

func TestWritePack(t *testing.T) {
	data := []byte("Hello, World!")
	p := &pack{&baseheader{243, 2, ReplyFlag}, nil, &replyheader{0}, data}
	buf := bytes.Buffer{}
	if err := p.write(&buf); err != nil {
		t.Errorf("err != nil: %v", err)
	}
	var id int32
	if err := binRead(&buf, &id); err != nil || id != 243 {
		t.Errorf("id == 1: id=%d: %v", id, err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

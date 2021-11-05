package fxk

import (
	"bytes"
	"encoding/binary"
)

type Global struct {
	Name      string `json:"name"`
	IP        string `json:"ip"`
	IPVersion string `json:"network"`
	Port      int    `json:"port"`
}

var GlobalVar *Global

func init() {
	GlobalVar = &Global{
		Name:      "zink",
		IP:        "0.0.0.0",
		IPVersion: "tcp",
		Port:      9000,
	}
}

func Encode(tag uint32, data string) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, tag); err != nil {
		return nil, err
	}
	dataBuf := []byte(data)
	if err := binary.Write(buf, binary.LittleEndian, uint32(len(dataBuf))); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, dataBuf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func Decode(b []byte) (uint32, string, error) {
	buf := bytes.NewBuffer(b)
	var tag, length uint32
	if err := binary.Read(buf, binary.LittleEndian, &tag); err != nil {
		return 0, "", err
	}
	if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return 0, "", err
	}
	dataBuf := make([]byte, length)
	if err := binary.Read(buf, binary.LittleEndian, &dataBuf); err != nil {
		return 0, "", err
	}
	return tag, string(dataBuf), nil
}

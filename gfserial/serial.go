package gfserial

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// Serializable 接口函数
type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

// Serializer 接口函数
type Serializer interface {
	Serialize(val interface{}) ([]byte, error)
	Deserialize(raw []byte, val interface{}) error
}

type JsonSerializer struct{}

func (JsonSerializer) Serialize(val interface{}) ([]byte, error) {
	return json.Marshal(val)
}

func (JsonSerializer) Deserialize(raw []byte, val interface{}) error {
	return json.Unmarshal(raw, val)
}

type GobSerializer struct{}

func (GobSerializer) Serialize(val interface{}) ([]byte, error) {
	buff := &bytes.Buffer{}
	enc := gob.NewEncoder(buff)
	err := enc.Encode(val)
	return buff.Bytes(), err
}

func (GobSerializer) Deserialize(raw []byte, val interface{}) error {
	buff := bytes.NewReader(raw)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(val)
	return err
}

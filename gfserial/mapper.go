package gfserial

import "encoding/json"

// Mapper 接口函数
type Mapper interface {
	ToMap() (map[string]interface{}, error)
	FromMap(map[string]interface{}) error
}

type Bytes []byte

func (b *Bytes) ToMap() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(*b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (b *Bytes) FromMap(m map[string]interface{}) error {
	*b, _ = json.Marshal(m)
	return nil
}

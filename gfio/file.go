package gfio

import "io"

func ReadBytes(body io.Reader, len int) ([]byte, error) {
	al := 0
	bs := make([]byte, len)
	for al < len {
		n, err := body.Read(bs[al:len])
		if err != nil {
			return nil, err
		}
		al = al + n
	}
	return bs, nil
}

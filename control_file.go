package go_lpd

import (
	"bytes"
	"fmt"
	"io"
)

type ControlFile map[ControlFileCommand]string

func (c *ControlFile) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	for cmd, value := range *c {
		_, err := buf.Write([]byte{byte(cmd)})
		if err != nil {
			return nil, err
		}

		_, err = buf.Write([]byte(value))
		if err != nil {
			return nil, err
		}

		_, err = buf.Write([]byte(LineEnding))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func NewControlFileDecoder(r io.Reader) *ControlFileDecoder {
	return &ControlFileDecoder{r}
}

type ControlFileDecoder struct {
	reader io.Reader
}

func (c *ControlFileDecoder) Decode(size int) (ControlFile, error) {
	data := make([]byte, size)
	readed, err := c.reader.Read(data)
	if err != nil {
		return nil, err
	}
	if readed != size {
		return nil, fmt.Errorf("could not read %d bytes from reader", size)
	}

	cf := ControlFile{}

	for _, line := range bytes.Split(data, []byte(LineEnding)) {
		cf[ControlFileCommand(line[0])] = string(line[1:])
	}

	return cf, nil
}

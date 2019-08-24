package lpd

import (
	"bytes"
	"fmt"
	"io"
)

type ControlFile map[ControlFileCommand]string

func (c *ControlFile) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	enc := NewControlFileCommandEncoder(buf)

	for cmd, value := range *c {
		if err := enc.Encode(cmd, value); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func NewControlFileCommandEncoder(w io.Writer) *ControlFileCommandEncoder {
	return &ControlFileCommandEncoder{w}
}

type ControlFileCommandEncoder struct {
	writer io.Writer
}

func (e *ControlFileCommandEncoder) Encode(cmd ControlFileCommand, value string) error {
	_, err := e.writer.Write([]byte{byte(cmd)})
	if err != nil {
		return err
	}

	_, err = e.writer.Write([]byte(value))
	if err != nil {
		return err
	}

	_, err = e.writer.Write([]byte(LineEnding))
	if err != nil {
		return err
	}

	return nil
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
		if len(line) > 0 {
			cf[ControlFileCommand(line[0])] = string(line[1:])
		}
	}

	return cf, nil
}

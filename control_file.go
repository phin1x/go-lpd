package go_lpd

import (
	"bytes"
	"fmt"
	"io"
)

type ControlFile map[ControlFileCommand]string

func NewControlFileEncoder(w io.Writer) *ControlFileEncoder {
	return &ControlFileEncoder{w}
}

type ControlFileEncoder struct {
	w io.Writer
}

func (c *ControlFileEncoder) Encode(cf ControlFile) (err error) {
	for cmd, value := range cf {
		_, err = c.w.Write([]byte{byte(cmd)})
		if err != nil {
			return
		}

		_, err = c.w.Write([]byte(value))
		if err != nil {
			return
		}

		_, err = c.w.Write([]byte(LineEnding))
		if err != nil {
			return
		}
	}

	return nil
}

func NewControlFileDecoder(r io.Reader) *ControlFileDecoder {
	return &ControlFileDecoder{r}
}

type ControlFileDecoder struct {
	r io.Reader
}

func (c *ControlFileDecoder) Decode(cf ControlFile, size int) (err error) {
	data := make([]byte, size)
	readed, err := c.r.Read(data)
	if err != nil {
		return
	}
	if readed != size {
		return fmt.Errorf("could not read %d bytes from reader", size)
	}

	cf = ControlFile{}

	for _, line := range bytes.Split(data, []byte(LineEnding)) {
		cf[ControlFileCommand(line[0])] = string(line[1:])
	}

	return nil
}

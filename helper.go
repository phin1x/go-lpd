package lpd

import (
	"errors"
	"io"
	"strings"
)

func SendCommandLine(w io.Writer, cmd byte, opts []string) (err error) {
	var line []byte

	line = append(line, cmd)

	optString := strings.Join(opts, Separator)
	line = append(line, []byte(optString)...)
	line = append(line, []byte(LineEnding)...)

	_, err = w.Write(line)
	return
}

func CheckAcknowledge(r io.Reader) error {
	buf := make([]byte, 1)

	if _, err := r.Read(buf); err != nil {
		return err
	}

	if buf[0] != Acknowledge {
		return errors.New("server not acknowledged the command")
	}

	return nil
}

func SendAbortOnError(w io.Writer, err error) error {
	if err != nil {
		_, inErr := w.Write([]byte{byte(AbortJob)})
		if inErr != nil {
			return inErr
		}

		_, inErr = w.Write([]byte(LineEnding))
		if inErr != nil {
			return inErr
		}
	}

	return nil
}

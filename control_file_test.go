package lpd

import (
	"bytes"
	"testing"
)

var encodingTestCases = []struct {
	Cmd   ControlFileCommand
	Value string
	Bytes []byte
}{
	{
		Cmd: Hostname,
		Value: "myhost",
		Bytes: []byte{72, 109, 121, 104, 111, 115, 116, 10},
	},
	{
		Cmd: UserID,
		Value: "testuser",
		Bytes: []byte{80, 116, 101, 115, 116, 117, 115, 101, 114, 10},
	},
	{
		Cmd: JobName,
		Value: "test job",
		Bytes: []byte{74, 116, 101, 115, 116, 32, 106, 111, 98, 10},
	},
	{
		Cmd: BannerClass,
		Value: "myhost",
		Bytes: []byte{67, 109, 121, 104, 111, 115, 116, 10},
	},
	{
		Cmd: PrintBanner,
		Value: "testuser",
		Bytes: []byte{76, 116, 101, 115, 116, 117, 115, 101, 114, 10},
	},
	{
		Cmd: UnlinkDataFile,
		Value: "dfA000myhost",
		Bytes: []byte{85, 100, 102, 65, 48, 48, 48, 109, 121, 104, 111, 115, 116, 10},
	},
	{
		Cmd: SourceFileName,
		Value: "test job",
		Bytes: []byte{78, 116, 101, 115, 116, 32, 106, 111, 98, 10},
	},
}

var decodingTestCases = []struct {
	ControlFile ControlFile
	Bytes       []byte
}{
	{
		ControlFile: ControlFile{
			Hostname:       "myhost",
			UserID:         "testuser",
			JobName:        "test job",
			BannerClass:    "myhost",
			PrintBanner:    "testuser",
			UnlinkDataFile: "dfA000myhost",
			SourceFileName: "test job",
		},
		Bytes: []byte{72, 109, 121, 104, 111, 115, 116, 10, 80, 116, 101, 115, 116, 117, 115, 101, 114, 10, 74, 116, 101, 115, 116, 32, 106, 111, 98, 10, 67, 109, 121, 104, 111, 115, 116, 10, 76, 116, 101, 115, 116, 117, 115, 101, 114, 10, 85, 100, 102, 65, 48, 48, 48, 109, 121, 104, 111, 115, 116, 10, 78, 116, 101, 115, 116, 32, 106, 111, 98, 10},
	},
}

func TestControlFileEncoding(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := NewControlFileCommandEncoder(buf)

	for _, c := range encodingTestCases {
		if err := enc.Encode(c.Cmd, c.Value); err != nil {
			t.Errorf("error while encoding controfile command: %v", err)
		}

		result := buf.Bytes()

		if !bytes.Equal(result, c.Bytes) {
			t.Errorf("encoding result is not correct, expected %v, got %v", c.Bytes, result)
		}

		buf.Reset()
	}
}

func TestControlFileeDecoding(t *testing.T) {
	buf := new(bytes.Buffer)
	dec := NewControlFileDecoder(buf)

	for _, c := range decodingTestCases {
		buf.Write(c.Bytes)

		cf, err := dec.Decode(len(c.Bytes))
		if err != nil {
			t.Errorf("error while decoding bytes %v: %v", c.Bytes, err)
		}

		for cmd, value := range c.ControlFile {
			cfValue, exists := cf[cmd]

			if !exists {
				t.Errorf("attribute %v not found in decoded controlfile", cmd)
			}

			if cfValue != value {
				t.Errorf("decoded value is not correct, expected %v, got %v", cfValue, value)

			}
		}

		buf.Reset()
	}
}

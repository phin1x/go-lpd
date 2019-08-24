package go_lpd

import (
	"io"
	"net"
	"os"
	"os/user"
	"strconv"
)

type Document struct {
	Document io.Reader
	Size     int
	Name     string
}

type Client struct {
	// dest format is host:port
	dest string
}

func (c *Client) PrintFile(printer string,  doc Document, cf ControlFile) (err error) {
	// get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return
	}

	//get current user
	currentUser, err:= user.Current()
	if err != nil {
		return
	}

	controlFileName := "cfA000" + hostname
	dataFileName := "dfA000" + hostname

	// build control file
	controlFile := make(ControlFile)
	controlFile[Hostname] = hostname
	controlFile[UserID] = currentUser.Username
	controlFile[JobName] = doc.Name
	controlFile[BannerClass] = hostname
	controlFile[PrintBanner] = currentUser.Username
	controlFile[PlainTextFile] = dataFileName
	controlFile[UnlinkDataFile] = dataFileName
	controlFile[SourceFileName] = doc.Name

	// append custom cf params
	for cmd, value := range cf {
		controlFile[cmd] = value
	}

	// open connection
	conn, err := net.Dial("tcp", c.dest)
	if err != nil {
		return err
	}
	defer conn.Close()

	// send receive job command
	err = SendCommandLine(conn, byte(ReceiveJob), []string{printer})
	if err != nil {
		return
	}
	if err = CheckAcknowledge(conn); err != nil {
		return
	}

	// ensure the we send abort if we return with error
	defer SendAbortOnError(conn, err)

	// write controlfile to buffer, so we can capture the size
	encodedControlFile, err := controlFile.Encode()
	if err != nil {
		return
	}

	// send controlfile sub command
	err = SendCommandLine(conn, byte(SendControlFile), []string{strconv.Itoa(len(encodedControlFile)), controlFileName})
	if err != nil {
		return
	}
	if err = CheckAcknowledge(conn); err != nil {
		return
	}

	// send controlfile
	_, err = conn.Write(encodedControlFile)
	if err != nil {
		return
	}
	if err = CheckAcknowledge(conn); err != nil {
		return
	}

	// send datafile sub command
	err = SendCommandLine(conn, byte(SendDataFile), []string{strconv.Itoa(doc.Size), dataFileName})
	if err != nil {
		return
	}
	if err = CheckAcknowledge(conn); err != nil {
		return
	}

	// send spool file
	if _, err = io.Copy(conn, doc.Document); err != nil {
		return
	}
	if err = CheckAcknowledge(conn); err != nil {
		return
	}

	return nil
}




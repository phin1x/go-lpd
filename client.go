package go_lpd

import (
	"net"
	"os"
	"os/user"
	"path"
	"strconv"
)

type Client struct {
	// dest format is host:port
	dest string
}

func (c *Client) PrintFile(printer string, cf ControlFile, file string) (err error) {
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

	// check file
	fileStat, err := os.Stat(file)
	if os.IsNotExist(err) {
		return err
	}

	fileName := path.Base(file)
	controlFileName := "cfA000" + hostname
	dataFileName := "dfA000" + hostname

	// build control file
	controlFile := make(ControlFile)
	controlFile[Hostname] = hostname
	controlFile[UserID] = currentUser.Username
	controlFile[JobName] = fileName
	controlFile[BannerClass] = hostname
	controlFile[PrintBanner] = currentUser.Username
	controlFile[PlainTextFile] = dataFileName
	controlFile[UnlinkDataFile] = dataFileName
	controlFile[SourceFileName] = fileName

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
	err = EncodeCommandLine(conn, byte(ReceiveJob), []string{printer})
	if err != nil {
		return
	}
	if err = CheckAcknowledge(conn); err != nil {
		return
	}

	// ensure the we send abort if we return with error
	defer SendAbortOnError(conn, err)

	// send controlfile subcommand
	// FIXME: we have to send the size of the control file, not the file to be printed
	err = EncodeCommandLine(conn, byte(SendControlFile), []string{strconv.FormatInt(fileStat.Size(), 10), controlFileName})
	if err != nil {
		return
	}
	if err = CheckAcknowledge(conn); err != nil {
		return
	}

	// send controlfile
	if err = NewControlFileEncoder(conn).Encode(controlFile); err != nil {
		return
	}
	

	return nil
}




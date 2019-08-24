package lpd

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
	"path"
	"strconv"
)

type Document struct {
	Document io.Reader
	Size     int
	Name     string
}

func NewClient(remote string, port int) *Client {
	return &Client{
		dest: fmt.Sprintf("%s:%d", remote, port),
	}
}

type Client struct {
	// dest format is host:port
	dest string
}

func (c *Client) PrintFile(filePath, queue string, cf ControlFile) error {
	fileStats, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return err
	}

	fileName := path.Base(filePath)

	document, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer document.Close()

	return c.PrintDocument(Document{
			Document: document,
			Name:     fileName,
			Size:     int(fileStats.Size()),
	}, queue, cf, PlainTextFile)
}

func (c *Client) PrintDocument(doc Document, queue string, cf ControlFile, of OutputFormat) (err error) {
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
	controlFile[UnlinkDataFile] = dataFileName
	controlFile[SourceFileName] = doc.Name

	controlFile[ControlFileCommand(of)] = dataFileName

	// append custom cf params
	if cf != nil {
		for cmd, value := range cf {
			controlFile[cmd] = value
		}
	}

	// open connection
	conn, err := net.Dial("tcp", c.dest)
	if err != nil {
		return err
	}
	defer conn.Close()

	// send receive job command
	err = SendCommandLine(conn, byte(ReceiveJob), []string{queue})
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

func (c *Client) PrintWaitingJobs(queue string) (err error) {
	// open connection
	conn, err := net.Dial("tcp", c.dest)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = SendCommandLine(conn, byte(PrintJobs), []string{queue})
	if err != nil {
		return
	}
	err = CheckAcknowledge(conn)

	return
}

func (c *Client) GetQueueStateShort(queue string, jobNumbers, usernames []string) (err error) {
	conn, err := net.Dial("tcp", c.dest)
	if err != nil {
		return err
	}
	defer conn.Close()

	// TODO: send username and job number if given
	err = SendCommandLine(conn, byte(QueueStatsShort), []string{queue})
	if err != nil {
		return
	}

	// TODO: parse response

	return
}

func (c *Client) GetQueueStateLong(queue string, jobNumbers, usernames []string) (err error) {
	conn, err := net.Dial("tcp", c.dest)
	if err != nil {
		return err
	}
	defer conn.Close()

	// TODO: send username and job number if given
	err = SendCommandLine(conn, byte(QueueStatsLong), []string{queue})
	if err != nil {
		return
	}

	// TODO: parse response

	return
}

//magent is the username making the request
func (c *Client) RemoveJobs(queue, agent string, jobNumbers, usernames []string) (err error) {
	conn, err := net.Dial("tcp", c.dest)
	if err != nil {
		return err
	}
	defer conn.Close()

	// TODO: send username and job number if given
	err = SendCommandLine(conn, byte(RemoveJobs), []string{queue, agent})
	if err != nil {
		return
	}
	err = CheckAcknowledge(conn)

	return
}


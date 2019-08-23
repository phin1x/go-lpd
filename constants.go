package go_lpd

// the descriptions copied from https://www.ietf.org/rfc/rfc1179.txt

var (
	Separator = " "
	LineEnding = "\n"
)

var (
	// the server has to acknowledge all commands with a octet of zero bytes
	Acknowledge byte = 0x0
)

type DaemonCommand uint8

const (
	/*
	      +----+-------+----+
	      | 01 | Queue | LF |
	      +----+-------+----+
	      Command code - 1
	      Operand - Printer queue name

	   This command starts the printing process if it not already running.
	*/
	PrintJobs DaemonCommand = 0x01

	/*
	      +----+-------+----+
	      | 02 | Queue | LF |
	      +----+-------+----+
	      Command code - 2
	      Operand - Printer queue name

	   Receiving a job is controlled by a second level of commands.  The
	   daemon is given commands by sending them over the same connection.
	   The commands are described in the next section (6).

	   After this command is sent, the client must read an acknowledgement
	   octet from the daemon.  A positive acknowledgement is an octet of
	   zero bits.  A negative acknowledgement is an octet of any other
	   pattern.
	 */
	ReceiveJob DaemonCommand = 0x02

	/*
	      +----+-------+----+------+----+
	      | 03 | Queue | SP | List | LF |
	      +----+-------+----+------+----+
	      Command code - 3
	      Operand 1 - Printer queue name
	      Other operands - User names or job numbers

	   If the user names or job numbers or both are supplied then only those
	   jobs for those users or with those numbers will be sent.

	   The response is an ASCII stream which describes the printer queue.
	   The stream continues until the connection closes.  Ends of lines are
	   indicated with ASCII LF control characters.  The lines may also
	   contain ASCII HT control characters.
	 */
	QueueStatsShort DaemonCommand = 0x03

	/*
	      +----+-------+----+------+----+
	      | 04 | Queue | SP | List | LF |
	      +----+-------+----+------+----+
	      Command code - 4
	      Operand 1 - Printer queue name
	      Other operands - User names or job numbers

	   If the user names or job numbers or both are supplied then only those
	   jobs for those users or with those numbers will be sent.

	   The response is an ASCII stream which describes the printer queue.
	   The stream continues until the connection closes.  Ends of lines are
	   indicated with ASCII LF control characters.  The lines may also
	   contain ASCII HT control characters.
	 */
	QueueStatsLong DaemonCommand = 0x04

	/*
	      +----+-------+----+-------+----+------+----+
	      | 05 | Queue | SP | Agent | SP | List | LF |
	      +----+-------+----+-------+----+------+----+
	      Command code - 5
	      Operand 1 - Printer queue name
	      Operand 2 - User name making request (the agent)
	      Other operands - User names or job numbers

	   This command deletes the print jobs from the specified queue which
	   are listed as the other operands.  If only the agent is given, the
	   command is to delete the currently active job.  Unless the agent is
	   "root", it is not possible to delete a job which is not owned by the
	   user.  This is also the case for specifying user names instead of
	   numbers.  That is, agent "root" can delete jobs by user name but no
	   other agents can.
	*/
	RemoveJobs DaemonCommand = 0x05
)

type SubCommand byte

const (
	/*
	      Command code - 1
	      +----+----+
	      | 01 | LF |
	      +----+----+

	   No operands should be supplied.  This subcommand will remove any
	   files which have been created during this "Receive job" command.
	 */
	AbortJob SubCommand = 0x01

	/*
	      +----+-------+----+------+----+
	      | 02 | Count | SP | Name | LF |
	      +----+-------+----+------+----+
	      Command code - 2
	      Operand 1 - Number of bytes in control file
	      Operand 2 - Name of control file

	   The control file must be an ASCII stream with the ends of lines
	   indicated by ASCII LF.  The total number of bytes in the stream is
	   sent as the first operand.  The name of the control file is sent as
	   the second.  It should start with ASCII "cfA", followed by a three
	   digit job number, followed by the host name which has constructed the
	   control file.  Acknowledgement processing must occur as usual after
	   the command is sent.

	   The next "Operand 1" octets over the same TCP connection are the
	   intended contents of the control file.  Once all of the contents have
	   been delivered, an octet of zero bits is sent as an indication that
	   the file being sent is complete.  A second level of acknowledgement
	   processing must occur at this point.
	 */
	SendControlFile SubCommand = 0x02

	/*
	      +----+-------+----+------+----+
	      | 03 | Count | SP | Name | LF |
	      +----+-------+----+------+----+
	      Command code - 3
	      Operand 1 - Number of bytes in data file
	      Operand 2 - Name of data file

	   The data file may contain any 8 bit values at all.  The total number
	   of bytes in the stream may be sent as the first operand, otherwise
	   the field should be cleared to 0.  The name of the data file should
	   start with ASCII "dfA".  This should be followed by a three digit job
	   number.  The job number should be followed by the host name which has
	   constructed the data file.  Interpretation of the contents of the
	   data file is determined by the contents of the corresponding control
	   file.  If a data file length has been specified, the next "Operand 1"
	   octets over the same TCP connection are the intended contents of the
	   data file.  In this case, once all of the contents have been
	   delivered, an octet of zero bits is sent as an indication that the
	   file being sent is complete.  A second level of acknowledgement
	   processing must occur at this point.
	 */

	SendDataFile SubCommand = 0x03
)

// the control file commands are defined as string, i defined them as ascii hex values
type ControlFileCommand byte

const (
	/*
	      +---+-------+----+
	      | C | Class | LF |
	      +---+-------+----+
	      Command code - 'C'
	      Operand - Name of class for banner pages

	   This command sets the class name to be printed on the banner page.
	   The name must be 31 or fewer octets.  The name can be omitted.  If it
	   is, the name of the host on which the file is printed will be used.
	   The class is conventionally used to display the host from which the
	   printing job originated.  It will be ignored unless the print banner
	   command ('L') is also used.
	 */
	BannerClass ControlFileCommand = 0x43

	/*
	      +---+------+----+
	      | H | Host | LF |
	      +---+------+----+
	      Command code - 'H'
	      Operand - Name of host

	   This command specifies the name of the host which is to be treated as
	   the source of the print job.  The command must be included in the
	   control file.  The name of the host must be 31 or fewer octets.
	 */
	Hostname ControlFileCommand = 0x48

	/*
	      +---+-------+----+
	      | I | count | LF |
	      +---+-------+----+
	      Command code - 'I'
	      Operand - Indenting count

	   This command specifies that, for files which are printed with the
	   'f', of columns given.  (It is ignored for other output generating
	   commands.)  The identing count operand must be all decimal digits.
	 */
	Indent ControlFileCommand = 0x49

	/*
	      +---+----------+----+
	      | J | Job name | LF |
	      +---+----------+----+
	      Command code - 'J'
	      Operand - Job name

	   This command sets the job name to be printed on the banner page.  The
	   name of the job must be 99 or fewer octets.  It can be omitted.  The
	   job name is conventionally used to display the name of the file or
	   files which were "printed".  It will be ignored unless the print
	   banner command ('L') is also used.
	 */
	JobName ControlFileCommand = 0x4a

	/*
	      +---+------+----+
	      | L | User | LF |
	      +---+------+----+
	      Command code - 'L'
	      Operand - Name of user for burst pages

	   This command causes the banner page to be printed.  The user name can
	   be omitted.  The class name for banner page and job name for banner
	   page commands must precede this command in the control file to be
	   effective.
	 */
	PrintBanner ControlFileCommand = 0x4c

	/*
	      +---+------+----+
	      | M | user | LF |
	      +---+------+----+
	      Command code - 'M'
	      Operand - User name

	   This entry causes mail to be sent to the user given as the operand at
	   the host specified by the 'H' entry when the printing operation ends
	   (successfully or unsuccessfully).
	*/
	MailWhenPrinted ControlFileCommand = 0x4d

	/*
	      +---+------+----+
	      | N | Name | LF |
	      +---+------+----+
	      Command code - 'N'
	      Operand - File name

	   This command specifies the name of the file from which the data file
	   was constructed.  It is returned on a query and used in printing with
	   the 'p' command when no title has been given.  It must be 131 or
	   fewer octets.
	 */
	SourceFileName ControlFileCommand = 0x4e

	/*
	      +---+------+----+
	      | P | Name | LF |
	      +---+------+----+
	      Command code - 'P'
	      Operand - User id

	   This command specifies the user identification of the entity
	   requesting the printing job.  This command must be included in the
	   control file.  The user identification must be 31 or fewer octets.
	 */
	UserID ControlFileCommand = 0x50

	/*
	      +---+--------+----+-------+----+
	      | S | device | SP | inode | LF |
	      +---+--------+----+-------+----+
	      Command code - 'S'
	      Operand 1 - Device number
	      Operand 2 - Inode number

	   This command is used to record symbolic link data on a Unix system so
	   that changing a file's directory entry after a file is printed will
	   not print the new file.  It is ignored if the data file is not
	   symbolically linked.
	 */

	SymbolicLink ControlFileCommand = 0x53

	/*
	      +---+-------+----+
	      | T | title | LF |
	      +---+-------+----+
	      Command code - 'T'
	      Operand - Title text

	   This command provides a title for a file which is to be printed with
	   either the 'p' command.  (It is ignored by all of the other printing
	   commands.)  The title must be 79 or fewer octets.
	 */
	Title ControlFileCommand = 0x54

	/*
	      +---+------+----+
	      | U | file | LF |
	      +---+------+----+
	      Command code - 'U'
	      Operand - File to unlink

	   This command indicates that the specified file is no longer needed.
	   This should only be used for data files.
	 */
	UnlinkDataFile ControlFileCommand = 0x55

	/*
	      +---+-------+----+
	      | W | width | LF |
	      +---+-------+----+
	      Command code - 'W'
	      Operand - Width count

	   This command limits the output to the specified number of columns for
	   the 'f', 'l', and 'p' commands.  (It is ignored for other output
	   generating commands.)  The width count operand must be all decimal
	   digits.  It may be silently reduced to some lower value.  The default
	   value for the width is 132.
	 */
	WidthOfOutput ControlFileCommand = 0x57

	/*
	      +---+------+----+
	      | 1 | file | LF |
	      +---+------+----+
	      Command code - '1'
	      Operand - File name

	   This command specifies the file name for the troff R font.  [1] This
	   is the font which is printed using Times Roman by default.
	 */
	TroffRFont ControlFileCommand = 0x31

	/*
	      +---+------+----+
	      | 2 | file | LF |
	      +---+------+----+
	      Command code - '2'
	      Operand - File name

	   This command specifies the file name for the troff I font.  [1] This
	   is the font which is printed using Times Italic by default.
	 */
	TroffIFont ControlFileCommand = 0x32

	/*
	      +---+------+----+
	      | 3 | file | LF |
	      +---+------+----+
	      Command code - '3'
	      Operand - File name

	   This command specifies the file name for the troff B font.  [1] This
	   is the font which is printed using Times Bold by default.
	 */
	TroffBFont ControlFileCommand = 0x33

	/*
	      +---+------+----+
	      | 4 | file | LF |
	      +---+------+----+
	      Command code - '4'
	      Operand - File name

	   This command specifies the file name for the troff S font.  [1] This
	   is the font which is printed using Special Mathematical Font by
	   default.
	 */
	TroffSFont ControlFileCommand = 0x34

	/*
	      +---+------+----+
	      | c | file | LF |
	      +---+------+----+
	      Command code - 'c'
	      Operand - File to plot

	   This command causes the data file to be plotted, treating the data as
	   CIF (CalTech Intermediate Form) graphics language. [2]
	 */
	CIFFile ControlFileCommand = 0x63

	/*
	      +---+------+----+
	      | d | file | LF |
	      +---+------+----+
	      Command code - 'd'
	      Operand - File to print

	   This command causes the data file to be printed, treating the data as
	   DVI (TeX output). [3]
	 */
	DVIFile ControlFileCommand = 0x64

	/*
	      +---+------+----+
	      | f | file | LF |
	      +---+------+----+
	      Command code - 'f'
	      Operand - File to print

	   This command cause the data file to be printed as a plain text file,
	   providing page breaks as necessary.  Any ASCII control characters
	   which are not in the following list are discarded: HT, CR, FF, LF,
	   and BS.
	 */
	PlainTextFile ControlFileCommand = 0x66

	/*
	      +---+------+----+
	      | g | file | LF |
	      +---+------+----+
	      Command code - 'g'
	      Operand - File to plot

	   This command causes the data file to be plotted, treating the data as
	   output from the Berkeley Unix plot library. [1]
	 */
	PlotFile ControlFileCommand = 0x67

	/*
	      +---+------+----+
	      | l | file | LF |
	      +---+------+----+
	      Command code - 'l' (lower case L)
	      Operand - File to print

	   This command causes the specified data file to printed without
	   filtering the control characters (as is done with the 'f' command).
	 */
	PrintWithLeavingControlCharacters ControlFileCommand = 0x6c

	/*
	      +---+------+----+
	      | n | file | LF |
	      +---+------+----+
	      Command code - 'n'
	      Operand - File to print

	   This command prints the data file to be printed, treating the data as
	   ditroff output. [4]
	 */
	DitroffFile ControlFileCommand = 0x6e

	/*
	      +---+------+----+
	      | o | file | LF |
	      +---+------+----+
	      Command code - 'o'
	      Operand - File to print

	   This command prints the data file to be printed, treating the data as
	   standard Postscript input.
	 */
	PostscriptFile ControlFileCommand = 0x6f

	/*
	      +---+------+----+
	      | p | file | LF |
	      +---+------+----+
	      Command code - 'p'
	      Operand - File to print

	   This command causes the data file to be printed with a heading, page
	   numbers, and pagination.  The heading should include the date and
	   time that printing was started, the title, and a page number
	   identifier followed by the page number.  The title is the name of
	   file as specified by the 'N' command, unless the 'T' command (title)
	   has been given.  After a page of text has been printed, a new page is
	   started with a new page number.  (There is no way to specify the
	   length of the page.)
	 */
	PRFormat ControlFileCommand = 0x70

	/*
	      +---+------+----+
	      | r | file | LF |
	      +---+------+----+
	      Command code - 'r'
	      Operand - File to print

	   This command causes the data file to be printed, interpreting the
	   first column of each line as FORTRAN carriage control.  The FORTRAN
	   standard limits this to blank, "1", "0", and "+" carriage controls.
	   Most FORTRAN programmers also expect "-" (triple space) to work as
	   well.
	 */
	FortranCarriageControlFormat ControlFileCommand = 0x72

	/*
	      +---+------+----+
	      | t | file | LF |
	      +---+------+----+
	      Command code - 't'
	      Operand - File to print

	   This command prints the data file as Graphic Systems C/A/T
	   phototypesetter input.  [5] This is the standard output of the Unix
	   "troff" command.
	 */
	TroffFormat ControlFileCommand = 0x74

	/*
	      +---+------+----+
	      | v | file | LF |
	      +---+------+----+
	      Command code - 'v'
	      Operand - File to print

	   This command prints a Sun raster format file. [6]
	 */
	RasterFormat ControlFileCommand = 0x76
)


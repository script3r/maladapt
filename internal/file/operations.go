package file

type Fs interface {
	Copy(src string, remote string) error
	Move(src string, remote string) error
	Stat(src string, remote string) error
}

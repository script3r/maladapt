package file

type Fs interface {
	Copy(src string, dst string) error
	Move(src string, dst string) error
	Stat(src string, dst string) error
}

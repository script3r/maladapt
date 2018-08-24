package file

type Quarantiner interface {
	RenderInert()
	RenderAlive()
}

type QuarantineAdmin struct {
	Location string
}

func NewQuarantineAdmin(location string) *QuarantineAdmin {
	return &QuarantineAdmin{location}
}

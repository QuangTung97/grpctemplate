package backend

// IPort interface for incoming port
type IPort interface {
	Hello(HelloInput) error
}

// Port incoming port
type Port struct {
}

// NewPort create a Port
func NewPort() *Port {
	return &Port{}
}

var _ IPort = &Port{}

// Hello do hello
func (p *Port) Hello(input HelloInput) error {
	return nil
}

package message

type Message interface {
	GetFrom() string
	GetTarget() string
	GetContent() string
}

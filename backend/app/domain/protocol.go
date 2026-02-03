package domain

type ProtocolType string

const (
	ProtocolTypeCustom = ProtocolType("custom")
)

type Protocol struct {
	Type        ProtocolType
	Title       string
	Description string
}

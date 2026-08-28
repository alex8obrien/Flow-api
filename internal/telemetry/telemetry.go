package telemetry

// Attribute is a lightweight key/value pair.
type Attribute struct {
	Key   string
	Value string
}

type Status string

const (
	StatusUnset   Status = "unset"
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

func String(key, value string) Attribute {
	return Attribute{Key: key, Value: value}
}

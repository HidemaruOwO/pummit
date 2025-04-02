package variable

import _ "embed"

//go:embed config.json
var DEFAULT_CONFIG []byte

const (
	VERSION = "2.0.0"
)

package variable

import _ "embed"

//go:embed config.json
var DEFAULT_CONFIG string

//go:embed discord-emojis.flat.json
var EMOJIS_JSON string

const (
	VERSION = "2.0.0"
)

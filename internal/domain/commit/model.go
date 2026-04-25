package commit

type Request struct {
	EmojiInput string
	Message    string
	AutoEmoji  bool
}

type Result struct {
	Committed bool
	Message   string
}

func (r Result) GetCommitted() bool {
	return r.Committed
}

func (r Result) GetMessage() string {
	return r.Message
}

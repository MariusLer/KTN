package messages

type ClientPayload struct {
	Request string
	Content string
}

type ServerPayload struct {
	Timestamp string
	Sender    string
	Response  string
	Content   string
}

type HistoryPayload struct {
	Timestamp string
	Sender    string
	Response  string
	Content   [][]byte
}

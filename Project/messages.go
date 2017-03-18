package messages

type ClientPayload struct {
	Request string
	Content string
}

type ServerPayload struct {
	Timespamp string
	Sender    string
	Response  string
	Content   []string
}

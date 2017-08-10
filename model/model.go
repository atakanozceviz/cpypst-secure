package model

type Data struct {
	Action  string `json:"action"`
	From    string `json:"from"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

type Settings struct {
	IncomingClip bool
	OutgoingClip bool
	IncomingFile bool
	OutgoingFile bool
}

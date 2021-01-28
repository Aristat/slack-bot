package db

// Message example Message table
type Message struct {
	Channel   string
	Timestamp string
	Users     []User
}

// Messages list
var Messages = map[string]*Message{}

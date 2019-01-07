package internal

import (
	"time"
)

type DataUnit struct {
	Date		time.Time
	Subject		string
	From 		string
	Sender 		string
	ReplyTo 	string
	To			string
	CC			string
	BCC			string
	Body		string
	InReplyTo	string
	MessageID 	string
}

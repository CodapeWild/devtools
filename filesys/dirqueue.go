package filesys

import (
	"devtools/msgque"
)

type DirTicket struct {
	Path  string
	Count int
}

type DirTicketQueue struct {
	*msgque.SimpleTicketQueue
}

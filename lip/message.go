package lip

import (
	"errors"
	"strings"
)

type Operation string

const (
	Execute Operation = "#"
	Monitor Operation = "~"
	Query   Operation = "?"
)

type CmdType string

const (
	Device CmdType = "DEVICE"
	Output CmdType = "OUTPUT"
)

type LIPMessage struct {
	Operation Operation
	CmdType   CmdType
	ID        string
	Params    []string
}

func ParseLIPMessage(msg string) (*LIPMessage, error) {
	if strings.HasPrefix(msg, "GNET> ") {
		msg = strings.Trim(msg, "GNET> ")
	}

	if msg == "" {
		return nil, errors.New("empty string")
	}

	msg = strings.Trim(msg, "\r\n")
	l := new(LIPMessage)
	l.Operation = Operation(msg[0])

	c := strings.Split(msg[1:], ",")
	l.CmdType = CmdType(c[0])
	l.ID = c[1]

	if len(c) > 1 {
		l.Params = c[2:]
	}
	return l, nil
}

func (l *LIPMessage) String() string {
	s := string(l.Operation) + string(l.CmdType) + "," + l.ID

	if l.Params != nil {
		s += "," + strings.Join(l.Params, ",")
	}

	return s
}

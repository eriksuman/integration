package integration

import (
	"github.com/eriksuman/integration/lip"
)

func Start(host, user, pass string) {

	conn, err := lip.New(host, user, pass)
	if err != nil {
		panic(err)
	}

	go handleExhaustFan(conn, conn.RegisterObserver(lip.Registration{
		Op:      lip.Monitor,
		CmdType: lip.Output,
		ID:      "9",
	}))

	<-make(chan struct{})
}

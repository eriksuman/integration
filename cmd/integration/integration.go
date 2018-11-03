package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/eriksuman/integration"
)

type ProgrammingHandle func(chan *integration.LIPMessage)

func main() {
	ip := flag.String("ip", "", "IP address of integration server")
	user := flag.String("username", "", "Telnet username for integration server")
	pass := flag.String("password", "", "Telnet password for integration server")
	flag.Parse()

	c, err := integration.New(*ip, *user, *pass)
	if err != nil {
		fmt.Println(err)
		return
	}

	go handleExhaustFan(c, c.RegisterObserver(integration.Registration{
		Op:      integration.Monitor,
		CmdType: integration.Output,
		ID:      "9",
	}))

	<-make(chan struct{})
}

func handleExhaustFan(conn *integration.LIPConn, c chan *integration.LIPMessage) {
	for {
		m := <-c
		if m.Params == nil {
			fmt.Printf("no params: %s\n", m)
			continue
		}

		if len(m.Params) != 2 {
			fmt.Printf("bad: %s\n", m)
			continue
		}

		l := m.Params[1]
		if l != "0.00" {
			select {
			case <-time.After(10 * time.Minute):
				conn.IssueCommand(&integration.LIPMessage{
					Operation: integration.Execute,
					CmdType:   integration.Output,
					ID:        "9",
					Params:    []string{"1", "0"},
				})
			case <-c:
			}
		}
	}
}

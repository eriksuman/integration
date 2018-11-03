package integration

import (
	"fmt"
	"time"

	"github.com/eriksuman/integration/lip"
)

func handleExhaustFan(conn *lip.LIPConn, c chan *lip.LIPMessage) {
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
				conn.IssueCommand(&lip.LIPMessage{
					Operation: lip.Execute,
					CmdType:   lip.Output,
					ID:        "9",
					Params:    []string{"1", "0"},
				})
			case <-c:
			}
		}
	}
}

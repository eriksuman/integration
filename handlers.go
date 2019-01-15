package integration

import (
	"fmt"
	"time"

	"github.com/eriksuman/integration/lip"
)

func handleHeater(conn *lip.LIPConn, c chan *lip.LIPMessage) {
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
			conn.IssueCommand(&lip.LIPMessage{
				Operation: lip.Execute,
				CmdType:   lip.Output,
				ID:        "9",
				Params:    []string{"1", "0"},
			})
			select {
			case <-time.After(10 * time.Minute):
				conn.IssueCommand(&lip.LIPMessage{
					Operation: lip.Execute,
					CmdType:   lip.Output,
					ID:        "10",
					Params:    []string{"1", "0"},
				})
			case <-c:
			}
		}
	}
}

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
			conn.IssueCommand(&lip.LIPMessage{
				Operation: lip.Execute,
				CmdType:   lip.Output,
				ID:        "10",
				Params:    []string{"1", "0"},
			})
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

func handleBedroomPico(conn *lip.LIPConn, c chan *lip.LIPMessage) {
	for {
		m := <-c
		if m.Params == nil {
			fmt.Printf("no params: %s", m)
			continue
		}
		if len(m.Params) != 2 {
			fmt.Printf("wrong number of params: %s", m)
			continue
		}
		if m.Params[0] != "4" {
			fmt.Println("bedroom pico non off pressed")
			continue
		}

		switch m.Params[1] {
		case "3":
			conn.IssueCommand(&lip.LIPMessage{
				Operation: lip.Execute,
				CmdType:   lip.Device,
				ID:        "1",
				Params:    []string{"4", "3"},
			})
		case "4":
			conn.IssueCommand(&lip.LIPMessage{
				Operation: lip.Execute,
				CmdType:   lip.Output,
				ID:        "12",
				Params:    []string{"1", "0"},
			})
		}
	}

}

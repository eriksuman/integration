package lip

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type LIPConn struct {
	conn      net.Conn
	lock      sync.Locker
	cmdChan   chan *LIPMessage
	observers map[Registration]chan *LIPMessage
}

func New(host, usrname, passwd string) (*LIPConn, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, "23"))
	if err != nil {
		return nil, err
	}

	if err := login(conn, usrname, passwd); err != nil {
		return nil, err
	}

	lc := &LIPConn{
		conn:      conn,
		lock:      &sync.Mutex{},
		cmdChan:   make(chan *LIPMessage),
		observers: make(map[Registration]chan *LIPMessage),
	}

	go lc.processCmds()
	go lc.listen()

	return lc, nil
}

type Registration struct {
	Op      Operation
	CmdType CmdType
	ID      string
}

func (l *LIPConn) RegisterObserver(r Registration) chan *LIPMessage {
	c := make(chan *LIPMessage)
	l.observers[r] = c
	return c
}

func (l *LIPConn) IssueCommand(m *LIPMessage) {
	fmt.Printf("sending %s\n", m)
	l.cmdChan <- m
}

func (l *LIPConn) listen() {
	fmt.Println("listening")
	for {
		msg, err := l.read()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(msg)

		r := Registration{
			Op:      msg.Operation,
			CmdType: msg.CmdType,
			ID:      msg.ID,
		}

		c, ok := l.observers[r]
		if !ok {
			continue
		}

		select {
		case <-c:
			c <- msg
		case c <- msg:
		}
	}
}

func (l *LIPConn) processCmds() {
	for {
		msg := <-l.cmdChan
		fmt.Printf("command: %s\n", msg)

		_, err := l.conn.Write([]byte(msg.String() + "\r\n"))
		if err != nil {
			panic(err)
		}
	}
}

func (l *LIPConn) read() (*LIPMessage, error) {
	reader := bufio.NewReader(l.conn)

	msg, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return ParseLIPMessage(msg)
}

func login(conn net.Conn, u, p string) error {
	_, err := conn.Write([]byte(u + "\r\n"))
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(p + "\r\n"))
	return err
}

func read(conn net.Conn) (string, error) {
	buf := make([]byte, 256)

	_, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

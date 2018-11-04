package main

import (
	"flag"

	"github.com/eriksuman/integration"
)

func main() {
	ip := flag.String("ip", "", "IP address of integration server")
	user := flag.String("username", "", "Telnet username for integration server")
	pass := flag.String("password", "", "Telnet password for integration server")
	flag.Parse()

	integration.Start(*ip, *user, *pass)
}

package socres

import (
	"flag"
	"gopkg.in/routeros.v2"
	"gopkg.in/routeros.v2/proto"
	"log"
	"strings"
)

type Client struct {
	Address  string
	Username string
	Password string
	Async    bool
	UseTLS   bool
}

func (client *Client) SendMikrotik(command string) []*proto.Sentence {
	flag.Parse()

	c, err := client.dial()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if client.Async {
		c.Async()
	}

	log.Println(command)
	r, err := c.RunArgs(strings.Split(command, " "))
	if err != nil {
		log.Fatal(err)
	}
	return r.Re
}

func (client *Client)dial() (*routeros.Client, error) {
	if client.UseTLS {
		return routeros.DialTLS(client.Address, client.Username, client.Password, nil)
	}
	return routeros.Dial(client.Address, client.Username, client.Password)
}

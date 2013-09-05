package nntp

import (
	"crypto/tls"
	"fmt"
	"net/textproto"
	"log"
	"strings"
)

type NntpClient struct {
	conn *textproto.Conn
}

func Dial(addr string) (*NntpClient, error) {
	conn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}
	client := &NntpClient{conn: textproto.NewConn(conn)}
	_, err = client.read(200) // get welcome message
	if err != nil {
		return nil, err
	}
	return client, nil
}

func DialAuth(addr string, user string, pass string) (*NntpClient, error) {
	client, err := Dial(addr)
	if err != nil {
		return nil, err
	}

	err = client.Auth(user, pass)

	return client, err
}

func (client *NntpClient) read(expectCode int) (resp string, err error) {
	_, resp, err = client.conn.ReadCodeLine(expectCode)
	log.Println("read:", resp)
	return
}

func (client *NntpClient) write(format string, args ...interface{}) error {
	log.Printf("wrote: " + format, args...)
	return client.conn.PrintfLine(format, args...)
}

func (client *NntpClient) ExecuteCommand(expectCode int, format string,
										 args ...interface{}) (resp string, err error) {
	err = client.write(format, args...)
	if err != nil {
		return
	}
	resp, err = client.read(expectCode)
	if err != nil {
		return
	}

	return
}

func (client *NntpClient) ReadResults() (lines []string, err error) {
	return client.conn.ReadDotLines()
}

func (client *NntpClient) Auth(user string, pass string) (err error) {
	_, err = client.ExecuteCommand(381, "AUTHINFO USER %s", user)
	if err != nil {
		return
	}
	_, err = client.ExecuteCommand(281, "AUTHINFO PASS %s", pass)
	return
}


func (client *NntpClient) XOver(group string, overRange string) (err error) {
	_, err = client.ExecuteCommand(211, "GROUP %s", group)
	if err != nil {
		return
	}

	_, err = client.ExecuteCommand(224, "XOVER %s", overRange)
	if err != nil {
		return
	}

	lines, err := client.ReadResults()
	if err != nil {
		return
	}

	for _, line := range lines {
		words := strings.Split(line, "\t")
		fmt.Println(words)
	}
	return
}

func (client *NntpClient) List(filter string) (items []ListItem, err error) {
	_, err = client.ExecuteCommand(215, "LIST ACTIVE %s", filter)
	if err != nil {
		return
	}

	lines, err := client.ReadResults()
	if err != nil {
		return
	}
	items = make([]ListItem, 0, len(lines))
	for _, line := range lines {
		item := ListItem{}
		_, err := fmt.Sscanf(line, "%s %d %d %s", &item.Name, &item.High, &item.Low, &item.Status)
		if err != nil {
			log.Println(err)
			continue
		}
		items = append(items, item)
	}
	return
}

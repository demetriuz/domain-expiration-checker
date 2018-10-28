package whois_backends

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type RawWhoisBackend struct{}


func (b RawWhoisBackend) Fetch(domain string) (string, error){

	conn, err := net.Dial("tcp", "whois.tcinet.ru:43") // TODO: address to cfg

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn.Write([]byte(domain))

	var outBuilder strings.Builder

	connbuf := bufio.NewReader(conn)
	for{
		str, err := connbuf.ReadString('\n')
		outBuilder.WriteString(str)
		if err!= nil {
			break
		}
	}

	if err != nil{
		return "", err
	}else{
		return string(outBuilder.String()), nil
	}
}

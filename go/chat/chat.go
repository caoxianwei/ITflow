package chat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7000")
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan bool)
	go userRead(conn)
	go userWrite(conn)
	<-ch
}

func userRead(conn net.Conn) {
	for {
		output := bufio.NewScanner(conn)
		for output.Scan() {
			fmt.Println(output.Text())
		}
	}

}

func userWrite(conn net.Conn) {
	for {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Print("enter an string:")
		input, err := inputReader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		conn.Write([]byte(input))
	}

}
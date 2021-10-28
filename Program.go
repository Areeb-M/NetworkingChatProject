package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Server Variables
	var port int64 = 43980
	var service string = ":" + strconv.FormatInt(port, 10)
	var network string = "tcp"
	var tcpListener net.Listener
	var err error

	var connectionsRecord map[int]net.Conn

	tcpListener, err = net.Listen(network, service)
	checkServerError(err)

	defer tcpListener.Close()

	connectionsRecord = make(map[int]net.Conn)

	for {
		// Accept new connections
		conn, err := tcpListener.Accept()
		checkServerError(err)

		// Store the conn in the record with a randomly assigned ID
		var randInt int
		for randInt == 0 {
			randInt = generateRandom5DigitNum()

			if _, ok := connectionsRecord[randInt]; ok {
				randInt = 0
			}
		}

		connectionsRecord[randInt] = conn

		// Pass them off to a goroutine for processing
		go handleServerClient(conn, randInt, connectionsRecord)
	}

}

func checkServerError(err error) {
	if err != nil {
		fmt.Println("Houston, we've got a problem.")
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateRandom5DigitNum() int {
	return 10000 + rand.Intn(99999-10000)
}

func handleServerClient(conn net.Conn, connId int, connectionsRecord map[int]net.Conn) {
	defer delete(connectionsRecord, connId)
	defer conn.Close()

	fmt.Println(conn.RemoteAddr())

	var buf [2048]byte
	var recipient int64 = 0
	var reply string
	var message string

	for {
		n, err := conn.Read(buf[:])

		if err != nil {
			fmt.Print("[Error]")
			fmt.Print(err)
			fmt.Printf(" -- Closing connection %d.\n", connId)

			return
		}

		// Check to see if there's a mention at the beginning of the message
		if buf[0] == '@' {
			// If there is, only return the message to the recipient and the sender
			recipient, err = strconv.ParseInt(strings.ReplaceAll(string(buf[1:7]), " ", ""), 10, 0)
			_, ok := connectionsRecord[int(recipient)]

			if err != nil || !ok {
				// There was an error, tell the sender that their recipient doesn't exist
				message = string(buf[:n])
				fmt.Printf("[Error] Message (%s) from %d had an invalid recipient.", strings.ReplaceAll(message, "\n", ""), connId)
				conn.Write([]byte(fmt.Sprintf("[Error] Message (%s) had an invalid recipient. It was not forwarded to anyone.\n", message)))
				continue
			}

			message = string(buf[7:n])
		} else {
			recipient = 0
			message = string(buf[:n])
		}

		// Format and retransmit the message
		if recipient == 0 {
			reply = fmt.Sprintf("[%d]: %s", connId, message)
			fmt.Printf("[%s]"+reply, conn.RemoteAddr().String())
			byteReply := []byte(reply)

			for id := range connectionsRecord {
				connectionsRecord[id].Write(byteReply)
			}
		} else {
			reply = fmt.Sprintf("[%d] (to %d): %s", connId, recipient, message)
			fmt.Printf("[%s]"+reply, conn.RemoteAddr().String())

			conn.Write([]byte(reply))
			if connId != int(recipient) {
				connectionsRecord[int(recipient)].Write([]byte(reply))
			}
		}
	}
}

package main

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
	"github.com/juju/loggo"
	"github.com/vmihailenco/msgpack"
	"log"
	"net"
	"os"
	"time"
)

var logger = loggo.GetLogger("collector")

type Message struct {
	length int
	buffer []byte
}

func NewMessage(length int, buffer []byte) *Message {
	return &Message{length, buffer}
}

func (m *Message) Body() []byte {
	return m.buffer[0:m.length]
}

func (m *Message) Decode() (interface{}, error) {
	//var out []byte
	dec := msgpack.NewDecoder(bytes.NewBuffer(m.buffer))
	return dec.DecodeInterface()
	//return out, err
}

func main() {

	serverAddr, _ := net.ResolveUDPAddr("udp", "mt2:3000")
	logger.Debugf("saddr : %v", serverAddr)

	conn, err := net.DialUDP("udp6", nil, serverAddr)
	logger.Debugf("conn : %v", conn)

	if err != nil {
		logger.Errorf("err : %s", err)
		os.Exit(1)
	}

	c := make(chan *Message, 1)

	go ReadData(conn, c)

	sendBuf := []byte{0}

	for {

		logger.Debugf("sending : %v", sendBuf)
		_, err = conn.Write(sendBuf)

		if err != nil {
			logger.Errorf("err : %s", err)
			os.Exit(1)
		}

		select {
		case msg := <-c:
			// use err and reply
			//log.Printf("recv : %v", msg.Body())
			data, err := msg.Decode()
			logger.Debugf("err : %v, data : %s", err, spew.Sdump(data))

		case <-time.After(10e8):
			// call timed out
			log.Printf("timed out")

		}
	}
}

func ReadData(conn *net.UDPConn, rec chan *Message) {
	var buffer [1500]byte

	for {
		// Read from server
		n, saddr, err := conn.ReadFrom(buffer[0:])
		if err != nil {
			logger.Debugf("err : ", err)
		}

		logger.Debugf("Recieved from server to %s.\n", saddr.String())

		rec <- NewMessage(n, buffer[0:])
	}
}

package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
)

const terminationSequence = "\x00"

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:27015")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	packet := NewPacket(SERVERDATA_AUTH, "hello")

	err = writePacket(&conn, packet)
	if err != nil {
		panic(err)
	}
	packet, err = readPacket(&conn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", packet)
}

func writePacket(conn *net.Conn, packet *Packet) error {
	buffer := packet.ToBuffer()

	bytes := buffer.Bytes()

	n, err := (*conn).Write(bytes)
	if n != len(bytes) {
		return errors.New("Invalid read")
	}
	return err
}

func readPacket(conn *net.Conn) (response *Packet, err error) {
	var packet Packet
	if err = binary.Read(*conn, binary.LittleEndian, &packet.Size); err != nil {
		return
	}
	if err = binary.Read(*conn, binary.LittleEndian, &packet.Id); err != nil {
		return
	}
	if err = binary.Read(*conn, binary.LittleEndian, &packet.Type); err != nil {
		return
	}

	var n int
	bytesRead := 0
	bytesTotal := int(packet.Size - 4)
	body := make([]byte, bytesTotal)

	for bytesTotal > bytesRead {
		n, err = (*conn).Read(body[bytesRead:])
		if err != nil {
			return
		}
		bytesRead += n
	}

	packet.Body = strings.TrimRight(string(body), terminationSequence)

	return &packet, nil
}

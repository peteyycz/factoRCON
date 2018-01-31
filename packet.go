package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
)

type Packet struct {
	Size int32
	Id   int32
	Type int32
	Body string
}

const (
	SERVERDATA_AUTH           int32 = 3
	SERVERDATA_AUTH_RESPONSE  int32 = 2
	SERVERDATA_EXECCOMMAND    int32 = 2
	SERVERDATA_RESPONSE_VALUE int32 = 0
)

func NewPacket(typename int32, body string) *Packet {
	var size, id int32

	// Since the only one of these values that can change in length is the body, an easy way to calculate the size of a packet is to find the byte-length of the packet body, then add 10 to it.
	// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol
	size = int32(len(body) + 10)

	binary.Read(rand.Reader, binary.LittleEndian, &id)

	return &Packet{
		Type: typename,
		Body: body,
		Size: size,
		Id:   id,
	}
}

func (packet *Packet) ToBuffer() *bytes.Buffer {
	size := int32(packet.Size + 4 /* Header 4 byte integer header */)
	buffer := bytes.NewBuffer(make([]byte, 0, size))

	binary.Write(buffer, binary.LittleEndian, packet.Size)
	binary.Write(buffer, binary.LittleEndian, packet.Id)
	binary.Write(buffer, binary.LittleEndian, packet.Type)

	buffer.WriteString(packet.Body)
	binary.Write(buffer, binary.LittleEndian, byte(0))

	binary.Write(buffer, binary.LittleEndian, byte(0))

	return buffer
}

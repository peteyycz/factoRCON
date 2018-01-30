package main

import (
        "encoding/binary"
        "bytes"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:2020")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

        message := "hellok\n"
        var message_bytes [20]byte
        message_bytes_slice := message_bytes[:]
        copy(message_bytes_slice, message)

        buffer := bytes.NewBuffer(make([]byte, 0, 14))
        buffer.WriteString(message)
        binary.Write(buffer, binary.LittleEndian, message_bytes_slice)

        conn.Write(buffer.Bytes())
}

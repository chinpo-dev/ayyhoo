package main

import (
	"fmt"
	"net"
	"os"
//	"strconv"
	"bytes"
	"bufio"
	"io"
)

const (
	CONN_HOST = ""
	CONN_PORT = "5050"
	CONN_TYPE = "tcp"
	PROTOCOL_VER = 10
)

func main() {
	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Connected:", conn.RemoteAddr())

		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func decodeData(data []byte) map[string]string {
	dat := make(map[string]string)
	splat := bytes.Split(data, []byte{192,128})
	for i := 0 ; i<len(splat) -1 ; i+=2 {
		dat[string(splat[i])] = string(splat[i + 1])
	}
	return dat
}

func encodeData(dataarray []string) []byte {
	payload := []byte{}
	sep := []byte{192,128}
	for i := 0 ; i<len(dataarray) ; i++ {
		entry := []byte(dataarray[i])
		entry = append(entry, sep...)
		payload = append(payload, entry...)
	}
	return payload
}

func ymsg(service byte, status byte, session_id byte, data []byte) []byte {
	msg := make([]byte, 20)
	msg[0], msg[1], msg[2], msg[3] = 89, 77, 83, 71 // YMSG protocol
	msg[4] = PROTOCOL_VER // version 10
	msg[11] = service
	msg[15] = status
	msg[19] = session_id
	data_s := len(data)
	if data_s > 0 {
		msg[8] = byte(data_s / 256)
		msg[9] = byte(data_s % 256)
		msg = append(msg, data...)
	}
	return msg
}

func handleRequest(conn net.Conn) {
	c := bufio.NewReader(conn)
	for {
		buf := make([]byte, 20)
		_, err := io.ReadFull(c, buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			fmt.Println("Closing connection")
			conn.Close()
			break
		}

		data_s := 256*int(buf[8]) + int(buf[9])
		data := make([]byte, data_s)
		if data_s > 0 {
			_, err := io.ReadFull(c,data)
			if err != nil {
				fmt.Println("Error reading data:", err.Error())
			}
		}
		datamap := make(map[string]string)
		if data_s > 0 {
			datamap = decodeData(data)
		}

		message := make([]byte, 20)
		if buf[11] == 76 {
			//handshake
			message = ymsg(76, 1, 1, []byte{})
		} else if buf[11] == 87 {
			//auth
			// FIXME: actually implement authorization
			username := datamap["1"]
			message = ymsg(87, 1, 1, encodeData([]string{"1", username,
				 "94", "AA.BBCCDDEEGGHHIIH_HII--"}))
		} else if buf[11] == 84 {
			//authresp
			// FIXME: actually implement authorization
			username := datamap["1"]
			respdata := make([]string, 0)
			respdata = append(respdata, []string{"87","\x0A"}...)
			respdata = append(respdata, []string{"88",""}...)
			respdata = append(respdata, []string{"89",""}...)
			//respdata = append(respdata, []string{"59","Y\x09"}...)
			//respdata = append(respdata, []string{"59","T\x09"}...)
			//respdata = append(respdata, []string{"59","C\x09mg=1"}...)
			respdata = append(respdata, []string{"3" , username}...)
			respdata = append(respdata, []string{"90","1"}...)
			respdata = append(respdata, []string{"100","0"}...)
			respdata = append(respdata, []string{"101",""}...)
			respdata = append(respdata, []string{"102",""}...)
			respdata = append(respdata, []string{"93","86400"}...)
			message = ymsg(1, 0, 1, encodeData([]string{"7",""}))
			fmt.Printf("Sent back: ")
			fmt.Println(message)
			conn.Write([]byte(message));
			message = ymsg(85, 0, 1, encodeData(respdata))
			// 85: LIST
		} else if buf[11] == 18 {
			message = ymsg(18, 1, 1, encodeData([]string{}))
		}

		//message := buf
		n := bytes.Index(buf, []byte{0})
		_ = n
		fmt.Printf("Received: ")
		fmt.Println(buf)
		if data_s > 0 {
			fmt.Println("Enclosed data: ", datamap)
		}
		fmt.Printf("Sent back: ")
		fmt.Println(message)
		fmt.Println()

		// Write the message in the connection channel.
		conn.Write([]byte(message));
	}
}

package main

import (
	"fmt"
	"net"
	"os"
//	"strconv"
	"bytes"
	"bufio"
	"io"
//    _ "github.com/mattn/go-sqlite3"
)

const (
	CONN_HOST = ""
	CONN_PORT = "5050"
	CONN_TYPE = "tcp"
	PROTOCOL_VER = 10
)

type connection struct {
    sess int
    conn net.Conn
}

var queue []byte
var usernames map[string]int //username -> session no
var get_usernames map[int]string //session no -> username
var connections map[int]net.Conn //session no -> connection
var sessioncount = 0

func main() {
    connections = make(map[int]net.Conn)
    usernames = make(map[string]int)
    get_usernames = make(map[int]string)
	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
        for k := range connections {
            fmt.Println(k)
        }
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

func ymsg(service byte, status byte, session_id int, data []byte) []byte {
	msg := make([]byte, 20)
	msg[0], msg[1], msg[2], msg[3] = 89, 77, 83, 71 // YMSG protocol
	msg[4] = PROTOCOL_VER // version 10
	msg[11] = service
	msg[15] = status
	msg[19] = byte(session_id)
	data_s := len(data)
	if data_s > 0 {
		msg[8] = byte(data_s / 256)
		msg[9] = byte(data_s % 256)
		msg = append(msg, data...)
	}
	return msg
}

func send(conn net.Conn, service byte, status byte, session_id int, dataarray []string) {
    msg := ymsg(service, status, session_id, encodeData(dataarray))
    conn.Write(msg)
	fmt.Printf("Sent: ")
	fmt.Println(msg)
    if len(dataarray) > 0 {
	    fmt.Printf("Sent data: ")
	    fmt.Println(dataarray)
    }
	fmt.Println()
}

func handleRequest(conn net.Conn) {
    sessioncount += 1
    sess := sessioncount
    connections[sess] = conn
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

		if buf[11] == 76 {
			//handshake
			go send(conn, 76, 1, sess, []string{})
		} else if buf[11] == 87 {
			//auth
			// FIXME: actually implement authorization
			username := datamap["1"]
            usernames[username] = sess
            get_usernames[sess] = username
			go send(conn, 87, 1, sess, []string{"1", username,
				 "94", "AA.BBCCDDEEGGHHIIH_HII--"})
		} else if buf[11] == 84 {
			//authresp
			// FIXME: actually implement authorization
			// FIXME: send real list of friends
			// FIXME: send online status of friends
			go send(conn, 1, 0, 1, []string{"7",""})
			username := datamap["1"]
			respdata := []string{ "87","(No Group):ab\n",
			                      "88","",
			                      "89","",
			                      "59","C\x09mg=1",
			                      //"59","Y\x09",
			                      //"59","T\x09",
			                      "3" , username,
			                      "90","1",
			                      "100","0",
			                      "101","",
			                      "102","",
			                      "93","86400"}
			go send(conn, 85, 1, sess, respdata) // 85: LIST
		} else if buf[11] == 18 {
			go send(conn, 18, 1, sess, []string{})
		} else if buf[11] == 6 {
            target := datamap["5"]
            // look up in the map for online usernames and get session number
            if to_sess, ok := usernames[target]; ok {
                msgcontent := datamap["14"]
                if msgcontent[0:14] == "<font face=\"\">" {
                    msgcontent = msgcontent[14:]
                }
			    msgdata := []string{"5",target, // to
			                        "4",datamap["1"], // from
			                        "14",msgcontent, // message content
			                        "63",";0",
			                        "64","0",
			                        "97","1"}
                to_conn := connections[to_sess]
                go send(to_conn, 6, 1, to_sess, msgdata)
            }
        }

		n := bytes.Index(buf, []byte{0})
		_ = n
		fmt.Printf("Received: ")
		fmt.Println(buf)
		if data_s > 0 {
			fmt.Println("Enclosed data: ", datamap)
		}
		fmt.Println()
	}
}

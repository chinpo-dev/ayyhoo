package main

import (
    "fmt"
    "net"
    "os"
    "strconv"
    "bytes"
)

const (
    CONN_HOST = ""
    CONN_PORT = "5050"
    CONN_TYPE = "tcp"
)

func main() {
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }

        //logs an incoming message
        fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

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
    fmt.Println(dat)
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
    msg[4] = 10 // version 10
    msg[11] = service
    msg[15] = status
    data_s := len(data)
    if data_s > 0 {
        msg[8] = byte(data_s / 256)
        msg[9] = byte(data_s % 256)
        msg = append(msg, data...)
    }
    return msg
}

func handleRequest(conn net.Conn) {
    for {
        buf := make([]byte, 20)
        reqLen, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Error reading:", err.Error())
            fmt.Println("Closing connection")
            conn.Close()
            break
        }

        data_s := 256*int(buf[8]) + int(buf[9])
        data := make([]byte, data_s)
        if data_s > 0 {
            reqLen, err := conn.Read(data)
            if err != nil {
                fmt.Println("Error reading data:", err.Error())
            }
            _ = reqLen
        }
        datamap := make(map[string]string)
        if data_s > 0 {
            datamap = decodeData(data)
        }
        _ = datamap

        message := make([]byte, 20)
        if buf[11] == 76 {
            //handshake
            message = ymsg(76, 1, 0, []byte{})
        } else if buf[11] == 87 {
            //auth
            // FIXME: actually implement authorization
            username := datamap["1"]
            message = ymsg(87, 1, 1, encodeData([]string{"1", username, "94", "AA.BBCCDDEEGGHHIIH_HII--"}))
        } else if buf[11] == 84 {
            //authresp
            // FIXME: actually implement authorization
            username := datamap["1"]
            respdata := make([]string, 0)
            respdata = append(respdata, "87")
            respdata = append(respdata, "a\x0A")
            respdata = append(respdata, "88")
            respdata = append(respdata, "")
            respdata = append(respdata, "89")
            respdata = append(respdata, "")
            respdata = append(respdata, "59")
            respdata = append(respdata, "Y\x09")
            respdata = append(respdata, "59")
            respdata = append(respdata, "T\x09")
            respdata = append(respdata, "59")
            respdata = append(respdata, "C\x09mg=1")
            respdata = append(respdata, "3")
            respdata = append(respdata, username)
            respdata = append(respdata, "90")
            respdata = append(respdata, "1")
            respdata = append(respdata, "100")
            respdata = append(respdata, "0")
            respdata = append(respdata, "101")
            respdata = append(respdata, "")
            respdata = append(respdata, "102")
            respdata = append(respdata, "")
            respdata = append(respdata, "93")
            respdata = append(respdata, "86400")
            message = ymsg(85, 0, 1, encodeData(respdata))
            // 85: LIST
        }

        //message := buf
        n := bytes.Index(buf, []byte{0})
        _ = n
        fmt.Println(buf)
        fmt.Println(strconv.Itoa(reqLen))
        fmt.Printf("Sent back: ")
        fmt.Println(message)

        // Write the message in the connection channel.
        conn.Write([]byte(message));
    }
}

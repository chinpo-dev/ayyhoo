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

func ymsg(service byte, status byte, session_id byte) []byte {
  msg := make([]byte, 20)
  msg[0], msg[1], msg[2], msg[3] = 89, 77, 83, 71 // YMSG protocol
  msg[4] = 10 // version 10
  msg[11] = service
  msg[15] = status
  return msg
}

func decodeData(data []byte) map[string][]byte {
  dat := make(map[string][]byte)
  splat := bytes.Split(data, []byte{192,128})
  for i := 0 ; i<len(splat) -1 ; i+=2 {
    dat[string(splat[i])] = splat[i + 1]
  }
  fmt.Println(dat)
  return dat
}

func handleRequest(conn net.Conn) {
  for {
    buf := make([]byte, 20)
    reqLen, err := conn.Read(buf)
    if err != nil {
      fmt.Println("Error reading:", err.Error())
    }

    data_s := 256*int(buf[8]) + int(buf[9])
    data := make([]byte, data_s)
    if data_s > 0 {
      reqLen, err := conn.Read(data)
      if err != nil {
        fmt.Println("Error reading:", err.Error())
      }
      _ = reqLen
    }
    if data_s > 0 {
      decodeData(data)
    }

    message := make([]byte, 20)
    if buf[11] == 76 {
      //handshake
      message = ymsg(76, 1, 0)
    }

    //message := buf
    n := bytes.Index(buf, []byte{0})
    _ = n
    fmt.Println(buf)
    fmt.Println(strconv.Itoa(reqLen))
    fmt.Println("Sent back:")
    fmt.Println(message)

    // Write the message in the connection channel.
    conn.Write([]byte(message));
  }
}
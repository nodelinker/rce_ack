package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"time"
	"webray/pwnpwn/db"
	"webray/pwnpwn/tcpserver"
)

var listenAddr string
var port int
var zeroCopy bool
var secure bool

func main() {

	dbInst := db.GetInstance()
	dbInst.Db.AutoMigrate(&db.PwnDataCache{})

	tfMap := make(map[bool]string)
	tfMap[true] = "on"
	tfMap[false] = "off"

	flag.StringVar(&listenAddr, "listen", "127.0.0.1:55535", "server listen addr")
	flag.BoolVar(&zeroCopy, "zerocopy", false, "use splice/sendfile zero copy")
	flag.BoolVar(&secure, "secure", false, "use TLS")
	flag.Parse()

	fmt.Printf("Running echo server on %s\n", listenAddr)
	fmt.Printf(" - zerocopy: %s\n", tfMap[zeroCopy])
	fmt.Printf(" - TLS secured: %s\n", tfMap[secure])

	server, _ := tcpserver.NewServer(listenAddr)

	if secure {
		cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
		if err != nil {
			panic("Error loading servert cert and key file: " + err.Error())
		}
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
		server.SetTLSConfig(tlsConfig)
	}

	server.SetListenConfig(&tcpserver.ListenConfig{
		SocketReusePort:   false,
		SocketFastOpen:    false,
		SocketDeferAccept: false,
	})
	server.SetRequestHandler(requestHandler)
	var err error
	if secure {
		err = server.ListenTLS()
	} else {
		err = server.Listen()
	}

	if err != nil {
		panic("Error listening on interface: " + err.Error())
	}

	err = server.Serve()
	if err != nil {
		panic("Error serving: " + err.Error())
	}
}

func requestHandler(conn tcpserver.Connection) {

	var rawData = ""
	// 记录client src ip, client src port
	// fmt.Println(connAddr)

	sha1Inst := sha1.New()
	// buf := make([]byte, 4096)
	// for {
	// 	n, err := conn.Read(buf)
	// 	if err != nil && n == 0 || n == 0 {
	// 		break
	// 	}

	// 	data := hex.EncodeToString(buf[:n])
	// 	rawData += data

	// 	// padding sdfsdf
	// 	sha1Inst.Write(buf[:n])

	// 	_, _ = conn.Write(buf[:n])
	// }
	buf := make([]byte, 256)
	n, _ := conn.Read(buf)
	data := hex.EncodeToString(buf[:n])
	rawData += data
	sha1Inst.Write(buf[:n])

	connAddr := conn.GetClientAddr()
	sha1Result := sha1Inst.Sum([]byte(""))

	rawDataHash := hex.EncodeToString(sha1Result)

	dataCache := db.PwnDataCache{
		ClientIP:    connAddr.IP.String(),
		ClientPort:  connAddr.Port,
		RawDataHash: rawDataHash,
		RawData:     rawData,
		Created:     0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now()}

	dbInst := db.GetInstance()
	dbInst.Db.Create(&dataCache)

	fmt.Println(connAddr, hex.EncodeToString(sha1Result))
	conn.Close()

}

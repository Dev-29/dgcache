package main

import (
	"github.com/Dev-29/dgcache/cache"
	"log"
	"net"
	"time"
)

func main() {
	opts := ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}

	go func() {
		time.Sleep(time.Second * 2)
		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.Write([]byte("SET Foo Bar 2500"))
		if err != nil {
			return
		}
	}()

	server := NewServer(opts, cache.New())
	err := server.Start()
	if err != nil {
		return
	}
}

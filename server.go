package main

import (
	"context"
	"fmt"
	"github.com/Dev-29/dgcache/cache"
	"log"
	"net"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts

	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s\n", err)
	}

	log.Printf("server listening on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.HandleConn(conn)
	}
}

func (s *Server) HandleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			return
		}
	}()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error: %s", err)
			break
		}

		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)
	if err != nil {
		fmt.Println("failed to parse command: ", err)
		// respond
		return
	}

	switch msg.Cmd {
	case CMDSet:
		if err := s.handleSetCmd(conn, msg); err != nil {
			// respond
			return
		}
	}
}

func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	// goroutine for sendToFollowers()
	go func() {
		err := s.sendToFollowers(context.TODO(), msg)
		if err != nil {
			return
		}
	}()

	return nil
}

// Distribute cache to follower server nodes.
func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	return nil
}

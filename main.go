package main

import "github.com/Dev-29/dgcache/cache"

func main() {
	opts := ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}
	server := NewServer(opts, cache.New())
	err := server.Start()
	if err != nil {
		return
	}
}

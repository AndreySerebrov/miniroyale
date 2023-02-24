package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"pow/internal/client"
	"pow/internal/mgr"
	"pow/internal/payload"
	"pow/internal/pow"
	"pow/internal/server"
	"syscall"
	"time"
)

func main() {
	fmt.Printf(`
	-----------------------------
	Hi, there!
	This is a simple PoW example!
	It uses Hashcash algorithm with %d target bits
	-----------------------------
	`, pow.GetTargetButs())
	time.Sleep(time.Second * 5)
	ctx, cancel := context.WithCancel(context.Background())
	timeout := time.Second * 2
	m := mgr.New(timeout)
	p, err := payload.New("./assets/Words-of-Wisdom.txt")
	if err != nil {
		log.Fatalf("error opening file: %s", err.Error())
	}

	server.NewServer(ctx, m, p)
	c, err := client.New(timeout * 5)
	if err != nil {
		log.Fatalf("error creating client: %s", err.Error())
	}
	go func() {
		for {
			_, _ = c.Session()
			fmt.Println("-------------------------------------")
			time.Sleep(time.Second * 3)

		}
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
	fmt.Println("all done")
}

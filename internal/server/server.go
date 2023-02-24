package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"pow/internal/mgr"
	"pow/internal/payload"
	"pow/internal/pow"
	"strconv"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/google/uuid"
)

const (
	HOST = "localhost"
	PORT = "9001"
	TYPE = "tcp"
)

type Server struct {
	ctx      context.Context
	listener net.Listener
	payload  *payload.Payload
	mgr      *mgr.Mgr
}

func NewServer(ctx context.Context, mgr *mgr.Mgr, payload *payload.Payload) *Server {
	s := Server{mgr: mgr, payload: payload, ctx: ctx}
	listener, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal("server falls with error: ", err)
		os.Exit(1)
	}
	s.listener = listener
	go s.run(ctx)
	return &s
}

func (s *Server) run(ctx context.Context) {
	fmt.Println("server starts on port", PORT)
	for {
		c, err := s.listener.Accept()
		if err != nil {
			fmt.Println("accepting data error: ", err)
			return
		}
		if ctx.Err() != nil {
			fmt.Println("context closed")
			s.listener.Close()
		}
		go s.handleConnection(c)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println("handleConnection", err)
			panic(err)
			continue
		}

		temp := strings.TrimSpace(string(netData))
		if strings.HasPrefix(temp, "ask") {
			fmt.Print(color.Colorize(color.Blue, fmt.Sprintf("---> SERVER received a task request: %s\n", temp)))
			id := uuid.New().String()
			s.mgr.Push(id)
			result := id + "\n"
			c.Write([]byte(result))
			fmt.Print(color.Colorize(color.Blue, fmt.Sprintf("---> SERVER sends a task: %s\n", id)))
			continue
		}
		if strings.HasPrefix(temp, "answer:") {
			fmt.Print(color.Colorize(color.Blue, fmt.Sprintf("---> SERVER received an answer:%s\n", temp)))
			l := strings.Split(temp, ":")
			if len(l) != 3 {
				result := `waiting for request format like: "answer:uuid:nonce"` + "\n"
				c.Write([]byte(result))
				continue
			}
			code := l[1]
			nonceStr := l[2]
			err := s.mgr.Use(code)
			if err != nil {
				result := err.Error() + "\n"
				c.Write([]byte(result))
				continue
			}
			nonce, err := strconv.ParseUint(nonceStr, 10, 64)
			if err != nil {
				result := "wrong nonce format: " + err.Error() + "\n"
				c.Write([]byte(result))
				continue
			}
			_, ok, err := pow.CheckHash([]byte(code), nonce)
			if err != nil || !ok {
				result := "you were trying to confuse me: " + err.Error() + "\n"
				c.Write([]byte(result))
				continue
			}
			result, err := s.payload.GetRandomQuote()
			if err != nil {
				result := "i got a problem with my wisdom source: " + err.Error() + "\n"
				c.Write([]byte(result))
				continue
			}
			c.Write([]byte(result + "\n"))
		}
	}
}

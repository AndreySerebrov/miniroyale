package client

import (
	"bufio"
	"fmt"
	"math/big"
	"pow/internal/pow"
	"strings"
	"time"

	"net"

	"github.com/TwiN/go-color"
)

const (
	HOST = "localhost"
	PORT = "9001"
	TYPE = "tcp"
)

type Client struct {
	c       net.Conn
	timeout time.Duration
}

func New(timeout time.Duration) (*Client, error) {
repeat:
	c, err := net.Dial(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Println("client error", err.Error())
		time.Sleep(time.Millisecond * 100)
		goto repeat
	}
	fmt.Println("client created")
	return &Client{c: c, timeout: timeout}, nil
}

func (c *Client) Session() (string, error) {
	task, err := c.ask()
	if err != nil {
		return "", err
	}
	nonce, err := c.getNonce(task)
	if err != nil {
		return "", err
	}
	quote, err := c.sendAnswer(task, nonce)
	if err != nil {
		return "", err
	}
	return quote, nil
}

func (c *Client) ask() (string, error) {
	fmt.Print(color.Colorize(color.Green, "---> CLIENT asks for a task\n"))
	_, err := c.c.Write([]byte("ask\n"))
	if err != nil {
		return "", err
	}
	netData, err := bufio.NewReader(c.c).ReadString('\n')
	if err != nil {
		return "", err
	}
	task := strings.TrimSuffix(netData, "\n")
	fmt.Print(color.Colorize(color.Green, fmt.Sprintf("<--- CLIENT received response with code %s\n", netData)))
	return task, err
}

func (c *Client) getNonce(task string) (uint64, error) {
	ts := time.Now()
	readyCh := make(chan struct{})
	var nonce uint64
	var err error
	var hash big.Int
	go func() {
		hash, nonce, err = pow.GetHash(readyCh, []byte(task))
		if err != nil {
			fmt.Println("hash error", err)
		}
		readyCh <- struct{}{}
	}()
	select {
	case <-readyCh:
		fmt.Print(color.Colorize(color.Green, fmt.Sprintf("     CLIENT got a hash %06x in %s\n", hash.Bytes(), time.Since(ts).String())))
	case <-time.After(c.timeout):
		fmt.Print(color.Colorize(color.Red, fmt.Sprintf("     CLIENT couldn't find a nonce in %s\n", c.timeout.String())))
		return 0, fmt.Errorf("CLIENT couldn't find a nonce in %s", c.timeout.String())

	}
	return nonce, err
}

func (c *Client) sendAnswer(task string, nonce uint64) (string, error) {
	answer := fmt.Sprintf("answer:%s:%d\n", task, nonce)
	_, err := c.c.Write([]byte(answer))
	if err != nil {
		fmt.Println("send answer error", err.Error())
		return "", err
	}
	fmt.Print(color.Colorize(color.Green, fmt.Sprintf("---> CLIENT sends an answer %s\n", answer)))
	netData, err := bufio.NewReader(c.c).ReadString('\n')
	if err != nil {
		return "", err
	}

	fmt.Print(color.Colorize(color.Green, "---> CLIENT received a quote:\n"))
	if strings.HasPrefix(netData, "code expired") {
		fmt.Print(color.Colorize(color.Red, fmt.Sprintf("\n     %s\n", netData)))
	} else {
		fmt.Print(color.Colorize(color.Cyan, fmt.Sprintf("     %s\n", netData)))
	}
	return netData, err
}

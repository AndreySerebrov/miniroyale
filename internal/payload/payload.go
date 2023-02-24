package payload

import (
	"bufio"
	"time"

	"math/rand"
	"os"
)

const fileOffset = 12

type Payload struct {
	linesNumber int
	file        string
	f           *os.File
}

func (p *Payload) lineCounter() (int, error) {
	rand.Seed(time.Now().Unix())
	f, err := os.Open(p.file)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(f)
	cnt := 0

	for scanner.Scan() {
		cnt++
	}
	return cnt, nil
}

func New(file string) (*Payload, error) {
	f, err := os.Open(file)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}
	p := &Payload{f: f, file: file}
	cnt, err := p.lineCounter()
	if err != nil {
		return nil, err
	}
	p.linesNumber = cnt - fileOffset
	return p, nil
}

func (p *Payload) GetRandomQuote() (string, error) {
	rand.Seed(time.Now().Unix())
	f, err := os.Open(p.file)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}
	lineNum := rand.Intn(p.linesNumber) + fileOffset - 1
	scanner := bufio.NewScanner(f)
	cnt := 0
	lastLine := ""
	for scanner.Scan() {
		if scanner.Text() == "" {
			cnt++
			continue
		}
		lastLine = scanner.Text()
		if cnt < lineNum {
			cnt++
			continue
		}

		return scanner.Text(), nil
	}
	return lastLine, nil
}

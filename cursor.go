package goparser

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"io"
	"strings"
)

type Cursor struct {
	reader  *bufio.Reader
	current int
	lines   Lines
}

func NewCursor(r io.Reader) Cursor {
	reader := bufio.NewReader(r)
	return Cursor{
		reader: reader,
	}
}

func (c *Cursor) Next() {
	switch {
	case len(c.lines) < 1:
		c.current = 0
		c.Read()
	case c.current == len(c.lines):
		c.current++
		c.Read()
	}

}

func (c *Cursor) Prev() {
	if c.current == 0 {
		return
	}
	c.current--
}

func (c *Cursor) Line() Line {
	return c.lines[c.current]
}

func (c *Cursor) Read() {
	str, err := c.reader.ReadString('\n')
	if err != nil {
		log.Error(err)
	}
	line := c.Parse(str)
	c.lines = append(c.lines, line)
}

func (c *Cursor) Parse(str string) Line {
	items := strings.Split(str, " ")
	line := make(map[string]interface{})
	for _, item := range items {
		if strings.Contains(item, "=") {
			item = StripColor(item)
			s := strings.Split(item, "=")
			line[s[0]] = s[1]
		}
	}
	return line
}

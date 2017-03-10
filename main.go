package goparser

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Line map[string]interface{}
type Lines []Line

type Parser struct {
	delimiter string
	path      string
	split     bool
	file      *os.File
	scanner   *bufio.Scanner
	lines     []Line
}

func NewParser(path string, split bool) *Parser {
	return &Parser{
		delimiter: "=",
		split:     split,
		path:      path,
	}
}

func (p *Parser) Open() *Parser {
	var err error
	p.file, err = os.Open(p.path)
	if err != nil {
		log.Error("Could not open file: ", err)
	}
	p.scanner = bufio.NewScanner(p.file)
	return p
}

func (p *Parser) Close() *Parser {
	err := p.file.Close()
	if err != nil {
		log.Error("Could not close file: ", err)
	}
	return p
}

func (p *Parser) SetDelimiter(del string) *Parser {
	p.delimiter = del
	return p
}

func (p *Parser) Max(key string) Line {
	c := 0
	lines := p.GetLines()
	for _, line := range lines {
		if val, ok := line[key]; ok {
			k, err := strconv.Atoi(val.(string))
			if err != nil {
				log.Error("Key not an integer: ", err)
				continue
			}
			if k > c {
				c = k
			}
		}
	}
	return lines[c-1]
}

func (p *Parser) Find(key string, out *string) *Parser {
	if p.lines == nil {
		p.lines = p.GetLines()
	}
	for _, line := range p.lines {
		if value, ok := line[key]; ok {
			val := value.(string)
			*out = val
		}
	}
	return p
}

func (p *Parser) GetLines() []Line {
	lines := make([]Line, 0)
	for p.scanner.Scan() {
		line := p.Parse(p.scanner.Text())
		lines = append(lines, line)
	}

	return lines
}

// Parses one line of string and convert it to map based on delimiter
func (p *Parser) Parse(str string) Line {
	items := []string{str}
	if p.split {
		items = strings.Split(str, " ")
	}
	line := make(map[string]interface{})
	for _, item := range items {
		if strings.Contains(item, p.delimiter) {
			item = StripColor(item)
			s := strings.Split(item, p.delimiter)
			if _, ok := line[s[0]]; ok {
				continue
			}
			line[s[0]] = s[1]
		}
	}
	return line
}

func StripColor(str string) string {
	re := regexp.MustCompile("\x1b\\[[^m]+m")
	return re.ReplaceAllString(str, "")
}

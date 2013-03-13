package mailthread

import (
	"fmt"
	"regexp"
)

type contentBuffer struct {
	lastLine    string
	content     string
	atHeadStart bool
	atHeadEnd   bool
	inHead      bool
	inFwHead    bool
	bType       string // fw | re
}

const (
	re_type = "re"
	fw_type = "fw"
)

func (buffer *contentBuffer) parseIn(line string) {
	buffer.content += buffer.lastLine
	buffer.lastLine = line
	buffer.parseLastLine()
}

func (buffer *contentBuffer) clear() {
	buffer.content = buffer.lastLine
	buffer.lastLine = ""
}

func (buffer *contentBuffer) parseLastLine() {
	buffer.atHeadStart = false
	buffer.atHeadEnd = false

	switch {
	case buffer.isFwHeadStart():
		buffer.bType = fw_type

		buffer.atHeadStart = true
		buffer.inHead = true
		buffer.inFwHead = true
	case !buffer.inFwHead && buffer.isReHeadStart():
		buffer.bType = re_type

		buffer.atHeadStart = true
		buffer.inHead = true
	case buffer.inHead && buffer.isHeadEnd():

		buffer.atHeadEnd = true
		buffer.inHead = false
		buffer.inFwHead = false
	}
}

var fwHeadStartExp *regexp.Regexp

func (buffer *contentBuffer) isFwHeadStart() bool {
	return fwHeadStartExp.MatchString(buffer.lastLine)
}

var reHeadStartExp *regexp.Regexp

func (buffer *contentBuffer) isReHeadStart() bool {
	return reHeadStartExp.MatchString(buffer.lastLine)
}

var headEndExp *regexp.Regexp

func (buffer *contentBuffer) isHeadEnd() bool {
	return headEndExp.MatchString(buffer.lastLine)
}

func init() {
	initMailComp()

	var err error

	fwHeadStartExp, err = regexp.Compile(headComp.fw)
	if err != nil {
		panic(err)
	}

	reHeadStartExp, err = regexp.Compile(fmt.Sprintf(`(%s|%s)`, headComp.from, headComp.re))
	if err != nil {
		panic(err)
	}

	headEndExp, err = regexp.Compile(`^\n`)
	if err != nil {
		panic(err)
	}
}

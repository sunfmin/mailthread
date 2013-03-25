package mailthread

import (
	"fmt"
	"regexp"
)

type contentBuffer struct {
	lastLine           string
	content            string
	atHeadStart        bool
	atHeadEnd          bool
	inHead             bool
	inFwHead           bool
	legalHeadContent   bool
	bType              string // fw | re
	headEndLineContent string
}

const (
	re_type = "re"
	fw_type = "fw"
)

var (
	fwHeadStartExp *regexp.Regexp
	reHeadStartExp *regexp.Regexp
	headEndExp     *regexp.Regexp
	legalFwCompExp *regexp.Regexp
)

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

	headEndExp, err = regexp.Compile(`^>*\n`)
	if err != nil {
		panic(err)
	}

	legalFwCompExp = regexp.MustCompile(headComp.legalFwComp)
}

func (buffer *contentBuffer) parseIn(line string) {
	buffer.content += buffer.lastLine
	buffer.lastLine = line
	buffer.parseLastLine()
}

func (buffer *contentBuffer) clean() {
	buffer.content = buffer.lastLine
	buffer.lastLine = ""
}

func (buffer *contentBuffer) rewind() {
	buffer.clean()
	buffer.inHead = false
	buffer.inFwHead = false
}

func (buffer *contentBuffer) parseLastLine() {
	switch {
	case buffer.isFwHeadStart():
		buffer.bType = fw_type

		buffer.atHeadStart = true
		buffer.atHeadEnd = false

		buffer.inHead = true
		buffer.inFwHead = true
		buffer.legalHeadContent = true
	case !buffer.inFwHead && buffer.isReHeadStart():
		buffer.bType = re_type

		buffer.atHeadStart = true
		buffer.atHeadEnd = false

		buffer.inHead = true
		buffer.legalHeadContent = true
	case buffer.inHead && buffer.isHeadEnd():
		buffer.atHeadStart = false
		buffer.atHeadEnd = true

		buffer.inHead = false
		buffer.inFwHead = false
		buffer.headEndLineContent = buffer.lastLine
	case buffer.inFwHead && !buffer.isLegalFwComp():
		buffer.legalHeadContent = false
	default:
		buffer.atHeadStart = false
		buffer.atHeadEnd = false
	}
}

func (buffer *contentBuffer) isFwHeadStart() bool {
	return fwHeadStartExp.MatchString(buffer.lastLine)
}

func (buffer *contentBuffer) isReHeadStart() bool {
	return reHeadStartExp.MatchString(buffer.lastLine)
}

func (buffer *contentBuffer) isHeadEnd() bool {
	return headEndExp.MatchString(buffer.lastLine)
}

func (buffer *contentBuffer) isLegalFwComp() bool {
	return legalFwCompExp.MatchString(buffer.lastLine)
}

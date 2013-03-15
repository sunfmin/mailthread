package test

import (
	"bytes"
	"github.com/sunfmin/mailthread"
	"io"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"strings"
)

type ProcessSuite struct{}

var _ = Suite(&ProcessSuite{})

var testFiles = []string{
	"gmail_style/simple_forwarding",
	"gmail_style/simply_replied_forwarding",
	"hotmail/nested_replied",
	"hotmail/fw and re",
	"yahoo mail/message",
	"japanese",
	"forward_havent_support1",
}

func (s *ProcessSuite) TestProcess(c *C) {
	for _, file := range testFiles {
		input, err := ioutil.ReadFile("input/" + file + ".eml")
		if err != nil {
			c.Fatal(err)
		}
		expectedOutput, err := ioutil.ReadFile("output/" + file + ".html")
		if err != nil {
			c.Fatal(err)
		}

		processedInput := mailthread.ProcessString(string(input))

		c.Log("TEST FILE: ", file)
		c.Check(processedInput, Equals, string(expectedOutput))
	}
}

func (s *ProcessSuite) TestGoThoughProcess(c *C) {
	for _, file := range testFiles {
		input, err := ioutil.ReadFile("input/" + file + ".eml")
		if err != nil {
			c.Fatal(err)
		}

		processedInput := mailthread.ProcessStringWithHandler(string(input), &mailthread.GoThroughHandler{})

		c.Log("TEST FILE: ", file)
		c.Check(processedInput, Equals, string(input))
	}
}

type CustomizedContentHandler struct {
	*mailthread.GoThroughHandler
	mainContent        bytes.Buffer
	otherContent       bytes.Buffer
	mainContentFilled  bool
	mainContentFilling bool
}

func (qch *CustomizedContentHandler) ForwardHeader(w io.Writer, header string) (err error) {
	println(qch.mainContentFilling, "ForwardHeader")
	if qch.mainContentFilling {
		qch.mainContentFilled = true
		qch.mainContentFilling = false
	}
	if qch.mainContentFilled {
		qch.otherContent.WriteString(header)
	} else {
		qch.mainContent.WriteString(header)
	}
	return
}

func (qch *CustomizedContentHandler) ReplyHeader(w io.Writer, header string) (err error) {
	println(qch.mainContentFilling, "ReplyHeader")
	err = qch.ForwardHeader(w, header)
	return
}

func (qch *CustomizedContentHandler) Text(w io.Writer, text string) (err error) {
	if qch.mainContentFilled {
		qch.otherContent.WriteString(text)
	} else {
		qch.mainContent.WriteString(text)
		if len(strings.TrimSpace(text)) > 0 {
			qch.mainContentFilling = true
		}
	}
	return
}

func (s *ProcessSuite) TestCustomizedProcess(c *C) {
	file := "japanese"
	input, err := ioutil.ReadFile("input/" + file + ".eml")
	if err != nil {
		c.Fatal(err)
	}
	maincontent, _ := ioutil.ReadFile("output/" + file + "_customized_main.html")
	othercontent, _ := ioutil.ReadFile("output/" + file + "_customized_other.html")

	ch := &CustomizedContentHandler{}
	mailthread.ProcessStringWithHandler(string(input), ch)

	c.Check(ch.mainContent.String(), Equals, string(maincontent))
	c.Check(ch.otherContent.String(), Equals, string(othercontent))
}

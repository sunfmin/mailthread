package test

import (
	"bytes"
	"fmt"
	"github.com/sunfmin/mailthread"
	"io"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"strings"
)

var fl = fmt.Println
var fp = fmt.Printf

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
	"left-arrowed-mail",
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

		newInput := string(input)
		if mailthread.HasLeadingNestedMailArrow(string(input)) {
			newInput, err = mailthread.RemoveLeadingNestedArrows(string(input))
			if err != nil {
				c.Fatal(err)
			}
		}

		c.Log("TEST FILE: ", file)
		c.Check(processedInput, Equals, newInput)
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
	if qch.mainContentFilling {
		qch.mainContentFilled = true
		qch.mainContentFilling = false
	}
	if qch.mainContentFilled {
		qch.otherContent.WriteString(header)
		err = mailthread.SkipParseLeftError
	} else {
		qch.mainContent.WriteString(header)
	}
	return
}

func (qch *CustomizedContentHandler) ReplyHeader(w io.Writer, header string) (err error) {
	err = qch.ForwardHeader(w, header)
	return
}

func (qch *CustomizedContentHandler) Skip(r io.Reader, w io.Writer) (err error) {
	io.Copy(&qch.otherContent, r)
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

	newMaincontent := string(maincontent)
	if mailthread.HasLeadingNestedMailArrow(string(maincontent)) {
		newMaincontent, err = mailthread.RemoveLeadingNestedArrows(string(maincontent))
		if err != nil {
			c.Fatal(err)
		}
	}
	newOthercontent := string(othercontent)
	if mailthread.HasLeadingNestedMailArrow(string(othercontent)) {
		newOthercontent, err = mailthread.RemoveLeadingNestedArrows(string(othercontent))
		if err != nil {
			c.Fatal(err)
		}
	}

	ch := &CustomizedContentHandler{}
	mailthread.ProcessStringWithHandler(string(input), ch)

	c.Check(ch.mainContent.String(), Equals, newMaincontent)
	c.Check(ch.otherContent.String(), Equals, newOthercontent)
}

func (s *ProcessSuite) TestCustomizedProcess2(c *C) {
	file := "left-arrowed-mail"
	input, err := ioutil.ReadFile("input/" + file + ".eml")
	if err != nil {
		c.Fatal(err)
	}
	maincontent, _ := ioutil.ReadFile("output/" + file + "-customized_main.html")
	othercontent, _ := ioutil.ReadFile("output/" + file + "-customized_other.html")

	newMaincontent := string(maincontent)
	newOthercontent := string(othercontent)

	ch := &CustomizedContentHandler{}
	mailthread.ProcessStringWithHandler(string(input), ch)

	// fmt.Println("+++++++++++++++++++============+++++++++++++++++++============")
	// fmt.Println(ch.mainContent.String())
	// fmt.Println("+++++++++++++++++++============+++++++++++++++++++============")
	// fmt.Println(ch.otherContent.String())

	c.Check(ch.mainContent.String(), Equals, newMaincontent)
	c.Check(ch.otherContent.String(), Equals, newOthercontent)
}

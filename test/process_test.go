package test

import (
	"github.com/sunfmin/mailthread"
	"io/ioutil"
	. "launchpad.net/gocheck"
)

type ProcessSuite struct{}

var _ = Suite(&ProcessSuite{})

var testFiles = []string{
	"gmail_style/simple_forwarding",
	"gmail_style/simply_replied_forwarding",
	"hotmail/nested_replied",
	"hotmail/fw and re",
	"yahoo mail/message",
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

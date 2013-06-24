package mailthread

import (
	"io/ioutil"
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type ProcessTest struct{}

var _ = Suite(&ProcessTest{})

func (s *ProcessTest) TestLeadingNestedMailArrow(c *C) {
	c.Check(getContent("> \n"), Equals, "\n")
	c.Check(getContent(">   \n"), Equals, "  \n")
	c.Check(getContent("> > >something\n"), Equals, ">something\n")
	c.Check(getContent("> > > something\n"), Equals, "something\n")
	c.Check(getContent("> >something in between > something\n"), Equals, ">something in between > something\n")

	c.Check(len(leadingNestedMailArrow.FindStringSubmatch("lorem ipsum > >\n")), Equals, 0)
	c.Check(len(leadingNestedMailArrow.FindStringSubmatch(" > lorem ipsum\n")), Equals, 0)
	c.Check(len(leadingNestedMailArrow.FindStringSubmatch("non match> >something in between > something\n")), Equals, 0)

	c.Check(getContent(">"), Equals, "")
	c.Check(getContent(">\n"), Equals, "\n")
	c.Check(getContent("> > >\n"), Equals, "\n")
	c.Check(getContent("> > > \n"), Equals, "\n")
}

func (s *ProcessTest) TestRemoveLeadingNestedArrows(c *C) {
	str, err := RemoveLeadingNestedArrows(">\n > >\n")
	c.Check(err, Equals, nil)
	c.Check(str, Equals, "\n > >\n")

	str, err = RemoveLeadingNestedArrows(">\n> >\n")
	c.Check(err, Equals, nil)
	c.Check(str, Equals, "\n\n")

	str, err = RemoveLeadingNestedArrows("something\n>test")
	c.Check(err, Equals, nil)
	c.Check(str, Equals, "something\n>test")

	str, err = RemoveLeadingNestedArrows("something\n> test")
	c.Check(err, Equals, nil)
	c.Check(str, Equals, "something\ntest")

	str, err = RemoveLeadingNestedArrows("something\n>  test")
	c.Check(err, Equals, nil)
	c.Check(str, Equals, "something\n test")

	str, err = RemoveLeadingNestedArrows("something> something\n> > > > test")
	c.Check(err, Equals, nil)
	c.Check(str, Equals, "something> something\ntest")

	str, err = RemoveLeadingNestedArrows("\nsomething>\n\n\n\n something\n> > > > test\n")
	c.Check(err, Equals, nil)
	c.Check(str, Equals, "\nsomething>\n\n\n\n something\ntest\n")

	input, err := ioutil.ReadFile("fixtures/left-arrowed-mail.eml")
	if err != nil {
		c.Fatal(err)
	}
	expectedOutput, err := ioutil.ReadFile("fixtures/left-arrowed-mail.html")
	if err != nil {
		c.Fatal(err)
	}
	processedInput, err := RemoveLeadingNestedArrows(string(input))
	if err != nil {
		c.Fatal(err)
	}
	c.Log(processedInput)
	c.Check(processedInput, Equals, string(expectedOutput))
}

package mailthread

import (
    . "launchpad.net/gocheck"
    "io/ioutil"
)

// Hook up gocheck into the gotest runner.
// func Test(t *testing.T) { TestingT(t) }

type DocumentMatch struct{}
var _ = Suite(&DocumentMatch{})

func (s *DocumentMatch) TestSimpleForwarding(c *C) {
    input, err := ioutil.ReadFile("test/input/gmail_style/simple_forwarding.eml")
    if err != nil { 
        c.Fatal(err)
    }
    expectedOutput, err:= ioutil.ReadFile("test/output/gmail_style/simple_forwarding.html")
    if err != nil { 
        c.Fatal(err)
    }
    
    processedInput, err := Process(string(input))
    if err != nil { 
        c.Fatal(err)
    }
    
    c.Check(processedInput, Equals, string(expectedOutput))
}

func (s *DocumentMatch) TestSimplyReplyedForwarding(c *C) {
    input, err := ioutil.ReadFile("test/input/gmail_style/simply_replied_forwarding.eml")
    if err != nil { 
        c.Fatal(err)
    }
    expectedOutput, err:= ioutil.ReadFile("test/output/gmail_style/simply_replied_forwarding.html")
    if err != nil { 
        c.Fatal(err)
    }
    
    processedInput, err := Process(string(input))
    if err != nil { 
        c.Fatal(err)
    }
    
    // c.Log(processedInput)
    c.Check(processedInput, Equals, string(expectedOutput))
}
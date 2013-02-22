package mailthread

import (
    . "launchpad.net/gocheck"
    "regexp"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type SimpleMatch struct{}
var _ = Suite(&SimpleMatch{})

func (s *SimpleMatch) TestIsForwardingBlockStart(c *C) {
    c.Check(isForwardingBlockStart("---------- Forwarded message ----------\n"), Equals, true)
}

func (s *SimpleMatch) TestIsForwardingBlockEnd(c *C) {
    c.Check(isForwardingBlockEnd("\n"), Equals, true)
    c.Check(isForwardingBlockEnd("From: BOM.D.Van <bom.d.van@gmail.com>\n"), Not(Equals), true)
}

func (s *SimpleMatch) TestEmailOfRlComp(c *C) {
    email, err := regexp.Compile(rlComp.email)
    if err != nil { 
        c.Fatal(err)
    }
    c.Check(email.MatchString("<bom.d.van@gmail.com>"), Equals, true)
    c.Check(email.MatchString("<bom.d.van@.gmail.com>"), Equals, false)
}

func (s *SimpleMatch) TestMonthAndWeekOfRlComp(c *C) {
    fullMonth, err := regexp.Compile(rlComp.fullMonth)
    if err != nil { 
        c.Fatal(err)
    }
    abbrMonth, err := regexp.Compile(rlComp.abbrMonth)
    if err != nil { 
        c.Fatal(err)
    }
    fullWeek, err := regexp.Compile(rlComp.fullWeek)
    if err != nil { 
        c.Fatal(err)
    }
    abbrWeek, err := regexp.Compile(rlComp.abbrWeek)
    if err != nil { 
        c.Fatal(err)
    }

    fullMonths := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
    abbrMonths := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sept", "Oct", "Nov", "Dec"}
    fullWeeks := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
    abbrWeeks := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
    
    // month matching
    for i, month := range fullMonths {
        // full month matching
        c.Log("Month Name: ", month)
        c.Check(fullMonth.MatchString(month), Equals, true)
        
        // short month matching
        c.Log("Month Name: ", abbrMonths[i])
        c.Check(abbrMonth.MatchString(abbrMonths[i]), Equals, true)
    }
    
    // week matching
    for i := 0; i < 7; i++ {
        // full week matching
        c.Log("Week Name: ", fullWeeks[i])
        c.Check(fullWeek.MatchString(fullWeeks[i]), Equals, true)
        
        // short week matching
        c.Log("Week Name: ", abbrWeeks[i])
        c.Check(abbrWeek.MatchString(abbrWeeks[i]), Equals, true)
    }
}

func (s *SimpleMatch) TestYYYYMMDDOfRlComp(c *C) {
    yyyymmdd, err := regexp.Compile(rlComp.yyyymmdd)
    if err != nil { 
        c.Fatal(err)
    }
    c.Check(yyyymmdd.MatchString("2013/2/20"), Equals, true)
    c.Check(yyyymmdd.MatchString("2013-02-20"), Equals, true)
}

func (s *SimpleMatch) TestNameOfRlComp(c *C) {
    name, err := regexp.Compile(rlComp.name)
    if err != nil { 
        c.Fatal(err)
    }
    c.Check(name.MatchString("bom.d.van dd"), Equals, true)
}

func (s *SimpleMatch) TestTwentyFourHourClockOfRlComp(c *C) {
    // testReg, err := regexp.Compile(`[0-5]`)
    // c.Log(testReg.MatchString("6"))
    twentyFourHourClock, err := regexp.Compile(rlComp.twentyFourHourClock)
    if err != nil { 
        c.Fatal(err)
    }
    c.Check(twentyFourHourClock.MatchString("at 20:00"), Equals, true)
    c.Check(twentyFourHourClock.MatchString("at 24:60"), Equals, false)
}

func (s *SimpleMatch) TestTwelveHourClockOfRlComp(c *C) {
    twelveHourClock, err := regexp.Compile(rlComp.twelveHourClock)
    if err != nil { 
        c.Fatal(err)
    }
    c.Check(twelveHourClock.MatchString("at 7:38 PM"), Equals, true)
    c.Check(twelveHourClock.MatchString("at 12:60 PM"), Equals, false)
}

func (s *SimpleMatch) TestIsReplyLine(c *C) {
    var matchR bool
    matchR = isReplyLine("2013/2/20 BOM.D.Van <bom.d.van@gmail.com>\n")
    c.Check(matchR, Equals, true)
    matchR = isReplyLine("On Wednesday, February 20, 2013, BOM.D.Van wrote:\n")
    c.Check(matchR, Equals, true)
    matchR = isReplyLine("On Wed, Feb 20, 2013 at 7:38 PM, BOM.D.Van <bom.d.van@gmail.com> wrote:\n")
    c.Check(matchR, Equals, true)
    matchR = isReplyLine("On 2013/2/20, at 20:00, BOM.D.Van <bom.d.van@gmail.com> wrote:\n")
    c.Check(matchR, Equals, true)
    
}
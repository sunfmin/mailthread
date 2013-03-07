package mailthread

import (
	. "launchpad.net/gocheck"
	"regexp"
)

type HeadComp struct{}

var _ = Suite(&HeadComp{})

func (s *HeadComp) TestEmail(c *C) {
	email, err := regexp.Compile(headComp.email)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(email.MatchString("<bom.d.van@gtime.com>"), Equals, true)
	c.Check(email.MatchString("<bom.d.van@hotmail.com>"), Equals, true)
	c.Check(email.MatchString("<bom.d.van@.gtime.com>"), Equals, false)
}

func (s *HeadComp) TestName(c *C) {
	name, err := regexp.Compile(headComp.name)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(name.MatchString("bom.d.van dd"), Equals, true)
}

func (s *HeadComp) TestFrom(c *C) {
	exp, err := regexp.Compile(headComp.from)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(exp.MatchString("From: BOM.D.Van <bom.d.van@gmail.com>\n"), Equals, true)
	c.Check(exp.MatchString("From: bom.d.van@hotmail.com\n"), Equals, true)
}

func (s *HeadComp) TestFw(c *C) {
	exp, err := regexp.Compile(headComp.fw)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(exp.MatchString("---------- Forwarded message ----------\n"), Equals, true)
	c.Check(exp.MatchString("----- Forwarded Message -----\n"), Equals, true)
	c.Check(exp.MatchString("From: BOM.D.Van <bom.d.van@gmail.com>\n"), Equals, false)
}

func (s *HeadComp) TestRe(c *C) {
	exp, err := regexp.Compile(headComp.re)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(exp.MatchString("2013/2/20 BOM.D.Van <bom.d.van@gmail.com>\n"), Equals, true)
	c.Check(exp.MatchString("On Wednesday, February 20, 2013, BOM.D.Van wrote:\n"), Equals, true)
	c.Check(exp.MatchString("On Wed, Feb 20, 2013 at 7:38 PM, BOM.D.Van <bom.d.van@gmail.com> wrote:\n"), Equals, true)
	c.Check(exp.MatchString("On 2013/2/20, at 20:00, BOM.D.Van <bom.d.van@gmail.com> wrote:\n"), Equals, true)
	c.Check(exp.MatchString("2013年2月13日 17:38 Anatole Varin a@theplant.jp:\n"), Equals, true)
	c.Check(exp.MatchString("2013/2/8 Maki Oka maki@theplant.jp\n"), Equals, true)
	c.Check(exp.MatchString("2013/2/14 ku充田久 男保 <m_kubota@bkkn.co.jp>\n"), Equals, true)
	c.Check(exp.MatchString("2013/2/13 Anatole Varin <a@theplant.jp>\n"), Equals, true)
	c.Check(exp.MatchString(">> 2013/2/13 ku充田久 男保 <m_kubota@bkkn.co.jp>\n"), Equals, true)
	c.Check(exp.MatchString(">>> 2013/2/13 ku充田久 男保 <m_kubota@bkkn.co.jp>\n"), Equals, true)
	c.Check(exp.MatchString(">>>> 2013年2月13日 17:38 Anatole Varin <a@theplant.jp>:\n"), Equals, true)
	c.Check(exp.MatchString(">>>>>> 2013年2月13日 12:02 Anatole Varin <a@theplant.jp>:\n"), Equals, true)
	c.Check(exp.MatchString(">>>>>>>> 2013/2/8 Maki Oka <maki@theplant.jp>\n"), Equals, true)
	c.Check(exp.MatchString(">>>>>>>>>> 2013年2月4日 20:07 ku充田久 男保 <m_kubota@bkkn.co.jp>:\n"), Equals, true)
	c.Check(exp.MatchString(">>>>>>>>>>>> 2013/2/1 柿沼宇成 <tkakinuma@fabricant.co.jp>\n"), Equals, true)
	c.Check(exp.MatchString("2013年2月21日 19:13 ku久保田 充男 <m_kubota@nkb.co.jp>:\n"), Equals, true)
	c.Check(exp.MatchString("On Wed, Feb 20, 2013 at 7:37 PM, BOM.D.Van <bom.d.van@gmail.com> wrote:\n"), Equals, true)
	c.Check(exp.MatchString("On 2013/02/20 at 7:37 PM, BOM.D.Van <bom.d.van@gmail.com> wrote:\n"), Equals, true)
	c.Check(exp.MatchString("On 2013-02-20 at 7:37 PM, BOM.D.Van <bom.d.van@gmail.com> wrote:\n"), Equals, true)
	c.Check(exp.MatchString("On 2013 02 20 at 7:37 PM, BOM.D.Van <bom.d.van@gmail.com> wrote:\n"), Equals, true)
	c.Check(exp.MatchString("2013/2/27 Van Hu <bom_d_van@yahoo.com>\n"), Equals, true)
	c.Check(exp.MatchString("1995:01:24T09:08:17.1823213 Van Hu <bom_d_van@yahoo.com>\n"), Equals, true)
}

type TimeComp struct{}

var _ = Suite(&TimeComp{})

func (s *TimeComp) TestMonthAndWeek(c *C) {
	fullMonth, err := regexp.Compile(timeComp.fullMonth)
	if err != nil {
		c.Fatal(err)
	}
	abbrMonth, err := regexp.Compile(timeComp.abbrMonth)
	if err != nil {
		c.Fatal(err)
	}
	fullWeek, err := regexp.Compile(timeComp.fullWeek)
	if err != nil {
		c.Fatal(err)
	}
	abbrWeek, err := regexp.Compile(timeComp.abbrWeek)
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

func (s *TimeComp) TestYYYYMMDD(c *C) {
	yyyymmdd, err := regexp.Compile(timeComp.yyyymmdd)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(yyyymmdd.MatchString("2013/2/20"), Equals, true)
	c.Check(yyyymmdd.MatchString("2013-02-20"), Equals, true)
	c.Check(yyyymmdd.MatchString("2013年2月13日"), Equals, true)
	c.Check(yyyymmdd.MatchString("2013/2/8"), Equals, true)
	c.Check(yyyymmdd.MatchString("2013年2月4日"), Equals, true)
	c.Check(yyyymmdd.MatchString("2013 02 23"), Equals, true)
}

func (s *TimeComp) TestTwentyFourHourClock(c *C) {
	twentyFourHourClock, err := regexp.Compile(timeComp.twentyFourHourClock)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(twentyFourHourClock.MatchString("20:00"), Equals, true)
	c.Check(twentyFourHourClock.MatchString("24:60"), Equals, false)
}

func (s *TimeComp) TestTwelveHourClock(c *C) {
	twelveHourClock, err := regexp.Compile(timeComp.twelveHourClock)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(twelveHourClock.MatchString("7:38 PM"), Equals, true)
	c.Check(twelveHourClock.MatchString("12:60 PM"), Equals, false)
}

func (s *TimeComp) TestTimeZoneOffset(c *C) {
	timeZoneOffset, err := regexp.Compile(timeComp.timeZoneOffset)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(timeZoneOffset.MatchString("+12:30"), Equals, true)
	c.Check(timeZoneOffset.MatchString("+1230"), Equals, true)
	c.Check(timeZoneOffset.MatchString("-14:01"), Equals, false)
	c.Check(timeZoneOffset.MatchString("-12:60"), Equals, false)
}

func (s *TimeComp) TestFullTimeClock(c *C) {
	fullTimeClock, err := regexp.Compile(timeComp.fullTimeClock)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(fullTimeClock.MatchString("00:03:05"), Equals, true)
	c.Check(fullTimeClock.MatchString("23:59:59"), Equals, true)
	c.Check(fullTimeClock.MatchString("00:60:60"), Equals, false)
	c.Check(fullTimeClock.MatchString("24:01:01"), Equals, false)
}

func (s *TimeComp) TestYearDigit(c *C) {
	yearDigit, err := regexp.Compile(timeComp.yearDigit)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(yearDigit.MatchString("2012"), Equals, true)
	c.Check(yearDigit.MatchString("00101"), Equals, false)
}

func (s *TimeComp) TestMonthDigit(c *C) {
	monthDigit, err := regexp.Compile(timeComp.monthDigit)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(monthDigit.MatchString("01"), Equals, true)
	c.Check(monthDigit.MatchString("1"), Equals, true)
	c.Check(monthDigit.MatchString("9"), Equals, true)
	c.Check(monthDigit.MatchString("10"), Equals, true)
	c.Check(monthDigit.MatchString("11"), Equals, true)
	c.Check(monthDigit.MatchString("00"), Equals, false)
	c.Check(monthDigit.MatchString("13"), Equals, false)
}

func (s *TimeComp) TestDateDigit(c *C) {
	dateDigit, err := regexp.Compile(timeComp.dateDigit)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(dateDigit.MatchString("01"), Equals, true)
	c.Check(dateDigit.MatchString("1"), Equals, true)
	c.Check(dateDigit.MatchString("8"), Equals, true)
	c.Check(dateDigit.MatchString("9"), Equals, true)
	c.Check(dateDigit.MatchString("10"), Equals, true)
	c.Check(dateDigit.MatchString("00"), Equals, false)
	c.Check(dateDigit.MatchString("32"), Equals, false)
}

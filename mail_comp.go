package mailthread

import "fmt"

// =================================================================================
//
// NOTE<IMPORTANT>:
// Comments here are not update-to-date, see mail_comp_test.go for each compositions
//
// =================================================================================

type headCompS struct {
	email             string // <bom.d.van@gmail.com>
	name              string // ' BOM.D.Van ', ' Van Hu ', etc(it can be any character)
	fw                string // '---------- Forwarded message ----------', '----- Forwarded Message -----'
	re                string // '2013/2/20 BOM.D.Van <bom.d.van@gmail.com>', 'On Wednesday, February 20, 2013, BOM.D.Van wrote:', 'On Wed, Feb 20, 2013 at 7:38 PM, BOM.D.Van <bom.d.van@gmail.com> wrote:', 'On 2013/2/20, at 20:00, BOM.D.Van <bom.d.van@gmail.com> wrote:'
	bareRe            string // re without ^
	from              string // From: bom.d.van@hotmail.com, From: BOM.D.Van <bom.d.van@gmail.com>
	to                string // To: bom.d.van@hotmail.com, To: BOM Van <bom.d.van@gmail.com>
	subject           string // Subject: RE: email test
	date              string // Date: Wed, 27 Feb 2013 00:03:05 +0000, Date: 2013/2/20
	cc                string // Cc: bom_d_van@yahoo.com, CC: bom_d_van@yahoo.com
	sent              string // Sent: Wednesday, February 27, 2013 9:45 AM
	legalFwComp       string // (fw|to|subject|date|cc|sent)
	leadingLeftArrows string // '>', '> >', '> > >', etc
}

var headComp headCompS

type timeCompS struct {
	yearDigit           string // 0000-9999
	monthDigit          string // 01-12
	dateDigit           string // 01-31
	yyyymmdd            string // 2013/2/20, 2013-02-20, etc
	fullMonth           string // July
	abbrMonth           string // Jul
	fullWeek            string // Sunday
	abbrWeek            string // Sun
	twelveHourClock     string // 7:38 PM, 07:38 PM, etc
	twentyFourHourClock string // 20:00
	fullTimeClock       string // 00:03:05
	timeZoneOffset      string // +0000, +00:00, etc(-14:00 through +14:00)
}

var timeComp timeCompS

func initMailComp() {
	yearDigit := `(\b\d{4}\b)`                          // 0000-9999
	monthDigit := `(0[1-9]|\b[1-9]\b|1[0-2])`           // 01-12, 1-12
	dateDigit := `(0[1-9]|\b[1-9]\b|[1-2][0-9]|3[0-1])` // 01-31, 1-31
	twelveHourClock := `(0?\d|1[0-1]):[0-5]?\d (AM|PM)` // 7:38 PM, 07:38 PM, etc
	twentyFourHourClock := `(([0-1]\d|2[0-3]):[0-5]\d)` // 20:00

	timeComp = timeCompS{
		yearDigit:           yearDigit,
		monthDigit:          monthDigit,
		dateDigit:           dateDigit,
		yyyymmdd:            fmt.Sprintf(`(%s[^\d]?%s[^\d]?%s)`, yearDigit, monthDigit, dateDigit),
		fullMonth:           `(January|February|March|April|May|June|July|August|September|October|November|December)`,
		abbrMonth:           `(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sept|Oct|Nov|Dec)`,
		fullWeek:            `(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday)`,
		abbrWeek:            `(Mon|Tue|Wed|Thu|Fri|Sat|Sun)`,
		twelveHourClock:     twelveHourClock,
		twentyFourHourClock: twentyFourHourClock,
		fullTimeClock:       fmt.Sprintf(`(24:00:00|%s:[0-5]\d)`, twentyFourHourClock),
		timeZoneOffset:      `([+-]((0\d|1[0-3]):?[0-5]\d|14:00))`,
	}

	name := `(.+)` // 'BOM.D.Van', 'Van Hu', etc(it can be any character)

	email := `([_a-zA-Z0-9-]+(\.[_a-zA-Z0-9-]+)*@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.(([0-9]{1,3})|([a-zA-Z]{2,3})|(aero|coop|info|museum|name)))` // bom.d.van@gmail.com

	headComp = headCompS{
		email:             fmt.Sprintf(`((\\)?\<(%s|\[%s\]\("?mailto:%s"?\))(\\)?\>)`, email, email, email), // see mail_comp_test.go#TestEmail
		name:              name,
		fw:                `^(((-{5,20}) Forwarded [Mm]essage (-{5,20}))|_{32})\n`,
		from:              `^( *(\*\*)?From:(\*\*)? .+)\n`,
		to:                `^( *(\*\*)?To:(\*\*)? .+)\n`,
		subject:           `^( *(\*\*)?Subject:(\*\*)? .+)\n`,
		cc:                `^( *(\*\*)?(Cc|CC):(\*\*)? .+)\n`,
		leadingLeftArrows: `^(>+( >)*).*\n`,

		// Date: Friday, May 31, 2013 15:08:42
		date: fmt.Sprintf(
			`^( *(\*\*)?Date:(\*\*)? (%s, %s %s %s %s %s|%s|%s, %s %s, %s at %s|%s, (%s|%s) %s, %s %s))\n`,
			timeComp.abbrWeek,
			timeComp.dateDigit,
			timeComp.abbrMonth,
			timeComp.yearDigit,
			timeComp.fullTimeClock,
			timeComp.timeZoneOffset,
			timeComp.yyyymmdd,
			timeComp.abbrWeek,
			timeComp.abbrMonth,
			timeComp.dateDigit,
			timeComp.yearDigit,
			timeComp.twelveHourClock,
			timeComp.fullWeek,
			timeComp.fullMonth,
			timeComp.abbrMonth,
			timeComp.dateDigit,
			timeComp.yearDigit,
			timeComp.fullTimeClock,
		),

		sent: fmt.Sprintf(
			`^(Sent: %s, %s %s, %s %s)\n`,
			timeComp.fullWeek,
			timeComp.fullMonth,
			timeComp.dateDigit,
			timeComp.yearDigit,
			timeComp.twelveHourClock,
		),
	}

	// 2013/2/20 BOM.D.Van <bom.d.van@gmail.com>
	re1 := fmt.Sprintf(
		`.*%s.*%s.*%s.*`,
		timeComp.yyyymmdd,
		headComp.name,
		email,
	)
	// On Wednesday, February 20, 2013, BOM.D.Van wrote:
	re2 := fmt.Sprintf(
		`On %s, %s %s, %s, %s wrote:`,
		timeComp.fullWeek,
		timeComp.fullMonth,
		timeComp.dateDigit,
		timeComp.yearDigit,
		headComp.name,
	)
	// On Wed, Feb 20, 2013 at 7:38 PM, BOM.D.Van <bom.d.van@gmail.com> wrote:
	re3 := fmt.Sprintf(
		`On %s, %s %s, %s at %s, %s %s wrote:`,
		timeComp.abbrWeek,
		timeComp.abbrMonth,
		timeComp.dateDigit,
		timeComp.yearDigit,
		timeComp.twelveHourClock,
		headComp.name,
		headComp.email,
	)
	// On 2013/2/20, at 20:00, BOM.D.Van <bom.d.van@gmail.com> wrote:
	re4 := fmt.Sprintf(
		`On %s, at %s, %s %s wrote:`,
		timeComp.yyyymmdd,
		timeComp.twentyFourHourClock,
		headComp.name,
		headComp.email,
	)
	// On May 27, 2013, at 4:21 PM, Kilian Muster wrote:
	re5 := fmt.Sprintf(
		`On (%s|%s) %s, %s, at %s, %s wrote:`,
		timeComp.fullMonth,
		timeComp.abbrMonth,
		timeComp.dateDigit,
		timeComp.yearDigit,
		timeComp.twelveHourClock,
		headComp.name,
	)
	// On Mon, Jul 8, 2013 at 4:24 PM, Finance \<[finance.van-test@qortex.theplant-dev.com]("mailto:finance.van-test@qortex.theplant-dev.com")\> wrote:
	re6 := fmt.Sprintf(
		`On %s, (%s|%s) %s, %s at %s, %s %s wrote:`,
		timeComp.abbrWeek,
		timeComp.fullMonth,
		timeComp.abbrMonth,
		timeComp.dateDigit,
		timeComp.yearDigit,
		timeComp.twelveHourClock,
		headComp.name,
		headComp.email,
	)
	headComp.re = fmt.Sprintf(`(^(%s|%s|%s|%s|%s|%s) *?\n)`, re1, re2, re3, re4, re5, re6)
	headComp.bareRe = fmt.Sprintf(`((%s|%s|%s|%s|%s|%s) *?\n)`, re1, re2, re3, re4, re5, re6)

	headComp.legalFwComp = fmt.Sprintf(
		`(^(%s|%s|%s|%s|%s|%s|%s))`,
		headComp.fw,
		headComp.from,
		headComp.to,
		headComp.subject,
		headComp.cc,
		headComp.date,
		headComp.sent,
	)
}

// Expose reply head matching for out-of-package using.
func GetBareReComp() string {
	return headComp.bareRe
}

package mailthread

import (
    "regexp"
    "fmt"
)

// ===== forwarding block regexpes ===== //
var forwardingBlockStart = func() *regexp.Regexp {
    expr, err := regexp.Compile(`-{10} Forwarded message -{10}\n`)
    if err != nil { 
        panic(err)
    }
    
    return expr
}()

func isForwardingBlockStart(str string) bool {
    return forwardingBlockStart.MatchString(str)
}


var forwardingBlockEnd = func() *regexp.Regexp {
    expr, err := regexp.Compile(`\n`)
    if err != nil { 
        panic(err)
    }
    
    return expr
}()

func isForwardingBlockEnd(str string) bool {
    return forwardingBlockEnd.MatchString(str)
}


// ===== reply line regexpes ===== //
var rlComp = struct{
    email string
    name string
    yyyymmdd string
    fullMonth string
    abbrMonth string
    fullWeek string
    abbrWeek string
    twelveHourClock string
    twentyFourHourClock string
}{
    email: `(\<[_a-zA-Z0-9-]+(\.[_a-zA-Z0-9-]+)*@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.(([0-9]{1,3})|([a-zA-Z]{2,3})|(aero|coop|info|museum|name))\>)`, // <bom.d.van@gmail.com>
    name: `([\w\. ]+)`, // BOM.D.Van, Van Hu, etc
    yyyymmdd: `(\d{4}[/|-]\d{1,2}[/|-]\d{1,2})`, // 2013/2/20, 2013-02-20, etc
    fullMonth: `(January|February|March|April|May|June|July|August|September|October|November|December)`, // July
    abbrMonth: `(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sept|Oct|Nov|Dec)`, // Jul
    fullWeek: `(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday)`, // Sunday
    abbrWeek: `(Mon|Tue|Wed|Thu|Fri|Sat|Sun)`, // Sun
    twelveHourClock: `(at [0-1]?\d:[0-5]?\d (AM|PM))`, // at 7:38 PM
    twentyFourHourClock: `(at [0-2]?\d:[0-5]?\d\b)`, // at 20:00
}

func initReplyLines() []*regexp.Regexp {
    // 2013/2/20 BOM.D.Van <bom.d.van@gmail.com>
    rl1, err := regexp.Compile(fmt.Sprintf(`^%s %s %s\n`, rlComp.yyyymmdd, rlComp.name, rlComp.email))
    if err != nil { 
        panic(err)
    }
    
    // On Wednesday, February 20, 2013, BOM.D.Van wrote:
    rl2, err := regexp.Compile(fmt.Sprintf(`^On %s, %s \d{1,2}, \d{4}, %s wrote:\n`, rlComp.fullWeek, rlComp.fullMonth, rlComp.name))
    if err != nil { 
        panic(err)
    }
    
    // On Wed, Feb 20, 2013 at 7:38 PM, BOM.D.Van <bom.d.van@gmail.com> wrote:
    rl3, err := regexp.Compile(fmt.Sprintf(`^On %s, %s \d{1,2}, \d{4} %s, %s %s wrote:\n`, rlComp.abbrWeek, rlComp.abbrMonth, rlComp.twelveHourClock, rlComp.name, rlComp.email))
    if err != nil { 
        panic(err)
    }
    
    // On 2013/2/20, at 20:00, BOM.D.Van <bom.d.van@gmail.com> wrote:
    rl4, err := regexp.Compile(fmt.Sprintf(`^On %s, %s, %s %s wrote:\n`, rlComp.yyyymmdd, rlComp.twentyFourHourClock, rlComp.name, rlComp.email))
    if err != nil { 
        panic(err)
    }
    
    return []*regexp.Regexp{rl1, rl2, rl3, rl4}
}

var replyLines = initReplyLines()
func isReplyLine(str string) (result bool) {
    for _, rl := range replyLines {
        if rl.MatchString(str) {
            result = true
            break
        }
    }
    return
}

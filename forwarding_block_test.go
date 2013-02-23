package mailthread

import (
    . "launchpad.net/gocheck"
)

type ForwardingBlock struct{}
var _ = Suite(&ForwardingBlock{})

func (s *ForwardingBlock) TestIsForwardingBlockStart(c *C) {
    c.Check(isForwardingBlockStart("---------- Forwarded message ----------\n"), Equals, true)
}

func (s *ForwardingBlock) TestIsForwardingBlockEnd(c *C) {
    c.Check(isForwardingBlockEnd("\n"), Equals, true)
    c.Check(isForwardingBlockEnd("From: BOM.D.Van <bom.d.van@gmail.com>\n"), Not(Equals), true)
}
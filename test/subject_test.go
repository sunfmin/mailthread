package test

import (
	"github.com/sunfmin/mailthread"
	"testing"
)

type subjectCase struct {
	input    string
	expected string
}

var subjectCases = []subjectCase{
	{"Re: paygentデータ", "paygentデータ"},
	{"Re: Fwd: Please refer to the attached file  ", "Please refer to the attached file"},
	{"Fwd: Re: [go-nuts] Digest for golang-nuts@googlegroups.com - 25 Messages in 12 Topics", "[go-nuts] Digest for golang-nuts@googlegroups.com - 25 Messages in 12 Topics"},
	{"Fwd: Fwd: Re: レスポンシブWebはCreaitve Cloudで。 Adobe Edge Reflow提供開始！", "レスポンシブWebはCreaitve Cloudで。 Adobe Edge Reflow提供開始！"},
	{"Reply:我爱你13165", "我爱你13165"},
	{"Forward:我爱你13165", "我爱你13165"},
	{"回复:我爱你13165", "我爱你13165"},
	{"回复:Re:Fwd:我爱你13165", "我爱你13165"},
}

func TestTrimSubject(t *testing.T) {
	for _, c := range subjectCases {
		trimed := mailthread.TrimSubject(c.input)
		if trimed != c.expected {
			t.Errorf("expected %+v, but was %+v", c.expected, trimed)
		}
	}
}

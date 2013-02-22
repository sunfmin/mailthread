package mailthread

import (
	"strings"
)

var SubjectTrimWords = []string{
	"Re:",
	"Fwd:",
	"Forward:",
	"Reply:",
	"回复:",
}

func TrimSubject(subject string) (r string) {
	r = subject
	for _, w := range SubjectTrimWords {
		r = strings.TrimSpace(r)
		i := strings.Index(r, w)
		if i == 0 {
			r = r[i+len(w):]
			r = strings.TrimSpace(r)
			r = TrimSubject(r)
		}
	}
	return
}

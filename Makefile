.DEFAULT_GOAL := test
.PHONY: test

test:
	go test -i github.com/sunfmin/mailthread
	go test github.com/sunfmin/mailthread
	go test -i github.com/sunfmin/mailthread/test
	go test github.com/sunfmin/mailthread/test

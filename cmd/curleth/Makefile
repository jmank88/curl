.PHONY: install
install:
	go install .

.PHONY: testscripts
testscripts: install
	go install github.com/rogpeppe/go-internal/cmd/testscript@latest
	find testdata -type d | xargs -I % \
		sh -c 'ls %/*.txtar > /dev/null 2>&1 || return 0 && testscript $(TS_FLAGS) %/*.txtar'

.PHONY: testscripts-update
testscripts-update:
	make testscripts TS_FLAGS="-u"

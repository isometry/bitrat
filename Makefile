GOOS		?= $(shell uname -s | tr "[:upper:]" "[:lower:]")
TEST_FILES	?= *.go */*.go
TEST_KEY	?= testhmac

debug:
	go build -v

.PHONY: release bitrat.darwin bitrat.linux bitrat.windows
release: bitrat.darwin bitrat.linux bitrat.windows

bitrat.darwin:
	env GOOS=darwin go build -v -ldflags="-s -w" -o $@ && upx $@

bitrat.linux:
	env GOOS=linux go build -v -ldflags="-s -w" -o $@ && upx $@

bitrat.windows:
	env GOOS=windows go build -v -ldflags="-s -w" -o $@ && upx $@

test: test-call-style test-sort test-blake2b test-sha1 test-sha256

test-call-style:
	bash -c "diff -u <(./bitrat --sort *.go */*.go) <(./bitrat -rs -n '*.go')"

test-sort:
	bash -c "diff -u <(./bitrat -r -j1 -n '*.go') <(./bitrat -rs -n '*.go')"

test-blake2b:
	bash -c "diff -u <(b2sum $(TEST_FILES) | sort) <(./bitrat --hash blake2b-512 $(TEST_FILES) | sort)"

test-blake3:
	bash -c "diff -u <(b3sum $(TEST_FILES) | sort) <(./bitrat --hash blake3 $(TEST_FILES) | sort)"

test-sha1:
	bash -c "diff -u <(shasum $(TEST_FILES) | sort) <(./bitrat --hash sha1 $(TEST_FILES) | sort)"

test-sha256:
	bash -c "diff -u <(sha256sum $(TEST_FILES) | sort) <(./bitrat --hash sha256 $(TEST_FILES) | sort)"

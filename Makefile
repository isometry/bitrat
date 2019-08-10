GOOS		?= $(shell uname -s | tr "[:upper:]" "[:lower:]")
TEST_FILES	?= *.go */*.go
TEST_KEY	?= testhmac

debug:
	go build -v

.PHONY: release hashbat.darwin hashbat.linux hashbat.windows
release: hashbat.darwin hashbat.linux hashbat.windows

hashbat.darwin:
	env GOOS=darwin go build -v -ldflags="-s -w" -o $@ && upx $@

hashbat.linux:
	env GOOS=linux go build -v -ldflags="-s -w" -o $@ && upx $@

hashbat.windows:
	env GOOS=windows go build -v -ldflags="-s -w" -o $@ && upx $@

test: test-call-style test-sort test-blake2b test-sha1 test-sha256

test-call-style:
	bash -c "diff -u <(./hashbat --sort *.go */*.go) <(./hashbat -rs -n '*.go')"

test-sort:
	bash -c "diff -u <(./hashbat -r -j1 -n '*.go') <(./hashbat -rs -n '*.go')"

test-blake2b:
	bash -c "diff -u <(b2sum $(TEST_FILES) | sort) <(./hashbat --hash blake2b-512 $(TEST_FILES) | sort)"

test-sha1:
	bash -c "diff -u <(shasum $(TEST_FILES) | sort) <(./hashbat --hash sha1 $(TEST_FILES) | sort)"

test-sha256:
	bash -c "diff -u <(sha256sum $(TEST_FILES) | sort) <(./hashbat --hash sha256 $(TEST_FILES) | sort)"

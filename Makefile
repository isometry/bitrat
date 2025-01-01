GOOS		?= $(shell uname -s | tr "[:upper:]" "[:lower:]")
TEST_PATH	?= .
TEST_FILES	?= *.go */*.go
TEST_KEY	?= testhmac

export BITRAT_HIDDEN_DIRS=true
export BITRAT_HIDDEN_FILES=true
export BITRAT_INCLUDE_GIT=true

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

test: test-call-style test-sort test-blake2b test-blake3 test-sha1 test-sha256

test-call-style:
	bash -c "shopt -s globstar; diff -u <(./bitrat --sort **/*.go) <(./bitrat -rs -n '*.go')"

test-recurse:
	bash -c "diff -u <(find * -type f | xargs -P9 b3sum | sort) <(./bitrat -r --hash blake3 | sort)"

test-sort:
	bash -c "diff -u <(./bitrat -r -j1 -n '*.go') <(./bitrat -rs -n '*.go')"

test-blake2b:
	bash -c "diff -u <(b2sum $(TEST_FILES) | sort) <(./bitrat --hash blake2b-512 $(TEST_FILES) | sort)"

test-blake3:
	bash -c "diff -u <(b3sum $(TEST_FILES) | sort) <(./bitrat --hash blake3 $(TEST_FILES) | sort)"

test-sha256-performance:
	hyperfine --warmup 1 \
		'bfs $(TEST_PATH) -type f -print0 | xargs -0 -n4 -P8 openssl sha256' \
		'bfs $(TEST_PATH) -type f -print0 | xargs -0 -n4 -P8 sha256sum' \
		'./bitrat -r $(TEST_PATH) --hash sha256' \
		'./bitrat -r $(TEST_PATH) --hash sha256-simd'

test-b3sum-performance:
	hyperfine --warmup 1 \
		'bfs  $(TEST_PATH) -type f -print0 | xargs -0 -P1  b3sum --num-threads=16' \
		'bfs  $(TEST_PATH) -type f -print0 | xargs -0 -P2  b3sum --num-threads=8' \
		'bfs  $(TEST_PATH) -type f -print0 | xargs -0 -P4  b3sum --num-threads=4' \
		'bfs  $(TEST_PATH) -type f -print0 | xargs -0 -P8  b3sum --num-threads=2' \
		'bfs  $(TEST_PATH) -type f -print0 | xargs -0 -P16 b3sum --num-threads=1' \

test-blake3-performance:
	hyperfine --warmup 1 \
		'bfs  $(TEST_PATH) -type f -print0 | xargs -0 -P16 b3sum --num-threads=4' \
		'./bitrat -r $(TEST_PATH) --hash blake3'

test-sha1:
	bash -c "diff -u <(shasum $(TEST_FILES) | sort) <(./bitrat --hash sha1 $(TEST_FILES) | sort)"

test-sha256:
	bash -c "diff -u <(sha256sum $(TEST_FILES) | sort) <(./bitrat --hash sha256 $(TEST_FILES) | sort)"

test-hmac-sha256:
	bash -c "diff -u <(openssl dgst -sha256 -hmac $(TEST_KEY) -r $(TEST_FILES) | tr '*' ' ' | sort) <(./bitrat --hash sha256 --hmac $(TEST_KEY) $(TEST_FILES) | sort)"

test-hmac-sha3-512:
	bash -c "diff -u <(openssl dgst -sha3-512 -hmac $(TEST_KEY) -r $(TEST_FILES) | tr '*' ' ' | sort) <(./bitrat --hash sha3-512 --hmac $(TEST_KEY) $(TEST_FILES) | sort)"

test-hmac-blake3:
	bash -c "diff -u <(b3sum --derive-key $(TEST_KEY) $(TEST_FILES) | sort) <(./bitrat --hash blake3 --hmac $(TEST_KEY) $(TEST_FILES) | sort)"

test-parallel:
	hyperfine --warmup 1 \
		'./bitrat -r -j1  $(TEST_PATH)' \
		'./bitrat -r -j2  $(TEST_PATH)' \
		'./bitrat -r -j3  $(TEST_PATH)' \
		'./bitrat -r -j4  $(TEST_PATH)' \
		'./bitrat -r -j5  $(TEST_PATH)' \
		'./bitrat -r -j6  $(TEST_PATH)' \
		'./bitrat -r -j7  $(TEST_PATH)' \
		'./bitrat -r -j8  $(TEST_PATH)' \
		'./bitrat -r -j9  $(TEST_PATH)' \
		'./bitrat -r -j10 $(TEST_PATH)' \
		'./bitrat -r -j11 $(TEST_PATH)' \
		'./bitrat -r -j12 $(TEST_PATH)' \
		'./bitrat -r -j13 $(TEST_PATH)' \
		'./bitrat -r -j14 $(TEST_PATH)' \
		'./bitrat -r -j15 $(TEST_PATH)' \
		'./bitrat -r -j16 $(TEST_PATH)'

test-hash:
	hyperfine --warmup 1 \
		'./bitrat -r $(TEST_PATH) --hash crc32' \
		'./bitrat -r $(TEST_PATH) --hash md5' \
		'./bitrat -r $(TEST_PATH) --hash sha1' \
		'./bitrat -r $(TEST_PATH) --hash sha256' \
		'./bitrat -r $(TEST_PATH) --hash sha3-512' \
		'./bitrat -r $(TEST_PATH) --hash blake2b-512' \
		'./bitrat -r $(TEST_PATH) --hash blake3'

COMMON_BZL_FLAGS := --test_output=all --cache_test_results=no

TEST := bazel test $(COMMON_BZL_FLAGS)

all:
	bazel run //:dora

build:
	bazel build //server:go_default_library

fmt:
	go fmt $(go list ./... | grep -v /vendor/)

bazel-update:
	bazel run //:gazelle
	python fix_update.py

update: bazel-update
	glide update

clean:
	glide cc
	bazel clean

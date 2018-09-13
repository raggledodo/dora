COMMON_BZL_FLAGS := --test_output=all --cache_test_results=no

TEST := bazel test $(COMMON_BZL_FLAGS)

all:
	bazel run //:dora

build_dora:
	bazel build //:dora

fmt:
	go fmt $(go list ./... | grep -v /vendor/)

bazel-update:
	bazel run //:gazelle
	python fix_update.py

glide-update:
	glide update

update: glide-update bazel-update

clean:
	glide cc
	bazel clean

COMMON_BZL_FLAGS := --test_output=all --cache_test_results=no

TEST := bazel test $(COMMON_BZL_FLAGS)

all:
	bazel run //server:main

build_dora:
	bazel build ...

fmt:
	gofmt -s -w .

update:
	gazelle fix

clean:
	bazel clean

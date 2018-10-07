COMMON_BZL_FLAGS := --test_output=all --cache_test_results=no

PLATFORM_FLAG := --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64

TEST := bazel test $(COMMON_BZL_FLAGS)

all:
	bazel run //server:main

docker:
	bazel run $(PLATFORM_FLAG) //server:dora_image

push:
	bazel run $(PLATFORM_FLAG) //server:dora_push

fmt:
	gofmt -s -w .

update:
	gazelle fix

clean:
	bazel clean

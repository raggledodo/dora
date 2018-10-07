load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def dependencies():
    # go dependency
    if "io_bazel_rules_go" not in native.existing_rules():
        http_archive(
            name = "io_bazel_rules_go",
            url = "https://github.com/bazelbuild/rules_go/releases/download/0.15.4/rules_go-0.15.4.tar.gz",
            sha256 = "7519e9e1c716ae3c05bd2d984a42c3b02e690c5df728dc0a84b23f90c355c5a1",
        )

    if "bazel_gazelle" not in native.existing_rules():
        http_archive(
            name = "bazel_gazelle",
            urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.14.0/bazel-gazelle-0.14.0.tar.gz"],
            sha256 = "c0a5739d12c6d05b6c1ad56f2200cb0b57c5a70e03ebd2f7b87ce88cabf09c7b",
        )

    if "io_bazel_rules_docker" not in native.existing_rules():
        http_archive(
            name = "io_bazel_rules_docker",
            sha256 = "29d109605e0d6f9c892584f07275b8c9260803bf0c6fcb7de2623b2bedc910bd",
            strip_prefix = "rules_docker-0.5.1",
            urls = ["https://github.com/bazelbuild/rules_docker/archive/v0.5.1.tar.gz"],
        )

    if "com_github_grpc_ecosystem_grpc_gateway" not in native.existing_rules():
        git_repository(
            name = "com_github_grpc_ecosystem_grpc_gateway",
            remote = "https://github.com/grpc-ecosystem/grpc-gateway",
            commit = "aeab1d96e0f1368d243e2e5f526aa29d495517bb",
        )

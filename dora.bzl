load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def dependencies():
    if "com_github_mingkaic_testify" not in native.existing_rules():
        git_repository(
            name = "com_github_mingkaic_testify",
            remote = "https://github.com/raggledodo/testify",
            commit = "efc070cbd323aeeb583edaaac20339552d2f25a7",
        )

    # go dependency
    if "io_bazel_rules_go" not in native.existing_rules():
        http_archive(
            name = "io_bazel_rules_go",
            urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.10.3/rules_go-0.10.3.tar.gz"],
            sha256 = "feba3278c13cde8d67e341a837f69a029f698d7a27ddbb2a202be7a10b22142a",
        )

    if "bazel_gazelle" not in native.existing_rules():
        http_archive(
            name = "bazel_gazelle",
            urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.13.0/bazel-gazelle-0.13.0.tar.gz"],
            sha256 = "bc653d3e058964a5a26dcad02b6c72d7d63e6bb88d94704990b908a1445b8758",
        )

    if "io_bazel_rules_docker" not in native.existing_rules():
        http_archive(
            name = "io_bazel_rules_docker",
            sha256 = "29d109605e0d6f9c892584f07275b8c9260803bf0c6fcb7de2623b2bedc910bd",
            strip_prefix = "rules_docker-0.5.1",
            urls = ["https://github.com/bazelbuild/rules_docker/archive/v0.5.1.tar.gz"],
        )

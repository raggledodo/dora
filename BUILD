load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_binary")

licenses(["notice"])

load("@bazel_gazelle//:def.bzl", "gazelle")

# ===== GO GENERATOR =====

# gazelle:prefix github.com/raggledodo/dora
gazelle(name = "gazelle")

go_binary(
    name = "dora",
    srcs = ["main.go"],
    deps = [
        "//server:go_default_library",
        "//vendor/github.com/sirupsen/logrus:go_default_library",
    ],
)

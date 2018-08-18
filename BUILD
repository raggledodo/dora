load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_binary")
load("@bazel_gazelle//:def.bzl", "gazelle")

licenses(["notice"])

package(default_visibility = ["//visibility:public"])

# ===== GO GENERATOR =====

# gazelle:prefix github.com/raggledodo/dora
gazelle(name = "gazelle")

go_binary(
    name = "dora",
    embed = [":go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/raggledodo/dora",
    deps = [
        "//server:go_default_library",
        "//vendor/github.com/sirupsen/logrus:go_default_library",
    ],
)

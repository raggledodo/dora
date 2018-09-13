load("@io_bazel_rules_docker//container:container.bzl", "container_push", "container_image")
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_binary")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
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

# docker image
go_image(
    name = "dora_image_base",
    embed = [":go_default_library"],
)

container_image(
    name = "dora_image",
    base = ":dora_image_base",
    ports = ["8581"],
)

container_push(
    name = "dora_push",
    format = "Docker",
    image = ":dora_image",
    registry = "index.docker.io",
    repository = "mkaichen/dora_server",
    tag = "latest",
)

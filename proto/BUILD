package(
    default_visibility = ["//visibility:public"],
)

licenses(["notice"])

load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

go_proto_library(
    name = "go_default_library",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/raggledodo/dora/proto",
    proto = "@com_github_mingkaic_testify//proto:proto",
)

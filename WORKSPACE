workspace(name = "com_github_raggledodo_dora")

load("//:dora.bzl", "dependencies")
dependencies()

load("@com_github_raggledodo_protos//:protos.bzl", "dependencies")
dependencies()

# go dependencies
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

load("@org_pubref_rules_protobuf//go:rules.bzl", "go_proto_repositories")
go_proto_repositories()

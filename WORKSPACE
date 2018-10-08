workspace(name = "com_github_raggledodo_dora")

# gazelle:repo bazel_gazelle

load("//:dora.bzl", "dependencies")

dependencies()

load("@com_github_grpc_ecosystem_grpc_gateway//:repositories.bzl", "repositories")

# go dependencies
load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_rules_dependencies",
    "go_register_toolchains",
)
load(
    "@bazel_gazelle//:deps.bzl",
    "gazelle_dependencies",
    "go_repository",
)
load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

repositories()

go_rules_dependencies()

go_register_toolchains()

gazelle_dependencies()

_go_image_repos()

go_repository(
    name = "org_golang_google_genproto",
    commit = "383e8b2c3b9e36c4076b235b32537292176bae20",
    importpath = "google.golang.org/genproto",
)

go_repository(
    name = "com_github_rogpeppe_fastuuid",
    commit = "6724a57986aff9bff1a1770e9347036def7c89f6",
    importpath = "github.com/rogpeppe/fastuuid",
)

go_repository(
    name = "com_github_go_resty_resty",
    commit = "f8815663de1e64d57cdd4ee9e2b2fa96977a030e",
    importpath = "github.com/go-resty/resty",
)

go_repository(
    name = "com_github_ghodss_yaml",
    commit = "0ca9ea5df5451ffdf184b4428c902747c2c11cd7",
    importpath = "github.com/ghodss/yaml",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    commit = "eb3733d160e74a9c7e442f435eb3bea458e1d19f",
    importpath = "gopkg.in/yaml.v2",
)

load("@org_pubref_rules_protobuf//cpp:rules.bzl", "cpp_proto_repositories")
cpp_proto_repositories()

load("@org_pubref_rules_protobuf//python:rules.bzl", "py_proto_repositories")
py_proto_repositories()

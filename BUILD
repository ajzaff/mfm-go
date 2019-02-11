load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_binary(
    name = "mfmgo",
    srcs = ["main.go"],
    visibility = ["//visibility:public"],
    deps = [
        "//experimental/computation/mfm-go:go_default_library",
        "//experimental/computation/mfm-go/hook:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"]),
    importpath = "ajz_xyz/experimental/computation/mfm-go",
    visibility = ["//visibility:public"],
    deps = [
        "//numerics/random/xorshift64star:go_default_library",
        "@com_github_fatih_color//:color",
        "@com_github_nsf_termbox-go//:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["bench_test.go"],
    deps = [
        "//experimental/computation/mfm-go:go_default_library",
        "//experimental/computation/mfm-go/atom:go_default_library",
    ],
)

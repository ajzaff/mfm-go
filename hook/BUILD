load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"]),
    importpath = "ajz_xyz/experimental/computation/mfm-go/hook",
    visibility = ["//visibility:public"],
    deps = [
        "//experimental/computation/mfm-go:go_default_library",
        "//experimental/computation/mfm-go/atom:go_default_library",
        "@com_github_nsf_termbox-go//:go_default_library",
    ],
)

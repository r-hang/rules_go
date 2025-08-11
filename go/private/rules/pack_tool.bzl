# Copyright 2025 The Bazel Go Rules Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load(
    "//go/private:common.bzl",
    "GO_TOOLCHAIN",
)
load(
    "//go/private:context.bzl",
    "go_context",
)
load(
    "//go/private:providers.bzl",
    "GoConfigInfo",
)
load(
    "//go/private/actions:pack_tool.bzl",
    "emit_pack_tool",
)

def _pack_tool_impl(ctx):
    go = go_context(ctx, include_deprecated_properties = False)
    return [emit_pack_tool(go)]

go_pack_tool = rule(
    implementation = _pack_tool_impl,
    cfg = "exec",
    attrs = {
        "cgo_context_data": attr.label(),
        "_go_config": attr.label(
            default = "//:go_config",
            providers = [GoConfigInfo],
        ),
    },
    doc = """go_pack_tool builds the pack command for the execution platform.

This rule builds the cmd/pack command from Go source code for the execution
platform. This ensures that the pack tool is always available and compatible
with the execution platform where compilation actions run.""",
    toolchains = [GO_TOOLCHAIN],
)
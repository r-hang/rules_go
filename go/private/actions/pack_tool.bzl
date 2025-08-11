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
    "GO_TOOLCHAIN_LABEL",
    "SUPPORTS_PATH_MAPPING_REQUIREMENT",
)
load(
    "//go/private:providers.bzl",
    "GoPackTool",
)

def emit_pack_tool(go):
    """Builds the pack tool for the execution platform.

    Args:
        go: GoContext object.

    Returns:
        GoPackTool provider with the pack executable.
    """
    sdk = go.sdk
    pack_tool = go.declare_file(go, path = "bin/pack")
    args = go.builder_args(go, "packtool", use_path_mapping = True)
    
    # Use a file rather than pack_tool.dirname as the latter is just a string and thus
    # not subject to path mapping.
    args.add("-out", pack_tool.path)

    inputs_direct = [sdk.go, sdk.package_list, sdk.root_file]
    inputs_transitive = [sdk.headers, sdk.srcs, sdk.tools]

    go.actions.run(
        inputs = depset(direct = inputs_direct, transitive = inputs_transitive),
        outputs = [pack_tool],
        mnemonic = "GoPackTool",
        executable = go.toolchain._builder,
        arguments = [args],
        env = _build_env(go),
        toolchain = GO_TOOLCHAIN_LABEL,
        execution_requirements = SUPPORTS_PATH_MAPPING_REQUIREMENT,
    )
    
    return GoPackTool(
        pack = pack_tool,
    )

def _build_env(go):
    """Build environment for pack tool, simplified from stdlib env."""
    env = go.env
    
    # Set basic Go environment for building pack
    env.update({
        "CGO_ENABLED": "0",  # Pack tool doesn't need CGO
    })
    
    return env
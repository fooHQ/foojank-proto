#!/usr/bin/env bash

set -euo pipefail

CAPNP_JAVA_INCLUDE="$(nix path-info nixpkgs#capnproto-java)/include"

install() {
    install_capnproto_go
}

install_capnproto_go() {
    if [ ! -d "./build/go-capnp" ]; then
        git clone -b v3.1.0-alpha.2 --depth 1 https://github.com/capnproto/go-capnp ./build/go-capnp
    fi

    cd ./build/go-capnp || exit 1
    go build -modfile go.mod -o ../capnpc-go ./capnpc-go
    cd - || exit 1
}

build() {
    local schemas=(
        "./schema/agent.capnp"
    )
    for schema in "${schemas[@]}"; do
        build_schema_go "$schema"
        build_schema_java "$schema"
    done
}

build_schema_go() {
    local output_dir="./go/$(basename "$1" ".capnp")"
    mkdir -p "$output_dir"
    capnp compile --src-prefix schema -I ./build/go-capnp/std/ -I "$CAPNP_JAVA_INCLUDE" -o ./build/capnpc-go:"$output_dir" "$1"

    local module_path="$(grep '\$Go\.import' "$1" | awk -F'"' '{print $2}')"
    rm -f "$output_dir/go.mod" "$output_dir/go.sum"
    go -C "$output_dir" mod init "$module_path"
    go -C "$output_dir" mod tidy
}

build_schema_java() {
    local output_dir="./java/io/github/foohq/foojank/$(basename "$1" ".capnp")"
    mkdir -p "$output_dir"
    capnp compile --src-prefix schema -I ./build/go-capnp/std/ -I "$CAPNP_JAVA_INCLUDE" -o java:"$output_dir" "$1"
}

eval $@

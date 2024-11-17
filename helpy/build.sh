#!/bin/bash

# clean up
rm -rf ./bin
rm -rf ./target
mkdir -p ./bin

# run Rust build with Cargo
cargo build --release

# copy all binaries
cp target/release/scw-* ./bin/
cp target/release/bunny-* ./bin/

# remove all .d files
rm -rf ./bin/*.d

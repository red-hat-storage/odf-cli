#!/bin/sh
# Download a nightly build of odf cli executable for the current platform and
# architecture.

set -e

version="$1"

if [ -z "$version" ]; then
    cat <<EOF
Usage: $0 VERSION

Examples:

  $0 v4.19

EOF
    exit 1
fi

image="quay.io/rhceph-dev/odf4-odf-cli-rhel9:$version"
platform=$(uname)

case $platform in
Darwin)
    platform=macosx
    executable=odf
    ;;
Linux)
    platform=linux
    machine=$(uname -m)
    case $machine in
    x86_64)
        executable=odf-amd64
        ;;
    aarch64)
        executable=odf-arm64
        ;;
    ppc64le,s390x)
        executable=odf-$machine
        ;;
    *)
        echo "unsupported machine $machine"
        exit 1
        ;;
    esac
    ;;
*)
    echo "unsupported platform $platform"
    exit 1
    ;;
esac

# We have only linux/amd64, linux/s390x, and linux/ppc64le images. The
# linux/amd64 image includes builds for all platforms.
oc image extract "$image" \
    --filter-by-os "^linux/amd64" \
    --file "/usr/share/odf/$platform/$executable"

chmod +x "$executable"

if [ "$executable" != "odf" ]; then
    mv "$executable" odf
fi

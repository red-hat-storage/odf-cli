#!/bin/sh
# Download a nightly build of odf cli executable for the current platform and
# architecture.

set -e

image="quay.io/rhceph-dev/odf4-odf-cli-rhel9:v4.19"
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

oc image extract "$image" --file "/usr/share/odf/$platform/$executable"

chmod +x "$executable"

if [ "$executable" != "odf" ]; then
    mv "$executable" odf
fi

#!/usr/bin/env bash

set -e

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
SDIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

export CGO_ENABLED=0

source ${SDIR}/env.sh
MGDIR_GOPATH=${MGDIR}

pushd ${MGDIR_GOPATH}/cmd/mailguard

#BUILD_TIME="$(date -u '+%Y-%m-%d_%I:%M:%S%p')"
#TAG="current"
#REVISION="current"
#if hash git 2>/dev/null && [ -e ${MGDIR_GOPATH}/.git ]; then
#  TAG="$(git describe --tags)"
#  REVISION="$(git rev-parse HEAD)"
#fi

#LD_FLAGS="-X github.com/mailguard/mailguard/pkg/version.appVersionTag=${TAG} -X github.com/mailguard/mailguard/pkg/version.appVersionRev=${REVISION} -X github.com/mailguard/malguard/pkg/version.appVersionTime=${BUILD_TIME}"
#gox -osarch="linux/amd64" -ldflags "${LD_FLAGS}" -output "${MGDIR_GOPATH}/bin/linux/mailguard"
#gox -osarch="darwin/amd64" -ldflags "${LD_FLAGS}" -output "${MGDIR_GOPATH}/bin/mac/mailguard"
#gox -osarch="linux/arm" -output "$MGDIR_GOPATH/bin/linux_arm/mailguard"
go build
mv mailguard* ${MGDIR_GOPATH}/bin/
popd
#rm -rfv ${MGDIR_GOPATH}/bin
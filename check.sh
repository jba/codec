#!/usr/bin/env bash
# Copyright 2021 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

EXIT_CODE=0

# Support ** in globs.
shopt -s globstar

main() {
  verify_header **/*.go

  go vet -all ./... | warnout

  ensure_go_binary honnef.co/go/tools/cmd/staticcheck
  staticcheck ./... | warnout

  ensure_go_binary github.com/client9/misspell/cmd/misspell
  misspell **/*.{go,sh,md} | warnout

  ensure_go_binary mvdan.cc/unparam
  unparam ./... | warnout

  go test ./...

  exit $EXIT_CODE
}

if [ -t 1 ] && which tput >/dev/null 2>&1; then
  RED="$(tput setaf 1)"
  GREEN="$(tput setaf 2)"
  YELLOW="$(tput setaf 3)"
  NORMAL="$(tput sgr0)"
else
  RED=""
  GREEN=""
  YELLOW=""
  NORMAL=""
fi

info() { echo -e "${GREEN}$@${NORMAL}" 1>&2; }
warn() { echo -e "${YELLOW}$@${NORMAL}" 1>&2; }
err() { echo -e "${RED}$@${NORMAL}" 1>&2; EXIT_CODE=1; }

die() {
  err $@
  exit 1
}

warnout() {
  while read line; do
    warn "$line"
  done
}

# ensure_go_binary verifies that a binary exists in $PATH corresponding to the
# given go-gettable URI. If no such binary exists, it is fetched via `go get`.
ensure_go_binary() {
  local binary=$(basename $1)
  if ! [ -x "$(command -v $binary)" ]; then
    info "Installing: $1"
    # Run in a subshell for convenience, so that we don't have to worry about
    # our PWD.
    (set -x; cd && env GO111MODULE=on go get -u $1)
  fi
}

# verify_header checks that all given files contain the standard header for Go
# projects.
verify_header() {
  if [[ "$@" != "" ]]; then
    for FILE in $@
    do
        # Allow for the copyright header to start on either of the first two
        # lines, to accommodate conventions for CSS and HTML.
        line="$(head -3 $FILE)"
        if [[ ! $line == *"The Go Authors. All rights reserved."* && ! $line == *"Code generated by"* ]]; then
              err "missing license header: $FILE"
        fi
    done
  fi
}


main $@
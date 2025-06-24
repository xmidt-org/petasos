#!/usr/bin/env sh
## SPDX-FileCopyrightText: 2016 Comcast Cable Communications Management, LLC
## SPDX-License-Identifier: Apache-2.0
set -e

# check arguments for an option that would cause /petasos to stop
# return true if there is one
_want_help() {
    local arg
    for arg; do
        case "$arg" in
            -'?'|--help|-v)
                return 0
                ;;
        esac
    done
    return 1
}

_main() {
    # if command starts with an option, prepend petasos
    if [ "${1:0:1}" = '-' ]; then
        set -- /petasos "$@"
    fi

    # skip setup if they aren't running /petasos or want an option that stops /petasos
    if [ "$1" = '/petasos' ] && ! _want_help "$@"; then
        echo "Entrypoint script for petasos Server ${VERSION} started."

        if [ ! -s /etc/petasos/petasos.yaml ]; then
            echo "Building out template for file"
            /bin/spruce merge /tmp/petasos_spruce.yaml > /etc/petasos/petasos.yaml
        fi
    fi

    exec "$@"
}

_main "$@"

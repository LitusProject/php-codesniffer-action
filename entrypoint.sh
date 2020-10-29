#! /bin/sh

cd "$GITHUB_WORKSPACE"

phpcs \
    -q \
    --basepath="$GITHUB_WORKSPACE" \
    --report=json \
    --runtime-set ignore_errors_on_exit true \
    --runtime-set ignore_warnings_on_exit true \
    | php-codesniffer-action

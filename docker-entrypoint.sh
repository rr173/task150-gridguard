#!/bin/sh
set -eu

if [ "$#" -eq 0 ] || [ "$1" = "bash" ]; then
  exec bash "$@"
fi

exec /app/gridguard "$@"

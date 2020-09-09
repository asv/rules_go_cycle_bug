#! /usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

readonly TOPLEVEL="$(git rev-parse --show-toplevel)"
cd "${TOPLEVEL}"

if ! hash bazelisk 2>/dev/null ; then
  echo '*** Ошибка: Не установлен bazelisk.' 1>&2
  echo '    Для установки скачай последнюю версию https://github.com/bazelbuild/bazelisk/releases.' 1>&2
  exit 1
fi

bazelisk run //:gazelle -- update-repos -from_file=go.mod -prune -to_macro=godeps.bzl%go_deps
bazelisk run //:gazelle

exit 0

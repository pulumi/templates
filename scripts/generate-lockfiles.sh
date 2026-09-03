#!/usr/bin/env bash

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRIPT="${REPO_ROOT}/scripts/$(basename "${BASH_SOURCE[0]}")"
CONCURRENCY="${LOCKFILE_CONCURRENCY:-6}"

REPO_TOOLING_DIRS=(generator metadata scripts tests)

node_project_dirs() {
    local excludes=()
    for dir in "${REPO_TOOLING_DIRS[@]}"; do
        excludes+=(-path "./${dir}" -prune -o)
    done

    find . "${excludes[@]}" \
        \( -name node_modules -o -name .git -o -name .jj \) -prune -o \
        -name package.json -print |
        sed 's|/package.json$||' |
        grep -v '^\.$' |
        sort
}

project_runtime() {
    local project="$1/Pulumi.yaml" runtime
    [[ -f "${project}" ]] || return 0

    runtime="$(tr -d '\r' <"${project}" | awk '
        /^runtime:[[:space:]]*[^[:space:]]/ {
            sub(/^runtime:[[:space:]]*/, "")
            sub(/[[:space:]]+$/, "")
            print
            exit
        }
        /^runtime:[[:space:]]*$/ { in_runtime_block = 1; next }
        in_runtime_block && /^[^[:space:]]/ { exit }
        in_runtime_block && /^[[:space:]]+name:[[:space:]]*/ {
            sub(/^[[:space:]]+name:[[:space:]]*/, "")
            sub(/[[:space:]]+$/, "")
            print
            exit
        }
    ')"

    runtime="${runtime#\"}"
    runtime="${runtime%\"}"
    runtime="${runtime#\'}"
    runtime="${runtime%\'}"
    printf '%s\n' "${runtime}"
}

declares_bun_runtime() {
    [[ "$(project_runtime "$1")" == "bun" ]]
}

if [[ "${1:-}" == "--generate-one" ]]; then
    dir="${REPO_ROOT}/$2"
    cd "${dir}"
    if declares_bun_runtime "${dir}"; then
        rm -f package-lock.json
        bun install --lockfile-only
    else
        rm -f bun.lock bun.lockb
        npm install --package-lock-only --no-audit --no-fund --loglevel=error
    fi
    exit 0
fi

cd "${REPO_ROOT}"

for tool in npm bun; do
    if ! command -v "${tool}" >/dev/null 2>&1; then
        echo "${tool} is required but was not found on \$PATH" >&2
        exit 1
    fi
done

dirs_file="$(mktemp)"
trap 'rm -f "${dirs_file}"' EXIT
node_project_dirs >"${dirs_file}"
count=$(wc -l <"${dirs_file}" | tr -d ' ')

if [[ "${count}" -eq 0 ]]; then
    echo "found no Node projects, which means the discovery logic is broken" >&2
    exit 1
fi

echo "Regenerating lockfiles for ${count} Node projects..."
xargs -P "${CONCURRENCY}" -I{} "${SCRIPT}" --generate-one {} <"${dirs_file}"

orphaned_lockfiles() {
    find . \( -name node_modules -o -name .git -o -name .jj \) -prune -o \
        \( -name package-lock.json -o -name bun.lock \) -print |
        while read -r lock; do
            [[ -f "$(dirname "${lock}")/package.json" ]] || echo "${lock}"
        done
}

orphans="$(orphaned_lockfiles)"
if [[ -n "${orphans}" ]]; then
    echo "found lockfiles with no package.json beside them:" >&2
    echo "${orphans}" >&2
    exit 1
fi

stray_lockfiles() {
    local dir name
    while read -r dir; do
        for name in yarn.lock pnpm-lock.yaml bun.lockb; do
            if [[ -f "${dir}/${name}" ]]; then
                echo "${dir}/${name}"
            fi
        done
    done <"${dirs_file}"
}

strays="$(stray_lockfiles)"
if [[ -n "${strays}" ]]; then
    echo "found lockfiles belonging to a package manager these templates do not declare:" >&2
    echo "${strays}" >&2
    exit 1
fi

echo "Done."

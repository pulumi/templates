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

declares_bun_runtime() {
    grep -qE '^runtime:[[:space:]]*bun[[:space:]]*$' "$1/Pulumi.yaml" 2>/dev/null
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

echo "Done."

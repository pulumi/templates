#!/usr/bin/env bash

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BASE_REF="${1:?usage: audit-lockfiles.sh --base <ref>}"
if [[ "${BASE_REF}" == "--base" ]]; then
    BASE_REF="${2:?--base needs a git ref}"
fi

cd "${REPO_ROOT}"

REPO_TOOLING_DIRS=(generator metadata scripts tests)

# Emits "<package manager>\t<directory>" so each project is audited by the tool that wrote its
# lockfile. npm cannot read bun.lock and bun cannot read package-lock.json.
lockfile_dirs() {
    local root="$1" excludes=() dir
    for dir in "${REPO_TOOLING_DIRS[@]}"; do
        excludes+=(-path "${root}/${dir}" -prune -o)
    done

    find "${root}" "${excludes[@]}" \
        \( -name node_modules -o -name .git -o -name .jj \) -prune -o \
        \( -name package-lock.json -o -name bun.lock \) -print |
        while read -r lock; do
            case "$(basename "${lock}")" in
            package-lock.json) printf 'npm\t%s\n' "$(dirname "${lock}")" ;;
            bun.lock) printf 'bun\t%s\n' "$(dirname "${lock}")" ;;
            esac
        done |
        sort
}

npm_advisories() {
    local dir="$1" label="$2" report
    report="$( (cd "${dir}" && npm audit --package-lock-only --json 2>/dev/null) || true)"

    if ! jq -e 'has("metadata")' >/dev/null 2>&1 <<<"${report}"; then
        echo "npm audit produced no usable report for ${dir}, refusing to report it as clean" >&2
        return 1
    fi

    jq -r --arg dir "${label}" '
        (.vulnerabilities // {}) | to_entries[] |
        .value.via[]? | select(type == "object") |
        "\($dir)\t\(.url // .source)\t\(.severity)\t\(.title)"
    ' <<<"${report}"
}

bun_advisories() {
    local dir="$1" label="$2" report
    report="$( (cd "${dir}" && bun audit --json 2>/dev/null) || true)"

    if ! jq -e 'type == "object"' >/dev/null 2>&1 <<<"${report}"; then
        echo "bun audit produced no usable report for ${dir}, refusing to report it as clean" >&2
        return 1
    fi

    jq -r --arg dir "${label}" '
        to_entries[] | .value[] |
        "\($dir)\t\(.url)\t\(.severity)\t\(.title)"
    ' <<<"${report}"
}

advisories_under() {
    local root="$1" manager dir
    while IFS=$'\t' read -r manager dir; do
        [[ -n "${dir}" ]] || continue
        case "${manager}" in
        npm) npm_advisories "${dir}" "${dir#"${root}"/}" ;;
        bun) bun_advisories "${dir}" "${dir#"${root}"/}" ;;
        esac
    done < <(lockfile_dirs "${root}") | sort -u
}

materialize_ref() {
    local ref="$1" dest="$2" lock dir
    for lock in $(git ls-tree -r --name-only "${ref}" | grep -E '(^|/)(package-lock\.json|bun\.lock)$' || true); do
        dir="$(dirname "${lock}")"
        mkdir -p "${dest}/${dir}"
        git show "${ref}:${lock}" >"${dest}/${lock}"
        git show "${ref}:${dir}/package.json" >"${dest}/${dir}/package.json" 2>/dev/null || true
    done
}

ids_of() {
    awk -F'\t' 'NF{print $1"\t"$2}' <<<"$1" | sort -u
}

found="$(advisories_under .)"

base_dir="$(mktemp -d)"
trap 'rm -rf "${base_dir}"' EXIT
materialize_ref "${BASE_REF}" "${base_dir}"
base_ids="$(ids_of "$(advisories_under "${base_dir}")")"

introduced="$(comm -23 <(ids_of "${found}") <(echo "${base_ids}"))"

if [[ -z "${introduced}" ]]; then
    echo "No vulnerabilities introduced against ${BASE_REF}."
    exit 0
fi

echo "::error::This change introduces vulnerabilities that are not in ${BASE_REF}" >&2
while IFS= read -r id; do
    [[ -n "${id}" ]] || continue
    grep -F "${id}" <<<"${found}" >&2 || echo "${id}" >&2
done <<<"${introduced}"
exit 1

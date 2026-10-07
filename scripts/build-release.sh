#!/bin/sh
set -eu

version="${1:-}"
output="${2:-}"
reference="${3:-HEAD}"

printf '%s\n' "$version" | grep -Eq \
	'^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)([-+][0-9A-Za-z.-]+)?$' || {
	printf 'release version must be a semantic v-prefixed version\n' >&2
	exit 1
}
case "$output" in
	''|/|.|..)
		printf 'release output directory is unsafe\n' >&2
		exit 1
		;;
esac

repository="$(git rev-parse --show-toplevel)"
prefix="$(git rev-parse --show-prefix)"
module_path="${prefix%/}"
commit="$(git -C "$repository" rev-parse "${reference}^{commit}")"
temporary="$(mktemp -d)"
trap 'rm -r "$temporary"' EXIT HUP INT TERM

mkdir -p "$output"
archive="prompts-${version}.tar.gz"
sbom="prompts-${version}.sbom.json"
test ! -e "$output/$archive" && test ! -e "$output/$sbom" || {
	printf 'release artifact already exists\n' >&2
	exit 1
}

tree="$commit"
if [ -n "$module_path" ]; then
	tree="${commit}:${module_path}"
fi
archive_time="$(git -C "$repository" show -s --format=%cI "$commit")"
git -C "$repository" archive --format=tar \
	--mtime="$archive_time" --prefix="prompts-${version}/" \
	"$tree" > "$temporary/source.tar"
go run ./scripts/rewrite-archive.go "$temporary/source.tar" \
	"$output/$archive" "prompts-${version}" "prompts-${version}"

GOWORK=off go run \
	github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@v1.10.0 \
	mod -json -licenses -type library -noserial -notimestamp \
	-output "$output/$sbom" .

(
	cd "$output"
	shasum -a 256 "$archive" "$sbom" > SHA256SUMS
)
gzip -t "$output/$archive"

#!/bin/sh
set -eu

output="$(mktemp -d)"
trap 'rm -rf "$output"' EXIT HUP INT TERM

for engine in go-prompts huh survey promptui bubbles; do
	name="$engine"
	if test "$engine" = go-prompts; then
		name="prompts"
	fi
	CGO_ENABLED=0 GOWORK=off go build -trimpath \
		-ldflags '-s -w -buildid=' -o "$output/$name" "./cmd/$engine"
	bytes="$(wc -c < "$output/$name" | tr -d ' ')"
	printf '%s\t%s\n' "$name" "$bytes"
	if test "$name" = prompts && test "$bytes" -gt 2500000; then
		printf 'prompts binary exceeds 2500000-byte budget\n' >&2
		exit 1
	fi
done

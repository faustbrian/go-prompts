# Security

Report vulnerabilities through
[GitHub private vulnerability reporting](https://github.com/faustbrian/go-prompts/security/advisories/new).
Do not include live credentials or secret prompt values in an issue, test
fixture, terminal capture, or proof of concept.

The published `v1.0.0` release begins the supported stable v1 line. The latest
stable v1 minor line receives security fixes. Security advisories describe
affected versions, safe upgrade versions, and whether secret disclosure or
terminal-state restoration is involved.

The package redacts its secret wrapper representations and safe errors. It
cannot control copies made after `Reveal`, application logging, clipboard
software, terminal emulators, swap, crash dumps, or Go string memory. See
`docs/security.md` and `docs/secrets.md` before handling credentials.

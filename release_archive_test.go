package prompts

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// The archive path uses real Git objects and the real archive writer. Only the
// unrelated external SBOM generator is replaced; release builds verify it live.
func TestReleaseArchiveUsesFrozenVersionedModuleSource(t *testing.T) {
	for _, moduleDir := range []string{"", "modules/nested"} {
		t.Run("QA-"+map[bool]string{true: "ROOT", false: "NESTED"}[moduleDir == ""], func(t *testing.T) {
			repository := t.TempDir()
			module := filepath.Join(repository, moduleDir)
			files := map[string]string{
				"go.mod":     "module example.com/release-fixture\n\ngo 1.27.0\n",
				"marker.txt": "frozen source\n",
			}
			for _, name := range []string{"build-release.sh", "rewrite-archive.go"} {
				data, err := os.ReadFile(filepath.Join("scripts", name))
				if err != nil {
					t.Fatal(err)
				}
				files["scripts/"+name] = string(data)
			}
			for name, data := range files {
				writeReleaseFixture(t, filepath.Join(module, name), data, 0o644)
			}
			goBinary, err := exec.LookPath("go")
			if err != nil {
				t.Fatal(err)
			}
			bin := t.TempDir()
			// Forward the owned Go writer to the actual toolchain. The provider
			// double writes an explicit fixture SBOM, not an archive or checksum.
			writeReleaseFixture(t, filepath.Join(bin, "go"), `#!/bin/sh
if [ "$1" = run ] && [ "$2" = github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@v1.10.0 ]; then
  while [ "$#" -gt 0 ]; do
    if [ "$1" = -output ]; then
      shift
      printf '%s\n' '{"bomFormat":"CycloneDX","metadata":{"component":{"name":"example.com/release-fixture"}}}' > "$1"
      exit 0
    fi
    shift
  done
  exit 2
fi
exec "$RELEASE_REAL_GO" "$@"
`, 0o755)
			env := append(os.Environ(), "GIT_CONFIG_GLOBAL="+os.DevNull, "GIT_CONFIG_NOSYSTEM=1",
				"GIT_AUTHOR_DATE=2000-01-01T00:00:00Z", "GIT_COMMITTER_DATE=2000-01-01T00:00:00Z",
				"GIT_AUTHOR_NAME=Release Fixture", "GIT_AUTHOR_EMAIL=fixture@example.com",
				"GIT_COMMITTER_NAME=Release Fixture", "GIT_COMMITTER_EMAIL=fixture@example.com",
				"PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"), "RELEASE_REAL_GO="+goBinary)
			git := func(args ...string) string {
				cmd := exec.Command("git", args...)
				cmd.Dir, cmd.Env = repository, env
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("git %v: %v\n%s", args, err, out)
				}
				return strings.TrimSpace(string(out))
			}
			git("init", "-q", "-b", "master")
			var paths []string
			for name := range files {
				paths = append(paths, filepath.Join(moduleDir, name))
			}
			if moduleDir != "" {
				for _, name := range []string{"root-marker.txt", "modules/sibling/marker.txt"} {
					writeReleaseFixture(t, filepath.Join(repository, name), "not selected\n", 0o644)
					paths = append(paths, name)
				}
			}
			git(append([]string{"add", "--"}, paths...)...)
			frozen := git("commit-tree", git("write-tree"), "-m", "test: freeze release fixture")
			git("update-ref", "refs/heads/master", frozen)
			writeReleaseFixture(t, filepath.Join(module, "marker.txt"), "later source\n", 0o644)
			git("add", "--", filepath.Join(moduleDir, "marker.txt"))
			later := git("commit-tree", git("write-tree"), "-p", frozen, "-m", "test: advance release fixture")
			git("update-ref", "refs/heads/master", later, frozen)
			writeReleaseFixture(t, filepath.Join(module, "marker.txt"), "dirty source\n", 0o644)
			writeReleaseFixture(t, filepath.Join(module, "untracked.txt"), "untracked source\n", 0o644)
			before := git("status", "--porcelain")
			build := func(output, reference string, wantFailure bool) {
				args := []string{"scripts/build-release.sh", "v1.1.3", output}
				if reference != "" {
					args = append(args, reference)
				}
				cmd := exec.Command("sh", args...)
				cmd.Dir, cmd.Env = module, env
				out, err := cmd.CombinedOutput()
				if (err != nil) != wantFailure {
					t.Fatalf("release build: %v\n%s", err, out)
				}
			}
			first, second := t.TempDir(), t.TempDir()
			build(first, frozen, false)
			build(second, frozen, false)
			archive := "prompts-v1.1.3.tar.gz"
			contents, bytes := readReleaseArchive(t, filepath.Join(first, archive))
			if !reflect.DeepEqual(contents, files) {
				t.Fatalf("frozen module export mismatch: got paths %v", reflect.ValueOf(contents).MapKeys())
			}
			_, repeated := readReleaseArchive(t, filepath.Join(second, archive))
			if string(bytes) != string(repeated) {
				t.Fatal("identical frozen source produced different archive bytes")
			}
			checkReleaseBundle(t, first)
			build(first, frozen, true) // QA-GUARDS: existing output is refused.
			_, unchanged := readReleaseArchive(t, filepath.Join(first, archive))
			if string(bytes) != string(unchanged) {
				t.Fatal("existing source artifact overwritten")
			}
			checkReleaseBundle(t, first)
			defaultOutput := t.TempDir()
			build(defaultOutput, "", false)
			current, _ := readReleaseArchive(t, filepath.Join(defaultOutput, archive))
			if current["marker.txt"] != "later source\n" {
				t.Fatal("default reference did not export HEAD")
			}
			if git("status", "--porcelain") != before || git("rev-parse", "HEAD") != later {
				t.Fatal("release build changed checkout, index, or refs")
			}
		})
	}
}

func writeReleaseFixture(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}

func readReleaseArchive(t *testing.T, path string) (map[string]string, []byte) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	z, err := gzip.NewReader(strings.NewReader(string(data)))
	if err != nil {
		t.Fatal(err)
	}
	defer z.Close()
	reader := tar.NewReader(z)
	files := make(map[string]string)
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if header.Typeflag == tar.TypeXGlobalHeader {
			continue // Git commit metadata is not an extracted member.
		}
		if header.ModTime.Unix() != 946684800 {
			t.Fatalf("archive member uses non-commit timestamp: %q at %v", header.Name, header.ModTime)
		}
		if !strings.HasPrefix(header.Name, "prompts-v1.1.3/") {
			t.Fatalf("archive member outside versioned root: %q", header.Name)
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		name := strings.TrimPrefix(header.Name, "prompts-v1.1.3/")
		if _, exists := files[name]; exists {
			t.Fatalf("duplicate member %q", name)
		}
		body, err := io.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}
		files[name] = string(body)
	}
	return files, data
}

func checkReleaseBundle(t *testing.T, output string) {
	t.Helper()
	sums, err := os.ReadFile(filepath.Join(output, "SHA256SUMS"))
	if err != nil {
		t.Fatal(err)
	}
	var expected strings.Builder
	for _, name := range []string{"prompts-v1.1.3.tar.gz", "prompts-v1.1.3.sbom.json"} {
		data, err := os.ReadFile(filepath.Join(output, name))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&expected, "%x  %s\n", sha256.Sum256(data), name)
		if strings.HasSuffix(name, ".json") {
			var bom struct {
				BOMFormat string `json:"bomFormat"`
				Metadata  struct{ Component struct{ Name string } }
			}
			if err := json.Unmarshal(data, &bom); err != nil {
				t.Fatal(err)
			}
			if bom.BOMFormat != "CycloneDX" || bom.Metadata.Component.Name != "example.com/release-fixture" {
				t.Fatal("wrong fixture SBOM identity")
			}
		}
	}
	if string(sums) != expected.String() {
		t.Fatal("checksum manifest does not bind both bundle assets")
	}
}

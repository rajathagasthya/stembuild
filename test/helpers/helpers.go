package helpers

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

func recursiveFileList(destDir, searchDir string) ([]string, []string, []string, error) {
	srcFileList := make([]string, 0)
	destFileList := make([]string, 0)
	dirList := make([]string, 0)
	leafSearchDir := searchDir
	lastSepIndex := strings.LastIndex(searchDir, string(filepath.Separator))
	if lastSepIndex >= 0 {
		leafSearchDir = searchDir[lastSepIndex:len(searchDir)]
	}

	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			dirList = append(dirList, filepath.Join(destDir, leafSearchDir, path[len(searchDir):len(path)]))
		} else {
			srcFileList = append(srcFileList, path)
			destFileList = append(destFileList, filepath.Join(destDir, leafSearchDir, path[len(searchDir):len(path)]))
		}
		return err
	})

	if e != nil {
		return nil, nil, nil, e
	}

	return destFileList, srcFileList, dirList, nil
}

func CopyRecursive(destRoot, srcRoot string) error {
	var err error
	destRoot, err = filepath.Abs(destRoot)
	if err != nil {
		return err
	}

	srcRoot, err = filepath.Abs(srcRoot)
	if err != nil {
		return err
	}

	destFileList, srcFileList, dirList, err := recursiveFileList(destRoot, srcRoot)
	if err != nil {
		return err
	}

	// create destination directory hierarchy
	for _, myDir := range dirList {
		if err = os.MkdirAll(myDir, os.ModePerm); err != nil {
			return err
		}
	}

	for i, _ := range srcFileList {
		srcFile, err := os.Open(srcFileList[i])
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destFileList[i])
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}

		if err = destFile.Sync(); err != nil {
			return err
		}
	}

	return nil
}

func extractArchive(archive io.Reader, dirname string) error {
	tr := tar.NewReader(archive)

	limit := 100
	for ; limit >= 0; limit-- {
		h, err := tr.Next()
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("tar: reading from archive: %s", err)
			}
			break
		}

		// expect a flat archive
		name := h.Name
		if filepath.Base(name) != name {
			return fmt.Errorf("tar: archive contains subdirectory: %s", name)
		}

		// only allow regular files
		mode := h.FileInfo().Mode()
		if !mode.IsRegular() {
			return fmt.Errorf("tar: unexpected file mode (%s): %s", name, mode)
		}

		path := filepath.Join(dirname, name)
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
		if err != nil {
			return fmt.Errorf("tar: opening file (%s): %s", path, err)
		}
		defer f.Close()

		if _, err := io.Copy(f, tr); err != nil {
			return fmt.Errorf("tar: writing file (%s): %s", path, err)
		}
	}
	if limit <= 0 {
		return errors.New("tar: too many files in archive")
	}
	return nil
}

// ExtractGzipArchive extracts the tgz archive name to a temp directory
// returning the filepath of the temp directory.
func ExtractGzipArchive(name string) (string, error) {
	tmpdir, err := os.MkdirTemp("", "test-")
	if err != nil {
		return "", err
	}

	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	if err := extractArchive(w, tmpdir); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return tmpdir, nil
}

func ReadFile(name string) (string, error) {
	b, err := os.ReadFile(name)
	return string(b), err
}

func BuildStembuild(version string) (string, error) {
	command := exec.Command("make", "build-integration")
	command.Env = AddOrReplaceEnvironment(os.Environ(), "STEMBUILD_VERSION", version)

	_, b, _, _ := runtime.Caller(0)
	root := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	command.Dir = root

	session, err := Start(
		command,
		NewPrefixedWriter(DebugOutPrefix, GinkgoWriter),
		NewPrefixedWriter(DebugErrPrefix, GinkgoWriter))
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 120*time.Second).Should(Exit(0))

	files, err := os.ReadDir(filepath.Join(root, "out"))
	Expect(err).NotTo(HaveOccurred())

	for _, f := range files {
		if strings.Contains(filepath.Base(f.Name()), "stembuild") {
			stem := filepath.Join(root, "out", f.Name())
			By(fmt.Sprintf("Stembuild: %s", stem))
			return stem, nil
		}
	}

	panic("Unable to find binary generated by 'make build'")
}

func EnvMustExist(variableName string) string {
	result := os.Getenv(variableName)
	if result == "" {
		Fail(fmt.Sprintf("%s must be set", variableName))
	}

	return result
}

package fileutil

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBinaryFile(t *testing.T) {
	assert.True(t, IsBinaryFile("./testfiles/test.png"), "test.png")
	assert.True(t, IsBinaryFile("./testfiles/test"), "test")
	assert.True(t, IsBinaryFile("./testfiles/test.blnk"), "test.blnk")

	assert.False(t, IsBinaryFile("./testfiles/test.dlnk"), "test.dlnk")
	assert.False(t, IsBinaryFile("./testfiles/test.lnk"), "test.lnk")
	assert.False(t, IsBinaryFile("./testfiles/test.md"), "test.md")
	assert.False(t, IsBinaryFile("./testfiles/test.php"), "test.php")
	assert.False(t, IsBinaryFile("./testfiles/test.txt"), "test.txt")
}

func TestIsSymlink(t *testing.T) {
	assert.True(t, IsSymlink("./testfiles/test.blnk"), "test.blnk")
	assert.True(t, IsSymlink("./testfiles/test.lnk"), "test.lnk")

	assert.False(t, IsSymlink("./testfiles/test.md"), "test.md")
	assert.False(t, IsSymlink("./testfiles/test.php"), "test.php")
	assert.False(t, IsSymlink("./testfiles/test.png"), "test.png")
	assert.False(t, IsSymlink("./testfiles/test.txt"), "test.txt")
	assert.False(t, IsSymlink("./testfiles/test"), "test")
}

func TestIsDir(t *testing.T) {
	assert.True(t, IsDir("./testfiles/"), "testfiles")
	assert.True(t, IsDir("./testfiles/test.dlnk"), "test.dlnk")

	assert.False(t, IsDir("./testfiles/test.blnk"), "test.blnk")
	assert.False(t, IsDir("./testfiles/test.lnk"), "test.lnk")
	assert.False(t, IsDir("./testfiles/test"), "test")
	assert.False(t, IsDir("./testfiles/test.png"), "test.png")
	assert.False(t, IsDir("./testfiles/test.md"), "test.md")
	assert.False(t, IsDir("./testfiles/test.php"), "test.php")
	assert.False(t, IsDir("./testfiles/test.txt"), "test.txt")
}

func TestDirToArray(t *testing.T) {
	expected := []string{
		"./testfiles/test",
		"./testfiles/test.blnk",
		"./testfiles/test.php",
		"./testfiles/testdir/test.nested",
		"./testfiles/test.dlnk/test.nested",
		"./testfiles/test.lnk",
		"./testfiles/test.md",
		"./testfiles/test.png",
		"./testfiles/test.txt",
	}
	files := DirToArray("./testfiles/", true, func(a, b string) bool { return true }, func(a, b string) bool { return true })

	sort.Strings(files)
	sort.Strings(expected)

	assert.Equal(t, expected, files)
}

func TestDefaultFileFilter(t *testing.T) {
	expected := []string{
		"./testfiles/test.php",
		"./testfiles/testdir/test.nested",
		"./testfiles/test.dlnk/test.nested",
		"./testfiles/test.md",
		"./testfiles/test.txt",
	}
	files := DirToArray("./testfiles/", true, DefaultFileFilter, func(a, b string) bool { return true })

	sort.Strings(files)
	sort.Strings(expected)

	assert.Equal(t, expected, files)
}

func TestFilterOutBinaries(t *testing.T) {
	expected := []string{
		"./testfiles/test.php",
		"./testfiles/testdir/test.nested",
		"./testfiles/test.dlnk/test.nested",
		"./testfiles/test.lnk",
		"./testfiles/test.md",
		"./testfiles/test.txt",
	}
	files := DirToArray("./testfiles/", true, func(a, b string) bool { return true }, func(a, b string) bool { return true })
	files = FilterOutBinaries(files)

	sort.Strings(files)
	sort.Strings(expected)

	assert.Equal(t, expected, files)
}

func TestDirToArrayNoFollow(t *testing.T) {
	expected := []string{
		"./testfiles/test",
		"./testfiles/test.php",
		"./testfiles/testdir/test.nested",
		"./testfiles/test.md",
		"./testfiles/test.png",
		"./testfiles/test.txt",
	}
	files := DirToArray("./testfiles/", false, func(a, b string) bool { return true }, func(a, b string) bool { return true })

	sort.Strings(files)
	sort.Strings(expected)

	assert.Equal(t, expected, files)
}

func TestFilterOutSymlinks(t *testing.T) {
	expected := []string{
		"./testfiles/test",
		"./testfiles/test.php",
		"./testfiles/testdir/test.nested",
		"./testfiles/test.dlnk/test.nested",
		"./testfiles/test.md",
		"./testfiles/test.png",
		"./testfiles/test.txt",
	}
	files := DirToArray("./testfiles/", true, func(a, b string) bool { return true }, func(a, b string) bool { return true })
	files = FilterOutSymlinks(files)

	sort.Strings(files)
	sort.Strings(expected)

	assert.Equal(t, expected, files)
}

func TestFilterExtWhitelist(t *testing.T) {
	expected := []string{
		"./testfiles/test.php",
		"./testfiles/test.png",
	}
	files := DirToArray("./testfiles/", true, func(a, b string) bool { return true }, func(a, b string) bool { return true })
	files = FilterExtWhitelist([]string{"php", "png"}, files)

	sort.Strings(files)
	sort.Strings(expected)

	assert.Equal(t, expected, files)
}

func TestFilterExtBlacklist(t *testing.T) {
	expected := []string{
		"./testfiles/test",
		"./testfiles/test.blnk",
		"./testfiles/testdir/test.nested",
		"./testfiles/test.dlnk/test.nested",
		"./testfiles/test.lnk",
		"./testfiles/test.md",
		"./testfiles/test.txt",
	}
	files := DirToArray("./testfiles/", true, func(a, b string) bool { return true }, func(a, b string) bool { return true })
	files = FilterExtBlacklist([]string{"php", "png"}, files)

	sort.Strings(files)
	sort.Strings(expected)

	assert.Equal(t, expected, files)
}

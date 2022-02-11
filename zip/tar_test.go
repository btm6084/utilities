package zip

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func getTarGZFromDisk(url string) (*http.Response, error) {
	file, err := os.Open("./TestData.tar.gz")
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	resp := &http.Response{
		ContentLength: stat.Size(),
		Status:        "OK",
		StatusCode:    http.StatusOK,
		Body:          file,
	}

	resp.Header = make(http.Header)
	resp.Header.Set("Content-Type", "application/gzip")

	return resp, nil
}

func TestGetTarFileFromDisk(t *testing.T) {
	r, err := GetTarGZFileFromDisk("./TestData.tar.gz")
	require.Nil(t, err)

	b, err := GetFileFromTarGZ(r, "readme.md")
	require.Nil(t, err)

	require.Equal(t, `This is just a simple file to test the zip/tar package.`, string(b))
}

func TestGetTarFileFromHTTP(t *testing.T) {
	r, err := DownloadTarGZFile(getTarGZFromDisk, "https://download.example.com/my/resource")
	require.Nil(t, err)

	b, err := GetFileFromTarGZ(r, "readme.md")
	require.Nil(t, err)

	require.Equal(t, `This is just a simple file to test the zip/tar package.`, string(b))
}

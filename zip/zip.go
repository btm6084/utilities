package zip

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var (
	ErrNotFound = errors.New("not found")
)

const (
	bodySampleSize = 45
)

// HTTPGetFunc is used to determine how remote URLs are downloaded.
// http.Get satisfies this type.
type HTTPGetFunc func(string) (*http.Response, error)

// DownloadZipFile downloads a zipfile from the given url and returns a zip.Reader
// The getFn matches the signature of http.Get
func DownloadZipFile(getFn HTTPGetFunc, url string) (*zip.Reader, error) {
	data, err := getFn(url)
	if err != nil {
		return nil, err
	}

	defer data.Body.Close()

	fileSize := data.ContentLength
	contentType := data.Header.Get("Content-Type")
	bodyData, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}

	if fileSize == -1 {
		fileSize = int64(len(bodyData))
	}

	if data.StatusCode != 200 || fileSize <= 0 || contentType != "application/zip" {
		// We just want a small sample of the body.
		comment := string(bodyData)[0:bodySampleSize]
		log.WithFields(log.Fields{
			"url":         url,
			"status":      data.Status,
			"contentType": contentType,
			"comment":     comment,
		}).Error("error downloading zipfile")
		return nil, errors.New("error downloading zipfile")
	}

	reader := bytes.NewReader(bodyData)

	zipFile, err := zip.NewReader(reader, fileSize)
	if err != nil {
		return nil, err
	}

	return zipFile, nil
}

// GetZipFileFromDisk opens a zipfile from the filesystem and returns a zip.Reader
func GetZipFileFromDisk(fileName string) (*zip.Reader, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	bodyData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(bodyData)

	zipFile, err := zip.NewReader(reader, fileSize)
	if err != nil {
		return nil, err
	}

	return zipFile, nil
}

func GetFileFromZip(r *zip.Reader, fileName string) ([]byte, error) {
	for _, file := range r.File {
		_, name := filepath.Split(file.Name)

		if name != fileName {
			continue // directory
		}

		f, err := file.Open()
		if err != nil {
			return nil, err
		}

		defer f.Close()

		return io.ReadAll(f)
	}

	return nil, ErrNotFound
}

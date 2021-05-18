package zip

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	bodySampleSize = 45
)

// DownloadZipFile downloads a zipfile from the given url and returns a zip.Reader
func DownloadZipFile(url string) (*zip.Reader, error) {
	http.DefaultClient.Timeout = 10 * time.Second
	data, err := http.Get(url)
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

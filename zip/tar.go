package zip

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// DownloadTarGZFile downloads a .tar.gz file from the given url and returns a zip.Reader
// The getFn matches the signature of http.Get
func DownloadTarGZFile(getFn HTTPGetFunc, url string) (*tar.Reader, error) {
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

	if data.StatusCode != 200 || fileSize <= 0 || contentType != "application/gzip" {
		// We just want a small sample of the body.
		comment := string(bodyData)[0:bodySampleSize]
		log.WithFields(log.Fields{
			"url":         url,
			"status":      data.Status,
			"contentType": contentType,
			"comment":     comment,
		}).Error("error downloading tar file")
		return nil, errors.New("error downloading tar file")
	}

	gzr, err := gzip.NewReader(bytes.NewReader(bodyData))
	if err != nil {
		return nil, err
	}

	return tar.NewReader(gzr), nil
}

// GetTarGZFileFromDisk opens a zipfile from the filesystem and returns a zip.Reader
func GetTarGZFileFromDisk(fileName string) (*tar.Reader, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	gzb, err := io.ReadAll(gzr)
	if err != nil {
		return nil, err
	}

	return tar.NewReader(bytes.NewReader(gzb)), nil
}

func GetFileFromTarGZ(r *tar.Reader, fileName string) ([]byte, error) {
	for {
		header, err := r.Next()

		if err == io.EOF {
			return nil, ErrNotFound
		}

		if err != nil {
			return nil, err
		}

		if header.Typeflag != tar.TypeReg || !strings.HasSuffix(header.Name, fileName) {
			continue
		}

		b := make([]byte, header.Size)
		_, err = io.ReadFull(r, b)
		if err != nil && err != io.EOF {
			return nil, err
		}

		return b, nil
	}
}

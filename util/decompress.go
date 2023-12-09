package util

import (
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func DecompressOpen(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(filePath, ".xml.gz") {
		return unzipGzip(file)
	} else if strings.HasSuffix(filePath, ".zip") {
		return unzipZip(file)
	} else if strings.HasSuffix(filePath, ".xml") {
		return file, nil
	}

	return nil, fmt.Errorf("Won't parse '%s'. It is neither of: .xml, .xml.gz, .zip", filePath)
}

func unzipGzip(gzipFile *os.File) (*os.File, error) {
	reader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	tempFile, err := createTempFile()
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(tempFile, reader)
	if err != nil {
		tempFile.Close()
		return nil, err
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		tempFile.Close()
		return nil, err
	}

	return tempFile, nil
}

func unzipZip(zipFile *os.File) (*os.File, error) {
	reader, err := zip.OpenReader(zipFile.Name())
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	if len(reader.File) != 1 {
		return nil, errors.New("zip file must contain exactly one XML file")
	}

	xmlFile, err := reader.File[0].Open()
	if err != nil {
		return nil, err
	}

	tempFile, err := createTempFile()
	if err != nil {
		xmlFile.Close()
		return nil, err
	}

	_, err = io.Copy(tempFile, xmlFile)
	if err != nil {
		tempFile.Close()
		xmlFile.Close()
		return nil, err
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		tempFile.Close()
		xmlFile.Close()
		return nil, err
	}

	return tempFile, nil
}

func createTempFile() (*os.File, error) {
	return os.CreateTemp("", "dmarc-report-*.xml")
}

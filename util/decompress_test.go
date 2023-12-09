package util

import (
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"testing"
)

func createTempFileWithContent(content string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "dmarc-report-testfile*.xml")
	if err != nil {
		return nil, err
	}

	_, err = tempFile.WriteString(content)
	if err != nil {
		return nil, err
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

func createTempGzipFileWithContent(content string) (*os.File, error) {
	tempFile, err := createTempFileWithContent(content)
	if err != nil {
		return nil, err
	}

	gzipFile, err := os.CreateTemp("", "dmarc-report-testfile*.gz")
	if err != nil {
		return nil, err
	}

	gzipWriter := gzip.NewWriter(gzipFile)
	_, err = io.Copy(gzipWriter, tempFile)
	if err != nil {
		return nil, err
	}
	gzipWriter.Close()
	tempFile.Close()
	os.Remove(tempFile.Name())

	_, err = gzipFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return gzipFile, nil
}

func TestUnzipGzip(t *testing.T) {
	content := "<foo>bar</foo>"

	gzipFile, err := createTempGzipFileWithContent(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(gzipFile.Name())
	defer gzipFile.Close()

	unzippedFile, err := unzipGzip(gzipFile)
	if err != nil {
		t.Fatal(err)
	}
	defer unzippedFile.Close()

	unzippedContent, err := io.ReadAll(unzippedFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(unzippedContent) != content {
		t.Errorf("Expected content %s, got %s", content, string(unzippedContent))
	}
}

func createTempZipFileWithContent(content string) (*os.File, error) {
	zipFile, err := os.CreateTemp("", "dmarc-report-testfile*.zip")
	if err != nil {
		return nil, err
	}

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileInZip, err := zipWriter.Create("file.txt")
	if err != nil {
		return nil, err
	}

	_, err = fileInZip.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	return zipFile, nil
}

func TestUnzipZip(t *testing.T) {
	content := "<foo>bar</foo>"

	zipFile, err := createTempZipFileWithContent(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(zipFile.Name())
	defer zipFile.Close()

	unzippedFile, err := unzipZip(zipFile)
	if err != nil {
		t.Fatal(err)
	}
	defer unzippedFile.Close()

	unzippedContent, err := io.ReadAll(unzippedFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(unzippedContent) != content {
		t.Errorf("Expected content %s, got %s", content, string(unzippedContent))
	}
}

func TestDecompressOpen_XML(t *testing.T) {
	content := "<foo>bar</foo>"

	xmlFile, err := createTempFileWithContent(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(xmlFile.Name())
	defer xmlFile.Close()

	openedFile, err := DecompressOpen(xmlFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer openedFile.Close()

	readContent, err := io.ReadAll(openedFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(readContent) != content {
		t.Errorf("Expected content %s, got %s", content, string(readContent))
	}
}

func TestDecompressOpen_Zip(t *testing.T) {
	content := "<foo>bar</foo>"

	zipFile, err := createTempZipFileWithContent(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(zipFile.Name())
	defer zipFile.Close()

	openedFile, err := DecompressOpen(zipFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer openedFile.Close()

	readContent, err := io.ReadAll(openedFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(readContent) != content {
		t.Errorf("Expected content %s, got %s", content, string(readContent))
	}
}

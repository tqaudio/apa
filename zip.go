package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func createZip(archivePath string, sourcePaths, destinationPaths []string) error {
	if len(sourcePaths) != len(destinationPaths) {
		return fmt.Errorf("mismatching source and destination paths")
	}

	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for i, sourcePath := range sourcePaths {
		destinationPath := destinationPaths[i]
		if err := appendFile(zipWriter, sourcePath, destinationPath); err != nil {
			return err
		}
	}
	return nil
}

func appendFile(zipWriter *zip.Writer, sourcePath, destinationPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = strings.Join(strings.Split(destinationPath, string(filepath.Separator)), "/")
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

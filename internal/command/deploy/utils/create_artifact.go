package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CreateArtifact(sourceDir, targetFileName string) error {
	// Create the .tar.gz file
	f, err := os.Create(targetFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a gzip writer on top of the file
	gw := gzip.NewWriter(f)
	defer gw.Close()

	// Create a tar writer on top of the gzip writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Walk through the source directory
	return filepath.Walk(sourceDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files/directories
		if strings.HasPrefix(fi.Name(), ".") {
			if fi.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Compute the relative path
		relPath, err := filepath.Rel(sourceDir, file)
		if err != nil {
			return err
		}

		// Skip the root folder itself
		if relPath == "." {
			return nil
		}

		// Create tar header
		header, err := tar.FileInfoHeader(fi, relPath)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If not a regular file, don't try to copy contents
		if !fi.Mode().IsRegular() {
			return nil
		}

		// Open and copy the file into the archive
		srcFile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		_, err = io.Copy(tw, srcFile)
		return err
	})
}

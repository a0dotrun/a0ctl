package main

import (
	"archive/zip"   // For creating zip files
	"io"            // For copying file content
	"log"           // For logging errors/info
	"os"            // For file operations
	"path/filepath" // For handling file paths
)

// main function — starting point of the program
func main() {
	// Compress current folder (".") into output.zip
	err := ZipFolder(".", "output.zip")
	if err != nil {
		log.Fatal(err) // If something goes wrong, exit with error
	}
	log.Println("Compressed current folder into output.zip")
}

// ZipFolder compresses the contents of the source folder into a target zip file
func ZipFolder(source, target string) error {
	// Create the output zip file
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close() // Make sure to close the file when we're done

	// Create a new zip writer on the file
	archive := zip.NewWriter(zipFile)
	defer archive.Close() // Close the archive once all files are written

	// Walk the entire source directory tree
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Return early if any error during walk
		}

		// Skip the output zip file itself if it's inside the source directory
		if filepath.Base(path) == target {
			return nil
		}

		// Create a zip file header based on the file/folder info
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Calculate the relative path (so we don't store absolute paths)
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		header.Name = relPath // Set relative path in zip

		// If it's a directory, add a trailing slash in the zip to mark it as a folder
		if info.IsDir() {
			header.Name += "/"
		} else {
			// Use compression for files
			header.Method = zip.Deflate
		}

		// Create a writer inside the zip for this file or folder
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		// If it's a file (not a directory), copy its contents into the zip
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Copy the file's contents into the zip writer
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		// Done processing this file or folder
		return nil
	})
}

package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
)

// unzup function
func ExtractTarGz(gzipStream io.Reader, cachDir string) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(filepath.Join(cachDir, header.Name), 0755); err != nil {
				log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(filepath.Join(cachDir, header.Name))
			if err != nil {
				log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

			err = os.Chmod(filepath.Join(cachDir, header.Name), 0755)
			if err != nil {
				log.Fatal(err)
			}

		default:
			log.Fatalf(
				"ExtractTarGz: uknown type: %b in %s",
				header.Typeflag,
				header.Name)
		}

	}
}

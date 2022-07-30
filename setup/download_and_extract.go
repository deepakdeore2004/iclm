package setup

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func Download(cacheDir, version string) {

	var localOS string
	if runtime.GOOS == "darwin" {
		localOS = "osx"
	} else {
		localOS = runtime.GOOS
	}

	file := "istio-" + version + "-" + localOS + ".tar.gz"
	link := "https://github.com/istio/istio/releases/download/" + version + "/" + file
	localFile := filepath.Join(cacheDir, file)

	if _, err := os.Stat(localFile); err == nil {
		log.Printf("%v already exists, skipping download\n", file)
	} else {
		log.Printf("downloading %v as %v\n", link, localFile)

		// create file
		out, err := os.Create(localFile)
		defer out.Close()

		// download
		resp, err := http.Get(link)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			_, err = io.Copy(out, resp.Body)
			log.Print("download successful: ", resp.StatusCode)
		} else {
			os.Remove(localFile)
			// fatal
			log.Println("download failed: ", resp.StatusCode)
		}
	}
	// unzip file
	log.Printf("Extracting %v\n", localFile)
	r, err := os.Open(localFile)
	if err != nil {
		// fatal
		log.Println(err.Error())
	} else {
		ExtractTarGz(r, cacheDir)
	}

}

func ExtractTarGz(gzipStream io.Reader, cachDir string) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		// fatal
		log.Println("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			// fatal
			log.Printf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(filepath.Join(cachDir, header.Name), 0755); err != nil {
				// fatal
				log.Printf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(filepath.Join(cachDir, header.Name))
			if err != nil {
				// fatal
				log.Printf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				// fatal
				log.Printf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

			err = os.Chmod(filepath.Join(cachDir, header.Name), 0755)
			if err != nil {
				// fatal
				log.Println(err)
			}

		default:
			log.Printf(
				"ExtractTarGz: uknown type: %b in %s",
				header.Typeflag,
				header.Name)
		}

	}
}

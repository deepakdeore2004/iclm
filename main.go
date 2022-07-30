package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	cachDir := flag.String("cacheDir", "/tmp", "Directory to download istio binaries")
	version := flag.String("version", "1.14.1", "istio version to install")
	flag.Parse()

	// Download istio if not downloaded already
	istioCtl := filepath.Join(*cachDir, "istio-"+*version, "bin", "istioctl")
	cmd := exec.Command(istioCtl, "version", "--remote=false")
	stdout, _ := cmd.Output()

	if strings.TrimSuffix(string(stdout), "\n") == *version {
		log.Printf("%v is present, skipping download\n", istioCtl)
	} else {
		fmt.Printf("%v isnt present, downloading\n", istioCtl)
		// set OS
		var localOS string
		if runtime.GOOS == "darwin" {
			localOS = "osx"
		} else {
			localOS = runtime.GOOS
		}

		file := "istio-" + *version + "-" + localOS + ".tar.gz"

		link := "https://github.com/istio/istio/releases/download/" + *version + "/" + file
		localFile := filepath.Join(*cachDir, file)

		if _, err := os.Stat(localFile); err == nil {
			log.Printf("%v already exists, skipping download\n", file)
		} else {
			download(link, localFile)
			// unzip file
			log.Printf("Extracting %v\n", localFile)
			r, err := os.Open(localFile)
			if err != nil {
				log.Fatal(err.Error())
			} else {
				ExtractTarGz(r, *cachDir)
			}
		}
	}
}

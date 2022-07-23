package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	cachDir := flag.String("cacheDir", "/tmp", "Directory to download istio binaries")
	version := flag.String("version", "1.14.1", "istio version to install")
	flag.Parse()

	// ensure cacheDir
	if _, err := os.Stat(*cachDir); os.IsNotExist(err) {
		err := os.Mkdir(*cachDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// set OS
	var localOS string
	if runtime.GOOS == "darwin" {
		localOS = "osx"
	} else {
		localOS = runtime.GOOS
	}

	// file := "istioctl-" + *version + "-" + localOS + ".tar.gz.sha256"
	file := "istioctl-" + *version + "-" + localOS + ".tar.gz"
	link := "https://github.com/istio/istio/releases/download/" + *version + "/" + file
	localFile := filepath.Join(*cachDir, file)

	if _, err := os.Stat(localFile); err == nil {
		log.Printf("%v already exists, skipping download\n", file)
	} else {
		download(link, localFile)
	}

	// unzip file
	log.Printf("Extracting %v\n", file)
	r, err := os.Open(localFile)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		ExtractTarGz(r, *cachDir)
	}
}

//https://github.com/istio/istio/releases/download/1.12.9/istioctl-1.12.9-osx.tar.gz.sha256
//https://github.com/istio/istio/releases/download/1.12.9/istioctl-1.12.9-osx-arm64.tar.gz.sha256

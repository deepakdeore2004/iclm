package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func download(link, localFile string) {

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
		log.Fatal("download failed: ", resp.StatusCode)
	}
}

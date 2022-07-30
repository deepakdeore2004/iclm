package main

import (
	"iclm/setup"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// var (
// 	// Info writes logs in the color blue with "INFO: " as prefix
// 	Info = log.New(os.Stdout, "\u001b[34mINFO: \u001B[0m", log.LstdFlags|log.Lshortfile)

// 	// Warning writes logs in the color yellow with "WARNING: " as prefix
// 	Warning = log.New(os.Stdout, "\u001b[33mWARNING: \u001B[0m", log.LstdFlags|log.Lshortfile)

// 	// Error writes logs in the color red with "ERROR: " as prefix
// 	Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)

// 	// Debug writes logs in the color cyan with "DEBUG: " as prefix
// 	Debug = log.New(os.Stdout, "\u001b[36mDEBUG: \u001B[0m", log.LstdFlags|log.Lshortfile)
// )

func main() {
	r := gin.Default()

	r.POST("/v1/setup", func(c *gin.Context) {
		cacheDir := c.DefaultQuery("cachedir", "/tmp")
		version := c.DefaultQuery("version", "1.14.1")

		c.String(http.StatusOK, "Downloading "+version+" in "+cacheDir)

		// Download istio if not downloaded already
		istioCtl := filepath.Join(cacheDir, "istio-"+version, "bin", "istioctl")
		cmd := exec.Command(istioCtl, "version", "--remote=false")
		stdout, _ := cmd.Output()

		if strings.TrimSuffix(string(stdout), "\n") == version {
			log.Printf("%v is present, skipping download\n", istioCtl)
		} else {
			setup.Download(cacheDir, version)
		}

	})

	r.Run(":8080")

}

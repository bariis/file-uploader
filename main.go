package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/oauth2"

	"github.com/putdotio/go-putio/putio"
)

// parameters to be taken
var (
	token         = flag.String("token", "", "user's token for application")
	rootpath      = flag.String("rootpath", "/bin/images", "folder path to be uploaded") // not must
	number0fFiles = 0
)

// colors for logging
var (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
	Blue  = "\033[34m"
)

func main() {
	paths := make(chan string)
	flag.Parse()

	if *token == "" {
		log.Fatal(Red + "Please pass a client token to use the service: -token={$token}")
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *token})
	oauthClient := oauth2.NewClient(context.TODO(), tokenSource)
	client := putio.NewClient(oauthClient)

	var wg = new(sync.WaitGroup)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go worker(paths, wg, client)
	}

	if err := filepath.Walk(*rootpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf(Red+"failed to walk directory: %v"+Reset, err)
		}
		if !info.IsDir() {
			paths <- path
		}
		return nil
	}); err != nil {
		log.Printf(Red+"failed walk: %v"+Reset, err)
	}

	close(paths)
	wg.Wait()
	log.Printf(Blue+"Uploaded files: %v."+Reset, number0fFiles)
}

func worker(paths <-chan string, wg *sync.WaitGroup, client *putio.Client) {
	defer wg.Done()
	for path := range paths {
		f, err := os.Open(path)
		if err != nil {
			log.Printf(Red+"Failed to open file %v for reading, err:"+Reset, f.Name(), err)
		}
		upload, err := client.Files.Upload(context.TODO(), f, path, 0)
		if err != nil {
			log.Printf(Red+"Failed to upload file %v: err:"+Reset, upload.File.Name, err)
		}
	
		number0fFiles++
		log.Printf(Green+"File %v has been uploaded succesfully"+Reset, upload.File.Name)
	}
}

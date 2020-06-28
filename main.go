package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/cavaliercoder/grab/grabui"
	"github.com/frezbo/rpminfo/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: rpminfo <url of yum repo>")
		os.Exit(1)
	}

	baseURL := os.Args[1]

	repomdURL := joinURL(baseURL, "repodata/repomd.xml")

	resp, err := http.Get(repomdURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	byteValue, _ := ioutil.ReadAll(resp.Body)

	Href := parser.PackageListURI(byteValue)
	if Href == "" {
		log.Fatalf("Cannot find primary package data")
	}
	log.Println(Href)
	dbfile := joinURL(baseURL, Href)
	respch, err := grabui.GetBatch(context.Background(), 0, ".", dbfile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if (<-respch).Err() != nil {
		log.Fatal("File downloaded failed")
	}

	gzFilePath := strings.Split(Href, "/")[1]

	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer gzFile.Close()

	zr, err := gzip.NewReader(gzFile)
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()
	err = parser.PackageListV2(zr)
	if err != nil {
		log.Fatal(err)
	}

}

func joinURL(baseURL string, path string) string {
	burl, _ := url.Parse(baseURL)
	burl.Path = filepath.Join(burl.Path, path)
	return burl.String()
}

package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"os"

	hel "github.com/thejini3/go-helper"
)

type xmlModel struct {
	Contents []xmlContent `xml:Contents`
}
type xmlContent struct {
	Key string `xml:Key`
}

var xmls xmlModel
var url string
var outDir string

var exts = map[string]string{
	"image/png":        "png",
	"image/jpg":        "jpg",
	"image/jpeg":       "jpeg",
	"application/json": "json",
}

func main() {
	flags()

	bytes := hel.URLContentMust(url, hel.UserAgentCrawler)

	// hel.Pl(string(bytes))

	err := xml.Unmarshal(bytes, &xmls)

	if err != nil {
		panic("Error xml Unmarshal - " + err.Error())
	}

	for _, c := range xmls.Contents {

		cURL := url + "/" + c.Key

		hel.Pl("Getting", cURL)

		if resp, err := hel.URLResponse(cURL, hel.UserAgentCrawler); err == nil {
			defer resp.Body.Close()

			if bytes, err := ioutil.ReadAll(resp.Body); err == nil {

				var fname = outDir + "/" + c.Key

				ext, ok := exts[resp.Header.Get("Content-Type")]

				if ok {
					fname += "." + ext
				}

				if err := hel.BytesToFile(fname, bytes); err == nil {
					hel.Pl("Wrote -", fname)
				}
			}
		}

	}

}

func flags() {

	flag.StringVar(&url, "u", "", "(required) Asset url, ex: https://s3.amazonaws.com/linksassetstest")
	flag.StringVar(&outDir, "o", "", "(required) Output directory, ex: out")

	flag.Parse()

	if len(url) == 0 || len(outDir) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	outDir = removeForwardSlashIfHave(outDir)
	url = removeForwardSlashIfHave(url)

	if err := hel.DirCreateIfNotExists(outDir); err != nil {
		panic("Error creating directory - " + err.Error())
	}
}

func removeForwardSlashIfHave(str string) string {
	c := len(str)
	have := string(str[c-1]) == "/"
	if have {
		return str[:c-1]
	}
	return str
}

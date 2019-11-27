package main

import (
	"encoding/json"
	"flag"
	"os"
	"path"
	"path/filepath"

	"github.com/shiguanghuxian/txai"
)

func imageFiles(root string) <-chan string {
	ret := make(chan string)
	go func() {
		defer close(ret)
		filepath.Walk(root, func(fp string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			if ext := path.Ext(fp); ext == ".jpg" || ext == ".png" {
				if abs, err := filepath.Abs(fp); err == nil {
					ret <- abs
				}
			}
			return nil
		})
	}()
	return ret
}

func main() {
	client := txai.New(config.ak, config.sk, false)
	for imgFile := range imageFiles(config.file) {
		result := map[string]interface{}{}
		if val, err := client.ImageFuzzyForPath(imgFile); err == nil {
			result["fuzzy.ind"] = val.Data.Fuzzy
			result["fuzzy.result"] = val.Data.Confidence
		}
		if it, err := client.ImageTagForPath(imgFile); err == nil {
			for _, tag := range it.Data.TagList {
				result["label."+tag.TagName] = tag.TagConfidence
			}
		}
		if vp, err := client.VisionPornForPath(imgFile); err == nil {
			result["vp.level"] = vp.Level
			for _, tag := range vp.Data.TagList {
				result["v.tag."+tag.TagName] = tag.TagConfidenceF
			}

		}
		result["uri"] = imgFile
		json.NewEncoder(os.Stdout).Encode(result)
	}
}

var config struct {
	ak   string
	sk   string
	file string
}

func init() {
	flag.StringVar(&config.ak, "ak", "1106841712", "")
	flag.StringVar(&config.sk, "sk", "WvfuiAt0tVgnE6ij", "")
	flag.StringVar(&config.file, "file", "", "")

	flag.Parse()
}

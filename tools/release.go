package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/cheggaaa/pb.v1"
)

type nsipConfig struct {
	Username string
	Password string
}

var cfg nsipConfig

// Release Information
type jsonRelease struct {
	Id        int    `json:"id"`
	TagName   string `json:"tag_name"`
	UploadURL string `json:"upload_url"`
}

type jsonUpload struct {
	DownloadURL string `json:"browser_download_url"`
}

// XXX Variables / Config
//	Base URL
//	Project Name
// 	How to get next release name/number

func getRelease(project string) jsonRelease {
	release := jsonRelease{}

	// Get the current Releases
	req, err := http.NewRequest("GET", "https://api.github.com/repos/nsip/"+project+"/releases/latest", nil)
	if err != nil {
		// handle err
	}
	// XXX DO NOT RELEASE !!!
	req.SetBasicAuth(cfg.Username, cfg.Password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	// log.Printf("Body: %s\n", body)

	err = json.Unmarshal(body, &release)
	if err != nil {
		log.Fatal(err)
	}

	return release
}

func uploadFile(release jsonRelease, name string, filename string) jsonUpload {
	upload := jsonUpload{}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	uploadURLs := strings.Split(release.UploadURL, "{")
	uploadURL := uploadURLs[0]

	versionSuffix := "-" + release.TagName + ".zip"
	versionedName := strings.Replace(name, ".zip", versionSuffix, 1)

	// log.Printf("Sending file to %s", uploadURL+"?name="+name)
	log.Printf("Sending file to %s", uploadURL+"?name="+versionedName)

	bar := pb.StartNew(int(fi.Size())).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	pr := bar.NewProxyReader(f)

	// req, err := http.NewRequest("POST", uploadURL+"?name="+name, pr)
	req, err := http.NewRequest("POST", uploadURL+"?name="+versionedName, pr)
	if err != nil {
		// handle err
		log.Printf("Error upload = %s", err)
	}
	req.SetBasicAuth(cfg.Username, cfg.Password)
	req.Header.Set("Content-Type", "application/zip")
	req.ContentLength = fi.Size()

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		// handle err
		log.Printf("Error upload = %s", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	bar.Finish()
	err = json.Unmarshal(body, &upload)
	if err != nil {
		log.Fatal(err)
	}

	if len(upload.DownloadURL) > 0 {
		log.Printf("Download URL = %s", upload.DownloadURL)
	} else {
		log.Printf("Body: %s\n", body)
	}
	return upload
}

// Arguments
//	1 = Project name
//	2 = File name to upload
//	3 = Local file path
// If called with just the project name, we generate versioning code
func main() {
	// Load configuration
	if _, err := toml.DecodeFile(os.Getenv("HOME")+"/.nsip.toml", &cfg); err != nil {
		log.Fatalln("Unable to read default config, aborting.", err)
	}
	// log.Printf("Username=%s, Password=%s", cfg.Username, cfg.Password)

	// TODO Check Username and password above

	// TODO Check os parameters

	// Get Latest Release
	release := getRelease(os.Args[1])

	log.Printf("Received release %d as %s", release.Id, release.TagName)

	if len(os.Args) == 2 {
		fmt.Printf("package version\nvar(\nId = %d\nTagName = \"%s\"\n)\n", release.Id, release.TagName)
	} else {
		uploadFile(release, os.Args[2], os.Args[3])
	}
}

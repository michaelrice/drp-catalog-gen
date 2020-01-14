package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/digitalrebar/provision/v4/models"
)

func main() {

	var catalog models.Content
	var body []byte
	useWeb := false

	if useWeb {
		body = getCatalogFromWeb("http://repo.rackn.io/")
	} else {
		body = loadCatalogFromFile("./data/catalog.json")
	}
	err := json.Unmarshal([]byte(body), &catalog)

	if err != nil {
		log.Fatal(err)
	}

	for _, v := range catalog.Sections {
		for items := range v {
			fmt.Println(items)
		}
	}
}

func getCatalogFromWeb(url string) []byte {
	drpClient :=  http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "drp-catalog-gen/0.1")

	res, getErr := drpClient.Do(req)

	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	defer res.Body.Close()
	return body
}

func loadCatalogFromFile(path string) []byte {
	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return jsonFile
}
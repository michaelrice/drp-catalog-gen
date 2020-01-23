/*
Copyright © 2020 RackN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a custom catalog for Digital Rebar",
	Long: `drpcli catalog create --catalog-src https://repo.rackn.io/ -p stable -o new_catalog.json 
Would create a new catalog in your cwd named new_catalog.json and would contain only stable 
versions of the packages.`,
	Run: func(cmd *cobra.Command, args []string) {
		command(cmd, args)
	},
}

var catalogSrc string
var pkgVer string
var outfile string
var catalog models.Content

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&catalogSrc, "catalog-src", "c", "", "catalog-src someCatalog (required) URL or File Path supported")
	createCmd.MarkFlagRequired("catalog-src")
	createCmd.Flags().StringVarP(&pkgVer, "pkg-version", "p", "", "pkg-version tip|stable (required)")
	createCmd.MarkFlagRequired("pkg-version")
	createCmd.Flags().StringVarP(&outfile, "outfile", "o", "","outfile new_catalog.json (required)")
	createCmd.MarkFlagRequired("outfile")
}

func command(cmd *cobra.Command, args []string) {
	var body []byte
	if isValidURl(catalogSrc) {
		body = loadCatalogFromWeb(catalogSrc)
	} else {
		body = loadCatalogFromFile(catalogSrc)
	}
	err := json.Unmarshal([]byte(body), &catalog)
	if err != nil {
		log.Fatal(err)
	}

	newCat := extractPackageSet(&catalog, pkgVer)
	writeCatalog(newCat, outfile)
}

func isValidURl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	} else {
		return true
	}
}

func loadCatalogFromWeb(url string) []byte {
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
	body, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func catalogToMap(cat *models.Content) map[string][]*models.CatalogItem {
	res := map[string][]*models.CatalogItem{}
	for _, v := range cat.Sections["catalog_items"] {
		item := &models.CatalogItem{}
		if err := models.Remarshal(v, &item); err != nil {
			continue
		}
		res[item.Name] = append(res[item.Name], item)
	}
	return res
}


func writeCatalog(cat *models.Content, filename string) {
	data,err := json.MarshalIndent(cat, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func extractPackageSet(cat *models.Content, version string) *models.Content {
	newMap := map[string]interface{}{}
	for k, v := range cat.Sections["catalog_items"] {
		item := &models.CatalogItem{}
		if err := models.Remarshal(v, &item); err != nil {
			continue
		}
		if item.Version == version {
			newMap[k] = item
		}
	}
	cat.Sections["catalog_items"] = newMap
	return cat
}

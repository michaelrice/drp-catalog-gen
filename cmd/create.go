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
	Short: "Create a drp catalog",
	Long: `Used to create a custom catalog for Digital Rebar.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringVarP(&catalogSrc, "catalog-src", "c", "", "catalog-src someCatalog (required)")
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

	//cataMap := catalogToMap(&catalog)

	//mainLoop()
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


//func createView() {
//	view := ui.AddWindow(5, 5, 80, 24, "RackN Catalog Generator")
//	view.SetTitleButtons(ui.ButtonMaximize | ui.ButtonClose)
//	view.SetActiveBackColor(ui.ColorBlack)
//
//	cataFrame := ui.CreateFrame(view, 10, 12, ui.BorderThin, ui.AutoSize)
//	cataFrame.SetActiveBackColor(ui.ColorBlack)
//
//
//	frmFull := ui.CreateFrame(cataFrame, 40, 5, ui.BorderThin, ui.Fixed)
//
//	frmFull.SetTitle("Catalog Items")
//	frmFull.SetActiveBackColor(ui.ColorBlack)
//	//frmHalf := ui.CreateFrame(frmViews, 8, 5, ui.BorderThin, ui.Fixed)
//	//frmHalf.SetPack(ui.Vertical)
//	//frmHalf.SetTitle("Half")
//	//frmNone := ui.CreateFrame(frmViews, 8, 5, ui.BorderThin, ui.Fixed)
//	//frmNone.SetPack(ui.Vertical)
//	//frmNone.SetTitle("None")
//	//
//	btnF1 := ui.CreateButton(frmFull, 2, 2, "Quit", ui.Fixed)
//	btnF1.SetShadowType(ui.ShadowNone)
//	btnF1.SetSize(1,1)
//	btnF1.SetTextColor(ui.ColorWhiteBold)
//	btnF1.SetBackColor(ui.ColorBlack)
//	btnF1.SetActiveBackColor(ui.ColorBlack)
//	btnF1.SetScale(10)
//	btnF1.SetPos(10, 20)
//
//	btnF1.OnClick(func(ev ui.Event) {
//		go ui.Stop()
//	})
//	//btnH3.OnClick(func(ev ui.Event) {
//	//	go ui.Stop()
//	//})
//	//btnN3.OnClick(func(ev ui.Event) {
//	//	go ui.Stop()
//	//})
//	//
//	ui.ActivateControl(view, btnF1)
//
//}
//
//func mainLoop() {
//	// Every application must create a single Composer and
//	// call its intialize method
//	ui.InitLibrary()
//	defer ui.DeinitLibrary()
//
//	createView()
//
//	// start event processing loop - the main core of the library
//	ui.MainLoop()
//}

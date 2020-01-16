/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// mikeCmd represents the mike command
var mikeCmd = &cobra.Command{
	Use:   "mike",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		body := loadCatalogFromFile("C:/Users/errr/go/src/github.com/digitalrebar/drp-catalog-gen/data/catalog.json")
		err := json.Unmarshal([]byte(body), &catalog)

		if err != nil {
			log.Fatal(err)
		}
		// str: [catalog_items]
		cataMap := catalogToMap(&catalog)
		for name, objs := range cataMap {
			fmt.Printf("%s -> %d versions available \n", name, len(objs))
			for _, cataItem := range objs {
				fmt.Println(" 	-> ", cataItem.Version)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mikeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mikeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mikeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

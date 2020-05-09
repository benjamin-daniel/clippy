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
	"errors"
	"strings"

	"github.com/benjamin-daniel/clippy/store"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

var like string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search your clipboard history",
	Long:  `Search your keyboard history with a keyword`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		db, err = gorm.Open("sqlite3", path+"/test.db")
		if err != nil {
			return err
			// panic("failed to connect database")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			return err
		}
		like = strings.TrimSpace(args[0])
		if len(like) < 2 {
			return errors.New("You need to proved a valid search string")
		}
		like = "%" + args[0] + "%"
		return searchClipBoardItems(limit, like)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		defer db.Close()
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			return (err)
		}
		mainQuery := func() {
			var clips store.ClipBoardItems
			db.Offset(listPage.Skip).Limit(limit).Order("id desc").Where("text LIKE ?", like).Find(&clips) //.Count(&listPage.count)
			clips.Print(pageN)
		}
		if !listPage.End() {
			return ask(limit, mainQuery)
		}
		return nil
	},
}

func searchClipBoardItems(limit int, like string) error {
	var clips store.ClipBoardItems
	listPage = &store.Page{Limit: float64(limit), Page: 1}
	db.Limit(limit).Order("id desc").Where("text LIKE ?", like).Find(&clips).Count(&listPage.Count)
	clips.Print(pageN)
	return nil
}

func init() {
	searchCmd.Flags().Int("limit", 20, "Limit the amount of clipboard items printed")
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

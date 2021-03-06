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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/benjamin-daniel/clippy/store"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
)

var listPage *store.Page

// opens a constant connection
var db *gorm.DB

var path string = "/usr/local/clippy"

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List your clipboard history",
	Long:    `List you clipboard history, you can specify a limit to the amount of things printed by using the limit flag.`,
	Aliases: []string{"ls"},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		db, err = gorm.Open("sqlite3", path+"/test.db")
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			return err
		}
		return listClipBoardItems(limit)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		defer db.Close()
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			return (err)
		}
		mainQuery := func() {
			var clips store.ClipBoardItems
			db.Offset(listPage.Skip).Limit(limit).Order("id desc").Find(&clips) //.Count(&listPage.count)
			clips.Print(pageN)
		}
		if !listPage.End() {
			return ask(limit, mainQuery)
		}
		return nil
	},
}

func init() {
	listCmd.Flags().Int("limit", 20, "Limit the amount of clipboard items printed")
	listCmd.Flags().Int("page", 1, "Specify the page you want to see")
	rootCmd.AddCommand(listCmd)
	// fmt.Println(cntxt)
	// cntxt = nil
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listClipBoardItems(limit int) error {
	var clips store.ClipBoardItems
	listPage = &store.Page{Limit: float64(limit), Page: 1}
	db.Limit(limit).Order("id desc").Find(&clips).Count(&listPage.Count)
	clips.Print(pageN)
	return nil
}
func pageN() {
	listPage.Init()
	fmt.Println(listPage.Pretty())
}
func printClip(clips store.ClipBoardItems) {
	for i := 0; i < len(clips); i++ {
		clip := clips[i]
		fmt.Printf("  %d\t%s\n", clip.ID, clip.TruncateText(50))
	}
	listPage.Init()
	fmt.Println(listPage.Pretty())
}

func ask(limit int, mainQuery func()) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintln(os.Stderr, listPage.Commands())
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			return (err)
		}
		s = strings.TrimSuffix(s, "\n")
		s = strings.ToLower(s)
		switch s {
		case "next":
			listPage.NextPage()
			mainQuery()
			if listPage.End() {
				fmt.Println(chalk.Magenta, "End of list", chalk.Reset)
				return nil
			}
			fmt.Fprintln(os.Stderr, listPage.Commands())
			continue
		case "prev":
			listPage.PrevPage()
			mainQuery()
			fmt.Fprintln(os.Stderr, listPage.Commands())
			continue
		case "last":
			listPage.Page = listPage.Max
			listPage.Init()
			mainQuery()
			return nil
		case "exit":
			return nil
		default:
			fmt.Fprintln(os.Stderr, listPage.Commands())
			continue
		}
	}
}

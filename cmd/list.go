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
	"math"
	"os"
	"strings"

	"github.com/benjamin-daniel/clippy/store"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
)

// Page represents a page in the document
type Page struct {
	page  int
	limit float64
	count float64
	skip  int
	max   int
}

var listPage *Page

func (p *Page) init() {
	limit := int(p.limit)
	p.skip = (p.page * limit) - limit
	p.max = int(math.Ceil(p.count / p.limit))
}

func (p *Page) nextPage() {
	p.page++
	p.init()
}

func (p *Page) prevPage() {
	if p.page == 1 {
		return
	}
	p.page--
	p.init()
}

func (p *Page) end() bool {
	return p.page >= p.max
}

func (p *Page) more() bool {
	return p.page < p.max
}

func (p *Page) String() string {
	return fmt.Sprintf("Page: %d\t Limit: %f\tCount:%f\t Skip: %d Max: %d More: %t End: %t\n", p.page, p.limit, p.count, p.skip, p.max, p.more(), p.end())
}
func (p *Page) pretty() string {
	return fmt.Sprintf("%sPage %d of %d%s", chalk.Blue, p.page, p.max, chalk.Reset)
}
func (p *Page) commands() string {
	return fmt.Sprintf("Please enter \nlast: to view last page\nnext: to view next page\nprev: to view previous page\nexit: to exit the shell")
}

// opens a constant connection
var db *gorm.DB

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		db, err = gorm.Open("sqlite3", "test.db")
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
		return listClipBoardItems(limit)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		defer db.Close()
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			panic(err)
		}
		if !listPage.end() {
			ask(limit)
		}
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
	// Observe how the b's and the d's, despite appearing in the
	// second cell of each line, belong to different columns.
	var clips []*store.ClipBoardItem
	listPage = &Page{limit: float64(limit), page: 1}
	db.Limit(limit).Order("id desc").Find(&clips).Count(&listPage.count)
	printClip(clips)
	return nil
}
func printClip(clips []*store.ClipBoardItem) {
	for i := 0; i < len(clips); i++ {
		clip := clips[i]
		fmt.Printf("  %d\t%s\n", clip.ID, clip.TruncateText(50))
	}
	listPage.init()
	fmt.Println(listPage.pretty())
}
func ask(limit int) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintln(os.Stderr, listPage.commands())
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		s = strings.TrimSuffix(s, "\n")
		s = strings.ToLower(s)
		switch s {
		case "next":
			listPage.nextPage()
			var clips []*store.ClipBoardItem
			db.Offset(listPage.skip).Limit(limit).Order("id desc").Find(&clips) //.Count(&listPage.count)
			printClip(clips)
			if listPage.end() {
				fmt.Println(chalk.Magenta, "End of list", chalk.Reset)
				return nil
			}
			fmt.Fprintln(os.Stderr, listPage.commands())
			continue
		case "prev":
			listPage.prevPage()
			var clips []*store.ClipBoardItem
			// fmt.Println(listPage)
			db.Offset(listPage.skip).Limit(limit).Order("id desc").Find(&clips) //.Count(&listPage.count)
			printClip(clips)
			continue
		case "last":
			listPage.page = listPage.max
			listPage.init()
			var clips []*store.ClipBoardItem
			// fmt.Println(listPage)
			db.Offset(listPage.skip).Limit(limit).Order("id desc").Find(&clips) //.Count(&listPage.count)
			printClip(clips)
			return nil
		case "exit":
			return nil
		default:
			fmt.Fprintln(os.Stderr, listPage.commands())
			continue
		}
	}
}
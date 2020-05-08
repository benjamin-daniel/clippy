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
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/benjamin-daniel/clippy/store"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			return err
		}
		return listClipBoardItems(limit)
	},
}

func init() {
	listCmd.Flags().Int("limit", 20, "Limit the amount of clipboard items printed")
	rootCmd.AddCommand(listCmd)
	// fmt.Println(cntxt)
	// cntxt = nil
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listClipBoardItems(limit int) error {
	// Observe how the b's and the d's, despite appearing in the
	// second cell of each line, belong to different columns.
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var clips []store.ClipBoardItem
	db.Limit(limit).Order("id desc").Find(&clips)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 30, ' ', tabwriter.Escape)
	fmt.Fprintln(w, chalk.Blue, fmt.Sprintf("%s\t%s", "ID", "Text"), chalk.Reset)
	for i := 0; i < len(clips); i++ {
		clip := clips[i]
		fmt.Fprintln(w, fmt.Sprintf("%d\t%s", clip.ID, clip.TruncateText(50)))
	}
	w.Flush()
	return nil
}

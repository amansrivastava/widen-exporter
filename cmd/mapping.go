/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/cheggaaa/pb/v3"
)

var filetype string

type result struct {
	Scroll_id string
	Count     int `json:"total_count"`
	Items     []asset
}
type asset struct {
	Id       string
	Filename string
	Metadata map[string]map[string][]string `json:"metadata"`
}

type Mapping []string

// mappingCmd represents the mapping command
var mappingCmd = &cobra.Command{
	Use:   "mapping",
	Short: "Export mapping sheet for media_acquiadam:2.x module.",
	Long:  `This command generates a csv mapping file which has webdam_id,widen_id mapping for replacing the old webdam ids with new webdamids in Drupal.`,

	Run: func(_ *cobra.Command, args []string) {
		query := "iwi:-(isEmpty)"
		url := baseUrl + "assets/search?expand=metadata&limit=100&scroll=true&query=" + query
		res, err := getData(url)
		f, err := os.Create(exportFile)
		defer f.Close()
		w := csv.NewWriter(f)
		defer w.Flush()
		if err != nil {
			log.Fatalln("failed to open file", err)
		}
		fmt.Println("Total assets found:-", res.Count)
		bar := pb.StartNew(res.Count / 100)
		r := record{"webdam_id", "widen_id"}
		writeToFile(r, w)
		items := res.Items
		for _, item := range items {
			r := record{
				item.Metadata["fields"]["webdam_id"][0],
				item.Id,
			}
			writeToFile(r, w)
		}
		if res.Scroll_id != "" {
			scrollAssets(res.Scroll_id, w, bar)
		}

	},
}

func scrollAssets(scroll_id string, w *csv.Writer, bar *pb.ProgressBar) {
	url := baseUrl + "assets/search/scroll?expand=metadata&scroll_id=" + scroll_id
	result, err := getData(url)
	if err != nil {
		return
	}
	items := result.Items
	for _, item := range items {
		r := record{
			item.Metadata["fields"]["webdam_id"][0],
			item.Id,
		}
		writeToFile(r, w)
	}
	if result.Scroll_id != "" && len(result.Items) > 0 {
		bar.Increment()
		scrollAssets(result.Scroll_id, w, bar)
	}
	bar.Finish()
}

func writeToFile(r record, w *csv.Writer) {
	if err := w.Write(r); err != nil {
		log.Fatalln("error writing record to file", err)
	}
}

func init() {
	rootCmd.AddCommand(mappingCmd)
}

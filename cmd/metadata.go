/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

type record []string

// metadataCmd represents the metadata command
var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "It will generate a csv with given metadata fields.",
	Long:  `This will generate a csv file with all the metadata fields added in the args like "widen-exporter metadata keywords webdam_version"`,

	Run: func(_ *cobra.Command, args []string) {
		url := baseUrl + "assets/search?expand=metadata,file_properties&limit=100&scroll=true&query=" + query
		res, err := getData(url)
		if err != nil {
			return
		}
		f, err := os.Create(exportFile)
		defer f.Close()
		w := csv.NewWriter(f)
		defer w.Flush()
		fmt.Println("Total assets found:-", res.Count)
		bar := pb.StartNew(res.Count / 100)
		if err != nil {
			log.Fatalln("failed to open file", err)
		}
		items := res.Items
		headers := generateHeaders(items, args)
		write(headers, w)
		for _, item := range items {
			r := record{
				item.Id,
			}
			if contains(args, "filename") {
				r = append(r, item.Filename)
			}
			for _, key := range args {
				v, exists := item.Metadata["fields"][key]
				if exists {
					r = append(r, strings.Join(v, ","))
				}
			}
			write(r, w)
		}
		if res.Scroll_id != "" {
			scrollMetadataAssets(res.Scroll_id, w, args, bar)
		}

	},
}

func scrollMetadataAssets(scroll_id string, w *csv.Writer, args []string, bar *pb.ProgressBar) {
	url := baseUrl + "assets/search/scroll?expand=metadata&scroll_id=" + scroll_id
	res, err := getData(url)
	if err != nil {
		return
	}
	items := res.Items
	for _, item := range items {
		r := record{
			item.Id,
		}
		if contains(args, "filename") {
			r = append(r, item.Filename)
		}
		for _, key := range args {
			v, exists := item.Metadata["fields"][key]
			if exists {
				r = append(r, strings.Join(v, ","))
			}
		}
		write(r, w)
	}
	if res.Scroll_id != "" && len(res.Items) > 0 {
		bar.Increment()
		scrollMetadataAssets(res.Scroll_id, w, args, bar)
	}
	bar.Finish()
}

func write(r record, w *csv.Writer) {
	if err := w.Write(r); err != nil {
		log.Fatalln("error writing record to file", err)
	}
}

func generateHeaders(items []asset, args []string) record {
	record := record{
		"id",
	}
	if contains(args, "filename") {
		record = append(record, "filename")
	}
	for _, key := range args {
		_, exists := items[0].Metadata["fields"][key]
		if exists {
			record = append(record, key)
		}
	}
	return record
}

func init() {
	metadataCmd.PersistentFlags().StringVar(&query, "query", "", "Your search query")
	rootCmd.AddCommand(metadataCmd)
}

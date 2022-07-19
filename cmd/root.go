/*
Copyright © 2022 Aman

*/
package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var exportFile string
var query string
var token string
var baseUrl string = "https://api.widencollective.com/v2/"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "widen-exporter",
	Short: "Export data using Widen APIs.",
	Long:  `This is a quick CLI tool to export widen mapping csv and metadata for different purposes.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		url := baseUrl + "user"
		method := "GET"

		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		req.Header.Add("Authorization", "Bearer "+viper.GetString("token"))
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode == 401 {
			return errors.New("Auth failed.")
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/.widen-exporter.yml)")
	rootCmd.PersistentFlags().StringVar(&exportFile, "filename", "export.csv", "exported filename (default is export.csv)")

	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Widen API Token")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		pwd, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in home directory with name ".widen-exporter" (without extension).
		viper.AddConfigPath(pwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".widen-exporter")
	}
	viper.ReadInConfig()
	viper.AutomaticEnv() // read in environment variables that match
	viper.BindPFlag("token", rootCmd.Flags().Lookup("token"))

}

func getData(url string) (result, error) {
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", "Bearer "+viper.GetString("token"))
	result := result{}
	if err != nil {
		return result, err
	}
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		return result, jsonErr
	}
	return result, nil
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

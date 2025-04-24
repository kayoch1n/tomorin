/*
Copyright © 2025 kayoch1n <hanayolawk@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"

	"github.com/kayoch1n/tomorin/revsh"
)

var depsOpts = struct {
	config string
}{}

// depsCmd represents the config command
var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Create a script to check dependencies",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(depsOpts.config)
		if err != nil {
			log.Fatalf("failed:%v\n", err)
		}
		var config revsh.Config
		yaml.Unmarshal(data, &config)

		var file *os.File
		file, err = os.OpenFile("check-deps", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Fatalf("failed:%v\n", err)
		}

		if err = revsh.Dependencies(&config, file); err != nil {
			log.Fatalf("failed to create configure file:%v\n", err)
		} else {
			log.Println("check-deps created.")
		}
	},
}

func init() {
	rootCmd.AddCommand(depsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// depsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// depsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	depsCmd.Flags().StringVarP(&depsOpts.config, "config", "c", "config.yml", "Path to the config file")
}

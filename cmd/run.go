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
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/kayoch1n/tomorin/revsh"
)

var runOpts = struct {
	config  string
	timeout int
	wait    int
}{}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run reverse shell samples from current host",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(runOpts.config)
		if err != nil {
			log.Fatalf("failed to open config file:%v\n", err)
		}
		var config revsh.Config
		yaml.Unmarshal(data, &config)
		log.Printf("%d samples loaded\n", len(config.Samples))

		if config.Timeout == 0 {
			config.Timeout = runOpts.timeout
		}
		if config.Wait == 0 {
			config.Wait = runOpts.timeout
		}

		results := revsh.Execute(&config)

		data, err = yaml.Marshal(results)
		if err != nil {
			log.Fatalf("failed to marshal: %v\n", err)
		}

		now := time.Now().Format("20060102150405")
		filename := now + ".yml"
		os.WriteFile(filename, data, 0644)
		log.Printf("configs saved to %s\n", filename)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().StringVarP(&runOpts.config, "config", "c", "config.yml", "Path to the config file")
	runCmd.Flags().IntVar(&runOpts.timeout, "timeout", 10, "Timeout of each sample")
	runCmd.Flags().IntVar(&runOpts.wait, "wait", 7, "Timeout until the next sample")
}

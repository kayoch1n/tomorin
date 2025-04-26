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
	"strings"
	"sync"

	"github.com/kayoch1n/tomorin/revsh"
	"github.com/spf13/cobra"
)

var serveOpts = struct {
	addresses []string
	cmd       string
	cmdExit   bool
	timeout   int
}{}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run reverse shell server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(serveOpts.addresses) == 0 {
			log.Fatal("at least one address is required")
		}
		var addresses [][2]string
		for _, address := range serveOpts.addresses {
			parts := strings.Split(address, ":")
			var proto, ip, port string
			switch len(parts) {
			case 1:
				proto, ip = "tcp", "0.0.0.0"
				port = parts[0]
			case 2:
				proto, port = parts[0], parts[1]
				ip = "0.0.0.0"
			case 3:
				proto, ip, port = parts[0], parts[1], parts[2]
			}
			proto = strings.ToLower(proto)
			if proto != "tcp" && proto != "udp" {
				log.Fatalf("Protocol %s in %s is not supported\n", proto, address)
			}
			addresses = append(addresses, [2]string{proto, ip + ":" + port})
		}

		command := serveOpts.cmd + "\n"
		if serveOpts.cmdExit {
			command = command + "exit\n"
		}

		var wg sync.WaitGroup
		for _, parts := range addresses {
			proto, address := parts[0], parts[1]
			wg.Add(1)
			go func() {
				defer wg.Done()
				log.Printf("listening %s on %s: cmd=\"%s\"\n", proto, address, revsh.LogEscape(command))
				switch proto {
				case "tcp":
					err := revsh.ServeTCP(address, command, serveOpts.timeout)
					log.Printf("exiting: %v ...\n", err)
				case "udp":
					err := revsh.ServeUDP(address, command)
					log.Printf("exiting: %v ...\n", err)
				}
			}()
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serveCmd.Flags().StringSliceVarP(&serveOpts.addresses, "address", "a", nil, "Addresses to listen, each of which has the format of [PROTO:[IP:]]PORT. If PROTO is not specified, then \"tcp\" will be used. If IP is not specified, then \"udp\" will be used.")
	serveCmd.Flags().StringVar(&serveOpts.cmd, "cmd", "whoami && sleep 2", "Command to be executed once a remove host is connected")
	serveCmd.Flags().IntVar(&serveOpts.timeout, "tcp-timeout", 10, "Timeout for each connection. Applied to TCP only")
	serveCmd.Flags().BoolVar(&serveOpts.cmdExit, "cmd-exit", true, "Whether to append an \"exit\" to the end of the provided cmd. This should always be true if you want a graceful exit on the remote shell.")
}

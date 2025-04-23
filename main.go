package main

import (
	"flag"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/kayoch1n/tomorin/revsh"
)

var (
	configFile    string
	scriptTimeout int
	waitTimeout   int
	// tcpPort       int
	// udpPort       int
	target string
)

func main() {
	flag.StringVar(&configFile, "c", "config.yml", "config file")
	flag.IntVar(&scriptTimeout, "timeout", 10, "timeout for each script")
	flag.IntVar(&waitTimeout, "wait", 10, "timeout until next script")
	flag.StringVar(&target, "t", "", "SSH target")
	flag.Parse()

	if target == "" {
		log.Fatalf("SSH target is required")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to open config file:%v\n", err)
	}
	var config revsh.Config
	yaml.Unmarshal(data, &config)
	log.Printf("%d samples loaded\n", len(config.Samples))

	if config.Timeout == 0 {
		config.Timeout = scriptTimeout
	}
	if config.Wait == 0 {
		config.Wait = waitTimeout
	}

	var results []revsh.Result
	results, err = revsh.Execute(target, &config)
	if err != nil {
		log.Fatalf("failed to execute: %v\n", err)
	}

	data, err = yaml.Marshal(results)
	if err != nil {
		log.Fatalf("failed to marshal: %v\n", err)
	}

	now := time.Now().Format("20060102150405")
	filename := now + ".yml"
	os.WriteFile(filename, data, 0644)
	log.Printf("configs saved to %s\n", filename)
}


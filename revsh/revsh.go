package revsh 

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"os/exec"

	"github.com/kballard/go-shellquote"
)

type properties struct {
	Timeout int `yaml:"timeout,omitempty"`
	Wait int `yaml:"wait,omitempty"`
}

type Sample struct {
	Name   string `yaml:"name,omitempty"`
	Script string `yaml:"script,omitempty"`
	properties `yaml:",inline"`
}

type Result struct {
	Sample `yaml:",inline"`
	Start  string `yaml:"start,omitempty"` // start time of the execution
	End    string `yaml:"end,omitempty"`  // end time of the execution
	Stderr string `yaml:"stderr,omitempty"` // captured stderr of the execution
	Stdout string `yaml:"stdout,omitempty"` // captured stdout of the execution
}

type Config struct {
	Samples []Sample `yaml:"samples,omitempty"`
	properties `yaml:",inline"`
}

func (c *Sample) content() string {
	return c.Script + "\nexit\n"
}


func Execute(target string, config *Config) (results []Result, err error) {
	for _, sample := range config.Samples {
		if sample.Timeout == 0 {
			sample.Timeout = config.Timeout
		}
		if sample.Wait == 0 {
			sample.Wait = config.Wait
		}
		result := Result{
			Sample: sample,
			Start: time.Now().Format("2006-01-02 15:04:05"),
		}
		log.Printf("execute script %s at %s\n", result.Name, result.Start)
		func() {
			var file *os.File
			file, err = os.CreateTemp("", "revsh")
			if err != nil {
				return
			}
			defer os.Remove(file.Name())

			if err = os.WriteFile(file.Name(), []byte(sample.content()), 0644); err != nil {
				return 
			}

			script := fmt.Sprintf("ssh %s -t -t bash < %s", target, file.Name())
			var stdout, stderr bytes.Buffer

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sample.Timeout)*time.Second)
			defer cancel()
			defer func() {
				result.End = time.Now().Format("2006-01-02 15:04:05")
				result.Stdout = stdout.String()
				result.Stderr = stderr.String()
				log.Printf("script %s finished at %s\n", result.Name, result.End)
				results = append(results, result)
			}()
			cmd := exec.CommandContext(ctx, "bash", "-c", script)
			log.Printf("cmdline: %s\n", shellquote.Join(cmd.Args...))
			cmd.WaitDelay = time.Duration(1)*time.Second
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				log.Printf("failed to execute: %v\n", err)
			}

			time.Sleep(time.Duration(sample.Wait) * time.Second)
		}()
	}
	return
}


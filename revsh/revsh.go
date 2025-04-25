package revsh

import (
	"bytes"
	"context"
	"log"
	"os"
	"time"

	"os/exec"

	"github.com/creack/pty"
)

type properties struct {
	Timeout int `yaml:"timeout,omitempty"`
	Wait    int `yaml:"wait,omitempty"`
}

type Sample struct {
	Name       string `yaml:"name,omitempty"`
	Script     string `yaml:"script,omitempty"`
	properties `yaml:",inline"`
}

type Result struct {
	Sample   `yaml:",inline"`
	Start    string `yaml:"start,omitempty"` // start time of the execution
	End      string `yaml:"end,omitempty"`   // end time of the execution
	Terminal string `yaml:"terminal,omitempty"`
}

type Config struct {
	Samples []Sample `yaml:"samples,omitempty"`
	Depends []string `yaml:"depends,omitempty"`
	//Require	   []string `yaml:"require,omitempty"`
	properties `yaml:",inline"`
}

func (c *Sample) content() string {
	return c.Script + "\nexit\n"
}

func executeInPty(ctx context.Context, script string) (output bytes.Buffer, err error) {
	cmd := exec.CommandContext(ctx, "bash")
	// https://stackoverflow.com/a/78429315/8706476
	cmd.WaitDelay = time.Duration(1) * time.Second

	var ptmx *os.File
	ptmx, err = pty.Start(cmd)
	if err != nil {
		return
	}
	defer ptmx.Close()

	if _, err = ptmx.WriteString(script); err != nil {
		return
	}

	if err = cmd.Wait(); err == nil {
		output.ReadFrom(ptmx)
	}

	return
}

func Execute(config *Config) (results []Result) {
	for _, sample := range config.Samples {
		if sample.Timeout == 0 {
			sample.Timeout = config.Timeout
		}
		if sample.Wait == 0 {
			sample.Wait = config.Wait
		}
		result := Result{
			Sample: sample,
			Start:  time.Now().Format("2006-01-02 15:04:05"),
		}
		log.Printf("execute script %s at %s\n", result.Name, result.Start)
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sample.Timeout)*time.Second)
			defer cancel()

			if buf, err := executeInPty(ctx, sample.content()); err == nil {
				result.Terminal = buf.String()
				result.End = time.Now().Format("2006-01-02 15:04:05")

				log.Printf("script %s finished at %s\n", result.Name, result.End)
				results = append(results, result)
				time.Sleep(time.Duration(sample.Wait) * time.Second)
			} else {
				log.Printf("failed to execute script: %v\n", err)
			}

		}()
	}
	return
}

package internal

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type GitHubLogLevel string

const (
	GitHubLogLevelDebug   GitHubLogLevel = "debug"
	GitHubLogLevelWarning GitHubLogLevel = "warning"
	GitHubLogLevelError   GitHubLogLevel = "error"
)

type Report struct {
	Totals struct {
		Errors   int `json:"errors"`
		Warnings int `json:"warnings"`
		Fixable  int `json:"fixable"`
	} `json:"totals"`
	Files map[string]struct {
		Errors   int `json:"errors"`
		Warnings int `json:"warnings"`
		Messages []struct {
			Message  string `json:"message"`
			Source   string `json:"source"`
			Severity int    `json:"severity"`
			Fixable  bool   `json:"fixable"`
			Type     string `json:"type"`
			Line     int    `json:"line"`
			Column   int    `json:"column"`
		} `json:"messages"`
	} `json:"files"`
}

func (r *Report) CreateMessages() ([]string, error) {
	var ms []string
	for k, v := range r.Files {
		for _, m := range v.Messages {
			if m.Type == "WARNING" && viper.GetBool("ignore-warnings") {
				continue
			}

			var l GitHubLogLevel
			switch m.Type {
			case "ERROR":
				l = GitHubLogLevelError
			case "WARNING":
				l = GitHubLogLevelWarning

			default:
				return nil, errors.New("invalid report message type")
			}

			ms = append(
				ms,
				fmt.Sprintf(
					"::%s file=%s,line=%d::%s",
					l,
					k,
					m.Line,
					m.Message,
				),
			)
		}
	}

	return ms, nil
}

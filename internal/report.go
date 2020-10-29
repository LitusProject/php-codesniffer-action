package internal

import (
	"errors"

	"github.com/google/go-github/v32/github"
	"github.com/spf13/viper"
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

func (r *Report) CreateCheckRunAnnotations() ([]*github.CheckRunAnnotation, error) {
	var as []*github.CheckRunAnnotation
	for k, v := range r.Files {
		for _, m := range v.Messages {
			if m.Type == "WARNING" && viper.GetBool("ignore-warnings") {
				continue
			}

			var l string
			switch m.Type {
			case "ERROR":
				l = "failure"
			case "WARNING":
				l = "warning"

			default:
				return nil, errors.New("invalid report message type")
			}

			a := &github.CheckRunAnnotation{
				Path:            github.String(k),
				StartLine:       github.Int(m.Line),
				EndLine:         github.Int(m.Line),
				AnnotationLevel: github.String(l),
				Message:         github.String(m.Message),
			}

			as = append(as, a)
		}
	}

	return as, nil
}

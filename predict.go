package main

import (
	"strings"

	"github.com/posener/complete"
)

type predictor []arg

func (p predictor) Predict(args complete.Args) []string {
	l := len(args.Completed)
	if l >= len(p) {
		return nil
	}

	return p[l].predict(args.Last)
}

type arg []string

var (
	predictBool     = arg([]string{"true", "false"})
	predictBlockNum = arg([]string{"earliest", "latest", "pending"})
	predictHex      = arg([]string{})
)

func (a arg) predict(s string) []string {
	if s == "" {
		return a
	}
	var p []string
	for _, o := range a {
		if strings.HasPrefix(o, s) {
			p = append(p, o)
		}
	}
	return p

}

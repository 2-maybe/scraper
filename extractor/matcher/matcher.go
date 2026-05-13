package matcher

import (
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
)

type CurrMatcherSignature = int

const (
	Href CurrMatcherSignature = iota
	Img
)

var (
	PrecompiledMatchers map[CurrMatcherSignature]goquery.Matcher
	once                sync.Once
)

func PreCompile() {
	once.Do(func() {
		PrecompiledMatchers = map[CurrMatcherSignature]goquery.Matcher{
			Href: cascadia.MustCompile("a[href]"),
			Img:  cascadia.MustCompile("img"),
		}
	})
}

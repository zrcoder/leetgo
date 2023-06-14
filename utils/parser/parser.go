package parser

import (
	"bytes"
	"strings"

	h2md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

type parser struct {
	data []byte
	err  error
}

func New(data []byte) *parser {
	return &parser{data: data}
}

func NewWithString(s string) *parser {
	return &parser{data: []byte(s)}
}

func (p *parser) PreRrgular() *parser {
	if p.err != nil {
		return p
	}
	replacer := strings.NewReplacer(
		"&nbsp;", " ",
		"\u00A0", " ",
		"\u200B", "",
	)
	p.data = []byte(replacer.Replace(string(p.data)))
	return p
}

func (p *parser) Regular() *parser {
	if p.err != nil {
		return p
	}

	p.data = bytes.TrimSpace(p.data)
	replacer := strings.NewReplacer(
		`\-`, "-",
		`\[`, "[",
		`\]`, `]`,
		`\\n`, "\n",
		`\n`, "\n",
		`\\`, `\`,
		`\*`, `*`,
		`\#`, `#`,
		"\\`", "`",
		`\\'`, `'`,
		`\_`, `_`,
		`\t`, "\t",
		`\\t`, "\t",
	)
	p.data = []byte(replacer.Replace(string(p.data)))
	return p
}

func (p *parser) ToMarkDown() *parser {
	if p.err != nil {
		return p
	}

	converter := h2md.NewConverter("", true, nil)
	replaceSub := h2md.Rule{
		Filter: []string{"sub"},
		Replacement: func(content string, selec *goquery.Selection, opt *h2md.Options) *string {
			selec.SetText(ReplaceSubscript(content))
			return nil
		},
	}
	replaceSup := h2md.Rule{
		Filter: []string{"sup"},
		Replacement: func(content string, selec *goquery.Selection, opt *h2md.Options) *string {
			selec.SetText(ReplaceSuperscript(content))
			return nil
		},
	}
	replaceEm := h2md.Rule{
		Filter: []string{"em"},
		Replacement: func(content string, selec *goquery.Selection, options *h2md.Options) *string {
			return h2md.String(content)
		},
	}
	converter.AddRules(replaceSub, replaceSup, replaceEm)
	p.data, p.err = converter.ConvertBytes(p.data)
	return p
}

func (p *parser) Data() ([]byte, error) {
	return p.data, p.err
}

func (p *parser) String() (string, error) {
	return string(p.data), p.err
}

func ReplaceSubscript(s string) string {
	return SubReplace().Replace(s)
}

func ReplaceSuperscript(s string) string {
	return SupReplace().Replace(s)
}

func SubReplace() *strings.Replacer {
	subscripts := map[string]string{
		"0": "\u2080",
		"1": "\u2081",
		"2": "\u2082",
		"3": "\u2083",
		"4": "\u2084",
		"5": "\u2085",
		"6": "\u2086",
		"7": "\u2087",
		"8": "\u2088",
		"9": "\u2089",
		"i": "\u1d62",
		"j": "\u2c7c",
		"a": "\u2090",
		"e": "\u2091",
		"o": "\u2092",
		"x": "\u2093",
		"y": "\u1d67",
		"h": "\u2095",
		"k": "\u2096",
		"l": "\u2097",
		"m": "\u2098",
		"n": "\u2099",
		"p": "\u209a",
		"s": "\u209b",
		"t": "\u209c",
	}
	args := make([]string, 0, len(subscripts)*2)
	for k, v := range subscripts {
		args = append(args, k, v)
	}
	return strings.NewReplacer(args...)
}

func SupReplace() *strings.Replacer {
	superscripts := map[string]string{
		"0": "\u2070",
		"1": "\u00b9",
		"2": "\u00b2",
		"3": "\u00b3",
		"4": "\u2074",
		"5": "\u2075",
		"6": "\u2076",
		"7": "\u2077",
		"8": "\u2078",
		"9": "\u2079",
		"i": "\u2071",
		"n": "\u207f",
	}
	args := make([]string, 0, len(superscripts)*2)
	for k, v := range superscripts {
		args = append(args, k, v)
	}
	return strings.NewReplacer(args...)
}

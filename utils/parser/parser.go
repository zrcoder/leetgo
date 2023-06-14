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
		"\u00A0", " ",
		"\u200B", "",
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
		"a": "\u2090",
		"e": "\u2091",
		"h": "\u2095",
		"i": "\u1d62",
		"j": "\u2c7c",
		"k": "\u2096",
		"l": "\u2097",
		"m": "\u2098",
		"n": "\u2099",
		"o": "\u2092",
		"p": "\u209a",
		"r": "\u1d63",
		"s": "\u209b",
		"t": "\u209c",
		"u": "\u1d64",
		"v": "\u1d65",
		"x": "\u2093",
		"y": "\u1d67",
		"+": "\u208A",
		"-": "\u208B",
		"=": "\u208C",
		"(": "\u208D",
		")": "\u208E",
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
		"a": "\u1D43",
		"b": "\u1D47",
		"c": "\u1D9C",
		"d": "\u1D48",
		"e": "\u1D49",
		"f": "\u1DA0",
		"g": "\u1D4D",
		"h": "\u02B0",
		"i": "\u2071",
		"j": "\u02B2",
		"k": "\u1D4F",
		"l": "\u02E1",
		"m": "\u1D50",
		"n": "\u207F",
		"o": "\u1D52",
		"p": "\u1D56",
		"q": "\u02A0",
		"r": "\u02B3",
		"s": "\u02E2",
		"t": "\u1D57",
		"u": "\u1D58",
		"v": "\u1D5B",
		"w": "\u02B7",
		"x": "\u02E3",
		"y": "\u02B8",
		"z": "\u1DBB",
		"+": "\u207A",
		"-": "\u207B",
		"=": "\u207C",
		"(": "\u207D",
		")": "\u207E",
	}
	args := make([]string, 0, len(superscripts)*2)
	for k, v := range superscripts {
		args = append(args, k, v)
	}
	return strings.NewReplacer(args...)
}

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"strings"
	"text/template"
	"time"

	"golang.org/x/tools/imports"

	log "github.com/golang/glog"
)

var (
	builtins = map[string]interface{}{
		"date": func() string {
			return time.Now().Format(time.RFC1123Z)
		},
		"to_lower": func(s string) string {
			return strings.ToLower(s)
		},
	}
	header = []byte(`
// auto-generated
// do not modify this file by hand

`)
)

func main() {
	var (
		in, out, cfg string
		c            map[string]interface{}
	)
	flag.StringVar(&in, "in", "", "")
	flag.StringVar(&cfg, "cfg", "", "")
	flag.StringVar(&out, "out", "", "")
	flag.Parse()

	if err := json.Unmarshal([]byte(cfg), &c); err != nil {
		log.Fatal(err)
	}

	t, err := template.
		New("").
		Funcs(builtins).
		Option("missingkey=error").
		ParseFiles(in)
	if err != nil {
		log.Fatal(err)
	}
	t = t.Templates()[0]

	var b bytes.Buffer
	b.Write(header)
	err = t.Execute(&b, c)
	if err != nil {
		log.Fatal(err)
	}

	outb, err := imports.Process(out, b.Bytes(), nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(out, outb, 0644); err != nil {
		log.Fatal(err)
	}
}

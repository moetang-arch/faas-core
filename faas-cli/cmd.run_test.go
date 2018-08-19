package main

import (
	"bytes"
	"html/template"
	"testing"
)

func TestExample(t *testing.T) {
	tttt, err := template.New("generatedSource").Parse(generatedSourceTemplate)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	err = tttt.Execute(buf, "hello")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buf.String())
}

package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseHTML(t *testing.T) {
	page := Page{}
	err := page.Load(bufio.NewReader(strings.NewReader(`
	<html>
		<head>
			<style>EEEE</style>
		</head>
		<body>
			AAAA<br/>
			BBBB

			<p>CCCC</p>

			<script>
			DDDD
			</script>
		</body>
	</html>
	`)))

	if err != nil {
		t.Error(err)
	}

	buf, err := page.GetText()

	if err != nil {
		t.Error(err)
	}

	text := buf.String()
	if strings.Contains(text, "EEEE") ||
		strings.Contains(text, "DDDD") {
		t.Error("found style or script in text")
	}
	if !strings.Contains(text, "AAAA") ||
		!strings.Contains(text, "BBBB") ||
		!strings.Contains(text, "CCCC") {
		t.Error("some text missing")
	}
}

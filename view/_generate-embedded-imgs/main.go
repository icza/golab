package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"go/format"
	"io/ioutil"

	"github.com/icza/golab/engine"
)

func main() {
	// Generate embedded-imgs.go source file, holding a map from file names to their base64 encoded contents.
	// Images are looked for in ../_images/.

	buf := &bytes.Buffer{}
	buf.WriteString(`// This file is generated by go generate.

package view

// Embedded images mapped from image (file) name to file content encoded in Base64 format.
// Whether these are used depends on the useEmbeddedImages const in images.go.
var base64Imgs = map[string]string{
`)

	var names []string
	for dir := engine.Dir(0); dir < engine.DirCount; dir++ {
		// Gopher images
		names = append(names, fmt.Sprintf("gopher-%s.png", dir))
		// Bulldog images
		names = append(names, fmt.Sprintf("bulldog-%s.png", dir))
	}

	names = append(names, "wall.png")
	names = append(names, "gopher-dead.png")
	names = append(names, "door.png")
	names = append(names, "marker.png")
	names = append(names, "won.png")

	// Generate map entries
	for _, name := range names {
		data, err := ioutil.ReadFile("../_images/" + name)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(buf, "\t\"%s\": \"%s\",\n", name, base64.StdEncoding.EncodeToString(data))
	}

	buf.WriteString("}")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("embedded-imgs.go", formatted, 0664); err != nil {
		panic(err)
	}
}

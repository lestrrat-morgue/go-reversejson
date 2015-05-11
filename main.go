package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
	"reflect"
	"regexp"
	"sort"
)

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func camelCase(src string) string {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		chunks[idx] = bytes.Title(val)
	}
	return string(bytes.Join(chunks, nil))
}

func main() {
	name := flag.String("name", "MyStruct", "The name of the resulting structure")
	flag.Parse()

	dec := json.NewDecoder(os.Stdin)

	thing := make(map[string]interface{})
	if err := dec.Decode(&thing); err != nil {
		log.Fatalf("Failed to decode JSON: %s", err)
	}

	out := &bytes.Buffer{}
	encode(out, thing, "", *name, 0)
	os.Stdout.Write(out.Bytes())
}

func encode(out *bytes.Buffer, thing interface{}, prefix, name string, nest int) {
	switch thing.(type) {
	case map[string]interface{}:
		// OK, no op
	default:
		log.Fatalf("Unsupported type: %s", reflect.TypeOf(thing))
	}

	if nest == 0 {
		out.WriteString("type " + name + " struct {\n")
	} else {
		out.WriteString("struct {\n")
	}

	encodeMapElems(out, thing.(map[string]interface{}), prefix + "  ", name, nest)
	if nest == 0 {
		out.WriteString(prefix + "}\n")
	} else {
		out.WriteString(prefix + "}")
	}
}

func encodeMapElems(out *bytes.Buffer, thing map[string]interface{}, prefix, name string, nest int) {
	keys := []string{}
	for k := range thing {
		keys = append(keys, k)
	}

	sort.StringSlice(keys).Sort()

	for _, k := range keys {
		v := thing[k]
		out.WriteString(prefix + camelCase(k) + " ")
		switch v.(type) {
		case string:
			out.WriteString("string")
		case bool:
			out.WriteString("bool")
		case float64:
			// By spec, a number is a float... not much else we can do here
			out.WriteString("float64")
		case []interface{}:
			// Assume list of strings
			out.WriteString("[]string")
		default:
			encode(out, v, prefix, name, nest + 1)
		}

		out.WriteString(" `json:\"" + k + "\"`\n")
	}
}

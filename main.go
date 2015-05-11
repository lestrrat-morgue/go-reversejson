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
	out.WriteString("type " + *name + " struct {\n")
	encode(out, thing, "  ")
	out.WriteString("}\n")
	os.Stdout.Write(out.Bytes())
}

func encode(out *bytes.Buffer, thing interface{}, prefix string) {
	switch thing.(type) {
	case map[string]interface{}:
		encodeMap(out, thing.(map[string]interface{}), prefix)
	default:
		log.Fatalf("Unsupported type: %s", reflect.TypeOf(thing))
	}
}

func encodeMap(out *bytes.Buffer, thing map[string]interface{}, prefix string) {
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
			out.WriteString("struct {\n")
			encode(out, v, prefix+"  ")
			out.WriteString(prefix + "}")
		}

		out.WriteString(" `json:\"" + k + "\"`\n")
	}
}

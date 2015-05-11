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
	"strconv"
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
	encode(out, thing, "  ", *name, 0)
	out.WriteString("}\n")
	os.Stdout.Write(out.Bytes())
}

func encode(out *bytes.Buffer, thing interface{}, prefix, name string, counter int) {
	switch thing.(type) {
	case map[string]interface{}:
		encodeMap(out, thing.(map[string]interface{}), prefix, name, counter)
	default:
		log.Fatalf("Unsupported type: %s", reflect.TypeOf(thing))
	}
}

func encodeMap(out *bytes.Buffer, thing map[string]interface{}, prefix, name string, counter int) {
	if counter == 0 {
		out.WriteString(prefix + "type " + name + " struct {\n")
	} else {
		out.WriteString(prefix + "type " + name + strconv.Itoa(counter) + " struct {\n")
	}

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
			encode(out, v, prefix+"  ", name, counter+1)
		}

		out.WriteString(" `json:\"" + k + "\"`\n")
	}

	out.WriteString(prefix + "}\n")
}

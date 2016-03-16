package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	elastigo "github.com/mattbaird/elastigo/lib"
)

type indexByName []elastigo.CatIndexInfo

func (i indexByName) Len() int           { return len(i) }
func (i indexByName) Swap(j, k int)      { i[j], i[k] = i[k], i[j] }
func (i indexByName) Less(j, k int) bool { return strings.Compare(i[j].Name, i[k].Name) == -1 }

var c *elastigo.Conn

func getIndicesName() []string {
	fmt.Print("Loading index list....\n")
	indices := c.GetCatIndexInfo("")
	sort.Sort(indexByName(indices))

	names := make([]string, len(indices))
	for i, v := range indices {
		names[i] = v.Name
	}

	return names
}

func getTypes(index string) elastigo.Mapping {
	//result := map[string]elastigo.MappingOptions
	var result elastigo.Mapping
	res, _ := c.DoCommand("GET", fmt.Sprintf("/%s/_mapping", index), nil, nil)
	var test map[string]json.RawMessage
	json.Unmarshal(res, &test)
	for _, mappings := range test {
		var mapping map[string]elastigo.Mapping
		json.Unmarshal(mappings, &mapping)

		result = mapping["mappings"]
	}
	return result
}

func selectIndicesToFix(indices []string) []string {
	toFix := make([]string, 0)
	for pos, index := range indices {
		fmt.Print("\r", pos+1, "/", len(indices), " : ", index)
		//fmt.Println("Index : ", index.Name, " ", index.Docs.Count, " docs / ", index.Store.Size, " bytes (", index.Status, ")")
		//fmt.Println("Mappings : ")
		types := getTypes(index)

		for docType, mapping := range types {
			for propertyName, _ := range mapping.Properties {
				//fmt.Println(docType, ".", propertyName)
				matched, err := regexp.MatchString("\\.", propertyName)
				if err != nil {
					fmt.Println(err)
				}
				if matched {
					fmt.Print(" property ", docType, ".", propertyName, " not compatible\n")
					toFix = append(toFix, index)
				}
			}
		}
	}
	fmt.Println()
	return toFix
}

func main() {
	c = elastigo.NewConn()
	c.Domain = "localhost"
	c.Port = "9300"

	indices := getIndicesName()

	toFix := selectIndicesToFix(indices)
	fmt.Println()
	fmt.Println("Indices to fix : ", strings.Join(toFix, " "))
}

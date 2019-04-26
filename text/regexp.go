package main

import (
	"fmt"
	"regexp"
)

func toMap (texts []string) map[string] int{
	m := make(map[string] int)
	
	for index,text := range texts {
		m[text] = index
	}
	return m
}

func main() {
	re := regexp.MustCompile("CLUSTER_(?P<port>\\d+)_URLPREFIX")
	groupNames := re.SubexpNames()
	groupIndexLookup := toMap ( groupNames )
	submatches := re.FindStringSubmatch("CLUSTER_80_URLPREFIX")
	
	fmt.Println(submatches[groupIndexLookup["port"]])
	
}

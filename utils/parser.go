package utils

import (
	"net/http"
	"regexp"
	"strconv"
)

// Regexp Pattern definition for:
// 	1. Structured Type
// 	2. Nested Structured Type
var (
	reMap   = regexp.MustCompile(`^([a-zA-Z_0-9]+)\[([a-zA-Z_0-9]+)\]$`)
	reArray = regexp.MustCompile(`^([a-zA-Z_0-9]+)\[([0-9]+)\]$`)

	reArrayOfMap = regexp.MustCompile(`^([a-zA-Z_0-9]+)\[([0-9]+)\]\[([a-zA-Z_0-9]+)\]$`)
	reMapOfArray = regexp.MustCompile(`^([a-zA-Z_0-9]+)\[([a-zA-Z_0-9]+)\]\[([0-9]+)\]$`)
)

// ParseFormData parses http request's FormData and returns a map[string]any.
// 	1. Notice that this is a PRELIMIARY implementation of a FormData parser,
// 		and `Content-Type:application/json` is a way better alternative.
// 	2. For structured type, parser ONLY supports array, map(objcet), array of map(array of object),
//		and map of array(object of which each fields is an array).
//		Allow letter, number, and underscore (`[a-zA-Z_0-9]+`) for the map key.
// 		Any other characters will cause the structured type being parsed as a basic type.
//	3. Any further nested structured type will also get parsed, however as a basic type.
// 		In other words, complex nested structured type MUST be passed by `Content-Type:application/json`.
func ParseFormData(r *http.Request) (ret map[string]any) {
	// init
	ret = map[string]any{}
	r.ParseForm()

	// iterate through all Form fields
	for key := range r.Form {
		matchesArray := reArray.FindStringSubmatch(key)
		matchesMap := reMap.FindStringSubmatch(key)
		matchesArrayOfMap := reArrayOfMap.FindStringSubmatch(key)
		matchesMapOfArray := reMapOfArray.FindStringSubmatch(key)

		// match for one and only one pattern
		switch {
		case len(matchesMapOfArray) == 4:
			mmkey := matchesMapOfArray[1] // main map key (ret)
			smkey := matchesMapOfArray[2] // sub map key
			index, _ := strconv.Atoi(matchesMapOfArray[3])
			// init
			var result map[string][]string
			if ret[mmkey] != nil {
				result = ret[mmkey].(map[string][]string)
			} else {
				result = make(map[string][]string)
			}
			// populate
			for index >= len(result[smkey]) {
				if result[smkey] == nil {
					result[smkey] = make([]string, index)
				}
				result[smkey] = append(result[smkey], "")
			}
			// save
			result[smkey][index] = r.Form[key][0]
			ret[mmkey] = result

		// matched as array of map (array of object)
		case len(matchesArrayOfMap) == 4:
			mmkey := matchesArrayOfMap[1] // main map key (ret)
			index, _ := strconv.Atoi(matchesArrayOfMap[2])
			smkey := matchesArrayOfMap[3] // sub map key

			// init
			var result []map[string]string
			if ret[mmkey] != nil {
				result = ret[mmkey].([]map[string]string)
			} else {
				result = make([]map[string]string, 0)
			}
			// populate
			for index >= len(result) {
				result = append(result, map[string]string{})
			}
			// save
			result[index][smkey] = r.Form[key][0]
			ret[mmkey] = result

		// matched as array
		case len(matchesArray) == 3:
			mmkey := matchesArray[1] // main map key (ret)
			index, _ := strconv.Atoi(matchesArray[2])
			// init
			var result []string
			if ret[mmkey] != nil {
				result = ret[mmkey].([]string)
			} else {
				result = make([]string, 0)
			}
			// populate
			for index >= len(result) {
				result = append(result, "")
			}
			// save
			result[index] = r.Form[key][0]
			ret[mmkey] = result

		// matched as map (object)
		case len(matchesMap) == 3:
			mmkey := matchesMap[1] // main map key (ret)
			smkey := matchesMap[2] // sub map key
			// init
			var result map[string]string
			if ret[mmkey] != nil {
				result = ret[mmkey].(map[string]string)
			} else {
				result = make(map[string]string)
			}
			// save
			result[smkey] = r.Form[key][0]
			ret[mmkey] = result

		// direct and basic type of k-v pairs
		default:
			ret[key] = r.Form[key][0]
		}

	}

	// RET
	return ret
}

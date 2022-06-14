package utils

import (
	"net/http"
	"regexp"
	"strconv"
)

// Regexp Pattern definition for:
// 	1. Structured Type: array, map(object)
// 	2. Nested Structured Type: array of map(array of object), map of array(object of which each fields is an array).
var (
	reMap   = regexp.MustCompile(`^([a-zA-Z_]+)\[([a-zA-Z_]+)\]$`)
	reArray = regexp.MustCompile(`^([a-zA-Z_]+)\[([0-9]+)\]$`)

	reArrayOfMap = regexp.MustCompile(`^([a-zA-Z_]+)\[([0-9]+)\]\[([a-zA-Z_]+)\]$`)
	reMapOfArray = regexp.MustCompile(`^([a-zA-Z_]+)\[([a-zA-Z_]+)\]\[([0-9]+)\]$`)
)

// ParseFormData parses http request's FormData and returns a map[string]any.
// 	1. Notice that this is a PRELIMIARY implementation of a FormData parser,
// 		and passing parameters through `Content-Type:application/json` is a way better alternative.
// 	2. For structured type, parser ONLY supports array, map(objcet),
//		array of map(array of object), and map of array(object of which each fields is an array).
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

		switch {
		case len(matchesMapOfArray) == 4:
			mainMapKey := matchesMapOfArray[1]
			subMapKey := matchesMapOfArray[2]
			index, _ := strconv.Atoi(matchesMapOfArray[3])
			var result map[string][]string
			if ret[mainMapKey] != nil {
				result = ret[mainMapKey].(map[string][]string)
			} else {
				result = make(map[string][]string)
			}
			for index >= len(result[subMapKey]) {
				if result[subMapKey] == nil {
					result[subMapKey] = make([]string, index)
				}
				result[subMapKey] = append(result[subMapKey], "")
			}
			result[subMapKey][index] = r.Form[key][0]
			ret[mainMapKey] = result

		case len(matchesArrayOfMap) == 4:
			mainMapKey := matchesArrayOfMap[1]
			index, _ := strconv.Atoi(matchesArrayOfMap[2])
			subMapKey := matchesArrayOfMap[3]

			var result []map[string]string
			if ret[mainMapKey] != nil {
				result = ret[mainMapKey].([]map[string]string)
			} else {
				result = make([]map[string]string, 0)
			}
			for index >= len(result) {
				result = append(result, map[string]string{})
			}
			result[index][subMapKey] = r.Form[key][0]
			ret[mainMapKey] = result

		case len(matchesArray) == 3:
			mainMapKey := matchesArray[1]
			index, _ := strconv.Atoi(matchesArray[2])
			var result []string
			if ret[mainMapKey] != nil {
				result = ret[mainMapKey].([]string)
			} else {
				result = make([]string, 0)
			}
			for index >= len(result) {
				result = append(result, "")
			}
			result[index] = r.Form[key][0]
			ret[mainMapKey] = result

		case len(matchesMap) == 3:
			mainMapKey := matchesMap[1]
			subMapKey := matchesMap[2]
			var result map[string]string
			if ret[mainMapKey] != nil {
				result = ret[mainMapKey].(map[string]string)
			} else {
				result = make(map[string]string)
			}
			result[subMapKey] = r.Form[key][0]
			ret[mainMapKey] = result

		default:
			ret[key] = r.Form[key][0]
		}

	}

	// RET
	return ret
}

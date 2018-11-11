package formatter

import (
	"fmt"
	"strings"
)

var formatType = map[string]int{
	"bool":         35,
	"key":          33,
	"num":          36,
	"objBracket":   1,
	"sliceBracket": 1,
	"string":       32,
}

// JSONToString formats and returns string converted JSON object
func JSONToString(rowData interface{}, formatting bool) string {
	var str string
	switch data := rowData.(type) {
	case int64, float64:
		str = numberToString(rowData, formatting)
	case bool:
		str = boolToString(rowData, formatting)
	case string:
		str = stringToString(rowData, formatting)
	case []interface{}:
		str = sliceToString(data, formatting, 1)
	case map[string]interface{}:
		str = objToString(data, formatting, 1)
	}
	return strings.TrimSuffix(str, ",")
}

func numberToString(num interface{}, formatting bool) string {
	return addComma(format(num, formatType["num"], formatting))
}

func stringToString(str interface{}, formatting bool) string {
	return addComma(format(fmt.Sprintf(`"%v"`, str), formatType["string"], formatting))
}

func boolToString(b interface{}, formatting bool) string {
	return addComma(format(b, formatType["bool"], formatting))
}

// sliceToString formats and returns string converted from Slice
func sliceToString(data []interface{}, formatting bool, nest int) string {
	lines := []string{format("[", formatType["sliceBracket"], formatting)}
	for _, value := range data {
		var str string
		switch convertedValue := value.(type) {
		case int64, float64:
			str = addSpace(numberToString(value, formatting), nest)
		case bool:
			str = addSpace(boolToString(value, formatting), nest)
		case string:
			str = addSpace(stringToString(value, formatting), nest)
		case []interface{}:
			str = sliceToString(convertedValue, formatting, nest+1)
		case map[string]interface{}:
			str = addSpace(objToString(convertedValue, formatting, nest+1), nest)
		}
		lines = append(lines, str)
	}
	lines[len(lines)-1] = strings.TrimSuffix(lines[len(lines)-1], ",")
	lines = append(lines, addSpace(addComma(format("]", formatType["sliceBracket"], formatting)), nest-1))
	return strings.Join(lines, "\n")
}

// objToString formats and returns string converted from Map
func objToString(data map[string]interface{}, formatting bool, nest int) string {
	lines := []string{format("{", formatType["objBracket"], formatting)}
	for key, value := range data {
		var str string
		fmtkey := format(fmt.Sprintf(`"%s"`, key), formatType["key"], formatting)
		switch convertedValue := value.(type) {
		case int64, float64:
			str = numberToString(value, formatting)
		case bool:
			str = boolToString(value, formatting)
		case string:
			str = stringToString(value, formatting)
		case []interface{}:
			str = sliceToString(convertedValue, formatting, nest+1)
		case map[string]interface{}:
			str = objToString(convertedValue, formatting, nest+1)
		}
		lines = append(lines, addSpace(fmt.Sprintf("%s: %s", fmtkey, str), nest))
	}
	lines[len(lines)-1] = strings.TrimSuffix(lines[len(lines)-1], ",")
	lines = append(lines, addSpace(addComma(format("}", formatType["objBracket"], formatting)), nest-1))
	return strings.Join(lines, "\n")
}

func format(s interface{}, typeNum int, formatting bool) string {
	if !formatting {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", typeNum, s)
}

func addSpace(str string, num int) string {
	return strings.Repeat("  ", num) + str
}

func addComma(str string) string {
	return str + ","
}

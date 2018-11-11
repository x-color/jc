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

// dataToString formats and returns string converted from all data type
func dataToString(data interface{}, formatting bool, nest int) string {
	switch convertedData := data.(type) {
	case int64, float64:
		return numberToString(data, formatting)
	case bool:
		return boolToString(data, formatting)
	case string:
		return stringToString(data, formatting)
	case []interface{}:
		return sliceToString(convertedData, formatting, nest)
	case map[string]interface{}:
		return objToString(convertedData, formatting, nest)
	}
	return ""
}

// JSONToString formats and returns string converted JSON object
func JSONToString(rowData interface{}, formatting bool) string {
	return strings.TrimSuffix(dataToString(rowData, formatting, 1), ",")
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
		lines = append(lines, addSpace(dataToString(value, formatting, nest+1), nest))
	}
	lines[len(lines)-1] = strings.TrimSuffix(lines[len(lines)-1], ",")
	lines = append(lines, addSpace(addComma(format("]", formatType["sliceBracket"], formatting)), nest-1))
	return strings.Join(lines, "\n")
}

// objToString formats and returns string converted from Map
func objToString(data map[string]interface{}, formatting bool, nest int) string {
	lines := []string{format("{", formatType["objBracket"], formatting)}
	for key, value := range data {
		fmtkey := format(fmt.Sprintf(`"%s"`, key), formatType["key"], formatting)
		fmtValue := dataToString(value, formatting, nest+1)
		lines = append(lines, addSpace(fmt.Sprintf("%s: %s", fmtkey, fmtValue), nest))
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

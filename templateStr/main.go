package templateStr

import (
	"bytes"
	"fmt"
	"reflect"
	regex "regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
	// "golang.org/x/text/cases"
)

type Any = interface{}
type VariableMap map[string]Any
type Func func([]Any) string
type FuncArray []Func

var regVariable = regex.MustCompile(REG_VARIABLE)
var regFunction = regex.MustCompile(REG_FUNCTION)
var regCondition = regex.MustCompile(REG_CONDITION)
var regSwitch = regex.MustCompile(REG_SWITCH)
var regTyping = regex.MustCompile(REG_STR + `|` + REG_BOOL + `|` + REG_INT + `|` + REG_FLOAT + `|` + REG_VAR + `|` + REG_LIST)

// Construtor
type TemplateStr struct {
	variableMap VariableMap
	funcArray   FuncArray
}

// `funcArray FuncArray` is a array of custom functions that can be used when you call a function with: `{{@myCustomFunction}}`
//
// `variableMap VariableMap` is a map of the values you want to use when you call: `{{$myVar}}`
//
// Typing:
//
//	keyVariable  : is the key of the value in the dictionary pass to the constructor (return the value)
//	<b:True>     : bool    (return true)
//	<n:123>      : int     (return 123)
//	<n:123.4>    : float64 (return 123.4)
//	"text"       : string  (return text)
func New(funcArray FuncArray, variableMap VariableMap) TemplateStr {

	return TemplateStr{
		variableMap,
		funcArray,
	}
}

// Function utility

func getNameFunc(function Func) string {
	spli := strings.Split(runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name(), ".")
	return spli[len(spli)-1]
}

func getVariable(key string, varMap VariableMap, index ...int) (Any, error) {

	var fvalue Any
	var ok bool

	if strings.Contains(key, ".") && !strings.Contains(key, " ") {
		var tempMap VariableMap

		keyArray := strings.Split(key, ".")
		lenKeyArray := len(keyArray)

		for index, keyMap := range keyArray {

			if index == 0 {
				tempMap, ok = varMap[keyMap].(VariableMap)
			} else if index == lenKeyArray-1 {
				fvalue, ok = tempMap[keyMap]
			} else {
				tempMap, ok = tempMap[keyMap].(VariableMap)
			}
			if !ok {
				return "", NotFoundVariableError("[key '[" + key + "]' not exist]")
			}
		}
	} else {
		fvalue, ok = varMap[key]
		if !ok {
			return "", NotFoundVariableError("[key '[" + key + "]' not exist]")
		}
	}

	if fvalue == nil {
		return "none", nil
	}

	if lenIndex := len(index); lenIndex != 0 && (reflect.TypeOf(fvalue).Kind() == reflect.Slice || reflect.TypeOf(fvalue).Kind() == reflect.Array) {
		if lenfvalue := len(fvalue.([]Any)); !(lenfvalue <= index[0]) {
			fvalue = fvalue.([]Any)[index[0]]
		} else {
			return "", IndexError("[index '[" + fmt.Sprint(index[0]) + "]' out of range]")
		}
	} else if lenIndex != 0 {
		return "", NotAArrayError("[key '[" + key + "]' is not array]")
	}

	return fvalue, nil
}

func checkExistFuncStr(functionArray FuncArray, compareStr string) (bool, int, string) {

	for index, function := range functionArray {
		if nameFunc := getNameFunc(function); nameFunc == compareStr {
			return true, index, nameFunc
		}
	}
	return false, 0, "None"
}

func swapCase(str string) string {
	b := new(bytes.Buffer)

	for _, elem := range str {
		if unicode.IsUpper(elem) {
			b.WriteRune(unicode.ToLower(elem))
		} else {
			b.WriteRune(unicode.ToUpper(elem))
		}
	}

	return b.String()
}

func upperCaseFirst(str string) string {

	arrayCase := strings.Split(fmt.Sprintf("%v", str), "")
	arrayCase[0] = strings.ToUpper(arrayCase[0])
	return strings.Join(arrayCase, "")
}

func findAllGroup(reg *regex.Regexp, str string) []map[string]string {

	group := reg.SubexpNames()
	arrayMatch := reg.FindAllStringSubmatch(str, -1)

	arrayMap := []map[string]string{}
	for _, valueArray := range arrayMatch {

		mapMatch := map[string]string{}
		mapMatch["match"] = valueArray[0]

		for i := 1; i < len(valueArray); i++ {

			mapMatch[group[i]] = valueArray[i]
		}

		arrayMap = append(arrayMap, mapMatch)
	}

	return arrayMap
}

func convertInterfaceToFloat(value1 Any, value2 Any) (value1F, value2F float64) {

	var b2i = map[bool]int8{false: 0, true: 1}

	switch value1 := value1.(type) {
	case int:
		value1F = float64(value1)
	case float64:
		value1F = value1
	case bool:
		value1F = float64(b2i[value1])
	case string:
		value1F = float64(len(value1))
	default:
		value1F = 0
	}

	switch value2 := value2.(type) {
	case int:
		value2F = float64(value2)
	case float64:
		value2F = value2
	case bool:
		value2F = float64(b2i[value2])
	case string:
		value2F = float64(len(value2))
	default:
		value2F = 0
	}

	return
}

func typing(str string, varMap VariableMap, types ...string) []Any {

	arrayTyping := []Any{}

	if len(types) == 0 {
		for _, groupParam := range findAllGroup(regTyping, str) {

			if groupParam["str_double"] != "" {
				arrayTyping = append(arrayTyping, groupParam["str_double"])
			} else if groupParam["str_single"] != "" {
				arrayTyping = append(arrayTyping, groupParam["str_single"])
			} else if groupParam["str_back"] != "" {
				arrayTyping = append(arrayTyping, groupParam["str_back"])
			} else if groupParam["bool"] != "" {
				bool, _ := strconv.ParseBool(groupParam["bool"])
				arrayTyping = append(arrayTyping, bool)
			} else if groupParam["int"] != "" {
				int, _ := strconv.Atoi(groupParam["int"])
				arrayTyping = append(arrayTyping, int)
			} else if groupParam["float"] != "" {
				float, _ := strconv.ParseFloat(groupParam["float"], 64)
				arrayTyping = append(arrayTyping, float)
			} else if groupParam["variable"] != "" {
				if groupParam["index"] == "" {
					value, _ := getVariable(groupParam["variable"], varMap)
					arrayTyping = append(arrayTyping, value)
				} else {
					i, _ := strconv.Atoi(groupParam["index"])
					value, _ := getVariable(groupParam["variable"], varMap, i)
					arrayTyping = append(arrayTyping, value)
				}
			} else if l := groupParam["list"]; l != "" {
				arrayTyping = append(arrayTyping, typing(l, varMap))
			}
		}
	} else if types[0] == "int" {
		int, _ := strconv.Atoi(str)
		arrayTyping = append(arrayTyping, int)
	} else if types[0] == "float" {
		float, _ := strconv.ParseFloat(str, 64)
		arrayTyping = append(arrayTyping, float)
	} else if types[0] == "str" {
		arrayTyping = append(arrayTyping, str)
	} else if types[0] == "bool" {
		bool, _ := strconv.ParseBool(str)
		arrayTyping = append(arrayTyping, bool)
	}

	return arrayTyping
}

func ternary(cond bool, val1 string, val2 string) string {
	if cond {
		return val1
	}
	return val2
}

// Method TemplateStr

// shortcuts to run all parsers
//
// return -> string
func (t TemplateStr) Parse(text string) (string, error) {

	var err error

	for t.HasOne(text) {


		text, err = t.ParseVariable(text)
		if err != nil {return "", err}

		text, err = t.ParseFunction(text)
		if err != nil {return "", err}

		text, err = t.ParseCondition(text)
		if err != nil {return "", err}

		text, err = t.ParseSwitch(text)
		if err != nil {return "", err}
	}

	return text, nil
}

// parse all the `${variable}` in the text give in
//
// return -> string
func (t TemplateStr) ParseVariable(text string) (string, error) {

	if !t.HasVariable(text) {
		return text, nil
	}

	for t.HasVariable(text) {
		for _, v := range findAllGroup(regVariable, text) {

			variable := v["variable"]
			var value Any
			var err error

			if index := v["index"]; index == "" {
				value, err = getVariable(variable, t.variableMap)
			} else {
				index, _ := strconv.Atoi(index)
				value, err = getVariable(variable, t.variableMap, index)
			}

			if err != nil {return "", err}

			key := fmt.Sprintf("%v", value)
			match := v["match"]

			text = strings.Replace(text, match, key, -1)
		}
	}

	return text, nil
}

// parse all the `{{@function param1 param2}}` or `{{@function}}` in the text give in
//
// return -> string
func (t TemplateStr) ParseFunction(text string) (string, error) {

	if !t.HasFunction(text) {
		return text, nil
	}

	// c := cases.Fold()

	for t.HasFunction(text) {

		for _, group := range findAllGroup(regFunction, text) {

			match := group["match"]
			parameters := group["parameters"]

			// var value string = "None"
			dateTime := time.Now()

			value := func(par string) (string, error) {

				v, err := getVariable(par, t.variableMap)
				if err == nil {
					if fmt.Sprint(v) != "" {
						return fmt.Sprint(v), nil
					}
				}

				return "", err
				
			}

			functionName := group["functionName"]

			switch functionName {
			case "uppercase":
				if v, err := value(parameters); err == nil {text = strings.Replace(text, match, strings.ToUpper(v), -1)} else { return "", err }
			case "uppercaseFirst":
				if v, err := value(parameters); err == nil {text = strings.Replace(text, match, upperCaseFirst(v), -1)} else { return "", err }
			case "lowercase":
				if v, err := value(parameters); err == nil {text = strings.Replace(text, match, strings.ToLower(v), -1)} else { return "", err }
			// case "casefold": text = strings.Replace(text, match, c.String(key), -1)
			case "swapcase":
				if v, err := value(parameters); err == nil {text = strings.Replace(text, match, swapCase(v), -1)} else { return "", err }
			case "time":
				text = strings.Replace(text, match, dateTime.Format("15:04:05"), -1)
			case "date":
				text = strings.Replace(text, match, dateTime.Format("02/01/2006"), -1)
			case "dateTime":
				text = strings.Replace(text, match, dateTime.Format("02/01/2006 15:04:05"), -1)
			default:
				if ok, index, customFuncstr := checkExistFuncStr(t.funcArray, functionName); ok {
					customFunc := t.funcArray[index]

					if functionName == customFuncstr {
						var resultTextfunc string
						if parameters != "" {
							resultTextfunc = customFunc(typing(parameters, t.variableMap))
						} else {
							resultTextfunc = customFunc([]Any{})
						}
						text = strings.Replace(text, match, resultTextfunc, -1)
					}
				} else {
					return "", NotFoundFunctionError("[Function " + functionName + " not exist]")
				}
			}
		}
	}

	return text, nil
}

// parse all the `{{#var1 == var2: value1 || value2}}` in the text give in
//
// return -> string
func (t TemplateStr) ParseCondition(text string) (string, error) {

	if !t.HasCondition(text) {
		return text, nil
	}

	for t.HasCondition(text) {

		for _, group := range findAllGroup(regCondition, text) {

			match := group["match"]
			conditionValue1 := group["conditionValue1"]
			conditionValue2 := group["conditionValue2"]
			conditionSymbol := group["conditionSymbol"]
			trueValue := group["trueValue"]
			falseValue := group["falseValue"]

			ArrayTyping := typing(conditionValue1+" "+conditionValue2, t.variableMap)

			if conditionSymbol == "==" {
				text = strings.Replace(text, match, ternary(ArrayTyping[0] == ArrayTyping[1], trueValue, falseValue), -1)
			} else if conditionSymbol == "!=" {
				text = strings.Replace(text, match, ternary(ArrayTyping[0] != ArrayTyping[1], trueValue, falseValue), -1)
			} else {
				v1, v2 := convertInterfaceToFloat(ArrayTyping[0], ArrayTyping[1])
				if conditionSymbol == "<=" {
					text = strings.Replace(text, match, ternary(v1 <= v2, trueValue, falseValue), -1)
				} else if conditionSymbol == ">=" {
					text = strings.Replace(text, match, ternary(v1 >= v2, trueValue, falseValue), -1)
				} else if conditionSymbol == "<" {
					text = strings.Replace(text, match, ternary(v1 < v2, trueValue, falseValue), -1)
				} else if conditionSymbol == ">" {
					text = strings.Replace(text, match, ternary(v1 > v2, trueValue, falseValue), -1)
				} else {
					return "", BadComparatorError("["+conditionSymbol+" is not valid comparator]")
				}
			}
		}
	}

	return text, nil
}

// parse all the `{{?var; value1=#0F0, 56=#00F, ..., default=#000}}` or
// `{{?var:int; 56=#0F0, 32=#00F, ..., default=#000}}` in the text give in
//
// return -> string
func (t TemplateStr) ParseSwitch(text string) (string, error) {

	if !t.HasSwitch(text) {
		return text, nil
	}

	for t.HasSwitch(text) {

		for _, group := range findAllGroup(regSwitch, text) {

			match := group["match"]

			mapTemp := map[string]string{}
			var result string

			for _, n := range strings.Split(group["values"], ", ") {
				keyValue := strings.Split(n, "::")
				mapTemp[keyValue[0]] = keyValue[1]
			}

			if group["type"] == "" {
				for key, value := range mapTemp {
					if key == t.variableMap[group["variable"]] {
						result = value
						break
					} else {
						result = group["defaultValue"]
					}
				}

			} else {
				keyVar := group["variable"]
				typeVar := group["type"]

				for key, value := range mapTemp {
					// println(fmt.Sprintf("%T", typing(key, t.variableMap, typeVar)[0]))
					if valVar, _ := getVariable(keyVar, t.variableMap); typing(key, t.variableMap, typeVar)[0] == valVar {
						result = value
						break
					} else {
						result = group["defaultValue"]
					}
				}
			}
			text = strings.Replace(text, match, result, -1)
		}
	}

	return text, nil
}

// Detects if there is the presence of min one syntaxe
//
// return -> bool
func (t TemplateStr) HasOne(text string) bool {

	if t.HasVariable(text) || t.HasFunction(text) || t.HasCondition(text) || t.HasSwitch(text) {
		return true
	}
	return false
}

// Detects if there is the presence of `{{$variable}}`
//
// return -> bool
func (t TemplateStr) HasVariable(text string) bool {
	return regVariable.MatchString(text)
}

// Detects if there is the presence of `{{@function param1 param2}}` or `{{@function}}`
//
// return -> bool
func (t TemplateStr) HasFunction(text string) bool {
	return regFunction.MatchString(text)
}

// Detects if there is the presence of `{{#var1 == var2: value1 || value2}}`
//
// return -> bool
func (t TemplateStr) HasCondition(text string) bool {
	return regCondition.MatchString(text)
}

// Detects if there is the presence of `{{?var: value1=#0F0, value2=#00F, ..., default=#000}}` or
// `{{?var:int; 56=#0F0, 32=#00F, ..., default=#000}}`
//
// return -> bool
func (t TemplateStr) HasSwitch(text string) bool {
	return regSwitch.MatchString(text)
}

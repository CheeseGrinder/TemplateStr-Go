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

    "golang.org/x/text/cases"
)

type Any = interface{}
type VariableMap map[string]Any
type Func func([]Any) string
type FuncArray []Func

var regVariable = regex.MustCompile(`(?P<match>{{\$(?P<key>[^{{$}}]+)}})`)
var regFunction = regex.MustCompile(`(?P<match>{{@(?P<function>[^{@}\s]+) ?(?P<key>[^{@}]+)?}})`)
var regCondition = regex.MustCompile(`(?P<match>{{#(?P<compValue1>[^{#}]+) (?P<compSymbol>[=!<>][=]?) (?P<compValue2>[^{#}]+): (?P<resultValue1>[^{}]+) \|\| (?P<resultValue2>[^{}]+)}})`)
var regSwitch = regex.MustCompile(`(?P<match>{{\?(?:(?P<key>[^{?}:]+)|(?P<keyTyped>[^{?}]+):(?P<type>str|int|float)); (?P<val>(?:[^{}]+)=(?:[^{}]+)), default=(?P<default>[^{}]+)}})`)
var regTyping = regex.MustCompile(`\"(?P<str_double>[^\"]+)\"|\'(?P<str_single>[^\']+)\'|\x60(?P<str_back>[^\x60]+)\x60|<b:(?P<bool>True|False)>|<n:(?P<number>[0-9_.]+)>|(?P<variable>[^<>\" ]+)`)

// Construtor
type TemplateStr struct {
    variableMap VariableMap
    funcArray FuncArray
}

// `funcArray FuncArray` is a array of custom functions that can be used when you call a function with: `{{@myCustomFunction}}`
//
// `variableMap VariableMap` is a map of the values you want to use when you call: `{{$myVar}}`
//
// Typing:
//     keyVariable  : is the key of the value in the dictionary pass to the constructor (return the value)
//     <b:True>     : bool    (return true)
//     <n:123>      : int     (return 123)
//     <n:123.4>    : float64 (return 123.4)
//     "text"       : string  (return text)
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

func getVariable(key string, varMap VariableMap, ) (Any, bool) {

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
        }
    } else {
        fvalue, ok = varMap[key]
    }

    if fvalue == nil { fvalue = "None" }

    return fvalue, ok
}

func checkExistFuncStr(functionArray FuncArray, compareStr string) bool {

    for _, function := range functionArray {
        if getNameFunc(function) == compareStr {
            return true
        }
    }
    return false
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
    return strings.Join(arrayCase,"")
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

func typing(str string, varMap VariableMap, typing ...string) []Any {

    arrayTyping := []Any{}

    if len(typing) == 0 {
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
            } else if groupParam["number"] != "" {
                if !strings.Contains(groupParam["number"], ".") {
                    int, _ := strconv.Atoi(groupParam["number"])
                    arrayTyping = append(arrayTyping, int)
                } else {
                    float, _ := strconv.ParseFloat(groupParam["number"], 64)
                    arrayTyping = append(arrayTyping, float)
                }
            } else if groupParam["variable"] != "" {
                value, _ := getVariable(groupParam["variable"], varMap)
                arrayTyping = append(arrayTyping, fmt.Sprintf("%v", value))
            }
        }
    } else if typing[0] == "int" {
        int, _ := strconv.Atoi(str)
        arrayTyping = append(arrayTyping, int)
    } else if typing[0] == "float" {
        float, _ := strconv.ParseFloat(str, 64)
        arrayTyping = append(arrayTyping, float)
    } else if typing[0] == "str" {
        arrayTyping = append(arrayTyping, str)
    } else if typing[0] == "bool" {
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
func (t TemplateStr) Parse(text string) string {

    text = t.ParseVariable(text)

    text = t.ParseFunction(text)

    text = t.ParseCondition(text)

    text = t.ParseSwitch(text)

    return text
}

// parse all the `{{$variable}}` in the text give in
//
// return -> string
func (t TemplateStr) ParseVariable(text string) string {

    if !t.HasVariable(text) { return text }

    for _, v := range findAllGroup(regVariable, text) {

        value, _ := getVariable(v["key"], t.variableMap)

        key := fmt.Sprintf("%v", value)
        match := v["match"]

        text = strings.Replace(text, match, key, -1)
    }

    return text
}

// parse all the `{{@function param1 param2}}` or `{{@function}}` in the text give in
//
// return -> string
func (t TemplateStr) ParseFunction(text string) string {

    if !t.HasFunction(text) { return text }

    c := cases.Fold()

    for _, group := range findAllGroup(regFunction, text) {

        match := group["match"]
        
        var key string
        dateTime := time.Now()

        if value, ok := getVariable(group["key"], t.variableMap); ok && fmt.Sprintf("%v", value) != ""{
            key = fmt.Sprintf("%v", value)
        } else {
            key = "None"
        }

        functionName := group["function"]

        switch functionName {
        case "uppercase": text = strings.Replace(text, match, strings.ToUpper(key), -1)
        case "uppercaseFirst": text = strings.Replace(text, match, upperCaseFirst(key), -1)
        case "lowercase": text = strings.Replace(text, match, strings.ToLower(key), -1)
        case "casefold": text = strings.Replace(text, match, c.String(key), -1)
        case "swapcase": text = strings.Replace(text, match, swapCase(key), -1)
        case "time": text = strings.Replace(text, match, dateTime.Format("15:04:05"), -1)
        case "date": text = strings.Replace(text, match, dateTime.Format("02/01/2006"), -1)
        case "dateTime": text = strings.Replace(text, match, dateTime.Format("02/01/2006 15:04:05"), -1)
        default:
            if checkExistFuncStr(t.funcArray, functionName) {

                for _, customFunc := range t.funcArray {
    
                    customFuncstr := getNameFunc(customFunc)
    
                    if functionName == customFuncstr{
                        var resultTextfunc string
                        
                        if group["key"] != "" {
                            
                            resultTextfunc = customFunc(typing(group["key"], t.variableMap))
    
                        } else {
                            resultTextfunc = customFunc([]Any{})
                        }
    
                        text = strings.Replace(text, match, resultTextfunc, -1)
                    } 
                }
            } else {
                text = "(NoFunction " + functionName + ")"
            }
        }
    }

    return text
}

// parse all the `{{#var1 == var2: value1 || value2}}` in the text give in
//
// return -> string
func (t TemplateStr) ParseCondition(text string) string {

    if !t.HasCondition(text) { return text }

    for _, group := range findAllGroup(regCondition, text) {

        match := group["match"]
        compValue1 := group["compValue1"]
        compValue2 := group["compValue2"]
        compSymbol := group["compSymbol"]
        resultValue1 := group["resultValue1"]
        resultValue2 := group["resultValue2"]

        ArrayTyping := typing(compValue1 + " " + compValue2, t.variableMap)
        
        if compSymbol == "==" {
            text = strings.Replace(text, match, ternary(ArrayTyping[0] == ArrayTyping[1], resultValue1, resultValue2), -1)
        } else if compSymbol == "!=" {
            text = strings.Replace(text, match, ternary(ArrayTyping[0] != ArrayTyping[1], resultValue1, resultValue2), -1)
        } else {
            v1, v2 := convertInterfaceToFloat(ArrayTyping[0], ArrayTyping[1])
            if compSymbol == "<=" {
                text = strings.Replace(text, match, ternary(v1 <= v2, resultValue1, resultValue2), -1)
            } else if compSymbol == ">=" {
                text = strings.Replace(text, match, ternary(v1 >= v2, resultValue1, resultValue2), -1)
            } else if compSymbol == "<" {
                text = strings.Replace(text, match, ternary(v1 < v2, resultValue1, resultValue2), -1)
            } else if compSymbol == ">" {
                text = strings.Replace(text, match, ternary(v1 > v2, resultValue1, resultValue2), -1)
            }
        }
    }
    return text
}

// parse all the `{{?var; value1=#0F0, 56=#00F, ..., default=#000}}` or 
//`{{?var:int; 56=#0F0, 32=#00F, ..., default=#000}}` in the text give in
//
// return -> string
func (t TemplateStr) ParseSwitch(text string) string {

    if !t.HasSwitch(text) { return text }

    for _, group := range findAllGroup(regSwitch, text) {

        match := group["match"]

        mapTemp := map[string]string{}
        var result string

        for _, n := range strings.Split(group["val"], ", ") {
            keyValue := strings.Split(n, "=")
            mapTemp[keyValue[0]] = keyValue[1]
        }

        if group["key"] != "" {
            for key, value := range mapTemp {
                if key == t.variableMap[group["key"]] {
                    result = value
                    break
                } else {
                    result = group["default"]
                }
            }

        } else if group["keyTyped"] != ""{
            keyVar := group["keyTyped"]
            typeVar := group["type"]
            
            for key, value := range mapTemp {
                // println(fmt.Sprintf("%T", typing(key, t.variableMap, typeVar)[0]))
                if valVar, _ := getVariable(keyVar, t.variableMap); typing(key, t.variableMap, typeVar)[0] == valVar {
                    result = value
                    break
                } else {
                    result = group["default"]
                }
            }
        }

        text = strings.Replace(text, match, result, -1)
    }
    return text
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
//`{{?var:int; 56=#0F0, 32=#00F, ..., default=#000}}`
//
// return -> bool
func (t TemplateStr) HasSwitch(text string) bool {
    return regSwitch.MatchString(text)
}



package templateStr

import (
	"fmt"
	"reflect"
	"testing"
	ti "time"

    tem "github.com/CheeseGrinder/TemplateStr-Go/templateStr"
)

var Reset  = "\033[0m"
var Red    = "\033[31m"
var Green  = "\033[32m"
var Yellow = "\033[33m"

type Any = interface{}
type VarMap = tem.VariableMap
type FuncArray = tem.FuncArray

func test([]Any) string {
    return "Test1"
}

func testType(list []Any) string {

    if reflect.TypeOf(list[0]).Kind() != reflect.String { return "list[0] != String"}
    if reflect.TypeOf(list[1]).Kind() != reflect.String { return "list[1] != String"}
    if reflect.TypeOf(list[2]).Kind() != reflect.String { return "list[2] != String"}
    if reflect.TypeOf(list[3]).Kind() != reflect.Bool { return "list[3] != Bool"}
    if reflect.TypeOf(list[4]).Kind() != reflect.Int { return "list[4] != Int"}
    if reflect.TypeOf(list[5]).Kind() != reflect.Float64 { return "list[5] != Float64"}
    if reflect.TypeOf(list[6]).Kind() != reflect.String { return "list[6] != String"}

    var finalStr string

    for index, v := range list {
        if index != 0 {
            finalStr = finalStr + " " + fmt.Sprintf("%v", v)
        } else {
            finalStr = fmt.Sprintf("%v", v)
        }
    }
    return finalStr
}

var arrayFunc FuncArray = FuncArray{testType, test}
var varMap VarMap = VarMap{
    "Build": "Succes",
    "var": "int",
    "str": "Jame",
    "int": 32,
    "float": 4.2,
    "bool": true,
    "lower": "azerty", 
    "upper": "AZERTY", 
    "swap": "AzErTy",
    // "cfold": "grüßen",
    "Map": VarMap{
        "value": "Map in Map",
    },
    "MasterMap": VarMap{
        "SecondMap": VarMap{
            "value": "Map in Map in Map",
        },
    },
}

func TestAll(t *testing.T) {
    
    testAll_1 := []string{
        "Name is @{{uppercase; str}}, ${{int}} years old. Map: ${{Map.value}}. my keyboard: #{{lower == 'azerty'; azerty | qwerty}}, ?{{lower; azerty:yes, AZERTY:no, _:anyway}}",
         "Name is JAME, 32 years old. Map: Map in Map. my keyboard: azerty, yes",
        }

    testAll_2 := []string{"test var in var ${{${{var}}}}", "test var in var 32"}
    testAll_3 := []string{"test func in func @{{lowercase; @{{uppercase; str}}}}", "test func in func none"}
    testAll_4 := []string{"test if in if #{{lower == 'azerty2'; azerty | #{{lower == 'querty'; yes | no}}}}","test if in if no"}
    testAll_5 := []string{"test switch in switch ?{{str; Jame:?{{Build; Succes:#0, Failed:#1, _:#default}}, Tony:#1, Marco:#2, _:#default}}", "test switch in switch #0"}
    testAll_6 := []string{"test wtf ?{{str; Jame:?{{int/${{var}}; 32:#0, 36:#1, _:#default}}, Tony:#1, Marco:#2, _:#default2}}", "test wtf #0"}

    parser := tem.New(arrayFunc, varMap)

    if text := parser.Parse(testAll_1[0]); text != testAll_1[1] {
        t.Fatalf("testAll_1 : '" + Red + text + Reset + "' != '" + Yellow + testAll_1[1] + Reset + "'")
    }

    if text := parser.Parse(testAll_2[0]); text != testAll_2[1] {
        t.Fatalf("testAll_2 : '" + Red + text + Reset + "' != '" + Yellow + testAll_2[1] + Reset + "'")
    }
    
    if text := parser.Parse(testAll_3[0]); text != testAll_3[1] {
        t.Fatalf("testAll_3 : '" + Red + text + Reset + "' != '" + Yellow + testAll_3[1] + Reset + "'")
    }
    
    if text := parser.Parse(testAll_4[0]); text != testAll_4[1] {
        t.Fatalf("testAll_4 : '" + Red + text + Reset + "' != '" + Yellow + testAll_4[1] + Reset + "'")
    }
    
    if text := parser.Parse(testAll_5[0]); text != testAll_5[1] {
        t.Fatalf("testAll_5 : '" + Red + text + Reset + "' != '" + Yellow + testAll_5[1] + Reset + "'")
    }
    
    if text := parser.Parse(testAll_6[0]); text != testAll_6[1] {
        t.Fatalf("testAll_6 : '" + Red + text + Reset + "' != '" + Yellow + testAll_6[1] + Reset + "'")
    }
}

func TestVariable(t *testing.T) {

    text_1 := []string{"var bool = ${{bool}} and name = ${{str}}", "var bool = true and name = Jame"}
    text_2 := []string{"${{Map.value}}", "Map in Map"}
    text_3 := []string{"${{MasterMap.SecondMap.value}}", "Map in Map in Map"}
    text_4 := []string{"${{word}}", "None"}
    text_5 := []string{"${{dict.dict1.value}}", "None"}

    parser := tem.New(FuncArray{}, varMap)

    if text := parser.ParseVariable(text_1[0]); text != text_1[1] {t.Fatalf("text_1 : '" + Red + text + Reset + "' != '" + Yellow + text_1[1] + Reset + "'")}
    if text := parser.ParseVariable(text_2[0]); text != text_2[1] {t.Fatalf("text_2 : '" + Red + text + Reset + "' != '" + Yellow + text_2[1] + Reset + "'")}
    if text := parser.ParseVariable(text_3[0]); text != text_3[1] {t.Fatalf("text_3 : '" + Red + text + Reset + "' != '" + Yellow + text_3[1] + Reset + "'")}
    if text := parser.ParseVariable(text_4[0]); text != text_4[1] {t.Fatalf("text_4 : '" + Red + text + Reset + "' != '" + Yellow + text_4[1] + Reset + "'")}
    if text := parser.ParseVariable(text_5[0]); text != text_5[1] {t.Fatalf("text_5 : '" + Red + text + Reset + "' != '" + Yellow + text_5[1] + Reset + "'")}

}
func TestInternFunction(t *testing.T) {

    uppercase := []string{"@{{uppercase; lower}}", "AZERTY"}
    uppercase2 := []string{"@{{uppercase; word}}", "NONE"}
    uppercaseFirst := []string{"@{{uppercaseFirst; lower}}", "Azerty"}
    lowercase := []string{"@{{lowercase; upper}}", "azerty"}
    // casefold := []string{"@{{casefold cfold}}", "grüssen"}
    swapcase := []string{"@{{swapcase; swap}}", "aZeRtY"}
    time := "@{{time}}"
    date := "@{{date}}"
    dateTime := "@{{dateTime}}"

    parser := tem.New(FuncArray{}, varMap)

    if text := parser.ParseFunction(uppercase[0]); text != uppercase[1] {t.Fatalf("uppercase : '" + Red + text + Reset + "' != '" + Yellow + uppercase[1] + Reset + "'")}
    if text := parser.ParseFunction(uppercase2[0]); text != uppercase2[1] {t.Fatalf("uppercase2 : '" + Red + text + Reset + "' != '" + Yellow + uppercase2[1] + Reset + "'")}
    if text := parser.ParseFunction(uppercaseFirst[0]); text != uppercaseFirst[1] {t.Fatalf( "uppercaseFirst : '" + Red + text + Reset + "' != '" + Yellow + uppercaseFirst[1] + Reset + "'")}
    if text := parser.ParseFunction(lowercase[0]); text != lowercase[1] {t.Fatalf("lowercase : '" + Red + text + Reset + "' != '" + Yellow + lowercase[1] + Reset + "'")}
    // if text := parser.ParseFunction(casefold[0]); text != casefold[1] {t.Fatalf("casefold : '" + Red + text + Reset + "' != '" + Yellow + casefold[1] + Reset + "'")}
    if text := parser.ParseFunction(swapcase[0]); text != swapcase[1] {t.Fatalf("swapcase : '" + Red + text + Reset + "' != '" + Yellow + swapcase[1] + Reset + "'")}

    dateTime1 := ti.Now()
    if text := parser.ParseFunction(time); text != dateTime1.Format("15:04:05") {
        t.Fatalf("time : '" + Red + text + Reset + "' != '" + Yellow + dateTime1.Format("15:04:05") + Reset + "'")
    }
    if text := parser.ParseFunction(date); text != dateTime1.Format("02/01/2006") {
        t.Fatalf("date : '" + Red + text + Reset + "' != '" + Yellow + dateTime1.Format("02/01/2006") + Reset + "'")
    }
    dateTime2 := ti.Now()
    if text := parser.ParseFunction(dateTime); text != dateTime2.Format("02/01/2006 15:04:05") {
        t.Fatalf("dateTime : '" + Red + text + Reset + "' != '" + Yellow + dateTime2.Format("02/01/2006 15:04:05") + Reset + "'")
    }
}

func TestCustomFunction(t *testing.T) {
    
    test := []string{"@{{test}}", "Test1"}
    testType := []string{"@{{testType; \"text\" 'text' `text` b/True i/123 f/123.4 int}}", "text text text true 123 123.4 32"}

    parser := tem.New(arrayFunc, varMap)

    if text := parser.ParseFunction(test[0]); text != test[1] {
        t.Fatalf( "'" + Red + text + Reset + "' != '" + Yellow + test[1] + Reset + "'")
    }
    if text := parser.ParseFunction(testType[0]); text != testType[1] {
        t.Fatalf( "'" + Red + text + Reset + "' != '" + Yellow + testType[1] + Reset + "'")
    }
}

func TestConditionEqual(t *testing.T) {
    
    str_Equal_Str := []string{"#{{'text' == 'text'; yes | no}}", "yes"}
    str_Equal2_Str := []string{"#{{'text' == 'texte'; yes | no}}", "no"}
    int_Equal_Str := []string{"#{{i/4 == 'text'; yes | no}}", "no"}
    float_Equal_Str := []string{"#{{f/4.5 == 'texte'; yes | no}}", "no"}
    bool_Equal_Str := []string{"#{{b/True == 'texte'; yes | no}}", "no"}

    parser := tem.New(FuncArray{}, varMap)

    if text := parser.ParseCondition(str_Equal_Str[0]); text != str_Equal_Str[1] {t.Fatalf("str_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Equal_Str[1] + Reset + "'")}
    if text := parser.ParseCondition(str_Equal2_Str[0]); text != str_Equal2_Str[1] {t.Fatalf("str_Equal2_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Equal2_Str[1] + Reset + "'")}
    if text := parser.ParseCondition(int_Equal_Str[0]); text != int_Equal_Str[1] {t.Fatalf("int_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + int_Equal_Str[1] + Reset + "'")}
    if text := parser.ParseCondition(float_Equal_Str[0]); text != float_Equal_Str[1] {t.Fatalf("float_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + float_Equal_Str[1] + Reset + "'")}
    if text := parser.ParseCondition(bool_Equal_Str[0]); text != bool_Equal_Str[1] {t.Fatalf("bool_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + bool_Equal_Str[1] + Reset + "'")}
}

func TestConditionNotEqual(t *testing.T) {

    str_Not_Equal_Str := []string{"#{{'text' != 'text'; yes | no}}", "no"}
    str_Not_Equal2_Str := []string{"#{{'text' != 'texte'; yes | no}}", "yes"}
    int_Not_Equal_Str := []string{"#{{i/4 != 'text'; yes | no}}", "yes"}
    float_Not_Equal_Str := []string{"#{{f/4.5 != 'texte'; yes | no}}", "yes"}
    bool_Not_Equal_Str := []string{"#{{b/True != 'texte'; yes | no}}", "yes"}

    parser := tem.New(FuncArray{}, varMap)

    if text := parser.ParseCondition(str_Not_Equal_Str[0]); text != str_Not_Equal_Str[1] {t.Fatalf("str_Not_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Not_Equal_Str[1] + Reset + "'")}
    if text := parser.ParseCondition(str_Not_Equal2_Str[0]); text != str_Not_Equal2_Str[1] {t.Fatalf("str_Not_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Not_Equal2_Str[1] + Reset + "'")}
    if text := parser.ParseCondition(int_Not_Equal_Str[0]); text != int_Not_Equal_Str[1] {t.Fatalf("int_Not_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + int_Not_Equal_Str[1] + Reset + "'")}
    if text := parser.ParseCondition(float_Not_Equal_Str[0]); text != float_Not_Equal_Str[1] {t.Fatalf("float_Not_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + float_Not_Equal_Str[1] + Reset + "'")}
    if text := parser.ParseCondition(bool_Not_Equal_Str[0]); text != bool_Not_Equal_Str[1] {t.Fatalf("bool_Not_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + bool_Not_Equal_Str[1] + Reset + "'")}
}

func TestConditionSuperiorEqual(t *testing.T) {

    parser := tem.New(FuncArray{}, varMap)

    // String
    str_Superior_Equal_Str := []string{"#{{'text' >= 'text'; yes | no}}", "yes"}
    str_Superior_Equal2_Str := []string{"#{{'text' >= 'texte'; yes | no}}", "no"}
    str_Superior_Equal_Int := []string{"#{{'text' >= i/4; yes | no}}", "yes"}
    str_Superior_Equal2_Int := []string{"#{{'text' >= i/123; yes | no}}", "no"}
    str_Superior_Equal_Float := []string{"#{{'text' >= f/4.5; yes | no}}", "no"}
    str_Superior_Equal2_Float := []string{"#{{'text' >= f/3.5; yes | no}}", "yes"}
    str_Superior_Equal_Bool := []string{"#{{'text' >= b/True; yes | no}}", "yes"}
    str_Superior_Equal2_Bool := []string{"#{{'text' >= b/False; yes | no}}", "yes"}

    if text := parser.ParseCondition(str_Superior_Equal_Str[0]); text != str_Superior_Equal_Str[1] {
        t.Fatalf("str_Superior_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Equal_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior_Equal2_Str[0]); text != str_Superior_Equal2_Str[1] {
        t.Fatalf("str_Superior_Equal2_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Equal2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior_Equal_Int[0]); text != str_Superior_Equal_Int[1] {
        t.Fatalf("str_Superior_Equal_Int : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Equal_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior_Equal2_Int[0]); text != str_Superior_Equal2_Int[1] {
        t.Fatalf("str_Superior_Equal2_Int : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Equal2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior_Equal_Float[0]); text != str_Superior_Equal_Float[1] {
        t.Fatalf("str_Superior_Equal_Float : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Equal_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior_Equal2_Float[0]); text != str_Superior_Equal2_Float[1] {
        t.Fatalf("str_Superior_Equal2_Float : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Equal2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior_Equal_Bool[0]); text != str_Superior_Equal_Bool[1] {
        t.Fatalf("str_Superior_Equal_Bool : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Equal_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior_Equal2_Bool[0]); text != str_Superior_Equal2_Bool[1] {
        t.Fatalf("str_Superior_Equal2_Bool : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Equal2_Bool[1] + Reset + "'")
    }

    // Int
    int_Superior_Equal_Str := []string{"#{{i/4 >= 'text'; yes | no}}", "yes"}
    int_Superior_Equal2_Str := []string{"#{{i/4 >= 'texte'; yes | no}}", "no"}
    int_Superior_Equal_Int := []string{"#{{i/4 >= i/4; yes | no}}", "yes"}
    int_Superior_Equal2_Int := []string{"#{{i/4 >= i/5; yes | no}}", "no"}
    int_Superior_Equal_Float := []string{"#{{i/4 >= f/3.5; yes | no}}", "yes"}
    int_Superior_Equal2_Float := []string{"#{{i/4 >= f/4.5; yes | no}}", "no"}
    int_Superior_Equal_Bool := []string{"#{{i/4 >= b/True; yes | no}}", "yes"}
    int_Superior_Equal2_Bool := []string{"#{{i/4 >= b/False; yes | no}}", "yes"}

    if text := parser.ParseCondition(int_Superior_Equal_Str[0]); text != int_Superior_Equal_Str[1] {
        t.Fatalf("int_Superior_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Equal_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior_Equal2_Str[0]); text != int_Superior_Equal2_Str[1] {
        t.Fatalf("int_Superior_Equal2_Str : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Equal2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior_Equal_Int[0]); text != int_Superior_Equal_Int[1] {
        t.Fatalf("int_Superior_Equal_Int : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Equal_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior_Equal2_Int[0]); text != int_Superior_Equal2_Int[1] {
        t.Fatalf("int_Superior_Equal2_Int : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Equal2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior_Equal_Float[0]); text != int_Superior_Equal_Float[1] {
        t.Fatalf("int_Superior_Equal_Float : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Equal_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior_Equal2_Float[0]); text != int_Superior_Equal2_Float[1] {
        t.Fatalf("int_Superior_Equal2_Float : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Equal2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior_Equal_Bool[0]); text != int_Superior_Equal_Bool[1] {
        t.Fatalf("int_Superior_Equal_Bool : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Equal_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior_Equal2_Bool[0]); text != int_Superior_Equal2_Bool[1] {
        t.Fatalf("int_Superior_Equal2_Bool : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Equal2_Bool[1] + Reset + "'")
    }

    // Float
    float_Superior_Equal_Str := []string{"#{{f/4.5 >= 'text'; yes | no}}", "yes"}
    float_Superior_Equal2_Str := []string{"#{{f/4.5 >= 'texte'; yes | no}}", "no"}
    float_Superior_Equal_Int := []string{"#{{f/4.5 >= i/4; yes | no}}", "yes"}
    float_Superior_Equal2_Int := []string{"#{{f/4.5 >= i/5; yes | no}}", "no"}
    float_Superior_Equal_Float := []string{"#{{f/4.5 >= f/4.5; yes | no}}", "yes"}
    float_Superior_Equal2_Float := []string{"#{{f/4.5 >= f/4.6; yes | no}}", "no"}
    float_Superior_Equal_Bool := []string{"#{{f/4.5 >= b/True; yes | no}}", "yes"}
    float_Superior_Equal2_Bool := []string{"#{{f/4.5 >= b/False; yes | no}}", "yes"}

    if text := parser.ParseCondition(float_Superior_Equal_Str[0]); text != float_Superior_Equal_Str[1] {
        t.Fatalf("float_Superior_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Equal_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior_Equal2_Str[0]); text != float_Superior_Equal2_Str[1] {
        t.Fatalf("float_Superior_Equal2_Str : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Equal2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior_Equal_Int[0]); text != float_Superior_Equal_Int[1] {
        t.Fatalf("float_Superior_Equal_Int : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Equal_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior_Equal2_Int[0]); text != float_Superior_Equal2_Int[1] {
        t.Fatalf("float_Superior_Equal2_Int : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Equal2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior_Equal_Float[0]); text != float_Superior_Equal_Float[1] {
        t.Fatalf("float_Superior_Equal_Float : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Equal_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior_Equal2_Float[0]); text != float_Superior_Equal2_Float[1] {
        t.Fatalf("float_Superior_Equal2_Float : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Equal2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior_Equal_Bool[0]); text != float_Superior_Equal_Bool[1] {
        t.Fatalf("float_Superior_Equal_Bool : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Equal_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior_Equal2_Bool[0]); text != float_Superior_Equal2_Bool[1] {
        t.Fatalf("float_Superior_Equal2_Bool : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Equal2_Bool[1] + Reset + "'")
    }

    // Bool
    bool_Superior_Equal_Str := []string{"#{{b/True >= 'text'; yes | no}}", "no"}
    bool_Superior_Equal2_Str := []string{"#{{b/False >= 'texte'; yes | no}}", "no"}
    bool_Superior_Equal_Int := []string{"#{{b/True >= i/4; yes | no}}", "no"}
    bool_Superior_Equal2_Int := []string{"#{{b/False >= i/5; yes | no}}", "no"}
    bool_Superior_Equal_Float := []string{"#{{b/True >= f/4.5; yes | no}}", "no"}
    bool_Superior_Equal2_Float := []string{"#{{b/False >= f/4.6; yes | no}}", "no"}
    bool_Superior_Equal_Bool := []string{"#{{b/True >= b/True; yes | no}}", "yes"}
    bool_Superior_Equal2_Bool := []string{"#{{b/False >= b/False; yes | no}}", "yes"}

    if text := parser.ParseCondition(bool_Superior_Equal_Str[0]); text != bool_Superior_Equal_Str[1] {
        t.Fatalf("bool_Superior_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Equal_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior_Equal2_Str[0]); text != bool_Superior_Equal2_Str[1] {
        t.Fatalf("bool_Superior_Equal2_Str : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Equal2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior_Equal_Int[0]); text != bool_Superior_Equal_Int[1] {
        t.Fatalf("bool_Superior_Equal_Int : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Equal_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior_Equal2_Int[0]); text != bool_Superior_Equal2_Int[1] {
        t.Fatalf("bool_Superior_Equal2_Int : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Equal2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior_Equal_Float[0]); text != bool_Superior_Equal_Float[1] {
        t.Fatalf("bool_Superior_Equal_Float : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Equal_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior_Equal2_Float[0]); text != bool_Superior_Equal2_Float[1] {
        t.Fatalf("bool_Superior_Equal2_Float : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Equal2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior_Equal_Bool[0]); text != bool_Superior_Equal_Bool[1] {
        t.Fatalf("bool_Superior_Equal_Bool : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Equal_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior_Equal2_Bool[0]); text != bool_Superior_Equal2_Bool[1] {
        t.Fatalf("bool_Superior_Equal2_Bool : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Equal2_Bool[1] + Reset + "'")
    }
}

func TestConditionSuperior(t *testing.T) {

    parser := tem.New(FuncArray{}, varMap)

    // String
    str_Superior_Str := []string{"#{{'text' > 'text'; yes | no}}", "no"}
    str_Superior2_Str := []string{"#{{'text' > 'texte'; yes | no}}", "no"}
    str_Superior_Int := []string{"#{{'text' > i/4; yes | no}}", "no"}
    str_Superior2_Int := []string{"#{{'text' > i/123; yes | no}}", "no"}
    str_Superior_Float := []string{"#{{'text' > f/4.5; yes | no}}", "no"}
    str_Superior2_Float := []string{"#{{'text' > f/3.5; yes | no}}", "yes"}
    str_Superior_Bool := []string{"#{{'text' > b/True; yes | no}}", "yes"}
    str_Superior2_Bool := []string{"#{{'text' > b/False; yes | no}}", "yes"}

    if text := parser.ParseCondition(str_Superior_Str[0]); text != str_Superior_Str[1] {
        t.Fatalf("str_Superior_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior2_Str[0]); text != str_Superior2_Str[1] {
        t.Fatalf("str_Superior2_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Superior2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior_Int[0]); text != str_Superior_Int[1] {
        t.Fatalf("str_Superior_Int : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior2_Int[0]); text != str_Superior2_Int[1] {
        t.Fatalf("str_Superior2_Int : '" + Red + text + Reset + "' != '" + Yellow + str_Superior2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior_Float[0]); text != str_Superior_Float[1] {
        t.Fatalf("str_Superior_Float : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior2_Float[0]); text != str_Superior2_Float[1] {
        t.Fatalf("str_Superior2_Float : '" + Red + text + Reset + "' != '" + Yellow + str_Superior2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior_Bool[0]); text != str_Superior_Bool[1] {
        t.Fatalf("str_Superior_Bool : '" + Red + text + Reset + "' != '" + Yellow + str_Superior_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Superior2_Bool[0]); text != str_Superior2_Bool[1] {
        t.Fatalf("str_Superior2_Bool : '" + Red + text + Reset + "' != '" + Yellow + str_Superior2_Bool[1] + Reset + "'")
    }

    // Int
    int_Superior_Str := []string{"#{{i/4 > 'text'; yes | no}}", "no"}
    int_Superior2_Str := []string{"#{{i/4 > 'texte'; yes | no}}", "no"}
    int_Superior_Int := []string{"#{{i/4 > i/4; yes | no}}", "no"}
    int_Superior2_Int := []string{"#{{i/4 > i/5; yes | no}}", "no"}
    int_Superior_Float := []string{"#{{i/4 > f/3.5; yes | no}}", "yes"}
    int_Superior2_Float := []string{"#{{i/4 > f/4.5; yes | no}}", "no"}
    int_Superior_Bool := []string{"#{{i/4 > b/True; yes | no}}", "yes"}
    int_Superior2_Bool := []string{"#{{i/4 > b/False; yes | no}}", "yes"}

    if text := parser.ParseCondition(int_Superior_Str[0]); text != int_Superior_Str[1] {
        t.Fatalf("int_Superior_Str : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior2_Str[0]); text != int_Superior2_Str[1] {
        t.Fatalf("int_Superior2_Str : '" + Red + text + Reset + "' != '" + Yellow + int_Superior2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior_Int[0]); text != int_Superior_Int[1] {
        t.Fatalf("int_Superior_Int : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior2_Int[0]); text != int_Superior2_Int[1] {
        t.Fatalf("int_Superior2_Int : '" + Red + text + Reset + "' != '" + Yellow + int_Superior2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior_Float[0]); text != int_Superior_Float[1] {
        t.Fatalf("int_Superior_Float : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior2_Float[0]); text != int_Superior2_Float[1] {
        t.Fatalf("int_Superior2_Float : '" + Red + text + Reset + "' != '" + Yellow + int_Superior2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior_Bool[0]); text != int_Superior_Bool[1] {
        t.Fatalf("int_Superior_Bool : '" + Red + text + Reset + "' != '" + Yellow + int_Superior_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Superior2_Bool[0]); text != int_Superior2_Bool[1] {
        t.Fatalf("int_Superior2_Bool : '" + Red + text + Reset + "' != '" + Yellow + int_Superior2_Bool[1] + Reset + "'")
    }

    // Float
    float_Superior_Str := []string{"#{{f/4.5 > 'text'; yes | no}}", "yes"}
    float_Superior2_Str := []string{"#{{f/4.5 > 'texte'; yes | no}}", "no"}
    float_Superior_Int := []string{"#{{f/4.5 > i/4; yes | no}}", "yes"}
    float_Superior2_Int := []string{"#{{f/4.5 > i/5; yes | no}}", "no"}
    float_Superior_Float := []string{"#{{f/4.5 > f/4.4; yes | no}}", "yes"}
    float_Superior2_Float := []string{"#{{f/4.5 > f/4.6; yes | no}}", "no"}
    float_Superior_Bool := []string{"#{{f/4.5 > b/True; yes | no}}", "yes"}
    float_Superior2_Bool := []string{"#{{f/4.5 > b/False; yes | no}}", "yes"}

    if text := parser.ParseCondition(float_Superior_Str[0]); text != float_Superior_Str[1] {
        t.Fatalf("float_Superior_Str : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior2_Str[0]); text != float_Superior2_Str[1] {
        t.Fatalf("float_Superior2_Str : '" + Red + text + Reset + "' != '" + Yellow + float_Superior2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior_Int[0]); text != float_Superior_Int[1] {
        t.Fatalf("float_Superior_Int : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior2_Int[0]); text != float_Superior2_Int[1] {
        t.Fatalf("float_Superior2_Int : '" + Red + text + Reset + "' != '" + Yellow + float_Superior2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior_Float[0]); text != float_Superior_Float[1] {
        t.Fatalf("float_Superior_Float : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior2_Float[0]); text != float_Superior2_Float[1] {
        t.Fatalf("float_Superior2_Float : '" + Red + text + Reset + "' != '" + Yellow + float_Superior2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior_Bool[0]); text != float_Superior_Bool[1] {
        t.Fatalf("float_Superior_Bool : '" + Red + text + Reset + "' != '" + Yellow + float_Superior_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Superior2_Bool[0]); text != float_Superior2_Bool[1] {
        t.Fatalf("float_Superior2_Bool : '" + Red + text + Reset + "' != '" + Yellow + float_Superior2_Bool[1] + Reset + "'")
    }

    // Bool
    bool_Superior_Str := []string{"#{{b/True > 'text'; yes | no}}", "no"}
    bool_Superior2_Str := []string{"#{{b/False > 'texte'; yes | no}}", "no"}
    bool_Superior_Int := []string{"#{{b/True > i/4; yes | no}}", "no"}
    bool_Superior2_Int := []string{"#{{b/False > i/5; yes | no}}", "no"}
    bool_Superior_Float := []string{"#{{b/True > f/4.5; yes | no}}", "no"}
    bool_Superior2_Float := []string{"#{{b/False > f/4.6; yes | no}}", "no"}
    bool_Superior_Bool := []string{"#{{b/True > b/True; yes | no}}", "no"}
    bool_Superior2_Bool := []string{"#{{b/False > b/False; yes | no}}", "no"}

    if text := parser.ParseCondition(bool_Superior_Str[0]); text != bool_Superior_Str[1] {
        t.Fatalf("bool_Superior_Str : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior2_Str[0]); text != bool_Superior2_Str[1] {
        t.Fatalf("bool_Superior2_Str : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior_Int[0]); text != bool_Superior_Int[1] {
        t.Fatalf("bool_Superior_Int : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior2_Int[0]); text != bool_Superior2_Int[1] {
        t.Fatalf("bool_Superior2_Int : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior_Float[0]); text != bool_Superior_Float[1] {
        t.Fatalf("bool_Superior_Float : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior2_Float[0]); text != bool_Superior2_Float[1] {
        t.Fatalf("bool_Superior2_Float : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior_Bool[0]); text != bool_Superior_Bool[1] {
        t.Fatalf("bool_Superior_Bool : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Superior2_Bool[0]); text != bool_Superior2_Bool[1] {
        t.Fatalf("bool_Superior2_Bool : '" + Red + text + Reset + "' != '" + Yellow + bool_Superior2_Bool[1] + Reset + "'")
    }
}

func TestConditionInferiorEqual(t *testing.T) {

    parser := tem.New(FuncArray{}, varMap)

    // String
    str_Inferior_Equal_Str := []string{"#{{'text' <= 'text'; yes | no}}", "yes"}
    str_Inferior_Equal2_Str := []string{"#{{'text' <= 'texte'; yes | no}}", "yes"}
    str_Inferior_Equal_Int := []string{"#{{'text' <= i/4; yes | no}}", "yes"}
    str_Inferior_Equal2_Int := []string{"#{{'text' <= i/123; yes | no}}", "yes"}
    str_Inferior_Equal_Float := []string{"#{{'text' <= f/4.5; yes | no}}", "yes"}
    str_Inferior_Equal2_Float := []string{"#{{'text' <= f/3.5; yes | no}}", "no"}
    str_Inferior_Equal_Bool := []string{"#{{'text' <= b/True; yes | no}}", "no"}
    str_Inferior_Equal2_Bool := []string{"#{{'text' <= b/False; yes | no}}", "no"}

    if text := parser.ParseCondition(str_Inferior_Equal_Str[0]); text != str_Inferior_Equal_Str[1] {
        t.Fatalf("str_Inferior_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Equal_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior_Equal2_Str[0]); text != str_Inferior_Equal2_Str[1] {
        t.Fatalf("str_Inferior_Equal2_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Equal2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior_Equal_Int[0]); text != str_Inferior_Equal_Int[1] {
        t.Fatalf("str_Inferior_Equal_Int : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Equal_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior_Equal2_Int[0]); text != str_Inferior_Equal2_Int[1] {
        t.Fatalf("str_Inferior_Equal2_Int : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Equal2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior_Equal_Float[0]); text != str_Inferior_Equal_Float[1] {
        t.Fatalf("str_Inferior_Equal_Float : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Equal_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior_Equal2_Float[0]); text != str_Inferior_Equal2_Float[1] {
        t.Fatalf("str_Inferior_Equal2_Float : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Equal2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior_Equal_Bool[0]); text != str_Inferior_Equal_Bool[1] {
        t.Fatalf("str_Inferior_Equal_Bool : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Equal_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior_Equal2_Bool[0]); text != str_Inferior_Equal2_Bool[1] {
        t.Fatalf("str_Inferior_Equal2_Bool : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Equal2_Bool[1] + Reset + "'")
    }

    // Int
    int_Inferior_Equal_Str := []string{"#{{i/4 <= 'text'; yes | no}}", "yes"}
    int_Inferior_Equal2_Str := []string{"#{{i/4 <= 'texte'; yes | no}}", "yes"}
    int_Inferior_Equal_Int := []string{"#{{i/4 <= i/4; yes | no}}", "yes"}
    int_Inferior_Equal2_Int := []string{"#{{i/4 <= i/5; yes | no}}", "yes"}
    int_Inferior_Equal_Float := []string{"#{{i/4 <= f/3.5; yes | no}}", "no"}
    int_Inferior_Equal2_Float := []string{"#{{i/4 <= f/4.5; yes | no}}", "yes"}
    int_Inferior_Equal_Bool := []string{"#{{i/4 <= b/True; yes | no}}", "no"}
    int_Inferior_Equal2_Bool := []string{"#{{i/4 <= b/False; yes | no}}", "no"}

    if text := parser.ParseCondition(int_Inferior_Equal_Str[0]); text != int_Inferior_Equal_Str[1] {
        t.Fatalf("int_Inferior_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Equal_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior_Equal2_Str[0]); text != int_Inferior_Equal2_Str[1] {
        t.Fatalf("int_Inferior_Equal2_Str : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Equal2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior_Equal_Int[0]); text != int_Inferior_Equal_Int[1] {
        t.Fatalf("int_Inferior_Equal_Int : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Equal_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior_Equal2_Int[0]); text != int_Inferior_Equal2_Int[1] {
        t.Fatalf("int_Inferior_Equal2_Int : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Equal2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior_Equal_Float[0]); text != int_Inferior_Equal_Float[1] {
        t.Fatalf("int_Inferior_Equal_Float : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Equal_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior_Equal2_Float[0]); text != int_Inferior_Equal2_Float[1] {
        t.Fatalf("int_Inferior_Equal2_Float : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Equal2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior_Equal_Bool[0]); text != int_Inferior_Equal_Bool[1] {
        t.Fatalf("int_Inferior_Equal_Bool : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Equal_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior_Equal2_Bool[0]); text != int_Inferior_Equal2_Bool[1] {
        t.Fatalf("int_Inferior_Equal2_Bool : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Equal2_Bool[1] + Reset + "'")
    }

    // Float
    float_Inferior_Equal_Str := []string{"#{{f/4.5 <= 'text'; yes | no}}", "no"}
    float_Inferior_Equal2_Str := []string{"#{{f/4.5 <= 'texte'; yes | no}}", "yes"}
    float_Inferior_Equal_Int := []string{"#{{f/4.5 <= i/4; yes | no}}", "no"}
    float_Inferior_Equal2_Int := []string{"#{{f/4.5 <= i/5; yes | no}}", "yes"}
    float_Inferior_Equal_Float := []string{"#{{f/4.5 <= f/4.4; yes | no}}", "no"}
    float_Inferior_Equal2_Float := []string{"#{{f/4.5 <= f/4.6; yes | no}}", "yes"}
    float_Inferior_Equal_Bool := []string{"#{{f/4.5 <= b/True; yes | no}}", "no"}
    float_Inferior_Equal2_Bool := []string{"#{{f/4.5 <= b/False; yes | no}}", "no"}

    if text := parser.ParseCondition(float_Inferior_Equal_Str[0]); text != float_Inferior_Equal_Str[1] {
        t.Fatalf("float_Inferior_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Equal_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior_Equal2_Str[0]); text != float_Inferior_Equal2_Str[1] {
        t.Fatalf("float_Inferior_Equal2_Str : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Equal2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior_Equal_Int[0]); text != float_Inferior_Equal_Int[1] {
        t.Fatalf("float_Inferior_Equal_Int : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Equal_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior_Equal2_Int[0]); text != float_Inferior_Equal2_Int[1] {
        t.Fatalf("float_Inferior_Equal2_Int : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Equal2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior_Equal_Float[0]); text != float_Inferior_Equal_Float[1] {
        t.Fatalf("float_Inferior_Equal_Float : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Equal_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior_Equal2_Float[0]); text != float_Inferior_Equal2_Float[1] {
        t.Fatalf("float_Inferior_Equal2_Float : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Equal2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior_Equal_Bool[0]); text != float_Inferior_Equal_Bool[1] {
        t.Fatalf("float_Inferior_Equal_Bool : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Equal_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior_Equal2_Bool[0]); text != float_Inferior_Equal2_Bool[1] {
        t.Fatalf("float_Inferior_Equal2_Bool : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Equal2_Bool[1] + Reset + "'")
    }

    // Bool
    bool_Inferior_Equal_Str := []string{"#{{b/True <= 'text'; yes | no}}", "yes"}
    bool_Inferior_Equal2_Str := []string{"#{{b/False <= 'texte'; yes | no}}", "yes"}
    bool_Inferior_Equal_Int := []string{"#{{b/True <= i/4; yes | no}}", "yes"}
    bool_Inferior_Equal2_Int := []string{"#{{b/False <= i/5; yes | no}}", "yes"}
    bool_Inferior_Equal_Float := []string{"#{{b/True <= f/4.5; yes | no}}", "yes"}
    bool_Inferior_Equal2_Float := []string{"#{{b/False <= f/4.6; yes | no}}", "yes"}
    bool_Inferior_Equal_Bool := []string{"#{{b/True <= b/True; yes | no}}", "yes"}
    bool_Inferior_Equal2_Bool := []string{"#{{b/False <= b/False; yes | no}}", "yes"}

    if text := parser.ParseCondition(bool_Inferior_Equal_Str[0]); text != bool_Inferior_Equal_Str[1] {
        t.Fatalf("bool_Inferior_Equal_Str : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Equal_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior_Equal2_Str[0]); text != bool_Inferior_Equal2_Str[1] {
        t.Fatalf("bool_Inferior_Equal2_Str : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Equal2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior_Equal_Int[0]); text != bool_Inferior_Equal_Int[1] {
        t.Fatalf("bool_Inferior_Equal_Int : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Equal_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior_Equal2_Int[0]); text != bool_Inferior_Equal2_Int[1] {
        t.Fatalf("bool_Inferior_Equal2_Int : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Equal2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior_Equal_Float[0]); text != bool_Inferior_Equal_Float[1] {
        t.Fatalf("bool_Inferior_Equal_Float : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Equal_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior_Equal2_Float[0]); text != bool_Inferior_Equal2_Float[1] {
        t.Fatalf("bool_Inferior_Equal2_Float : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Equal2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior_Equal_Bool[0]); text != bool_Inferior_Equal_Bool[1] {
        t.Fatalf("bool_Inferior_Equal_Bool : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Equal_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior_Equal2_Bool[0]); text != bool_Inferior_Equal2_Bool[1] {
        t.Fatalf("bool_Inferior_Equal2_Bool : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Equal2_Bool[1] + Reset + "'")
    }
}

func TestConditionInferior(t *testing.T) {

    parser := tem.New(FuncArray{}, varMap)

    // String
    str_Inferior_Str := []string{"#{{'text' < 'text'; yes | no}}", "no"}
    str_Inferior2_Str := []string{"#{{'text' < 'texte'; yes | no}}", "yes"}
    str_Inferior_Int := []string{"#{{'text' < i/4; yes | no}}", "no"}
    str_Inferior2_Int := []string{"#{{'text' < i/123; yes | no}}", "yes"}
    str_Inferior_Float := []string{"#{{'text' < f/4.5; yes | no}}", "yes"}
    str_Inferior2_Float := []string{"#{{'text' < f/3.5; yes | no}}", "no"}
    str_Inferior_Bool := []string{"#{{'text' < b/True; yes | no}}", "no"}
    str_Inferior2_Bool := []string{"#{{'text' < b/False; yes | no}}", "no"}

    if text := parser.ParseCondition(str_Inferior_Str[0]); text != str_Inferior_Str[1] {
        t.Fatalf("str_Inferior_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior2_Str[0]); text != str_Inferior2_Str[1] {
        t.Fatalf("str_Inferior2_Str : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior_Int[0]); text != str_Inferior_Int[1] {
        t.Fatalf("str_Inferior_Int : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior2_Int[0]); text != str_Inferior2_Int[1] {
        t.Fatalf("str_Inferior2_Int : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior_Float[0]); text != str_Inferior_Float[1] {
        t.Fatalf("str_Inferior_Float : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior2_Float[0]); text != str_Inferior2_Float[1] {
        t.Fatalf("str_Inferior2_Float : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior_Bool[0]); text != str_Inferior_Bool[1] {
        t.Fatalf("str_Inferior_Bool : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(str_Inferior2_Bool[0]); text != str_Inferior2_Bool[1] {
        t.Fatalf("str_Inferior2_Bool : '" + Red + text + Reset + "' != '" + Yellow + str_Inferior2_Bool[1] + Reset + "'")
    }

    // Int
    int_Inferior_Str := []string{"#{{i/4 < 'text'; yes | no}}", "no"}
    int_Inferior2_Str := []string{"#{{i/4 < 'texte'; yes | no}}", "yes"}
    int_Inferior_Int := []string{"#{{i/4 < i/4; yes | no}}", "no"}
    int_Inferior2_Int := []string{"#{{i/4 < i/5; yes | no}}", "yes"}
    int_Inferior_Float := []string{"#{{i/4 < f/3.5; yes | no}}", "no"}
    int_Inferior2_Float := []string{"#{{i/4 < f/4.5; yes | no}}", "yes"}
    int_Inferior_Bool := []string{"#{{i/4 < b/True; yes | no}}", "no"}
    int_Inferior2_Bool := []string{"#{{i/4 < b/False; yes | no}}", "no"}

    if text := parser.ParseCondition(int_Inferior_Str[0]); text != int_Inferior_Str[1] {
        t.Fatalf("int_Inferior_Str : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior2_Str[0]); text != int_Inferior2_Str[1] {
        t.Fatalf("int_Inferior2_Str : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior_Int[0]); text != int_Inferior_Int[1] {
        t.Fatalf("int_Inferior_Int : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior2_Int[0]); text != int_Inferior2_Int[1] {
        t.Fatalf("int_Inferior2_Int : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior_Float[0]); text != int_Inferior_Float[1] {
        t.Fatalf("int_Inferior_Float : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior2_Float[0]); text != int_Inferior2_Float[1] {
        t.Fatalf("int_Inferior2_Float : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior_Bool[0]); text != int_Inferior_Bool[1] {
        t.Fatalf("int_Inferior_Bool : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(int_Inferior2_Bool[0]); text != int_Inferior2_Bool[1] {
        t.Fatalf("int_Inferior2_Bool : '" + Red + text + Reset + "' != '" + Yellow + int_Inferior2_Bool[1] + Reset + "'")
    }

    // Float
    float_Inferior_Str := []string{"#{{f/4.5 < 'text'; yes | no}}", "no"}
    float_Inferior2_Str := []string{"#{{f/4.5 < 'texte'; yes | no}}", "yes"}
    float_Inferior_Int := []string{"#{{f/4.5 < i/4; yes | no}}", "no"}
    float_Inferior2_Int := []string{"#{{f/4.5 < i/5; yes | no}}", "yes"}
    float_Inferior_Float := []string{"#{{f/4.5 < f/4.5; yes | no}}", "no"}
    float_Inferior2_Float := []string{"#{{f/4.5 < f/4.6; yes | no}}", "yes"}
    float_Inferior_Bool := []string{"#{{f/4.5 < b/True; yes | no}}", "no"}
    float_Inferior2_Bool := []string{"#{{f/4.5 < b/False; yes | no}}", "no"}

    if text := parser.ParseCondition(float_Inferior_Str[0]); text != float_Inferior_Str[1] {
        t.Fatalf("float_Inferior_Str : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior2_Str[0]); text != float_Inferior2_Str[1] {
        t.Fatalf("float_Inferior2_Str : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior_Int[0]); text != float_Inferior_Int[1] {
        t.Fatalf("float_Inferior_Int : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior2_Int[0]); text != float_Inferior2_Int[1] {
        t.Fatalf("float_Inferior2_Int : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior_Float[0]); text != float_Inferior_Float[1] {
        t.Fatalf("float_Inferior_Float : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior2_Float[0]); text != float_Inferior2_Float[1] {
        t.Fatalf("float_Inferior2_Float : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior_Bool[0]); text != float_Inferior_Bool[1] {
        t.Fatalf("float_Inferior_Bool : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(float_Inferior2_Bool[0]); text != float_Inferior2_Bool[1] {
        t.Fatalf("float_Inferior2_Bool : '" + Red + text + Reset + "' != '" + Yellow + float_Inferior2_Bool[1] + Reset + "'")
    }

    // Bool
    bool_Inferior_Str := []string{"#{{b/True < 'text'; yes | no}}", "yes"}
    bool_Inferior2_Str := []string{"#{{b/False < 'texte'; yes | no}}", "yes"}
    bool_Inferior_Int := []string{"#{{b/True < i/4; yes | no}}", "yes"}
    bool_Inferior2_Int := []string{"#{{b/False < i/5; yes | no}}", "yes"}
    bool_Inferior_Float := []string{"#{{b/True < f/4.5; yes | no}}", "yes"}
    bool_Inferior2_Float := []string{"#{{b/False < f/4.6; yes | no}}", "yes"}
    bool_Inferior_Bool := []string{"#{{b/True < b/True; yes | no}}", "no"}
    bool_Inferior2_Bool := []string{"#{{b/False < b/False; yes | no}}", "no"}

    if text := parser.ParseCondition(bool_Inferior_Str[0]); text != bool_Inferior_Str[1] {
        t.Fatalf("bool_Inferior_Str : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior2_Str[0]); text != bool_Inferior2_Str[1] {
        t.Fatalf("bool_Inferior2_Str : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior2_Str[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior_Int[0]); text != bool_Inferior_Int[1] {
        t.Fatalf("bool_Inferior_Int : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior2_Int[0]); text != bool_Inferior2_Int[1] {
        t.Fatalf("bool_Inferior2_Int : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior2_Int[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior_Float[0]); text != bool_Inferior_Float[1] {
        t.Fatalf("bool_Inferior_Float : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior2_Float[0]); text != bool_Inferior2_Float[1] {
        t.Fatalf("bool_Inferior2_Float : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior2_Float[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior_Bool[0]); text != bool_Inferior_Bool[1] {
        t.Fatalf("bool_Inferior_Bool : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior_Bool[1] + Reset + "'")
    }
    if text := parser.ParseCondition(bool_Inferior2_Bool[0]); text != bool_Inferior2_Bool[1] {
        t.Fatalf("bool_Inferior2_Bool : '" + Red + text + Reset + "' != '" + Yellow + bool_Inferior2_Bool[1] + Reset + "'")
    }
}

func TestSwitch(t *testing.T) {
    
    text_Switch_1 := []string{"?{{str; Jame:#0, Tony:#1, Marco:#2, _:#default}}", "#0"}
    text_Switch_2 := []string{"?{{int/int; 56:#0, 36:#1, 32:#2, _:#default}}", "#2"}

    parser := tem.New(arrayFunc, varMap)

    if text := parser.ParseSwitch(text_Switch_1[0]); text != text_Switch_1[1] {
        t.Fatalf("text_Switch_1 : '" + Red + text + Reset + "' != '" + Yellow + text_Switch_1[1] + Reset + "'")
    }

    if text := parser.Parse(text_Switch_2[0]); text != text_Switch_2[1] {
        t.Fatalf("text_Switch_2 : '" + Red + text + Reset + "' != '" + Yellow + text_Switch_2[1] + Reset + "'")
    }
}

func TestHasOne(t *testing.T) {
    
    text_Has_One_1 := []Any{"?{{age:int; 56:#0, 36:#1, 32:#2, _:#default}} and ?{{age:int; 56:#0, 36:#1, 32:#2, _:#default}}", true}
    text_Has_One_2 := []Any{"?{{age:int; 56:#0, 36:#1, 32:#2, _:#default}} and ${{bool}}", true}
    text_Has_One_3 := []Any{"@{{uppercase ${{var}}}} and ${{name}}", true}
    text_Has_One_4 := []Any{"text", false}
    text_Has_One_5 := []Any{"%{{bool}}", false}

    parser := tem.New(FuncArray{}, VarMap{})

    if text := parser.HasOne(fmt.Sprintf("%v",text_Has_One_1[0])); text != text_Has_One_1[1] {
        t.Fatalf("text_Has_One_1 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_One_1[1]) + Reset + "'")
    }

    if text := parser.HasOne(fmt.Sprintf("%v",text_Has_One_2[0])); text != text_Has_One_2[1] {
        t.Fatalf("text_Has_One_2 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_One_2[1]) + Reset + "'")
    }

    if text := parser.HasOne(fmt.Sprintf("%v",text_Has_One_3[0])); text != text_Has_One_3[1] {
        t.Fatalf("text_Has_One_3 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_One_3[1]) + Reset + "'")
    }
    
    if text := parser.HasOne(fmt.Sprintf("%v",text_Has_One_4[0])); text != text_Has_One_4[1] {
        t.Fatalf("text_Has_One_4 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_One_4[1]) + Reset + "'")
    }
    
    if text := parser.HasOne(fmt.Sprintf("%v",text_Has_One_5[0])); text != text_Has_One_5[1] {
        t.Fatalf("text_Has_One_5 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_One_5[1]) + Reset + "'")
    }
}

func TestHasVariable(t *testing.T) {
    
    text_Has_Variable_1 := []Any{"${{bool}} and ${{name}}", true}
    text_Has_Variable_2 := []Any{"${{bool}} and @{{uppercase lower}}", true}
    text_Has_Variable_3 := []Any{"@{{uppercaseFirst bool}} and @{{uppercase lower}}", false}

    parser := tem.New(FuncArray{}, VarMap{})

    if text := parser.HasVariable(fmt.Sprintf("%v",text_Has_Variable_1[0])); text != text_Has_Variable_1[1] {
        t.Fatalf("text_Has_Variable_1 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Variable_1[1]) + Reset + "'")
    }

    if text := parser.HasVariable(fmt.Sprintf("%v",text_Has_Variable_2[0])); text != text_Has_Variable_2[1] {
        t.Fatalf("text_Has_Variable_2 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Variable_2[1]) + Reset + "'")
    }

    if text := parser.HasVariable(fmt.Sprintf("%v",text_Has_Variable_3[0])); text != text_Has_Variable_3[1] {
        t.Fatalf("text_Has_Variable_3 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Variable_3[1]) + Reset + "'")
    }
}

func TestHasFunction(t *testing.T) {
    
    text_Has_Function_1 := []Any{"@{{uppercase; lower}} and @{{uppercaseFirst; lower}}", true}
    text_Has_Function_2 := []Any{"@{{uppercase; lower}} and #{{'text' > 'text'; yes | no}}", true}
    text_Has_Function_3 := []Any{"#{{'text' > 'text'; yes | no}} and #{{'text' < 'text'; yes | no}}", false}
    
    parser := tem.New(FuncArray{}, VarMap{})
    
    if text := parser.HasFunction(fmt.Sprintf("%v",text_Has_Function_1[0])); text != text_Has_Function_1[1] {
        t.Fatalf("text_Has_Function_1 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Function_1[1]) + Reset + "'")
    }
    
    if text := parser.HasFunction(fmt.Sprintf("%v",text_Has_Function_2[0])); text != text_Has_Function_2[1] {
        t.Fatalf("text_Has_Function_2 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Function_2[1]) + Reset + "'")
    }
    
    if text := parser.HasFunction(fmt.Sprintf("%v",text_Has_Function_3[0])); text != text_Has_Function_3[1] {
        t.Fatalf("text_Has_Function_3 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Function_3[1]) + Reset + "'")
    }
}

func TestHasCondition(t *testing.T) {
    
    text_Has_Condition_1 := []Any{"#{{'text' > 'text'; yes | no}} and #{{'text' < 'text'; yes | no}}", true}
    text_Has_Condition_2 := []Any{"#{{'text' > 'text'; yes | no}} and ?{{age:int; 56:#0, 36:#1, 32:#2, _:#default}}", true}
    text_Has_Condition_3 := []Any{"?{{age:int; 56:#0, 36:#1, 32:#2, _:#default}} and ?{{age:int; 56:#0, 36:#1, 32:#2, _:#default}}", false}
    
    parser := tem.New(FuncArray{}, VarMap{})
    
    if text := parser.HasCondition(fmt.Sprintf("%v",text_Has_Condition_1[0])); text != text_Has_Condition_1[1] {
        t.Fatalf("text_Has_Condition_1 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Condition_1[1]) + Reset + "'")
    }
    
    if text := parser.HasCondition(fmt.Sprintf("%v",text_Has_Condition_2[0])); text != text_Has_Condition_2[1] {
        t.Fatalf("text_Has_Condition_2 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Condition_2[1]) + Reset + "'")
    }
    
    if text := parser.HasCondition(fmt.Sprintf("%v",text_Has_Condition_3[0])); text != text_Has_Condition_3[1] {
        t.Fatalf("text_Has_Condition_3 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Condition_3[1]) + Reset + "'")
    }
}

func TestHasSwitch(t *testing.T) {
    
    text_Has_Switch_1 := []Any{"?{{age:int; 56:#0, 36:#1, 32:#2, _:#default}} and ?{{age:int; 56:#0, 36:#1, 32:#2, _:#default}}", true}
    text_Has_Switch_2 := []Any{"?{{age:int; 56:#0, 36:#1, 32:#2, _:#default}} and ${{bool}}", true}
    text_Has_Switch_3 := []Any{"${{bool}} and ${{name}}", false}

    parser := tem.New(FuncArray{}, VarMap{})

    if text := parser.HasSwitch(fmt.Sprintf("%v",text_Has_Switch_1[0])); text != text_Has_Switch_1[1] {
        t.Fatalf("text_Has_Switch_1 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Switch_1[1]) + Reset + "'")
    }

    if text := parser.HasSwitch(fmt.Sprintf("%v",text_Has_Switch_2[0])); text != text_Has_Switch_2[1] {
        t.Fatalf("text_Has_Switch_2 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Switch_2[1]) + Reset + "'")
    }

    if text := parser.HasSwitch(fmt.Sprintf("%v",text_Has_Switch_3[0])); text != text_Has_Switch_3[1] {
        t.Fatalf("text_Has_Switch_3 : '" + Red + fmt.Sprintf("%v", text) + Reset + "' != '" + Yellow + fmt.Sprintf("%v", text_Has_Switch_3[1]) + Reset + "'")
    }
}
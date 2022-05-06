
<div align="center">
    <h1>TemplateStr-Go</h1>
    <h3>TemplateStr allows to add variable, function, condition and switch in a string.</h3>
    <img src="https://img.shields.io/static/v1?label=Go&message=1.11%5E&color=22CFFA&style=flat-square&logo=Go&logoColor=00ADD8"/>
    <a href="https://github.com/CheeseGrinder/TemplateStr-Go/actions/workflows/python-app.yml">
        <img src="https://img.shields.io/github/workflow/status/CheeseGrinder/TemplateStr-Go/Go Test?label=Test&style=flat-square"/>
    </a>
</div>

#### Install :
```
go get -u github.com/CheeseGrinder/TemplateStr-Go/templateStr@latest
```

#### Import :

```go
import (
    "github.com/CheeseGrinder/TemplateStr-Go/templateStr"
)
// to simplify the types
type VarMap = templateStr.VariableMap
type FuncArray = templateStr.FuncArray
```

#### Construtor :

```go
parser := templateStr.New(funcArray, varMap)
```

<ul>
<li>
<details>
<summary><code>funcArray</code>: is a array of Functions you want to pass to be called in your text</summary><br>

```go
var funcArray FuncArray = FuncArray{meCustomFunc, otherCustomFunc}
```

</details>
</li>
<li>
<details>
<summary><code>varMap</code>: is a map of the Variables you want to pass to be called in your text</summary><br>

```go
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
    "Map": VarMap{
        "value": "Map in Map",
    },
    "MasterMap": VarMap{
        "SecondMap": VarMap{
            "value": "Map in Map in Map",
        },
    },
}
```

</details>
</li>
</ul>

#### Function :

```go
parser.Parse(text)
```

- `Parse(text: string) string` : parse all (variable, function, condition and switch)
- `ParseVariable(text: string) string` : parse Variable ; ${{variable}}
- `ParseFunction(text: string) string` : parse Function and Custom Function ; @{{function}}
- `ParseCondition(text: string) string` : parse Condition ; #{{var1 == var2: value1 || value2}}
- `ParseSwitch(text: string) string` : parse Switch ; ?{{var; value1=#0F0, value2=#00F, ..., default=#000}}
- `HasVariable(text: string) bool` : check if there are any Variable
- `HasFunction(text: string) bool` : check if there are any Function
- `HasCondition(text: string) bool` : check if there are any Condition
- `HasSwitch(text: string) bool` : check if there are any Switch
- `HasOne(text: string) bool` : check if there are one syntaxe

#### Exemple Syntaxe :

<details>
<summary><strong>Variable</strong></summary>
</br>

The syntax of the Variables is like if : 
- `${{variable}}`
- `${{Map.value}}`
- `${{MasterMap.SecondMap.value. ...}}`

if the value does not exist then `None` is return

```go
var varMap = VarMap{
    "variable": "yes",
}

text := "are you a variable : ${{variable}}"

parser := templateStr.New(FuncArray{}, varMap)

println(parser.Parse(text))
```

```go
var varMap = VarMap{
    "variable": VarMap{
        "value": "yes",
    },
}

text := "are you a variable : ${{variable.value}}"

parser := templateStr.New(FuncArray{}, varMap)

println(parser.Parse(text))
```

```go
variable := "yes"

println("are you a variable : " + variable)
```

The 3 codes will return

```text
are you a variable : yes
```

</details>

<details>
<summary><strong>Function</strong></summary>
</br>

The syntax of the Function is like if : `@{{function variable}}`

list of basic functions : 

- `@{{uppercase variable}}`
- `@{{uppercaseFirst variable}}`
- `@{{lowercase variable}}`
- `@{{swapcase variable}}`
- `@{{time}}`
- `@{{date}}`
- `@{{dateTime}}`

```go
var varMap = VarMap{
    "variable": "no",
}

text := "is lower case : @{{uppercase variable}}"

parser := templateStr.New(FuncArray{}, varMap)

println(parser.Parse(text))
```

```go
variable := "no"

println("is lower case : " + strings.ToUpper(variable))
```

The two codes will return

```text
is lower case : NO
```
</details>

<details>
<summary><strong>Custom Function</strong></summary>
</br>

The syntax of the Custom Function is like if : `@{{customFunction param1 param2 ...}}`

`Typing` can be used at the parameter level of custom functions

parameters to be passed in a list

the custom function must necessarily return a str

```go
func customFunc(list []Any) string{
    return strings.Replace(list[0], "no", "maybe", -1)
}

text := "are you a customFunction : @{{customFunc 'no'}}"

parser := templateStr.New(FuncArray{customFunc}, varMap)

println(parser.Parse(text))
```
The codes will return

```text
are you a customFunction : maybe
```

</details>

<details>
<summary><strong>Condition</strong></summary>
</br>

The syntax of the Condition is like if : 
- `#{{var1 == var2: value1 || value2}}`

comparator:
- `==`
- `!=`
- `<=`*
- `<`*
- `>=`*
- `>`*

*for this comparator the type `string` and `bool` are modified :
- `string` it's the number of characters that is compared ('text' = 4)
- `bool` it's the value in int that is compared (True = 1)


`var1` is compared with `var2`

`Typing` can be used at `var1` and `var2` level

```go
var varMap = VarMap{
    "var1": "no",
    "var2": "o2",
}

text := "are you a variable : #{{'test' == var2: yes || no}}"

parse := templateStr.New(FuncArray{}, varMap)

println(parser.Parse(text))
```
```go
var1 := "no"
var2 := "o2"

var text string

if "test" == var2 {
    text = "yes"
} else {
    text = "no"
}

println(text)
```

The 2 codes will return

```text
are you a variable : no
```

</details>

<details>
<summary><strong>Switch</strong></summary>
</br>

The syntax of the Switch is like if : 
- `?{{var; value1=#0F0, value2=#00F, ..., default=#000}}`
- `?{{var:type; value1=#0F0, value2=#00F, ..., default=#000}}`

`var` can be typed, if it is typed then all the `values` will be typed of the same type

type accept :
- `str`
- `int`
- `float`

```go
var varMap = VarMap{
    "variable": "yes",
}

text := "=( ?{{variable; yes=#A, no=#B, maybe=#C, default=#000}} )="

parse := templateStr.New(FuncArray{}, varMap)

println(parser.Parse(text))
```
```go
var varMap = VarMap{
    "variable": 42,
}

text := "=( ?{{variable:int; 42=#A, 32=#B, 22=#C, default=#000}} )="

parse := templateStr.New(FuncArray{}, varMap)

println(parser.Parse(text))
```
```go
var result string
variable := "yes"

switch variable {
case "yes":
    result = "#A"
case "no":
    result = "#B"
case "maybe":
    result = "#C"
default:
    result = "#000"
}

println(result)
```

The 3 codes will return

```text
=( #A )=
```

</details>

#### Typing :

| format                       | type    | description                                                       | return                 |
|------------------------------|---------|-------------------------------------------------------------------|------------------------|
| keyVariable                  | `*`     | is the key of the value in the dictionary pass to the constructor | value of `keyVariable` |
| \<b:True>                    | `bool`  |                                                                   | True                   |
| \<n:123>                     | `int`   |                                                                   | 123                    |
| \<n:123.4>                   | `float` |                                                                   | 123.4                  |
| "text" or 'text' or \`text\` | `str`   |                                                                   | text                   |

### TODO

- [ ] : Add exemple
- [x] : Add test

 

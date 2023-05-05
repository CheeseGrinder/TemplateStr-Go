
<div align="center">
    <h1>TemplateStr-Go</h1>
    <h2>TemplateStr allows to add variable, function, condition and switch in a string.</h2>
    <img src="https://img.shields.io/static/v1?label=Go&message=1.16%5E&color=22CFFA&style=flat-square&logo=Go&logoColor=00ADD8"/>
    <a href="https://github.com/CheeseGrinder/TemplateStr-Go/actions/workflows/python-app.yml">
        <img src="https://img.shields.io/github/actions/workflow/status/CheeseGrinder/TemplateStr-Go/go_test.yml?label=Test&style=flat-square"/>
    </a>
</div>

### Install :
```
go get -u github.com/CheeseGrinder/TemplateStr-Go/templateStr@latest
```

### Import :

```go
import (
    "github.com/CheeseGrinder/TemplateStr-Go/templateStr"
)
// to simplify the types
type VarMap = templateStr.VariableMap
type FuncArray = templateStr.FuncArray
```

### Construtor :

```go
parser := templateStr.New(funcArray, varMap)
```

<ul>
<li>
<details>
<summary><code>funcArray</code>: is an array of functions passed to the constructor that can be called in the parsed text</summary><br>

```go
var funcArray FuncArray = FuncArray{myCustomFunc, otherCustomFunc}
```

</details>
</li>

<li>
<details>
<summary><code>varMap</code>: is a map of variables passed to the constructor that can be used in the parsed text</summary><br>

```go
var varMap VarMap = VarMap{
    "foo": "bar",
    "str": "Jame",
    "int": 32,
    "float": 4.2,
    "bool": true,
    "array": []Any{"foo", 42},
    "Map": VarMap{
        "value": "Map in Map",
    },
    "Map1": VarMap{
        "Map2": VarMap{
            "value": "Map in Map in Map",
        },
    },
}
```

</details>
</li>
</ul>

### Function :

```go
result, err := parser.Parse(text)
```

- `Parse(text: string) (string, error)` : parse all (variable, function, condition and switch)
- `ParseVariable(text: string) (string, error)` : parse Variable ; ${variableName}
- `ParseFunction(text: string) (string, error)` : parse Function and Custom Function ; @{functionName}
- `ParseCondition(text: string) (string, error)` : parse Condition ; #{value1 == value2; trueValue | falseValue}
- `ParseSwitch(text: string) (string, error)` : parse Switch ; ?{var; value1::#0F0, value2::#00F, ..., _::#000}
- `HasOne(text: string) bool` : check if there are one syntaxe
- `HasVariable(text: string) bool` : check if there are any Variable
- `HasFunction(text: string) bool` : check if there are any Function
- `HasCondition(text: string) bool` : check if there are any Condition
- `HasSwitch(text: string) bool` : check if there are any Switch

### Exemple Syntaxe :

<ul>
<li>
<details>
<summary><strong>Variable</strong></summary>
</br>

The syntax of the Variables is like :
- `${variable}`
- `${Map.value}`
- `${MasterMap.SecondMap.value. ...}`
- `${variable[0]}`</br></br>

If the value does not exist an error is returned

<!-- V Be careful, it's not a "go" code, it's just to have some colour in the rendering -->
```go
//Example of parsing | is not code

name = "Jame"

"name is ${name}"
parse()
"name is Jame"
```

</details>
</li>

<li>
<details>
<summary><strong>Function</strong></summary>
</br>

The syntax of the Function is like : 
- `@{function; parameter}`
- `@{function}`</br></br>

Here is a list of the basic functions available  :

- `@{uppercase; variableName}`
- `@{uppercaseFirst; variableName}`
- `@{lowercase; variableName}`
- `@{swapcase; variableName}`
- `@{time}` HH/mm/ss
- `@{date}` DD/MM/YYYY
- `@{dateTime}` DD/MM/YYYY HH/mm/ss</br></br> 

<!-- V Be careful, it's not a "go" code, it's just to have some colour in the rendering -->
```go
//Example of parsing | is not code

name = "jame"

"name is @{uppercase; name}"
parse()
"name is JAME"
//=================================

"what time is it ? it's @{time}"
parse()
"what time is it ? it's 15:30:29"
```

</details>
</li>

<li>
<details>
<summary><strong>Custom Function</strong></summary>
</br>

The syntax of Custom function is the same as the basic functions, they can have 0,1 or more parameters : 
- `@{customFunction; param1 param2 variableName ...}`
- `@{customFunction}`</br></br>

The developer who adds his own function will have to document it

`Syntaxe Typing` can be used at the parameter level of custom functions

For developers :
- Parameters to be passed in a `list/vec/array`
- The custom function must necessarily return a `str/string`</br></br>

```go
func YourFunction(parameter []Any) string {

    //Your code

    return string
}
```


</details>
</li>

<li>
<details>
<summary><strong>Condition</strong></summary>
</br>

The syntax of the Condition is like :
- `#{value1 == value2; trueValue | falseValue}`</br></br>
  
comparator:
- `==`
- `!=`
- `<=` *
- `<` *
- `>=` *
- `>` *
</br></br>
<details>
<summary>* <i>for this comparator the type <code>string</code> and <code>bool</code> are modified</i> :</summary>

- `string` it's the number of characters that is compared ('text' = 4)
- `bool` it's the value in int that is compared (True = 1)

</details>

`value1` is compared with `value2`

`Syntaxe Typing` can be used at `value1` and `value2` level

<!-- V Be careful, it's not a "go" code, it's just to have some colour in the rendering -->
```go
//Example of parsing | is not code

name = "Jame"

"Jame is equal to James ? #{name == 'James'; Yes | No}"
parse()
"Jame is equal to James ? No"
```

</details>
</li>

<li>
<details>
<summary><strong>Switch</strong></summary>
</br>

The syntax of the Switch is like :
- `?{variableName; value1::#0F0, value2::#00F, ..., _::#000}`
- `?{type/variableName; value1::#0F0, value2::#00F, ..., _::#000}`</br></br>

The value of `variableName` is compared with all the `values*`,
if a `values*` is equal to the value of `variableName` then the value after the `::` will be returned.</br>
If no `values*` matches, the value after `_::` is returned

you can specify the type of `variableName`, but don't use `Syntaxe Typing`.</br>
If the type is specified then all `values*` will be typed with the same type.

syntax to specify the type of `variableName` :
- `str/variableName`
- `int/variableName`
- `float/variableName`</br></br>

<!-- V Be careful, it's not a "go" code, it's just to have some colour in the rendering -->
```go
//Example of parsing | is not code

name = "Jame"
yearsOld = 36

"how old is Jame ? ?{name; Jame::42 years old, William::36 years old, _::I don't know}"
parse()
"how old is Jame ? 42 years old"
//=================================

"who at 36 years old ? ?{int/yearsOld; 42::Jame !, 36::William !, _::I don't know}"
parse()
"who at 42 years old ? William !"
```

</details>
</li>
</ul>

#### !!! Warning if the syntax is not respected the text will not be parsed

### Syntaxe Typing :

Usable only in the <u>**custom function parameters**</u> and the two <u>**condition comparison values**</u>

| Format                       | Type    | Return                 | Note                                                                    |
|------------------------------|---------|------------------------|-------------------------------------------------------------------------|
| variableName                 | `*`     | value of `variableName`| Is the key of the value in the dictionary pass to the constructor       |
| b/True                       | `bool`  | True                   | Type the string True as `bool`                                          |
| i/123                        | `int`   | 123                    | Type the string 123 as type `int`                                       |
| f/123.4                      | `float` | 123.4                  | Type the string 123.4 as type `float`                                   |
| "text" or 'text' or \`text\` | `str`   | text                   | It just takes what's in quote, not to be interpreted as a variable name |
| ("test", i/56)               | `slice` | [test 56]              | Use typing for typed otherwise text will be used as variable name       |

```diff

This function takes as parameters a Bool, Int and String
+ @{MyCustomFunction; b/True i/15 "foo"}
- @{MyCustomFunction; True 15 foo}


+ #{"test" == "test"; Yes | No}
+ #{"56" == i/56; Yes | No}
- #{foo == 56; Yes | No}

```

### More
If you want another example you can look in the test file (`template_test.go`)

### TODO

- [ ] : Add exemple
- [x] : Add test

 

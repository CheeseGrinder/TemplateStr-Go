package templateStr

const REG_STR = `\"(?P<str_double>[^\"]+)\"|\'(?P<str_single>[^\']+)\'|\x60(?P<str_back>[^\x60]+)\x60`
const REG_BOOL = `b/(?P<bool>[Tt]rue|[Ff]alse)`
const REG_INT = `i/(?P<int>[0-9_]+)`
const REG_FLOAT = `f/(?P<float>[0-9_.]+)`
const REG_VAR = `(?P<variable>[\w._-]+)(?:\[(?P<index>[\d]+)])?`
const REG_LIST = `\((?P<list>[^\(\)]+)\)`

const REG_VARIABLE = `(?P<match>\${` + REG_VAR + `})`
const REG_FUNCTION = `(?P<match>@{(?P<functionName>[^{}\s]+)(?:; (?P<parameters>[^{}]+))?})`
const REG_CONDITION = `(?P<match>#{(?P<conditionValue1>[^{#}]+) (?P<conditionSymbol>==|!=|<=|<|>=|>) (?P<conditionValue2>[^{#}]+); (?P<trueValue>[^{}]+) \| (?P<falseValue>[^{}]+)})`
const REG_SWITCH = `(?P<match>\?{(?:(?P<type>str|int|float)/)?` + REG_VAR + `; (?P<values>(?:[^{}]+::[^{}]+){2,}), _::(?P<defaultValue>[^{}]+)})`

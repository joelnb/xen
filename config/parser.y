%{
package config

type pair struct {
  key string
  val interface{}
}

func setResult(l yyLexer, v map[string]interface{}) {
  l.(*lex).result = v
}
%}

%union{
  obj map[string]interface{}
  list []interface{}
  pair pair
  val interface{}
}

%token LexError
%token <val> String Number Literal

%type <obj> config config_items
%type <pair> pair
%type <val> array
%type <list> elements
%type <val> value config_value

%start config

%%

config: config_items
  {
    $$ = $1
    setResult(yylex, $$)
  }

config_items:
  {
    $$ = map[string]interface{}{}
  }
| pair
  {
    $$ = map[string]interface{}{
      $1.key: $1.val,
    }
  }
| config_items ';' pair
  {
    $1[$3.key] = $3.val
    $$ = $1
  }
| config_items '\n' pair
  {
    $1[$3.key] = $3.val
    $$ = $1
  }

pair: Literal '=' config_value
  {
    $$ = pair{key: $1.(string), val: $3}
  }

config_value:
  array
| value

array: '[' elements ']'
  {
    $$ = $2
  }

elements:
  {
    $$ = []interface{}{}
  }
| value
  {
    $$ = []interface{}{$1}
  }
| elements ',' value
  {
    $$ = append($1, $3)
  }

value:
  String
| Number
| Literal

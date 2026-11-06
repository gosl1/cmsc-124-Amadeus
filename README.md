# Amadeus

## Creators

- Ethan Sean T. Gapulan (Ethanerz)
- Marc Raven A. Sian (gosl1)

## Overview

Amadeus is a dynamically typed, general-purpose scripting language inspired by the visual novel and anime franchise Steins;Gate, built specifically to have execution-history awareness.

## Host language and build

- Host language: go 1.26.7
- Version metadata: go.mod
- Build: `./build.sh`

## Running it


| Command | What it does |
|---|---|
| `./run <file>` | [Executes a program. Available from Lab 4.] |
| `./run --tokenize <file>` | Scans the file and prints one token per line, including EOF.  |
| `./run --parse <file>` | - parses one expression per nonblank line and prints one AST per expression. |
| `./run --eval <file>` | [Evaluates each expression and prints its value.] |
| `./run` | Starts the parser REPL; parses each input line and prints its AST. |


Exit codes: 0 on successful run, 65 on a static error, 70 on a runtime error.

## File extension

`.ama` 

## Lexical structure


### Keywords

| Keyword | Purpose |
|---|---|
| res | variable declaration |
| dmail | print/output statement |
| div | if |
| con | else |
| timeleap | while |
| true | boolean literal — true |
| false | boolean literal — false |
| null | nil / absence-of-value literal |
| operation | function declaration |
| elpsy | return statement |
| readingsteiner | expose the execution-history record as a value (returns the log) |
| and | and logical operator |
| or | or logical operator |


### Operators


| Operator | Category | Operands | Associativity | Precedence |
|---|---|---|---|---|
| `=` | assignment | binary | right | 1 (lowest) |
| `or` | logical | binary | left | 2 |
| `and` | logical | binary | left | 3 |
| `== !=` | equality | binary | left | 4 |
| `< <= > >=` | comparison | binary | left | 5 |
| `+ -` | arithmetic | binary | left | 6 |
| `* / %` | arithmetic | binary | left | 7 |
| `! -` | unary | unary | right | 8 |

The parser currently supports only the operators listed here. Although the
scanner recognizes `[`, `]`, and `.`, indexing and field access are not yet
part of the expression grammar.

### Literals


| Kind | Syntax | Produces |
|---|---|---|
| number | `4`, `4.0` (integers and decimals; no leading-dot numbers, e.g. `.5` is invalid; a trailing dot with no following digit is its own token, so `3.toString` scans as NUMBER, DOT, IDENTIFIER) | numeric value |
| string | `"hello"`, double-quoted, single-line only (no multi-line strings); supports `\n`, `\"`, `\\` escape sequences | string value (decoded — lexeme keeps the raw escaped text, literal holds the decoded characters) |
| boolean | true / false | boolean value |
| nil | null | absence-of-value |


### Identifiers

- Start characters: letters (a-z, A-Z) and underscore (_)
- Continue characters: letters, digits (0-9), and underscore (_)
- Case-sensitive: yes
- Reserved patterns, length limits, or other restrictions: none beyond the reserved keyword list

### Comments

- Line comments: `//`, discarded and not counted, runs to end of line
- Block comments: not supported (documented decision — revisit only if the
  language later needs to annotate nested structure worth commenting on)
- Nesting: n/a
- Harness note: `comment_prefix` in `tests/lab*/manifest.json` is set to `//`

## Whitespace and termination

- Whitespace significant: no — spaces, tabs, and carriage returns are
  discarded and not counted as tokens
- The scanner discards newline tokens, but the parser uses token line numbers
  to separate expressions.
- The parser accepts one expression per nonblank line; statements and
  statement terminators are not implemented yet. Blank and comment-only lines
  are ignored. A file with no expressions is rejected.
- Block delimiters: braces `{ }`
- Grouping delimiters: parentheses `( )`
- Newlines are discarded but increment the line counter exactly once, at the
  point the newline character is consumed, so line numbers stay accurate
  even once multi-line constructs (strings, comments) are added later

## Token output format

```
Token(type=NUMBER, lexeme=42, literal=42, line=1)
Token(type=STRING, lexeme="hello\nworld", literal=hello\nworld, line=2)
Token(type=EOF, lexeme=, literal=null, line=3)
```

Each token prints as `Token(type=TYPE, lexeme=TEXT, literal=VALUE, line=N)`.
The fields are the token type, raw source lexeme, decoded/typed literal (or
`null` when there is no literal), and the 1-indexed line where the token
starts. String literal values are unquoted; backslashes, quotes, and embedded
newlines are escaped in the display so a token stays on one output line.
Numbers use fixed-point notation without unnecessary trailing zeroes. The
stream ends with one EOF token. Token lines are written to stdout only when
the whole file scans successfully.



## Grammar

```
expression  → assignment
assignment  → logic_or ( "=" assignment )?
logic_or    → logic_and ( "or" logic_and )*
logic_and   → equality ( "and" equality )*
equality    → comparison ( ( "==" | "!=" ) comparison )*
comparison  → term ( ( ">" | ">=" | "<" | "<=" ) term )*
term        → factor ( ( "-" | "+" ) factor )*
factor      → unary ( ( "/" | "*" | "%" ) unary )*
unary       → ( "!" | "-" ) unary
            | primary
primary     → "false"
            | "true"
            | "null"
            | NUMBER
            | STRING
            | IDENTIFIER
            | "(" expression ")"
```

Assignment is syntactically right-associative, but its left side must be a
variable reference; other assignment targets are rejected. Each nonblank,
non-comment line is parsed as one expression. Blank and comment-only lines are
ignored, but a file containing no expressions is rejected.

## Parse output format

```
Input:  1 + 2 * 3
Output: (+ 1 (* 2 3))
```

- Grouping expressions print as `(group expression)`.
- Numbers print using fixed-point notation without unnecessary trailing zeroes
  (for example, `4` and `3.14`). Values are parsed as `float64`.
- Strings print quoted, with special characters escaped.
- The parser prints one tree per expression line. If any line has a syntax
  error, it reports the errors to stderr, prints no trees to stdout, and exits
  with code 65.

## Semantics

### Values and types

[What runtime values exist, and how they are represented in the host
language.]

### Value printing

- Numbers: [e.g. 5 rather than 5.0]
- Nil: [spelling]
- Strings: [with or without quotes]

### Truthiness

[The complete rule. Which values are false in a condition; everything else is
true.]

### Operator semantics

- Arithmetic: [accepted operand types]
- `+` on strings: [concatenation, error, or coercion]
- Mixed types: [what happens]
- Comparison: [accepted operand types]
- Equality across types: [false, or an error]
- Division by zero: [value produced, or runtime error]

### Scope and bindings

- Redeclaration in the same scope: [allowed or an error]
- Uninitialized variable holds: [value]
- Shadowing: [behavior]
- Undefined name: [static error with exit 65, or runtime error with exit 70]

### Control flow and functions

- Logical operators return: [booleans, or the operand]
- Dangling else binds to: [which if]
- Closure capture of a loop variable: [per iteration, or shared]
- Function with no return statement produces: [value]
- Arity mismatch: [message and exit code]

## Native functions


| Name | Arguments | Returns | Notes |
|---|---|---|---|
| [name] | [count and types] | [type] | [caveats] |


## Errors and diagnostics

Scanner errors are reported to stderr in the format
`[Attractor Field Violation] line N: message`. Examples:

```
[Attractor Field Violation] line 1: Unexpected character '$'.
[Attractor Field Violation] line 2: Invalid escape sequence '\q'.
[Attractor Field Violation] line 2: Unterminated string.
[Attractor Field Violation] line 3: Malformed number '999999999999999999999999999999999999999'.
```

Parser errors are reported to stderr with the unexpected token, or `end` when
the parser reaches EOF. For example:

```
[line 1] Error at ')': Expect expression.
[line 1] Error at end: Expect expression.
```

The scanner continues scanning after lexical errors and reports each one. If a
file has any scanner error, it prints no tokens to stdout and exits 65; parsing
does not proceed. The parser reports a syntax error for each malformed
expression line and continues with later lines. If any parser error occurs, a
rejected file prints nothing to stdout and exits 65. A successful scan or
parse exits 0.

| Failure | Exit code |
|---|---|
| Unterminated string literal | 65 |
| Character that can't begin any lexeme | 65 |
| Invalid string escape | 65 |
| Malformed or out-of-range number | 65 |
| Parser syntax error | 65 |



## Testing conventions


| Folder | Activity | Mode | Flag |
|---|---|---|---|
| tests/lab1 | Scanner | sidecar | `--tokenize` |
| tests/lab2 | Parser | sidecar | `--parse` |
| tests/lab3 | Evaluator | inline | `--eval` |
| tests/lab4 | Context | inline | none |
| tests/lab5 | Functions | inline | none |


```
[specific tests]...
```

Run locally with:

```bash
curl -sSL https://raw.githubusercontent.com/WhiteLicorice/cmsc-124-harness/v1.1/run_tests.py -o run_tests.py
./build.sh
python run_tests.py tests/lab1
```

## Sample code

```
[a short program]
```

Output:

```
[its output]
```

## Design rationale

- Dynamically typed by requirement; no type-checking occurs at scan time.
- History-awareness (recording branch/loop/call outcomes and exposing them
  as a queryable value) was chosen over a purely thematic keyword layer
  after establishing that renamed control-flow keywords alone (e.g. an
  `if`/`else` pair under different names) would be functionally identical
  to ordinary conditionals and not a real language feature. This after all,
  is the choice of Steins;Gate.
- [keyword choices yet to be done]

## Known limitations

- No block comments; no multi-line strings.
- History-log growth is currently unbounded within a scope — no cap or
  clear mechanism decided yet.

## Changelog


| Activity | What changed in the language |
|---|---|
| Lab 1 | Initial lexical structure, run contract, dynamically-typed requirement, and history-awareness design direction established. Keyword names left open. |
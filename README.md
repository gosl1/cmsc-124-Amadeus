# Amadeus

## Creators

- Ethan Sean T. Gapulan (Ethanerz)
- Marc Raven A. Sian (gosl1)

## Overview

Amadeus is a dynamically typed, general-purpose scripting language inspired by the visual novel and anime franchise Steins;Gate, built specifically to have execution-history awareness. Amadeus natively treats the program's running execution path as an immutable, queryable "worldline." Rather than relying strictly on current variable states, developers can native-query the runtime history using custom syntax inspired by Steins;Gate to dynamically reroute logic based on what statements were executed and in what exact order.

## Host language and build

- Host language: go 1.26.7
- Version metadata: go.mod
- Build: `./build.sh`

## Running it


| Command | What it does |
|---|---|
| `./run <file>` | [Executes a program. Available from Lab 4.] |
| `./run --tokenize <file>` | [Prints the token stream.] |
| `./run --parse <file>` | [Prints the parsed tree.] |
| `./run --eval <file>` | [Evaluates each expression and prints its value.] |
| `./run` | [Starts the REPL.] |


Exit codes: 0 on successful run, 65 on a static error, 70 on a runtime error.

## File extension

`[.ama]` 

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
| `+` | arithmetic | binary | left | [TBD — set in Lab 2] |
| `-` | arithmetic | binary | left | 5 |
| `*` | arithmetic | binary | left | 6 |
| `/` | arithmetic | binary | left | 6 |
| `%` | arithmetic (modulo) | binary | left | 6 |
| `=` | assignment | binary | right | 1 |
| `== != < <= > >=` | comparison | binary | left | 4 |
| `!` | logical (not) | unary | right | 7 |
| `[` `]` | indexing (array/log access) | binary | left | 8 (tightest) |
| `.` | field access | binary | left | 8 (tightest) |

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
- Statement terminator: newline (no semicolons)
- Block delimiters: braces `{ }`
- Grouping delimiters: parentheses `( )`
- Newlines are discarded but increment the line counter exactly once, at the
  point the newline character is consumed, so line numbers stay accurate
  even once multi-line constructs (strings, comments) are added later

## Token output format

```
[TBD format name](type=[TBD], lexeme=[TBD], literal=[TBD], line=1)
```

Fields, left to right: token type (category, e.g. NUMBER), the lexeme (raw
source text), the literal value (decoded/typed value, or absent for
non-literal tokens), and the 1-indexed line number the token started on.
One token per line, EOF included as the final token of every stream. Format
must be frozen before the first `.expected` file is committed — see
changelog if it changes later.



## Grammar

```
[Your complete context-free grammar, current as of the latest activity.
Unambiguous, with precedence and associativity encoded in rule structure.]
```

## Parse output format

```
[one line of real --parse output, e.g. (+ 1.0 (* 2.0 3.0))]
```

- Groupings print as: [form]
- Numbers print as: [form]

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

Message format (draft, refine once real examples exist):

```
[TBD prefix] Unexpected character '$' at line 12
[TBD prefix] String has no closing quote, starting at line 7
```

Both are reported on stderr; scanning continues afterward so multiple
problems in one file are all reported; the file exits 65 once scanning
finishes. Nothing about a rejected file is printed to stdout. A clean scan
exits 0.

| Failure | Exit code |
|---|---|
| Unterminated string literal | 65 |
| Character that can't begin any lexeme | 65 |
| [Third rejection case — TBD once decided per assignment requirement] | 65 |



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
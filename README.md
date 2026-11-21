# Amadeus

## Creators

- Ethan Sean T. Gapulan (Ethanerz)
- Marc Raven A. Sian (gosl1)

## Overview

Amadeus is a small expression language inspired by Steins;Gate. Its scanner
turns source text into tokens, its recursive-descent parser builds expression
trees, and its evaluator computes values from those trees.

## Host language and build

- Host language: go 1.26.7
- Version metadata: go.mod
- Build: `./build.sh`

## Running it


| Command | What it does |
|---|---|
| `./run <file>` | Not implemented yet; plain file execution is reserved for a later activity. |
| `./run --tokenize <file>` | Scans the file and prints one token per line, including EOF. |
| `./run --parse <file>` | Parses one expression per nonblank, non-comment line and prints one AST per expression. |
| `./run --eval <file>` | Evaluates one expression per nonblank, non-comment line and prints one value per expression. |
| `./run` | Starts the REPL; evaluates each entered expression and prints its value. |

Exit codes:

| Code | Meaning |
|---|---|
| `0` | The requested scan, parse, or evaluation completed successfully. |
| `64` | Invalid command-line arguments. |
| `65` | File read, lexical, or syntax error. |
| `70` | Runtime evaluation error in `--eval` file mode. |

`--tokenize` and `--parse` print results only when the file is accepted.
`--parse` reports syntax errors to stderr, continues checking later expression
lines, and prints no ASTs if any line has a syntax error. `--eval` prints each
successful result as it evaluates the file and stops at the first runtime
error. The no-argument REPL reads one expression per entered line, prints its
value, and returns to the prompt after scanner, syntax, or runtime errors. EOF
(Ctrl+D on Unix-like systems or Ctrl+Z on Windows) ends the REPL.

## File extension

`.ama`

## Lexical structure


### Keywords

| Keyword | Token type | Use in the current expression grammar |
|---|---|---|---|
| `res` | `RES` | Reserved; not an expression form yet. |
| `dmail` | `DMAIL` | Reserved; not an expression form yet. |
| `div` | `DIV` | Reserved; not an expression form yet. |
| `con` | `CON` | Reserved; not an expression form yet. |
| `timeleap` | `TIMELEAP` | Reserved; not an expression form yet. |
| `true` | `TRUE` | Boolean literal. |
| `false` | `FALSE` | Boolean literal. |
| `null` | `NULL` | Null literal. |
| `operation` | `OPERATION` | Reserved; not an expression form yet. |
| `elpsy` | `ELPSY` | Reserved; not an expression form yet. |
| `readingsteiner` | `READINGSTEINER` | Reserved; not an expression form yet. |
| `and` | `AND` | Logical conjunction. |
| `or` | `OR` | Logical disjunction. |


### Operators


| Operator | Category | Operands | Associativity | Precedence |
|---|---|---:|---|---:|
| `=` | Assignment | Binary | Right | 1 (loosest) |
| `or` | Logical | Binary | Left | 2 |
| `and` | Logical | Binary | Left | 3 |
| `==`, `!=` | Equality | Binary | Left | 4 |
| `>`, `>=`, `<`, `<=` | Comparison | Binary | Left | 5 |
| `+`, `-` | Arithmetic | Binary | Left | 6 |
| `*`, `/`, `%` | Arithmetic | Binary | Left | 7 |
| `!`, `-` | Logical / arithmetic | Unary | Right | 8 |

The scanner also recognizes `(`, `)`, `{`, `}`, `[`, `]`, `.`, and `,`.
Parentheses group expressions. Braces, brackets, dot, and comma currently
have no expression grammar rule. The scanner uses lookahead to recognize
two-character operators such as `==`, `!=`, `<=`, and `>=`.

### Literals


| Kind | Syntax | Value produced by parsing |
|---|---|---|
| Number | `4`, `4.0`; digits, optionally followed by a decimal point and one or more digits | `float64` |
| String | Double-quoted and single-line; supports `\n`, `\"`, and `\\` | Decoded string |
| Boolean | `true`, `false` | `bool` |
| Null | `null` | `nil` |

Numbers do not include a leading decimal point. A decimal point is consumed as
part of a number only when followed by a digit; otherwise it is a separate
`DOT` token. `strconv.ParseFloat` errors, including out-of-range values, are
reported as scanner errors.


### Identifiers

- Character processing is byte-by-byte.
- Start characters: ASCII letters (`a-z`, `A-Z`) and underscore (`_`).
- Continue characters: ASCII letters, digits (`0-9`), and underscore.
- Case-sensitive; the complete identifier is consumed before keyword lookup.
- Non-ASCII characters are not identifier characters.
- There is no explicit identifier length limit.

### Comments

- Line comments: `//`, discarded through (but not including) the newline.
- The newline after a comment is handled by the main scanner loop and
  increments the line counter.
- Block comments are not supported.
- Comments may appear after code on the same line. The test harness uses `//`
  as the inline annotation prefix in inline-mode test folders.

## Whitespace and termination

- Spaces, tabs, and carriage returns are insignificant and discarded.
- Newlines are discarded by the scanner and increment the line counter. In
  `--parse` and `--eval` file mode, each nonblank, non-comment source line is a
  separate expression.
- Expressions do not require a terminator. Statements and block syntax are
  not part of the current parser.
- Parentheses delimit grouped expressions. Braces and brackets are tokenized
  but are not parsed as blocks or grouping forms.
- Strings cannot span lines; an unterminated or multiline string is reported.

## Token output format

```
Token(type=NUMBER, lexeme=42, literal=42, line=1)
Token(type=STRING, lexeme="hello\nworld", literal=hello\nworld, line=2)
Token(type=EOF, lexeme=, literal=null, line=3)
```

Each token prints as `Token(type=TYPE, lexeme=TEXT, literal=VALUE, line=N)`.
The fields are the token type, raw source lexeme, typed/decoded literal (or
`null` when there is no literal), and the 1-indexed line where the token
starts. String literals are unquoted in the literal field; `Token.String()`
escapes backslashes, quotes, and decoded newlines so each token stays on one
output line. Numbers use fixed-point notation without unnecessary trailing
zeroes. Keywords and operators have no literal payload, displayed as `null`.
The stream always ends with one EOF token, even after scanner errors.



## Grammar

The parser uses recursive descent, with a parsing function for each precedence
level. This grammar describes the expression forms currently accepted by
`--parse` and `--eval`:

```text
expression  -> assignment
assignment  -> logic_or ( "=" assignment )?
logic_or    -> logic_and ( "or" logic_and )*
logic_and   -> equality ( "and" equality )*
equality    -> comparison ( ( "==" | "!=" ) comparison )*
comparison  -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
term        -> factor ( ( "-" | "+" ) factor )*
factor      -> unary ( ( "/" | "*" | "%" ) unary )*
unary       -> ( "!" | "-" ) unary | primary
primary     -> "false" | "true" | "null"
             | NUMBER | STRING | IDENTIFIER
             | "(" expression ")"
```

Binary operators within each level are left-associative because the parser
accumulates them in a loop. Assignment is right-associative because its right
side recursively parses another assignment. Its left side must be a variable
reference; other assignment targets are syntax errors.

Precedence, from lowest to highest, is assignment, `or`, `and`, equality,
comparison, addition/subtraction, multiplication/division/modulo, unary, and
primary expressions. The scanner omits newline tokens but retains each
token's source line. File mode parses each nonblank, non-comment line as one
expression; blank and comment-only lines are skipped. A malformed line does
not prevent parsing later lines. A file with no expressions (empty or only
comments and blank lines) is rejected with a syntax error.

## Parse output format

The parser prints an AST in prefix form, with one tree per source expression:

```text
Input:  1 + 2 * 3
Output: (+ 1 (* 2 3))
```

- Grouping expressions print as `(group expression)`, for example
  `(1 + 2) * 3` prints `(* (group (+ 1 2)) 3)`.
- Numbers use fixed-point notation without unnecessary trailing zeroes:
  `4.0` prints as `4`.
- Strings print quoted, with special characters escaped.
- Variables print by name; assignment prints as `(= name value)`.

If any expression line has a syntax error, the CLI reports the errors to
stderr, prints no ASTs to stdout, and exits with code `65`. Scanner errors also
exit `65` and prevent parsing.

## Lab 1 scope

The scanner's job is to classify source characters into tokens. Parsing turns
those tokens into expression ASTs; evaluation is a separate stage that
computes values.


## Errors and diagnostics

Scanner errors are reported to stderr in the format
`[Attractor Field Violation] line N: message`. Examples:

```
[Attractor Field Violation] line 1: Unexpected character '$'.
[Attractor Field Violation] line 2: Invalid escape sequence '\q'.
[Attractor Field Violation] line 2: Unterminated string.
[Attractor Field Violation] line 3: Malformed number '999999999999999999999999999999999999999'.
```

The scanner continues scanning after lexical errors and reports each one. If a
file has any scanner error, `--tokenize` prints no tokens to stdout and exits
65. A successful file scan prints all tokens, including EOF, and exits 0.

Parser diagnostics include the 1-based source line and the unexpected token,
or `end` when the parser reaches EOF:

```text
[line 1] Error at end: Expect expression.
[line 2] Error at ')': Expect expression.
```

| Failure | Exit code |
|---|---|
| Invalid command-line arguments | 64 |
| File read error | 65 |
| Unterminated string, invalid escape, unexpected character, or malformed number | 65 |
| Parser syntax error (including invalid assignment target) | 65 |
| Runtime evaluation error in `--eval` file mode | 70 |

## Testing conventions

| Folder | Activity | Mode | Flag |
|---|---|---|---|
| `tests/lab1` | Scanner | sidecar | `--tokenize` |
| `tests/lab2` | Parser | sidecar | `--parse` |
| `tests/lab3` | Evaluator | sidecar | `--eval` |

Each sidecar test uses an `.ama` source file and a matching `.expected` stdout
file. Error cases also have an `.exit` file with the expected process exit
code. The shared harness compares stdout and exit status. Sidecar expectations
do not check the contents of stderr.

Lab 1 has 47 cases covering comments, identifiers and keywords, whitespace
and line tracking, numbers, operators, strings, and lexical errors. Lab 2 has
41 cases grouped into:

- `expressions`: precedence, associativity, grouping, literals, unary
  expressions, long/nested expressions, and multiple expressions per file.
- `identifiers`: identifier forms, keyword boundaries, and adjacent operator
  tokens.
- `literals`: number, string, boolean, and null AST output.
- `errors`: malformed expressions, invalid assignment targets, empty input,
  unmatched delimiters, lexical failures, and recovery across expression
  lines.

Lab 3 currently has 45 sidecar cases for literal values, grouping, logical and
unary operators, arithmetic, comparisons, and string addition.

Run a suite locally after building the interpreter:

```sh
curl -sSL https://raw.githubusercontent.com/WhiteLicorice/cmsc-124-harness/v1.1/run_tests.py -o run_tests.py
./build.sh
python3 run_tests.py tests/lab1
python3 run_tests.py tests/lab2
python3 run_tests.py tests/lab3
```

The GitHub Actions workflow runs all three suites (Lab 1, Lab 2, and Lab 3)
on each push. Lab suites can also be run locally with the commands above.

## Sample code

Source file:

```text
1 + 2 * 3
(1 + 2) * 3
a = b = 3
```

Command: `./run --parse sample.ama`

Output:

```text
(+ 1 (* 2 3))
(* (group (+ 1 2)) 3)
(= a (= b 3))
```

## Design rationale

- The recursive-descent functions mirror the grammar's precedence levels,
  making the binding order visible in the parser.
- Loops at binary precedence levels encode left associativity; recursive
  parsing on the right side of assignment encodes right associativity.
- Parentheses are retained as `Grouping` nodes so the AST output makes
  explicit grouping visible.
- The scanner records token line numbers even though it does not emit newline
  tokens. File parsing uses those lines to support one expression and one AST
  per line, and diagnostics can identify the source line.

## Known limitations

- No block comments or multi-line strings.
- Identifiers use ASCII character classes; Unicode identifiers are not
  supported.
- The parser accepts expressions only; declarations, statements, blocks,
  indexing, property access, and function calls are not in the current grammar.
- File parsing treats each nonblank source line as a complete expression;
  expressions cannot continue across lines.
- Diagnostics include line numbers but not columns.

## Changelog


| Activity | What changed |
|---|---|
| Lab 1 | Scanner token rules, `--tokenize` file scanning, scanner REPL, and scanner test suite. |
| Lab 2 | Added an expression grammar and recursive-descent parser, parenthesized AST output, `--parse` file mode, and line-based parse error recovery. |
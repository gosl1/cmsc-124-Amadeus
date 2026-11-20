# Amadeus

## Creators

- Ethan Sean T. Gapulan (Ethanerz)
- Marc Raven A. Sian (gosl1)

## Overview

Amadeus is a language project inspired by the visual novel and anime franchise Steins;Gate. This Lab 1 submission focuses on its scanner: converting source characters into tokens.

## Host language and build

- Host language: go 1.26.7
- Version metadata: go.mod
- Build: `./build.sh`

## Running it


| Command | What it does |
|---|---|
| `./run --tokenize <file>` | Scans the file and prints one token per line, including EOF.  |
| `./run --parse <file>` | Separate parser command; not part of the Lab 1 scanner demonstration. |
| `./run` | Starts the scanner REPL; scans each input line independently. |

For `--tokenize`, exit code 0 means scanning succeeded; exit code 65 indicates
a file read or scanner error.

The scanner REPL reads and scans one line at a time with a fresh scanner for
each line. It prints a line's tokens only if that scan succeeds, then continues
after lexical errors. EOF (Ctrl+D on Unix-like systems or Ctrl+Z on Windows)
ends the REPL cleanly.

## File extension

`.ama` 

## Lexical structure


### Keywords

| Keyword | Token type |
|---|---|
| `res` | `RES` |
| `dmail` | `DMAIL` |
| `div` | `DIV` |
| `con` | `CON` |
| `timeleap` | `TIMELEAP` |
| `true` | `TRUE` |
| `false` | `FALSE` |
| `null` | `NULL` |
| `operation` | `OPERATION` |
| `elpsy` | `ELPSY` |
| `readingsteiner` | `READINGSTEINER` |
| `and` | `AND` |
| `or` | `OR` |


### Operators


The scanner recognizes `(`, `)`, `{`, `}`, `[`, `]`, `.`, `,`, `+`, `-`, `*`,
`/`, `%`, `=`, `==`, `!`, `!=`, `<`, `<=`, `>`, and `>=`. It uses lookahead
to recognize the two-character operators. The scanner recognizes `[` `]` and
`.`; this only describes tokenization, not parser support or runtime behavior.

### Literals


| Kind | Syntax | Produces |
|---|---|---|
| number | `4`, `4.0`; digits, optionally followed by a decimal point and one or more digits | `float64` literal |
| string | Double-quoted and single-line; supports `\n`, `\"`, and `\\` | Decoded string literal; lexeme retains source spelling |
| boolean keyword | `true`, `false` | No scanner literal payload |
| null keyword | `null` | No scanner literal payload |

Numbers do not include a leading decimal point. A decimal point is consumed as
part of a number only when followed by a digit; otherwise it is a separate
`DOT` token. `strconv.ParseFloat` errors, including out-of-range values, are
reported as scanner errors.


### Identifiers

- Character processing is byte-by-byte.
- Start characters: ASCII letters (`a-z`, `A-Z`) and underscore (`_`).
- Continue characters: ASCII letters, digits (`0-9`), and underscore.
- The complete identifier is consumed before a case-sensitive keyword lookup.
- Non-ASCII characters are not identifier characters.

### Comments

- Line comments: `//`, discarded through (but not including) the newline.
- The newline after a comment is handled by the main scanner loop and
  increments the line counter.
- Block comments are not supported.

## Whitespace and termination

- Spaces, tabs, and carriage returns are discarded.
- Newlines are discarded and increment the line counter when consumed.
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
zeroes. The stream always ends with one EOF token, even after scanner errors.



## Lab 1 scope

The scanner classifies source text into tokens; it does not establish variable,
function, scope, truthiness, or runtime-value semantics. Those topics are
outside this scanner-focused submission. `--parse` remains a separate command
and is not part of the scanner REPL.


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

| Failure | Exit code |
|---|---|
| Unterminated string literal | 65 |
| Character that can't begin any lexeme | 65 |
| Invalid string escape | 65 |
| Malformed or out-of-range number | 65 |



## Testing conventions


## Testing conventions

## Testing conventions

| Folder | Activity | Mode | Flag |
|---|---|---|---|
| `tests/lab1` | Scanner | sidecar | `--tokenize` |

Lab 1 cases use `.ama` source files with matching `.expected` output files. Error cases also specify expected exit codes in `.exit` files.

Coverage includes comments, identifiers and keywords, whitespace and line tracking, numbers, operators, strings, and lexical errors.

The shared course harness checks standard output and process exit codes against the sidecar expectations. Error cases verify that rejected input produces no standard output and exits with code `65`. Scanner diagnostics are written to standard error; their exact contents are not checked by the sidecar harness.

Run the Lab 1 suite locally with:

```bash
curl -sSL https://raw.githubusercontent.com/WhiteLicorice/cmsc-124-harness/v1.1/run_tests.py -o run_tests.py
./build.sh
python3 run_tests.py tests/lab1
```

GitHub Actions automatically builds the project and runs the Lab 1 test suite on every push.
## Known limitations

- No block comments or multi-line strings.
- Identifiers use ASCII character classes; Unicode identifiers are not
  supported.

## Changelog


| Activity | What changed |
|---|---|
| Lab 1 | Scanner token rules, `--tokenize` file scanning, scanner REPL, and scanner test suite. |
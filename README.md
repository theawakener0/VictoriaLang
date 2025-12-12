# Victoria Programming Language

![Queen Victoria](media/QueenVictoria.jpg)

![Victoria Language](https://img.shields.io/badge/Victoria-Language-blueviolet?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen?style=for-the-badge)

**Victoria** is a modern, artistic programming language designed to be expressive and readable. It combines the structural clarity of C-style languages with the ease of use found in Python. Whether you're building simple scripts, data processing tools, or network servers, Victoria provides a clean and enjoyable syntax.

## Why Victoria?

Victoria is designed to be the language you **wish you had learned first**. Here's why Victoria stands out as a potential game-changer:

### The Best of All Worlds

| Aspect | Victoria's Approach |
|--------|---------------------|
| **Syntax** | Clean like Python, structured like JavaScript |
| **Learning Curve** | Gentle for beginners, powerful for experts |
| **Error Messages** | Rust-inspired diagnostics that teach, not frustrate |
| **Type System** | Dynamic typing with clear error guidance |
| **Functional Style** | First-class functions + lambdas + map/filter/reduce |

### Key Advantages

- **Zero Friction Start**: No complex toolchains, package managers, or configuration files needed
- **Self-Documenting Errors**: Every error message is a mini-tutorial explaining what went wrong and how to fix it
- **Expressive Syntax**: Write code that reads like pseudocode but runs like production code
- **Batteries Included**: Built-in JSON, Time, Math modules—no external dependencies for common tasks
- **Modern Features**: Lambda expressions, spread operators, string interpolation out of the box
- **Familiar Yet Fresh**: If you know JavaScript or Python, you're already 80% there

### Perfect For

- **Students** learning programming fundamentals without fighting cryptic errors
- **Hobbyists** prototyping ideas quickly without boilerplate
- **Educators** teaching programming concepts with clear, readable syntax
- **Developers** scripting automation tasks with minimal setup
- **Anyone** who wants a language that respects their time and intelligence

## What's New

- **Lambda Expressions**: Concise `x => x * 2` syntax for anonymous functions
- **JSON Module**: Parse, stringify, and validate JSON with `json.parse()`, `json.stringify()`, `json.valid()`
- **Time Module**: Get current time, format dates, and manipulate timestamps
- **Extended Math**: New `floor`, `ceil`, `round`, `min`, `max`, `random`, `tan`, `log`, `log10` functions

## Features

### Core Language
- **C-like Structure**: Familiar curly brace `{}` syntax for blocks.
- **Python-like Readability**: Clean, minimal syntax with optional semicolons.
- **Dynamic Typing**: No need to declare types explicitly.
- **Const Variables**: Immutable variables with `const` to prevent accidental reassignment.

### Functions & Lambdas
- **First-Class Functions**: Functions are values that can be passed around and returned.
- **Lambda Expressions**: Concise arrow function syntax `x => x * 2` and `(a, b) => a + b`.
- **Functional Programming**: `map`, `filter`, `reduce` with lambda support.

### Data Structures
- **Object-Oriented**: Support for `structs` and methods to organize your code.
- **Array Slicing**: Python-style slicing with `arr[1:3]`, `arr[:5]`, `arr[2:]`.
- **Spread Operator**: Merge and expand arrays with `[...arr1, ...arr2]`.

### Built-in Modules
- **Math Module**: `floor`, `ceil`, `round`, `min`, `max`, `random`, `sin`, `cos`, `tan`, `log`, and more.
- **JSON Module**: `json.parse()`, `json.stringify()`, `json.valid()` for JSON handling.
- **Time Module**: `time.now()`, `time.format()`, `time.parse()`, date/time manipulation.

### Control Flow & Operators
- **Versatile Loops**: Supports `for`, `for-in`, `while` loops with `break`/`continue`.
- **Modern Operators**: Ternary operator, compound assignments (`+=`, `-=`, `++`, `--`).
- **Switch Statements**: Clean pattern matching with `switch`/`case`/`default`.

### Developer Experience
- **Rust-Style Error Messages**: Beautiful, helpful error messages with code snippets and suggestions.
- **String Interpolation**: Embed expressions in strings with `${expr}`.
- **Rich Standard Library**: Built-in support for strings, arrays, hashes, and networking.

## Installation

To build Victoria from source, you need to have [Go](https://golang.org/) installed on your machine.

1.  Clone the repository:
    ```bash
    git clone https://github.com/theawakener0/VictoriaLang.git
    cd VictoriaLang 
    ```

2.  Build the interpreter:
    ```bash
    go build -o victoria cmd/victoria/main.go
    ```

3.  (Optional) Add `victoria` to your PATH for easy access.

## Usage

### Interactive REPL

Start the Read-Eval-Print Loop (REPL) by running the executable without arguments:

```bash
./victoria
```

```victoria
>> let name = "World"
>> print("Hello, ${name}!")
Hello, World!
```

### Running Scripts

Execute a Victoria script (`.vc` file) by passing the filename:

```bash
./victoria examples/hello.vc
```

## Quick Examples

### Lambda Expressions
```victoria
// Concise arrow function syntax
let double = x => x * 2
let add = (a, b) => a + b

print(double(5))    // 10
print(add(3, 4))    // 7

// Perfect for functional programming
let nums = [1, 2, 3, 4, 5]
let doubled = map(nums, x => x * 2)        // [2, 4, 6, 8, 10]
let evens = filter(nums, x => x % 2 == 0)  // [2, 4]
let sum = reduce(nums, (a, b) => a + b, 0) // 15
```

### JSON Module
```victoria
// Parse JSON strings
let data = json.parse('{"name": "Victoria", "version": 1}')
print(data["name"])  // Victoria

// Convert objects to JSON
let obj = {"users": ["Alice", "Bob"], "count": 2}
print(json.stringify(obj))

// Validate JSON
print(json.valid('{"valid": true}'))  // true
print(json.valid('invalid json'))     // false
```

### Time Module
```victoria
// Current time
print(time.now())    // 2024-01-15T10:30:45Z
print(time.nowMs())  // 1705315845123 (milliseconds)

// Formatting
let now = time.now()
print(time.format(now, "YYYY-MM-DD"))       // 2024-01-15
print(time.format(now, "HH:mm:ss"))         // 10:30:45
print(time.format(now, "MMMM DD, YYYY"))    // January 15, 2024

// Extract components
print(time.year(now))     // 2024
print(time.weekday(now))  // Monday
```

### Math Functions
```victoria
// Rounding functions
print(math.floor(3.7))  // 3
print(math.ceil(3.2))   // 4
print(math.round(3.5))  // 4

// Min/Max
print(math.min(5, 3, 8, 1))  // 1
print(math.max(5, 3, 8, 1))  // 8

// Random numbers
print(math.random())       // 0.0 to 1.0
print(math.random(1, 100)) // 1 to 100
```

### Variables & Data Types
```victoria
// String interpolation
let name = "Victoria"
print("Hello, ${name}!")

// Constant variables
const PI = 3.14159
const MAX_SIZE = 100
// PI = 3.0  // ERROR: cannot reassign constant

// Range-based loops
for i in 0..5 {
    print(i)  // 0, 1, 2, 3, 4
}

// Array slicing
let arr = [1, 2, 3, 4, 5]
print(arr[1:3])  // [2, 3]
print(arr[:2])   // [1, 2]
print(arr[3:])   // [4, 5]

// Spread operator
let a = [1, 2, 3]
let b = [4, 5, 6]
let merged = [...a, ...b]  // [1, 2, 3, 4, 5, 6]

// Functional programming with traditional functions
let nums = [1, 2, 3, 4, 5]
let doubled = map(nums, define(x) { x * 2 })
let evens = filter(nums, define(x) { x % 2 == 0 })

// Ternary operator
let status = age >= 18 ? "adult" : "minor"

// Switch statements
switch (day) {
    case 1: { print("Monday") }
    case 2: { print("Tuesday") }
    default: { print("Other day") }
}
```

## Documentation

For a deep dive into the language syntax and features, check out the [Language Documentation](docs/LANGUAGE.md).

Key topics include:
- [Variables & Data Types](docs/LANGUAGE.md#variables)
- [Control Flow](docs/LANGUAGE.md#control-flow)
- [Functions & Lambdas](docs/LANGUAGE.md#functions)
- [Structs](docs/LANGUAGE.md#structs)
- [Math Module](docs/LANGUAGE.md#math-module)
- [JSON Module](docs/LANGUAGE.md#json-module)
- [Time Module](docs/LANGUAGE.md#time-module)
- [Built-in Functions](docs/LANGUAGE.md#built-in-functions)
- [Error Handling](docs/LANGUAGE.md#error-handling)

## Rust-Inspired Error Messages

Victoria features **industry-leading error messages** inspired by the Rust programming language. We believe that clear, actionable error messages are essential for developer productivity and learning. Our error system draws from Rust's legendary compiler diagnostics to make debugging a pleasure, not a pain.

### Why Our Errors Are Different

Most languages show you cryptic messages like `undefined: x` or `TypeError`. Victoria shows you:

```
error[E0002]: identifier not found: unknownVar
 --> main.vc:4:9
  |
3 | // Using undefined variable
4 | let x = unknownVar + 5
  |         ^^^^^^^^^^ identifier not found: unknownVar
5 |
  |
  = help: did you mean to declare 'unknownVar' with 'let unknownVar = ...'?
```

### Error Enhancement Features

| Feature | Description |
|---------|-------------|
| **Color-Coded Output** | Errors in red, context in cyan, help in green—instantly parse severity |
| **Precise Location** | Line AND column numbers pinpoint the exact character causing issues |
| **Code Snippets** | See your actual code with context lines before and after |
| **Visual Markers** | Carets (`^^^^^`) underline exactly what's wrong |
| **Smart Suggestions** | Context-aware help tells you *how* to fix it, not just what's broken |
| **Error Codes** | Reference codes (E0001, E0002...) for documentation lookup |
| **Explanatory Notes** | Additional context explaining *why* the error occurred |
| **Typo Detection** | Wrote `println`? Victoria suggests `print`. Wrote `nil`? Try `null`. |

### Real Error Examples

**Type Mismatch with Conversion Hints:**
```
error[E0001]: type mismatch: STRING + INTEGER
 --> main.vc:2:18
  |
1 | let name = "hello"
2 | let result = name + 42
  |                  ^ type mismatch: STRING + INTEGER
3 |
  |
  = note: Victoria is a dynamically typed language, but operators require compatible types
  = help: use string() to convert integers to strings, or int() to convert strings to integers
```

**Smart Typo Suggestions:**
```
error[E0002]: identifier not found: println
 --> main.vc:1:1
  |
1 | println("Hello World")
  | ^^^^^^^ identifier not found: println
  |
  = help: did you mean 'print'? Victoria uses 'print' for output
```

**Array vs Function Call Confusion:**
```
error[E0005]: calling non-function with type ARRAY
 --> main.vc:2:1
  |
1 | let arr = [1, 2, 3]
2 | arr(0)
  | ^^^ calling non-function with type ARRAY
  |
  = help: to access array elements, use bracket notation: arr[0]
```

### Complete Error Code Reference

| Code | Error Type | Description |
|------|------------|-------------|
| `E0001` | Type Mismatch | Incompatible types in operation |
| `E0002` | Undefined Identifier | Variable or function not found |
| `E0003` | Unknown Operator | Unsupported operator for types |
| `E0004` | Unexpected Token | Parser encountered unexpected syntax |
| `E0005` | Not Callable | Trying to call a non-function value |
| `E0006` | Not Indexable | Using `[]` on non-array/hash/string |
| `E0007` | Division by Zero | Attempting to divide by zero |
| `E0008` | Property Not Found | Accessing undefined property/method |
| `E0009` | Struct Not Found | Using undefined struct type |
| `E0010` | Argument Count | Wrong number of function arguments |
| `E0011` | Slice Error | Invalid slice operation |
| `E0012` | Spread Error | Invalid spread operator usage |
| `E0013` | Hash Key Error | Invalid type used as hash key |
| `E0014` | Argument Type | Wrong argument type for function |
| `E0015` | Not Iterable | Iterating over non-iterable value |
| `E0016` | Range Error | Invalid range() arguments |
| `E0017` | Conversion Error | Type conversion failed |
| `E0018` | Assignment Error | Invalid assignment target |
| `E0019` | Operator Error | Invalid operator usage (++/--) |
| `E0020` | Member Access | Invalid dot notation access |
| `E0021` | Reduce Error | Reducing empty array without initial |
| `E0022` | Constant Reassign | Attempting to modify a constant |
| `E0100` | Parse Error | General syntax parsing failure |
| `E0101` | Illegal Character | Invalid character in source code |
| `E0102` | Unterminated String | String literal missing closing quote |


---

## Syntax Comparison

See how Victoria compares to languages you already know:

### Quick Reference Table

| Feature | Victoria | JavaScript | Python | Go |
|---------|----------|------------|--------|-----|
| **Variable** | `let x = 5` | `let x = 5` | `x = 5` | `x := 5` |
| **Constant** | `const PI = 3.14` | `const PI = 3.14` | `PI = 3.14` | `const PI = 3.14` |
| **Function** | `define add(a, b) { return a + b }` | `function add(a, b) { return a + b }` | `def add(a, b): return a + b` | `func add(a, b int) int { return a + b }` |
| **Lambda** | `x => x * 2` | `x => x * 2` | `lambda x: x * 2` | N/A |
| **Range Loop** | `for i in 0..5 { }` | `for (let i = 0; i < 5; i++)` | `for i in range(5):` | `for i := 0; i < 5; i++` |
| **For-Each** | `for item in arr { }` | `for (item of arr)` | `for item in arr:` | `for _, item := range arr` |
| **Print** | `print("Hello")` | `console.log("Hello")` | `print("Hello")` | `fmt.Println("Hello")` |
| **Interpolation** | `"Hi ${name}"` | `` `Hi ${name}` `` | `f"Hi {name}"` | `fmt.Sprintf("Hi %s", name)` |
| **Ternary** | `x ? a : b` | `x ? a : b` | `a if x else b` | N/A |
| **JSON Parse** | `json.parse(s)` | `JSON.parse(s)` | `json.loads(s)` | `json.Unmarshal(...)` |
| **Current Time** | `time.now()` | `new Date()` | `datetime.now()` | `time.Now()` |
| **Array Slice** | `arr[1:3]` | `arr.slice(1,3)` | `arr[1:3]` | `arr[1:3]` |
| **Spread** | `[...a, ...b]` | `[...a, ...b]` | `[*a, *b]` | N/A |

### Side-by-Side Examples

**Fibonacci Function:**

```victoria
// Victoria - Clean and readable
define fib(n) {
    if (n <= 1) { return n }
    return fib(n - 1) + fib(n - 2)
}
```

```javascript
// JavaScript - Similar structure
function fib(n) {
    if (n <= 1) { return n; }
    return fib(n - 1) + fib(n - 2);
}
```

```python
# Python - Different style
def fib(n):
    if n <= 1:
        return n
    return fib(n - 1) + fib(n - 2)
```

**Map/Filter with Lambdas:**

```victoria
// Victoria - Elegant functional style
let nums = [1, 2, 3, 4, 5]
let doubled = map(nums, x => x * 2)        // [2, 4, 6, 8, 10]
let evens = filter(nums, x => x % 2 == 0)  // [2, 4]
let sum = reduce(nums, (a, b) => a + b, 0) // 15
```

```javascript
// JavaScript - Very similar
let nums = [1, 2, 3, 4, 5];
let doubled = nums.map(x => x * 2);        // [2, 4, 6, 8, 10]
let evens = nums.filter(x => x % 2 === 0); // [2, 4]
let sum = nums.reduce((a, b) => a + b, 0); // 15
```

```python
# Python - Functional but verbose
nums = [1, 2, 3, 4, 5]
doubled = list(map(lambda x: x * 2, nums))        # [2, 4, 6, 8, 10]
evens = list(filter(lambda x: x % 2 == 0, nums))  # [2, 4]
from functools import reduce
sum_val = reduce(lambda a, b: a + b, nums, 0)     # 15
```

---

## Examples

We have a collection of examples in the `examples/` directory to help you get started:

- **[hello.vc](examples/hello.vc)**: Basic Hello World.
- **[fib.vc](examples/fib.vc)**: Fibonacci sequence using recursion.
- **[new_features_showcase.vc](examples/new_features_showcase.vc)**: Lambda expressions, JSON, Time, and Math modules.
- **[server.vc](examples/server.vc)**: A simple HTTP server.
- **[structs.vc](examples/structs.vc)**: Working with structs and methods.
- **[error_demo.vc](examples/error_demo.vc)**: Demonstrates the Rust-style error messages.
- **[new_features_test.vc](examples/new_features_test.vc)**: Demonstrates all modern features.

See the [Examples README](examples/README.md) for a full list.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests to improve the language, documentation, or tooling.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  <strong>Victoria</strong> — A language fit for a queen.
</p>

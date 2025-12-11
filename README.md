# Victoria Programming Language

![Queen Victoria](media/QueenVictoria.jpg)

![Victoria Language](https://img.shields.io/badge/Victoria-Language-blueviolet?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen?style=for-the-badge)

**Victoria** is a modern, artistic programming language designed to be expressive and readable. It combines the structural clarity of C-style languages with the ease of use found in Python. Whether you're building simple scripts, data processing tools, or network servers, Victoria provides a clean and enjoyable syntax.

## ✨ Features

- **C-like Structure**: Familiar curly brace `{}` syntax for blocks.
- **Python-like Readability**: Clean, minimal syntax with optional semicolons.
- **First-Class Functions**: Functions are values that can be passed around and returned.
- **Object-Oriented**: Support for `structs` and methods to organize your code.
- **Dynamic Typing**: No need to declare types explicitly.
- **Rust-Style Error Messages**: Beautiful, helpful error messages with code snippets and suggestions.
- **Rich Standard Library**: Built-in support for strings, arrays, hashes, and networking.
- **Versatile Loops**: Supports `for`, `for-in`, `while` loops with `break`/`continue`.
- **Modern Operators**: Ternary operator, compound assignments (`+=`, `-=`, `++`, `--`).
- **String Interpolation**: Embed expressions in strings with `${expr}`.
- **Functional Programming**: `map`, `filter`, `reduce` for array processing.
- **Switch Statements**: Clean pattern matching with `switch`/`case`/`default`.

## 🚀 Installation

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

## 📖 Usage

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

## 🎯 Quick Examples

```victoria
// String interpolation
let name = "Victoria"
print("Hello, ${name}!")

// Range-based loops
for i in 0..5 {
    print(i)  // 0, 1, 2, 3, 4
}

// Functional programming
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

## 📚 Documentation

For a deep dive into the language syntax and features, check out the [Language Documentation](docs/LANGUAGE.md).

Key topics include:
- [Variables & Data Types](docs/LANGUAGE.md#variables)
- [Control Flow](docs/LANGUAGE.md#control-flow)
- [Functions](docs/LANGUAGE.md#functions)
- [Structs](docs/LANGUAGE.md#structs)
- [Built-in Functions](docs/LANGUAGE.md#built-in-functions)

## 🎯 Rust-Inspired Error Messages

Victoria features beautiful, developer-friendly error messages **inspired by the Rust programming language**. We believe that clear, actionable error messages are essential for developer productivity and learning. Our error system draws from Rust's legendary compiler diagnostics to provide:

### Key Features

- **🎨 Color-coded output** - Errors in red, context in cyan, help in green
- **📍 Precise location tracking** - Line and column numbers pinpoint exactly where issues occur
- **📝 Inline code snippets** - See the problematic code directly in the terminal
- **👆 Visual markers** - Carets (`^`) point to the exact error location
- **💡 Actionable suggestions** - Context-aware help tells you how to fix the problem
- **🔢 Error codes** - Reference codes (like `E0002`) for easy documentation lookup
- **📝 Notes** - Additional context explaining why the error occurred

### Example Error Output

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

### Type Mismatch Example

```
error: type mismatch: STRING + INTEGER
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

### Supported Error Types

| Error Code | Description |
|------------|-------------|
| `E0001` | Type mismatch between operands |
| `E0002` | Undefined variable or identifier |
| `E0003` | Unknown or unsupported operator |
| `E0004` | Unexpected token during parsing |
| `E0005` | Attempting to call a non-function |
| `E0006` | Index operator on non-indexable type |
| `E0100` | General parse error |
| `E0101` | Illegal character in source |
| `E0102` | Unterminated string literal |

## 💡 Examples

We have a collection of examples in the `examples/` directory to help you get started:

- **[hello.vc](examples/hello.vc)**: Basic Hello World.
- **[fib.vc](examples/fib.vc)**: Fibonacci sequence using recursion.
- **[server.vc](examples/server.vc)**: A simple HTTP server.
- **[structs.vc](examples/structs.vc)**: Working with structs and methods.
- **[error_demo.vc](examples/error_demo.vc)**: Demonstrates the Rust-style error messages.
- **[new_features_test.vc](examples/new_features_test.vc)**: Demonstrates all modern features.

See the [Examples README](examples/README.md) for a full list.

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests to improve the language, documentation, or tooling.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

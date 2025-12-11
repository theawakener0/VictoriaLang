# Victoria Programming Language

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
- **Rich Standard Library**: Built-in support for strings, arrays, hashes, and networking.
- **Versatile Loops**: Supports `for`, `for-in`, and `while` loops.

## 🚀 Installation

To build Victoria from source, you need to have [Go](https://golang.org/) installed on your machine.

1.  Clone the repository:
    ```bash
    git clone https://github.com/yourusername/victoria.git
    cd victoria
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
>> print("Hello, " + name + "!")
Hello, World!
```

### Running Scripts

Execute a Victoria script (`.vc` file) by passing the filename:

```bash
./victoria examples/hello.vc
```

## 📚 Documentation

For a deep dive into the language syntax and features, check out the [Language Documentation](docs/LANGUAGE.md).

Key topics include:
- [Variables & Data Types](docs/LANGUAGE.md#variables)
- [Control Flow](docs/LANGUAGE.md#control-flow)
- [Functions](docs/LANGUAGE.md#functions)
- [Structs](docs/LANGUAGE.md#structs)
- [Built-in Functions](docs/LANGUAGE.md#built-in-functions)

## 💡 Examples

We have a collection of examples in the `examples/` directory to help you get started:

- **[hello.vc](examples/hello.vc)**: Basic Hello World.
- **[fib.vc](examples/fib.vc)**: Fibonacci sequence using recursion.
- **[server.vc](examples/server.vc)**: A simple HTTP server.
- **[structs.vc](examples/structs.vc)**: Working with structs and methods.

See the [Examples README](examples/README.md) for a full list.

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests to improve the language, documentation, or tooling.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

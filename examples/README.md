# Victoria Language Examples

This directory contains example scripts demonstrating various features of the Victoria programming language.

## Quick Start

```victoria
// Hello World with string interpolation
let name = "World"
print("Hello, ${name}!")

// Range-based loop
for i in range(5) {
    print(i)
}

// Functional programming
let nums = [1, 2, 3, 4, 5]
let doubled = map(nums, define(x) { x * 2 })
print(doubled)  // [2, 4, 6, 8, 10]
```

## Basic Examples

| File | Description |
|------|-------------|
| `hello.vc` | Simple "Hello, World!" program |
| `fib.vc` | Fibonacci sequence using recursion |
| `math_test.vc` | Demonstrates mathematical operations |
| `try_catch.vc` | Shows error handling capabilities |
| `new_features_test.vc` | Comprehensive test of all modern features |

## Data Structures

| File | Description |
|------|-------------|
| `structs.vc` | Demonstrates structs and methods |
| `include_test.vc` | Shows how to include other files |
| `mylib.vc` | Example library file for imports |

## Networking

| File | Description |
|------|-------------|
| `server.vc` | Simple HTTP server implementation |
| `html_server.vc` | HTTP server serving HTML content |
| `tcp_server.vc` | TCP server example |
| `tcp_client.vc` | TCP client example |
| `http_client_test.vc` | Making HTTP requests |

## Testing

| File | Description |
|------|-------------|
| `full_test.vc` | Comprehensive test suite for language features |
| `std_test.vc` | Tests for the standard library |
| `modules_test.vc` | Tests for module loading |

## How to Run

To run any of these examples, use the `victoria` interpreter from the root directory:

```bash
# Run a single example
./victoria examples/hello.vc

# Run the comprehensive feature test
./victoria examples/new_features_test.vc
```

## Featured: new_features_test.vc

This example demonstrates all modern language features:

- ✅ Pre/post increment/decrement (`++x`, `x++`, `--x`, `x--`)
- ✅ Compound assignments (`+=`, `-=`, `*=`, `/=`, `%=`)
- ✅ For loop with index (`for i, v in arr`)
- ✅ Range-based loops (`for i in 0..10`, `for i in range(10)`)
- ✅ Break and continue statements
- ✅ Logical operators (`&&`, `||`)
- ✅ Ternary operator (`condition ? a : b`)
- ✅ String interpolation (`"Hello ${name}"`)
- ✅ Multi-line strings (backticks)
- ✅ Functional methods (`map`, `filter`, `reduce`)
- ✅ Switch statements

# Victoria Language Examples

This directory contains example scripts demonstrating various features of the Victoria programming language.

## Quick Start

```victoria
// Hello World with string interpolation
let name = "World"
print("Hello, ${name}!")

// Lambda expressions
let double = x => x * 2
let add = (a, b) => a + b
print(double(5))  // 10

// Functional programming with lambdas
let nums = [1, 2, 3, 4, 5]
let doubled = map(nums, x => x * 2)
print(doubled)  // [2, 4, 6, 8, 10]

// JSON handling
let data = json.parse('{"name": "Victoria"}')
print(data["name"])

// Current time
print(time.now())
```

## Basic Examples

| File | Description |
|------|-------------|
| `hello.vc` | Simple "Hello, World!" program |
| `fib.vc` | Fibonacci sequence using recursion |
| `math_test.vc` | Demonstrates mathematical operations |
| `try_catch.vc` | Shows error handling capabilities |
| `new_features_test.vc` | Comprehensive test of all modern features |
| `new_features_showcase.vc` | **NEW:** Lambdas, JSON, Time, and Math modules |

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

## Featured: new_features_showcase.vc

This example demonstrates the latest language features:

### Lambda Expressions
- ✅ Single parameter: `x => x * 2`
- ✅ Multiple parameters: `(a, b) => a + b`
- ✅ No parameters: `() => 42`
- ✅ Works with `map`, `filter`, `reduce`

### JSON Module
- ✅ `json.parse(str)` - Parse JSON strings
- ✅ `json.stringify(obj)` - Convert objects to JSON
- ✅ `json.valid(str)` - Validate JSON strings

### Time Module
- ✅ `time.now()` - Current timestamp
- ✅ `time.format(t, fmt)` - Custom formatting
- ✅ `time.year/month/day/hour/minute/second(t)` - Extract components
- ✅ `time.weekday(t)` - Day of week
- ✅ `time.sleep(ms)` - Pause execution

### Math Module
- ✅ `math.floor/ceil/round(x)` - Rounding functions
- ✅ `math.min/max(...)` - Minimum and maximum
- ✅ `math.random()` - Random numbers
- ✅ `math.sin/cos/tan/log(x)` - Trigonometry and logarithms

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

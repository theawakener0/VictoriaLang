# Victoria Language Documentation

Victoria is a dynamic, interpreted programming language designed for readability and expressiveness. It combines features from C-style languages (like curly braces) with Python-like simplicity. Victoria supports **optional static typing** inspired by Go, making it great for learning and prototyping while enabling type safety when needed.

## Table of Contents

- [Variables](#variables)
  - [Constant Variables](#constant-variables)
  - [Type Annotations](#type-annotations)
- [Data Types](#data-types)
- [Type System](#type-system)
- [Operators](#operators)
- [Control Flow](#control-flow)
- [Functions](#functions)
  - [Typed Functions](#typed-functions)
  - [Lambda Functions (Arrow Functions)](#lambda-functions-arrow-functions)
- [Data Structures](#data-structures)
  - [Array Slicing](#array-slicing)
  - [Spread Operator](#spread-operator)
- [Structs](#structs)
- [Modules](#modules)
  - [Math Module](#math-module)
  - [JSON Module](#json-module)
  - [Time Module](#time-module)
- [Built-in Functions](#built-in-functions)
- [Error Handling](#error-handling)

## Variables

Variables are declared using the `let` keyword.

```victoria
let x = 5
let y = 10
let name = "Victoria"
let isActive = true
```

Variables can be reassigned:

```victoria
let x = 5
x = 10  // reassignment
```

### Constant Variables

Use `const` to declare immutable variables that cannot be reassigned:

```victoria
const PI = 3.14159
const MAX_SIZE = 100
const APP_NAME = "Victoria Lang"

print(PI)       // 3.14159

// Attempting to reassign a const will cause an error:
PI = 3.0        // ERROR: cannot reassign constant variable: PI
```

Constants are useful for values that should never change, like configuration values, mathematical constants, or API keys.

### Type Annotations

Victoria supports optional type annotations for variables, inspired by Go's type system:

```victoria
let x:int = 42
let name:string = "Victoria"
let pi:float = 3.14159
let isReady:bool = true
```

When a type is specified, Victoria will check that the assigned value matches the declared type:

```victoria
let x:int = "hello"  // ERROR: type mismatch: cannot assign string to variable of type int
```

Available types:
- `int` - Integer numbers
- `float` - Floating-point numbers
- `string` - Text strings
- `bool` - Boolean values (true/false)
- `char` - Single character (single-char string)
- `byte` - Single byte (0-255)
- `rune` - Unicode code point (like Go's rune)
- `array` - Arrays/lists
- `map` - Hash maps/dictionaries
- `any` - Any type (disables type checking)
- `void` - No value (for functions that don't return)

Type annotations are **optional** - you can mix typed and untyped code:

```victoria
let typed:int = 42     // typed variable
let untyped = "hello"  // untyped (dynamic) variable
```

## Preprocessor Directives

### #make (Compile-Time Constants)

The `#make` directive works like C's `#define` to create compile-time constants. These are especially useful for DSA (Data Structures and Algorithms) problems:

```victoria
#make MAX_SIZE 1000
#make MOD 1000000007
#make INF 999999999
#make EPSILON 0.00001

// Use in your code
let arr = []
for i in 0..MAX_SIZE {
    arr = push(arr, 0)
}

// Great for competitive programming
let result = (a * b) % MOD
```

`#make` creates immutable constants that cannot be reassigned, similar to `const` but semantically indicates a compile-time definition.

## Enums

Enums define a set of named integer constants, perfect for representing states, options, or categories:

```victoria
// Basic enum
enum Color {
    RED,        // 0
    GREEN,      // 1
    BLUE        // 2
}

// Enum with explicit values
enum HttpStatus {
    OK = 200,
    NOT_FOUND = 404,
    SERVER_ERROR = 500
}

// Using enums
let myColor = Color.RED
print(myColor)           // Color.RED

// Enums in switch statements
switch (myColor) {
    case Color.RED: {
        print("It's red!")
    }
    case Color.GREEN: {
        print("It's green!")
    }
    case Color.BLUE: {
        print("It's blue!")
    }
}

// Enum for graph algorithms
enum NodeState {
    UNVISITED,
    VISITING,
    VISITED
}

let states = []
for i in 0..n {
    states = push(states, NodeState.UNVISITED)
}
```

## Data Types

Victoria supports the following basic data types:

- **Integer**: `1`, `42`, `-10`
- **Float**: `3.14`, `0.001`
- **String**: `"Hello World"` or `` `multi-line string` ``
- **Char**: `'a'`, `'Z'`, `'\n'` (single character with single quotes)
- **Byte**: Single byte value (0-255)
- **Rune**: Unicode code point (32-bit)
- **Boolean**: `true`, `false`
- **Null**: `null` (implicit in some contexts)
- **Array**: `[1, 2, 3]`
- **Hash**: `{"key": "value"}`
- **Enum**: Named integer constants

### Character Literals

Use single quotes for character literals:

```victoria
let ch = 'A'           // Character literal
let newline = '\n'     // Escape sequences supported
let tab = '\t'

// Get ASCII/Unicode value
let ascii = ord('A')   // 65
print(ascii)

// Convert number to character
let letter = chr(65)   // "A"
print(letter)

// Character functions for DSA
if (isDigit(ch)) {
    print("It's a digit!")
}

if (isLetter(ch)) {
    print("It's a letter!")
}

// Useful for string parsing in algorithms
let str = "abc123"
for i in 0..len(str) {
    let c = str[i]
    if (isDigit(c)) {
        print("Found digit: ${c}")
    }
}
```

### String Interpolation

Embed expressions directly in strings using `${expr}`:

```victoria
let name = "Victoria"
let version = 1
print("Hello from ${name} version ${version}!")
// Output: Hello from Victoria version 1!

let a = 2
let b = 3
print("${a} + ${b} = ${a + b}")
// Output: 2 + 3 = 5
```

### Multi-line Strings

Use backticks for multi-line strings:

```victoria
let poem = `Roses are red,
Violets are blue,
Victoria is awesome,
And so are you!`
```

## Operators

### Arithmetic Operators

| Operator | Description |
|----------|-------------|
| `+` | Addition |
| `-` | Subtraction |
| `*` | Multiplication |
| `/` | Division |
| `%` | Modulo |

### Comparison Operators

| Operator | Description |
|----------|-------------|
| `==` | Equal |
| `!=` | Not equal |
| `<` | Less than |
| `>` | Greater than |
| `<=` | Less than or equal |
| `>=` | Greater than or equal |

### Logical Operators

| Operator | Description |
|----------|-------------|
| `&&` or `and` | Logical AND |
| `\|\|` or `or` | Logical OR |
| `!` | Logical NOT |

```victoria
if (x > 0 && x < 10) {
    print("x is between 0 and 10")
}

if (name == "admin" || isSuper) {
    print("Access granted")
}
```

### Compound Assignment Operators

| Operator | Description | Equivalent |
|----------|-------------|------------|
| `+=` | Add and assign | `x = x + y` |
| `-=` | Subtract and assign | `x = x - y` |
| `*=` | Multiply and assign | `x = x * y` |
| `/=` | Divide and assign | `x = x / y` |
| `%=` | Modulo and assign | `x = x % y` |

```victoria
let x = 10
x += 5   // x is now 15
x -= 3   // x is now 12
x *= 2   // x is now 24
x /= 4   // x is now 6
x %= 4   // x is now 2
```

### Increment/Decrement Operators

| Operator | Description | Returns |
|----------|-------------|---------|
| `++x` | Pre-increment | New value |
| `--x` | Pre-decrement | New value |
| `x++` | Post-increment | Old value |
| `x--` | Post-decrement | Old value |

```victoria
let x = 5
print(++x)  // 6 (increments first, then returns)
print(x++)  // 6 (returns first, then increments)
print(x)    // 7
```

### Ternary Operator

```victoria
let result = condition ? valueIfTrue : valueIfFalse

let age = 20
let status = age >= 18 ? "adult" : "minor"
print(status)  // "adult"
```

### Range Operator

Create ranges using `..`:

```victoria
for i in 0..5 {
    print(i)  // 0, 1, 2, 3, 4
}
```

## Control Flow

### If-Else

```victoria
let x = 10
if (x > 5) {
    print("x is greater than 5")
} else if (x == 5) {
    print("x equals 5")
} else {
    print("x is less than 5")
}
```

### Switch Statement

```victoria
let day = 3
switch (day) {
    case 1: {
        print("Monday")
    }
    case 2: {
        print("Tuesday")
    }
    case 3: {
        print("Wednesday")
    }
    default: {
        print("Other day")
    }
}
```

### Loops

**While Loop**

```victoria
let i = 0
while (i < 5) {
    print(i)
    i++
}
```

**For Loop (C-style)**

```victoria
for (let i = 0; i < 5; i++) {
    print(i)
}
```

**For-In Loop**

```victoria
let arr = [1, 2, 3]
for x in arr {
    print(x)
}
```

**For-In with Index**

```victoria
let fruits = ["apple", "banana", "cherry"]
for i, fruit in fruits {
    print("${i}: ${fruit}")
}
// Output:
// 0: apple
// 1: banana
// 2: cherry
```

**Range-based For Loop**

```victoria
// Using range operator
for i in 0..5 {
    print(i)  // 0, 1, 2, 3, 4
}

// Using range function
for i in range(5) {
    print(i)  // 0, 1, 2, 3, 4
}

for i in range(2, 8) {
    print(i)  // 2, 3, 4, 5, 6, 7
}

for i in range(0, 10, 2) {
    print(i)  // 0, 2, 4, 6, 8
}
```

### Break and Continue

```victoria
// Break exits the loop
for i in 0..10 {
    if (i == 5) {
        break
    }
    print(i)  // 0, 1, 2, 3, 4
}

// Continue skips to next iteration
for i in 0..5 {
    if (i == 2) {
        continue
    }
    print(i)  // 0, 1, 3, 4
}
```

## Functions

Functions are first-class citizens and are defined using the `define` keyword.

```victoria
define add(x, y) {
    return x + y
}

let result = add(5, 10)
print(result)  // 15
```

Functions can also be assigned to variables (anonymous functions):

```victoria
let multiply = define(x, y) {
    return x * y
}

print(multiply(3, 4))  // 12
```

### Typed Functions

Victoria supports Go-style typed function parameters and return types:

```victoria
// Function with typed parameters
define add(a:int, b:int) -> int {
    return a + b
}

// Multiple return values with different types
define greet(name:string, age:int) -> string {
    return "Hello, " + name + "! You are " + string(age) + " years old."
}

// Function with any type (accepts anything)
define printValue(val:any) -> void {
    print("Value:", val)
}
```

Type checking happens at runtime when functions are called:

```victoria
define add(a:int, b:int) -> int {
    return a + b
}

add(5, 3)        // OK - both arguments are int
add("hello", 5)  // ERROR: type mismatch for parameter 'a': expected int, got string
```

Return types are also checked:

```victoria
define getNumber() -> int {
    return "not a number"  // ERROR: return type mismatch: expected int, got string
}
```

You can mix typed and untyped parameters for flexibility:

```victoria
define log(message:string, data) {
    print(message, data)
}
```

### Higher-Order Functions

Functions can take other functions as arguments:

```victoria
define applyTwice(f, x) {
    return f(f(x))
}

let double = define(n) { n * 2 }
print(applyTwice(double, 5))  // 20
```

### Lambda Functions (Arrow Functions)

Victoria supports a concise lambda syntax using the `=>` operator, perfect for short functions:

```victoria
// Single parameter lambda
let double = x => x * 2
print(double(5))  // 10

// Multiple parameters (use parentheses)
let add = (x, y) => x + y
print(add(3, 4))  // 7

// No parameters
let getAnswer = () => 42
print(getAnswer())  // 42
```

Lambda functions are ideal for use with `map`, `filter`, and `reduce`:

```victoria
let numbers = [1, 2, 3, 4, 5]

// Map: transform each element
let doubled = map(numbers, x => x * 2)
print(doubled)  // [2, 4, 6, 8, 10]

// Filter: keep elements matching condition
let evens = filter(numbers, x => x % 2 == 0)
print(evens)  // [2, 4]

// Reduce: accumulate values
let sum = reduce(numbers, (acc, x) => acc + x, 0)
print(sum)  // 15

// Chain operations
let result = map(filter(numbers, x => x > 2), x => x * 10)
print(result)  // [30, 40, 50]
```

## Data Structures

### Arrays

Arrays are ordered lists of values.

```victoria
let arr = [1, 2, 3, "four", true]
print(arr[0])    // 1
print(len(arr))  // 5

// Array methods
let nums = [1, 2, 3, 4, 5]

// map - transform each element
let doubled = map(nums, define(x) { x * 2 })
print(doubled)  // [2, 4, 6, 8, 10]

// filter - keep elements matching condition
let evens = filter(nums, define(x) { x % 2 == 0 })
print(evens)  // [2, 4]

// reduce - accumulate values
let sum = reduce(nums, define(acc, x) { acc + x }, 0)
print(sum)  // 15
```

#### Array Slicing

Extract portions of arrays using slice syntax `[start:end]`:

```victoria
let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

// Get elements from index 2 to 5 (exclusive)
print(arr[2:5])   // [3, 4, 5]

// Get elements from start to index 3
print(arr[:3])    // [1, 2, 3]

// Get elements from index 5 to end
print(arr[5:])    // [6, 7, 8, 9, 10]

// Negative indexing (from end)
print(arr[-3:])   // [8, 9, 10]
print(arr[:-2])   // [1, 2, 3, 4, 5, 6, 7, 8]

// String slicing also works
let str = "Hello, World!"
print(str[0:5])   // "Hello"
print(str[7:])    // "World!"
print(str[-6:])   // "World!"
```

#### Spread Operator

Use the spread operator `...` to expand arrays:

```victoria
let arr1 = [1, 2, 3]
let arr2 = [4, 5, 6]

// Merge arrays
let merged = [...arr1, ...arr2]
print(merged)     // [1, 2, 3, 4, 5, 6]

// Add elements around spread
let extended = [0, ...arr1, 99, ...arr2, 100]
print(extended)   // [0, 1, 2, 3, 99, 4, 5, 6, 100]

// Copy an array
let original = [10, 20, 30]
let copy = [...original]
print(copy)       // [10, 20, 30]

// Combine with slicing
let arr = [1, 2, 3, 4, 5]
let combined = [...arr[:2], ...arr[3:]]
print(combined)   // [1, 2, 4, 5]
```

### Hashes

Hashes are key-value pairs (dictionaries).

```victoria
let person = {
    "name": "Alice",
    "age": 30,
    "city": "Wonderland"
}

print(person["name"])  // Alice

// Get all keys
print(keys(person))    // ["name", "age", "city"]

// Get all values
print(values(person))  // ["Alice", 30, "Wonderland"]
```

## Structs

Structs allow you to define custom data types with fields and methods.

```victoria
struct Person {
    name
    age
}

// Define a method for the Person struct
define Person.greet() {
    print("Hello, my name is ${self.name}")
}

define Person.birthday() {
    self.age = self.age + 1
    print("Happy birthday! Now ${self.age} years old.")
}

// Instantiate a struct
let p = Person {
    name: "Bob",
    age: 25
}

p.greet()     // Hello, my name is Bob
p.birthday()  // Happy birthday! Now 26 years old.
```

## Modules

You can include other Victoria files or built-in modules using `include`.

```victoria
include "net"        // Built-in network module
include "os"         // Built-in OS module  
include "math"       // Built-in math module
include "json"       // Built-in JSON module
include "time"       // Built-in time/date module
include "path"       // Built-in path module
include "mylib"      // Local file (mylib.vc)
```

Multiple modules can be included at once:

```victoria
include ("os", "math", "json", "time", "path")
```

### Available Modules

- `os`: File system operations (`readFile`, `writeFile`, `exists`, `exec`, `platform`, etc.)
- `net`: Networking (`http.get`, `http.post`, `tcp.connect`, etc.)
- `math`: Math functions (`abs`, `sqrt`, `pow`, `sin`, `cos`, `floor`, `ceil`, `round`, `min`, `max`, `random`, etc.)
- `json`: JSON parsing and serialization (`parse`, `stringify`, `valid`)
- `time`: Date/time operations (`now`, `format`, `parse`, `date`, `time`, etc.)
- `path`: Path manipulation (`join`, `base`, `dir`, `ext`, `abs`, `isAbs`)
- `std`: Standard utilities (common functions)

### Module Isolation

When you include a user-defined module (e.g., `include "mylib"`), Victoria evaluates it in an **isolated environment**. All top-level variables and functions defined in that module are exported as a Hash object named after the module.

Example `mylib.vc`:
```victoria
let version = "1.0.0"
let add = fn(a, b) { return a + b }
```

Example usage:
```victoria
include "mylib"
print(mylib.version)
print(mylib.add(1, 2))
```

### OS Module

The OS module provides functions for interacting with the operating system and file system.

```victoria
include "os"

// File operations
os.writeFile("test.txt", "Hello Victoria!")
let content = os.readFile("test.txt")
print(os.exists("test.txt")) // true
os.remove("test.txt")

// Directory operations
os.mkdir("data/logs")
let files = os.readDir(".")
os.chdir("data")
print(os.getwd())

// Environment and System
print(os.env("PATH"))
os.env("MY_VAR", "123")
print(os.platform) // "windows", "linux", etc.
print(os.arch)     // "amd64", "arm64", etc.
print(os.pid)      // Current process ID
print(os.hostname())
print(os.user().username)

// Command Execution
let output = os.exec("echo", "Hello from Shell")
print(output)

// Process
os.kill(1234) // Kill process by PID
os.exit(0)
```

| Function/Property | Description |
|-------------------|-------------|
| `os.readFile(path)` | Read file content as string |
| `os.writeFile(path, content)` | Write string to file |
| `os.remove(path)` | Delete a file |
| `os.exists(path)` | Check if file/dir exists |
| `os.mkdir(path)` | Create directory (recursive) |
| `os.readDir(path)` | List files in directory |
| `os.stat(path)` | Get file info (size, modTime, isDir) |
| `os.rename(old, new)` | Rename/move file or dir |
| `os.getwd()` | Get current working directory |
| `os.chdir(path)` | Change current working directory |
| `os.env()` | Get all environment variables as Hash |
| `os.env(key)` | Get environment variable |
| `os.env(key, val)` | Set environment variable |
| `os.args` | Array of command line arguments |
| `os.platform` | Operating system name |
| `os.arch` | CPU architecture |
| `os.pid` | Current process ID |
| `os.hostname()` | Get system hostname |
| `os.user()` | Get current user info (Hash) |
| `os.tempDir()` | Get system temp directory |
| `os.chmod(path, mode)` | Change file permissions |
| `os.exec(cmd, ...args)` | Execute external command and return output |
| `os.kill(pid)` | Kill a process |
| `os.exit(code)` | Exit the process |

### Path Module

The path module provides utilities for manipulating file and directory paths.

```victoria
include "path"

let fullPath = path.join("data", "users", "config.json")
print(path.base(fullPath)) // "config.json"
print(path.dir(fullPath))  // "data/users"
print(path.ext(fullPath))  // ".json"
print(path.abs("."))       // Absolute path to current dir
print(path.isAbs("/etc"))  // true (on Unix)
```

| Function | Description |
|----------|-------------|
| `path.join(...parts)` | Join path segments |
| `path.base(path)` | Get last element of path |
| `path.dir(path)` | Get directory part of path |
| `path.ext(path)` | Get file extension |
| `path.abs(path)` | Get absolute path |
| `path.isAbs(path)` | Check if path is absolute |

### Math Module

The math module provides mathematical constants and functions.

```victoria
include "math"

// Constants
print(math.pi)   // 3.141592653589793
print(math.e)    // 2.718281828459045

// Basic functions
print(math.abs(-5))      // 5
print(math.sqrt(16))     // 4
print(math.pow(2, 10))   // 1024

// Trigonometric functions
print(math.sin(0))       // 0
print(math.cos(0))       // 1
print(math.tan(0))       // 0

// Rounding functions
print(math.floor(3.7))   // 3
print(math.ceil(3.2))    // 4
print(math.round(3.5))   // 4

// Min/Max (accepts multiple arguments)
print(math.min(5, 3, 8, 1))   // 1
print(math.max(5, 3, 8, 1))   // 8

// Random numbers
print(math.random())          // Random float 0-1
print(math.random(10))        // Random int 0-9
print(math.random(1, 6))      // Random int 1-6 (dice roll)

// Logarithmic functions
print(math.log(2.718281828))  // ~1 (natural log)
print(math.log10(100))        // 2
```

| Function | Description |
|----------|-------------|
| `math.pi` | Pi constant (3.14159...) |
| `math.e` | Euler's number (2.71828...) |
| `math.abs(x)` | Absolute value |
| `math.sqrt(x)` | Square root |
| `math.pow(x, y)` | x raised to power y |
| `math.sin(x)` | Sine (radians) |
| `math.cos(x)` | Cosine (radians) |
| `math.tan(x)` | Tangent (radians) |
| `math.floor(x)` | Round down to integer |
| `math.ceil(x)` | Round up to integer |
| `math.round(x)` | Round to nearest integer |
| `math.min(a, b, ...)` | Minimum of values |
| `math.max(a, b, ...)` | Maximum of values |
| `math.random()` | Random float 0-1 |
| `math.random(n)` | Random integer 0 to n-1 |
| `math.random(min, max)` | Random integer min to max (inclusive) |
| `math.log(x)` | Natural logarithm |
| `math.log10(x)` | Base-10 logarithm |

### JSON Module

The JSON module provides functions for parsing and serializing JSON data.

```victoria
include "json"

// Parse JSON string to Victoria object
let jsonStr = `{"name": "Victoria", "version": 1.0}`
let obj = json.parse(jsonStr)
print(obj["name"])     // Victoria
print(obj["version"])  // 1

// Stringify Victoria object to JSON
let data = {
    "language": "Victoria",
    "year": 2024,
    "active": true,
    "tags": ["fast", "simple"]
}
print(json.stringify(data))
// {"active":true,"language":"Victoria","tags":["fast","simple"],"year":2024}

// Pretty print with indentation
print(json.stringify(data, 2))
// {
//   "active": true,
//   "language": "Victoria",
//   ...
// }

// Validate JSON string
print(json.valid(`{"valid": true}`))  // true
print(json.valid("not json"))         // false
```

| Function | Description |
|----------|-------------|
| `json.parse(str)` | Parse JSON string to Victoria object |
| `json.stringify(obj)` | Convert object to JSON string |
| `json.stringify(obj, indent)` | Convert with pretty printing (indent = spaces or string) |
| `json.valid(str)` | Check if string is valid JSON |

### Time Module

The time module provides functions for working with dates and times.

```victoria
include "time"

// Get current timestamp (Unix seconds)
let now = time.now()
print(now)  // e.g., 1702400000

// Get current timestamp in milliseconds
let nowMs = time.nowMs()
print(nowMs)  // e.g., 1702400000123

// Format timestamp to string
print(time.format(now))                    // "2024-12-12 10:30:00"
print(time.format(now, "YYYY-MM-DD"))      // "2024-12-12"
print(time.format(now, "HH:mm:ss"))        // "10:30:00"

// Parse date string to timestamp
let ts = time.parse("2024-12-25 10:30:00")
print(ts)  // Unix timestamp

// Get date/time components
print(time.year(now))     // 2024
print(time.month(now))    // 12
print(time.day(now))      // 12
print(time.hour(now))     // 10
print(time.minute(now))   // 30
print(time.second(now))   // 0
print(time.weekday(now))  // 0-6 (0=Sunday)

// Quick date/time strings
print(time.date())        // "2024-12-12"
print(time.time())        // "10:30:00"
print(time.date(now))     // Date of timestamp
print(time.time(now))     // Time of timestamp

// Sleep (pause execution)
time.sleep(1000)  // Sleep for 1000 milliseconds (1 second)
```

| Function | Description |
|----------|-------------|
| `time.now()` | Current Unix timestamp (seconds) |
| `time.nowMs()` | Current timestamp (milliseconds) |
| `time.format(ts)` | Format timestamp as "YYYY-MM-DD HH:mm:ss" |
| `time.format(ts, fmt)` | Format with custom format string |
| `time.parse(str)` | Parse date string to timestamp |
| `time.parse(str, fmt)` | Parse with custom format string |
| `time.year(ts)` | Get year from timestamp |
| `time.month(ts)` | Get month (1-12) |
| `time.day(ts)` | Get day of month (1-31) |
| `time.hour(ts)` | Get hour (0-23) |
| `time.minute(ts)` | Get minute (0-59) |
| `time.second(ts)` | Get second (0-59) |
| `time.weekday(ts)` | Get weekday (0=Sunday, 6=Saturday) |
| `time.date()` | Current date as "YYYY-MM-DD" |
| `time.time()` | Current time as "HH:mm:ss" |
| `time.sleep(ms)` | Pause execution for milliseconds |

#### Format Tokens

| Token | Description | Example |
|-------|-------------|---------|
| `YYYY` | 4-digit year | 2024 |
| `YY` | 2-digit year | 24 |
| `MM` | 2-digit month | 01-12 |
| `DD` | 2-digit day | 01-31 |
| `HH` | 24-hour hour | 00-23 |
| `hh` | 12-hour hour | 01-12 |
| `mm` | Minute | 00-59 |
| `ss` | Second | 00-59 |
| `A` | AM/PM | AM, PM |
| `a` | am/pm | am, pm |

### Net Module

The net module provides networking capabilities for both client and server applications.

```victoria
include "net"

// HTTP Client
let resp = net.get("https://api.github.com")
print(resp.status)
print(resp.body)

// Generic Request
let postResp = net.request("POST", "https://example.com/api", `{"id": 1}`, {"Content-Type": "application/json"})

// DNS and Interfaces
print(net.lookupHost("google.com"))
print(net.interfaces()[0].mac)

// HTTP Server (Simple)
net.listen(8080, fn(req) {
    return "Hello from Victoria!"
})

// HTTP Server (with Routing)
net.serve(8080, {
    "/": fn(req) { return "Home" },
    "/api": fn(req) { 
        return {
            "status": 200,
            "body": `{"status": "ok"}`,
            "headers": {"Content-Type": "application/json"}
        }
    }
})

// TCP Client
let conn = net.dial("localhost", 9000)
conn.write("Hello TCP\n")
print(conn.read())
conn.close()

// TCP Server
net.listenTcp(9000, fn(conn) {
    let msg = conn.read()
    print("Received: " + msg)
    conn.write("Echo: " + msg + "\n")
    conn.close()
})

// UDP Server
net.listenUdp(9001, fn(packet) {
    print("From " + packet.remote + ": " + packet.data)
})
```

| Function | Description |
|----------|-------------|
| `net.get(url)` | HTTP GET request |
| `net.post(url, body, [type])` | HTTP POST request |
| `net.request(method, url, [body], [headers])` | Generic HTTP request |
| `net.lookupHost(host)` | DNS lookup (returns array of IPs) |
| `net.interfaces()` | Get network interfaces info |
| `net.parseQuery(str)` | Parse URL query string to Hash |
| `net.listen(port, handler)` | Start simple HTTP server |
| `net.serve(port, routes)` | Start HTTP server with routing (Hash of path -> handler) |
| `net.dial(host, port)` | Connect to TCP server |
| `net.dialUdp(host, port)` | Connect to UDP server |
| `net.listenTcp(port, handler)` | Start TCP server |
| `net.listenUdp(port, handler)` | Start UDP server |

## Built-in Functions

Victoria comes with a standard library of built-in functions:

### General

| Function | Description |
|----------|-------------|
| `print(args...)` | Prints arguments to the console |
| `len(arg)` | Returns the length of a string or array |
| `type(arg)` | Returns the type of the argument as a string |
| `input(prompt)` | Reads input from the user |

### Type Conversion

| Function | Description |
|----------|-------------|
| `int(arg)` | Converts a value to an integer |
| `string(arg)` | Converts a value to a string |
| `float(arg)` | Converts a value to a float |
| `char(arg)` | Converts a value to a character |
| `byte(arg)` | Converts a value to a byte (0-255) |
| `rune(arg)` | Converts a value to a rune (Unicode code point) |

### Character Functions

These functions are essential for DSA problems involving string parsing, character classification, and manipulation:

| Function | Description |
|----------|-------------|
| `ord(char)` | Returns the ASCII/Unicode value of a character |
| `chr(int)` | Returns the character for an ASCII/Unicode value |
| `isDigit(char)` | Returns `true` if character is '0'-'9' |
| `isLetter(char)` | Returns `true` if character is 'a'-'z' or 'A'-'Z' |
| `isAlpha(char)` | Returns `true` if character is alphanumeric |
| `isSpace(char)` | Returns `true` if character is whitespace |
| `toUpper(char)` | Converts character to uppercase |
| `toLower(char)` | Converts character to lowercase |

```victoria
// Example: Count digits and letters in a string
let str = "Hello123World456"
let digits = 0
let letters = 0

for i in 0..len(str) {
    let c = str[i]
    if (isDigit(c)) {
        digits++
    }
    if (isLetter(c)) {
        letters++
    }
}
print("Digits: ${digits}, Letters: ${letters}")
// Output: Digits: 6, Letters: 10

// Example: Convert string to uppercase manually
define toUpperStr(s) {
    let result = ""
    for i in 0..len(s) {
        let c = s[i]
        if (ord(c) >= ord('a') && ord(c) <= ord('z')) {
            result = result + chr(ord(c) - 32)
        } else {
            result = result + c
        }
    }
    return result
}
print(toUpperStr("hello"))  // "HELLO"
```

### Range

| Function | Description |
|----------|-------------|
| `range(end)` | Returns `[0, 1, ..., end-1]` |
| `range(start, end)` | Returns `[start, start+1, ..., end-1]` |
| `range(start, end, step)` | Returns `[start, start+step, ..., <end]` |

```victoria
range(5)        // [0, 1, 2, 3, 4]
range(2, 7)     // [2, 3, 4, 5, 6]
range(0, 10, 2) // [0, 2, 4, 6, 8]
range(10, 0, -1) // [10, 9, 8, 7, 6, 5, 4, 3, 2, 1]
```

### String Functions

| Function | Description |
|----------|-------------|
| `format(str, args...)` | Formats a string (like printf) |
| `split(str, sep)` | Splits a string into an array |
| `join(arr, sep)` | Joins an array into a string |
| `upper(str)` | Converts to uppercase |
| `lower(str)` | Converts to lowercase |
| `contains(str, substr)` | Checks if string contains substring |
| `index(str, substr)` | Returns index of substring (-1 if not found) |

### Array Functions

| Function | Description |
|----------|-------------|
| `first(arr)` | Returns the first element |
| `last(arr)` | Returns the last element |
| `rest(arr)` | Returns all elements except first |
| `push(arr, elem)` | Returns new array with element added |
| `pop(arr)` | Returns new array with last element removed |
| `map(arr, fn)` | Transforms each element using function |
| `filter(arr, fn)` | Keeps elements where function returns true |
| `reduce(arr, fn, init)` | Reduces array to single value |
| `contains(arr, item)` | Checks if array contains item |
| `index(arr, item)` | Returns index of item (-1 if not found) |

### Hash Functions

| Function | Description |
|----------|-------------|
| `keys(hash)` | Returns array of all keys |
| `values(hash)` | Returns array of all values |

## Error Handling

Use `try`/`catch` for error handling:

```victoria
try {
    let result = riskyOperation()
    print(result)
} catch (err) {
    print("Error occurred: ${err}")
}
```

## Rust-Inspired Error Messages

Victoria provides beautiful, developer-friendly error messages **inspired by the Rust programming language**. When you make a mistake, Victoria helps you understand and fix it quickly with:

### Error Message Features

1. **Color-Coded Output**: Errors are highlighted in red, code context in cyan, and helpful suggestions in green.

2. **Precise Source Locations**: Every error shows the exact file, line number, and column where the issue occurred.

3. **Inline Code Snippets**: The actual source code is displayed with visual markers pointing to the problem.

4. **Actionable Suggestions**: Context-aware help messages suggest how to fix common mistakes.

5. **Error Codes**: Each error type has a unique code (e.g., `E0002`) for easy reference.

### Example: Undefined Variable

```
error[E0002]: identifier not found: myVariable
 --> script.vc:5:12
  |
4 | let x = 10
5 | let y = myVariable + x
  |         ^^^^^^^^^^ identifier not found: myVariable
6 |
  |
  = help: did you mean to declare 'myVariable' with 'let myVariable = ...'?
```

### Example: Type Mismatch

```
error: type mismatch: STRING + INTEGER
 --> script.vc:3:20
  |
2 | let greeting = "Hello"
3 | let result = greeting + 42
  |                      ^ type mismatch: STRING + INTEGER
4 |
  |
  = note: Victoria is a dynamically typed language, but operators require compatible types
  = help: use string() to convert integers to strings, or int() to convert strings to integers
```

### Example: Syntax Error

```
error[E0004]: expected '=' but found 'EOF'
 --> script.vc:2:1
  |
1 | let x
2 |
  | ^ expected '=' here
  |
  = note: reached end of file unexpectedly
  = help: variable declarations require an initial value: let name = value
```

### Common Error Codes

| Code | Description | Common Cause |
|------|-------------|-------------|
| `E0001` | Type mismatch | Mixing incompatible types in operations |
| `E0002` | Undefined identifier | Using a variable before declaring it |
| `E0003` | Unknown operator | Using an operator not defined for a type |
| `E0004` | Unexpected token | Syntax error or missing punctuation |
| `E0005` | Not a function | Trying to call something that isn't a function |
| `E0006` | Cannot index | Using `[]` on a type that doesn't support indexing |
| `E0007` | Division by zero | Dividing by zero |
| `E0008` | Property not found | Accessing a non-existent property |
| `E0009` | Struct not found | Using an undefined struct |
| `E0010` | Wrong argument count | Calling a function with wrong number of args |
| `E0011` | Invalid slice | Invalid slice indices or unsupported type |
| `E0012` | Spread error | Spread operator on non-array or wrong context |
| `E0013` | Invalid hash key | Using a non-hashable type as hash key |
| `E0014` | Argument type error | Passing wrong type to a function |
| `E0015` | Not iterable | Using for-in on non-iterable type |
| `E0016` | Range error | Invalid range parameters |
| `E0017` | Conversion error | Failed type conversion |
| `E0018` | Assignment error | Invalid assignment target |
| `E0019` | Operator error | Invalid use of ++/-- operators |
| `E0020` | Member access error | Dot notation on unsupported type |
| `E0021` | Empty reduce | reduce() on empty array without initial value |
| `E0022` | Join error | join() with non-string array elements |
| `E0100` | Parse error | General syntax/parsing error |
| `E0101` | Illegal character | Invalid character in source |
| `E0102` | Unterminated string | String literal missing closing quote |

### Smart Typo Detection

Victoria automatically suggests corrections for common mistakes from other languages:

| You typed | Victoria suggests |
|-----------|------------------|
| `println()` | `print()` |
| `str()` | `string()` |
| `len()` is called `length()` or `size()` | `len()` |
| `nil`, `None`, `undefined` | `null` |
| `fn`, `func`, `function`, `def` | `define` |
| `append()` | `push()` |
| `import`, `require` | `include` |
| `var` | `let` |
| `elif` | `else if` |

These error messages are designed to help developers—especially those learning to program—understand what went wrong and how to fix it, just like the Rust compiler does.

---

## DSA (Data Structures & Algorithms) Guide

Victoria is specifically designed to make learning and implementing data structures and algorithms easy and enjoyable. This section covers DSA-specific features and common patterns.

### Why Victoria for DSA?

| Feature | Victoria | Benefit for DSA |
|---------|----------|-----------------|
| `#make` directive | `#make MOD 1000000007` | C-like constants without includes |
| `enum` | `enum NodeState { UNVISITED, VISITING, VISITED }` | Clean state management for graphs |
| Character ops | `ord('a')`, `chr(97)`, `isDigit()` | Easy string parsing |
| Array slicing | `arr[1:5]`, `arr[:mid]`, `arr[mid:]` | Merge sort, divide & conquer |
| Hash assignment | `freq[c] = 1`, `memo[key] = result` | Frequency counting, memoization |
| First-class functions | `map()`, `filter()`, `reduce()` | Functional transformations |

---

### Common DSA Constants

```victoria
// Competitive programming essentials
#make MOD 1000000007        // Most common modulo
#make MOD2 998244353        // NTT-friendly prime
#make INF 999999999         // Infinity for shortest path
#make NINF -999999999       // Negative infinity
#make MAXN 100005           // Common array size
#make EPSILON 0.000001      // Float comparison threshold
```

---

### Pattern 1: Frequency Counting

```victoria
// Count character frequency (classic interview question)
define countFrequency(s) {
    let freq = {}
    for i in 0..len(s) {
        let c = s[i]
        if (freq[c] == null) {
            freq[c] = 1
        } else {
            freq[c] = freq[c] + 1
        }
    }
    return freq
}

// Usage
let freq = countFrequency("mississippi")
print(freq)  // {m: 1, i: 4, s: 4, p: 2}
```

---

### Pattern 2: Two Pointers

```victoria
// Check if palindrome
define isPalindrome(s) {
    let left = 0
    let right = len(s) - 1
    
    while (left < right) {
        if (toLower(s[left]) != toLower(s[right])) {
            return false
        }
        left = left + 1
        right = right - 1
    }
    return true
}

// Two Sum (sorted array)
define twoSum(arr, target) {
    let left = 0
    let right = len(arr) - 1
    
    while (left < right) {
        let sum = arr[left] + arr[right]
        if (sum == target) {
            return [left, right]
        } else if (sum < target) {
            left = left + 1
        } else {
            right = right - 1
        }
    }
    return []
}
```

---

### Pattern 3: Binary Search

```victoria
define binarySearch(arr, target) {
    let left = 0
    let right = len(arr) - 1
    
    while (left <= right) {
        // Use this form to avoid integer overflow
        let mid = left + (right - left) / 2
        
        if (arr[mid] == target) {
            return mid
        } else if (arr[mid] < target) {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return -1  // Not found
}

// Find lower bound (first >= target)
define lowerBound(arr, target) {
    let left = 0
    let right = len(arr)
    
    while (left < right) {
        let mid = left + (right - left) / 2
        if (arr[mid] < target) {
            left = mid + 1
        } else {
            right = mid
        }
    }
    return left
}
```

---

### Pattern 4: Memoization (Dynamic Programming)

```victoria
// Fibonacci with memoization
let fibMemo = {}

define fib(n) {
    if (n <= 1) {
        return n
    }
    
    // Check cache first
    if (fibMemo[n] != null) {
        return fibMemo[n]
    }
    
    // Compute and cache
    let result = fib(n - 1) + fib(n - 2)
    fibMemo[n] = result
    return result
}

// Now fib(40) runs instantly instead of taking seconds
print(fib(40))  // 102334155
```

---

### Pattern 5: Graph State Management with Enums

```victoria
// Perfect for cycle detection, topological sort
enum NodeState {
    UNVISITED,   // 0 - Not yet visited
    VISITING,    // 1 - Currently in DFS stack
    VISITED      // 2 - Completely processed
}

define hasCycle(graph, node, states) {
    if (states[node] == NodeState.VISITING) {
        return true  // Back edge found = cycle!
    }
    if (states[node] == NodeState.VISITED) {
        return false  // Already processed
    }
    
    states[node] = NodeState.VISITING
    
    for neighbor in graph[node] {
        if (hasCycle(graph, neighbor, states)) {
            return true
        }
    }
    
    states[node] = NodeState.VISITED
    return false
}
```

---

### Pattern 6: Divide and Conquer (Merge Sort)

```victoria
define mergeSort(arr) {
    if (len(arr) <= 1) {
        return arr
    }
    
    let mid = len(arr) / 2
    let left = mergeSort(arr[:mid])    // Array slicing!
    let right = mergeSort(arr[mid:])   // So clean!
    
    return merge(left, right)
}

define merge(left, right) {
    let result = []
    let i = 0
    let j = 0
    
    while (i < len(left) && j < len(right)) {
        if (left[i] <= right[j]) {
            result = push(result, left[i])
            i = i + 1
        } else {
            result = push(result, right[j])
            j = j + 1
        }
    }
    
    // Add remaining elements
    while (i < len(left)) {
        result = push(result, left[i])
        i = i + 1
    }
    while (j < len(right)) {
        result = push(result, right[j])
        j = j + 1
    }
    
    return result
}
```

---

### Pattern 7: Number Theory

```victoria
#make MOD 1000000007

// Greatest Common Divisor (Euclidean algorithm)
define gcd(a, b) {
    if (b == 0) {
        return a
    }
    return gcd(b, a % b)
}

// Least Common Multiple
define lcm(a, b) {
    return (a / gcd(a, b)) * b
}

// Modular exponentiation (a^b % MOD)
define modPow(base, exp, mod) {
    let result = 1
    base = base % mod
    
    while (exp > 0) {
        if (exp % 2 == 1) {
            result = (result * base) % mod
        }
        exp = exp / 2
        base = (base * base) % mod
    }
    return result
}

// Sieve of Eratosthenes
define sieve(n) {
    let isPrime = []
    for i in 0..n+1 {
        isPrime = push(isPrime, true)
    }
    isPrime[0] = false
    isPrime[1] = false
    
    for i in 2..n+1 {
        if (isPrime[i]) {
            let j = i * i
            while (j <= n) {
                isPrime[j] = false
                j = j + i
            }
        }
    }
    
    let primes = []
    for i in 0..n+1 {
        if (isPrime[i]) {
            primes = push(primes, i)
        }
    }
    return primes
}
```

---

### Pattern 8: String Parsing

```victoria
// Parse integers from a string
define parseNumbers(s) {
    let numbers = []
    let current = ""
    
    for i in 0..len(s) {
        let c = s[i]
        if (isDigit(c)) {
            current = current + c
        } else {
            if (len(current) > 0) {
                numbers = push(numbers, int(current))
                current = ""
            }
        }
    }
    
    // Don't forget the last number
    if (len(current) > 0) {
        numbers = push(numbers, int(current))
    }
    
    return numbers
}

// Validate parentheses
define isValidParentheses(s) {
    let count = 0
    for i in 0..len(s) {
        let c = s[i]
        if (c == "(") {
            count = count + 1
        } else if (c == ")") {
            count = count - 1
            if (count < 0) {
                return false
            }
        }
    }
    return count == 0
}
```

---

### Common Mistakes to Avoid

#### 1. Off-by-One Errors
```victoria
// WRONG: This misses the last element
for i in 0..len(arr) - 1 { }

// CORRECT: 0..n means 0 to n-1, which is all indices
for i in 0..len(arr) { }
```

#### 2. Integer Overflow
```victoria
// WRONG: May overflow for large values
let result = a * b

// CORRECT: Use modular arithmetic
#make MOD 1000000007
let result = ((a % MOD) * (b % MOD)) % MOD
```

#### 3. Missing Base Case
```victoria
// WRONG: No base case = infinite recursion
define fib(n) {
    return fib(n-1) + fib(n-2)
}

// CORRECT: Always have a base case
define fib(n) {
    if (n <= 1) { return n }
    return fib(n-1) + fib(n-2)
}
```

#### 4. Binary Search Boundaries
```victoria
// WRONG: May cause infinite loop
let mid = (left + right) / 2

// CORRECT: Prevents overflow and handles edge cases
let mid = left + (right - left) / 2
```

#### 5. Forgetting to Check Null/Empty
```victoria
// WRONG: May crash on empty array
return arr[0]

// CORRECT: Always check first
if (len(arr) == 0) {
    return null
}
return arr[0]
```

---

### Quick Reference Card

| Task | Victoria Code |
|------|---------------|
| Constant | `#make MOD 1000000007` |
| Enum | `enum State { A, B, C }` |
| Character to ASCII | `ord('A')` → 65 |
| ASCII to character | `chr(65)` → "A" |
| Is digit? | `isDigit('5')` → true |
| Is letter? | `isLetter('a')` → true |
| To uppercase | `toUpper("hello")` → "HELLO" |
| To lowercase | `toLower("HELLO")` → "hello" |
| Array slice | `arr[1:4]` → elements 1,2,3 |
| Slice from start | `arr[:3]` → first 3 |
| Slice to end | `arr[3:]` → from index 3 |
| Hash set | `freq[key] = 1` |
| Hash get | `freq[key]` |
| Hash check | `if (freq[key] != null)` |

---

Happy coding! ��� Victoria makes DSA fun and accessible.

# Victoria Language Documentation

Victoria is a dynamic, interpreted programming language designed for readability and expressiveness. It combines features from C-style languages (like curly braces) with Python-like simplicity.

## Table of Contents

- [Variables](#variables)
- [Data Types](#data-types)
- [Operators](#operators)
- [Control Flow](#control-flow)
- [Functions](#functions)
- [Data Structures](#data-structures)
- [Structs](#structs)
- [Modules](#modules)
- [Built-in Functions](#built-in-functions)

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

## Data Types

Victoria supports the following basic data types:

- **Integer**: `1`, `42`, `-10`
- **Float**: `3.14`, `0.001`
- **String**: `"Hello World"` or `` `multi-line string` ``
- **Boolean**: `true`, `false`
- **Null**: `null` (implicit in some contexts)
- **Array**: `[1, 2, 3]`
- **Hash**: `{"key": "value"}`

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

### Higher-Order Functions

Functions can take other functions as arguments:

```victoria
define applyTwice(f, x) {
    return f(f(x))
}

let double = define(n) { n * 2 }
print(applyTwice(double, 5))  // 20
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
include "./mylib.vc" // Local file
```

### Available Modules

- `os`: File system operations (`readFile`, `writeFile`, `exists`, etc.)
- `net`: Networking (`http.get`, `http.post`, `tcp.connect`, etc.)
- `math`: Math functions (`abs`, `sqrt`, `pow`, `sin`, `cos`, etc.)

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
| `E0003` | Unknown operator | Using an operator that doesn't exist for a type |
| `E0004` | Unexpected token | Syntax error or missing punctuation |
| `E0005` | Not a function | Trying to call something that isn't a function |
| `E0006` | Cannot index | Using `[]` on a type that doesn't support indexing |

These error messages are designed to help developers—especially those learning to program—understand what went wrong and how to fix it, just like the Rust compiler does.

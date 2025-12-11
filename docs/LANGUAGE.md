# Victoria Language Documentation

Victoria is a dynamic, interpreted programming language designed for readability and expressiveness. It combines features from C-style languages (like curly braces) with Python-like simplicity.

## Table of Contents

- [Variables](#variables)
- [Data Types](#data-types)
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

## Data Types

Victoria supports the following basic data types:

- **Integer**: `1`, `42`, `-10`
- **Float**: `3.14`, `0.001`
- **String**: `"Hello World"`
- **Boolean**: `true`, `false`
- **Null**: `null` (implicit in some contexts)

## Control Flow

### If-Else

```victoria
let x = 10
if (x > 5) {
    print("x is greater than 5")
} else {
    print("x is less than or equal to 5")
}
```

### Loops

**While Loop**

```victoria
let i = 0
while (i < 5) {
    print(i)
    let i = i + 1
}
```

**For Loop (C-style)**

```victoria
for (let i = 0; i < 5; i = i + 1) {
    print(i)
}
```

**For-In Loop (Python-style)**

```victoria
let arr = [1, 2, 3]
for (x in arr) {
    print(x)
}
```

## Functions

Functions are first-class citizens and are defined using the `define` keyword.

```victoria
define add(x, y) {
    return x + y
}

let result = add(5, 10)
print(result)
```

Functions can also be assigned to variables:

```victoria
let multiply = define(x, y) {
    return x * y
}
```

## Data Structures

### Arrays

Arrays are ordered lists of values.

```victoria
let arr = [1, 2, 3, "four", true]
print(arr[0]) // 1
print(len(arr)) // 5
```

### Hashes

Hashes are key-value pairs (dictionaries).

```victoria
let person = {
    "name": "Alice",
    "age": 30,
    "city": "Wonderland"
}

print(person["name"]) // Alice
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
    print("Hello, my name is " + self.name)
}

// Instantiate a struct
let p = Person {
    name: "Bob",
    age: 25
}

p.greet()
```

## Modules

You can include other Victoria files or built-in modules using `include`.

```victoria
include "net" // Built-in network module
include "./mylib.vc" // Local file
```

## Built-in Functions

Victoria comes with a standard library of built-in functions:

- `print(args...)`: Prints arguments to the console.
- `len(arg)`: Returns the length of a string or array.
- `range(start, end)`: Returns an array of integers from start to end (exclusive).
- `format(string, args...)`: Formats a string (similar to printf).
- `input(prompt)`: Reads input from the user.
- `int(arg)`: Converts a value to an integer.
- `string(arg)`: Converts a value to a string.
- `type(arg)`: Returns the type of the argument.
- `first(array)`: Returns the first element of an array.
- `last(array)`: Returns the last element of an array.
- `rest(array)`: Returns a new array containing all elements except the first.
- `push(array, element)`: Adds an element to the end of an array.
- `pop(array)`: Removes and returns the last element of an array.
- `split(string, separator)`: Splits a string into an array.
- `join(array, separator)`: Joins an array of strings into a single string.
- `contains(container, item)`: Checks if an array or hash contains an item.
- `keys(hash)`: Returns an array of keys in a hash.
- `values(hash)`: Returns an array of values in a hash.

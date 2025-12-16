# Victoria Language Levels

Victoria is designed to grow with you. Instead of throwing all features at you at once, the language naturally supports three levels of sophistication. You can write entire programs at any level, and mix levels freely within the same codebase.

## Level 1: Explorer

**Goal:** Write working code immediately, learn basic concepts.

At this level, you're learning what programming *is*. Victoria gets out of your way.

### What You Use

```victoria
// Variables
let name = "Alice"
let age = 25
let items = [1, 2, 3, 4, 5]

// Output
print("Hello, World!")
print("Name:", name)

// Basic control flow
if (age >= 18) {
    print("Adult")
} else {
    print("Minor")
}

// Simple loops
for item in items {
    print(item)
}

// Basic functions
define greet(name) {
    print("Hello, " + name + "!")
}

greet("Bob")
```

### What You Learn

- Variables hold values
- Programs run top to bottom
- Conditions let you make decisions
- Loops let you repeat things
- Functions let you reuse code

### Not Yet Needed

- Type annotations
- Return types
- Constants
- Advanced data structures
- Modules

---

## Level 2: Builder

**Goal:** Structure larger programs, understand data flow.

At this level, you're learning *how* to organize code. Victoria provides more tools.

### What You Use

```victoria
// Constants for values that shouldn't change
const MAX_USERS = 100
const API_URL = "https://api.example.com"

// Functions with return values
define calculateArea(width, height) {
    return width * height
}

define isEven(n) {
    return n % 2 == 0
}

// Using return values
let area = calculateArea(5, 10)
print("Area:", area)

// Structs for organizing data
struct User {
    name,
    email,
    age
}

let user = User{name: "Alice", email: "alice@example.com", age: 30}
print(user.name)

// Methods on structs
struct Rectangle {
    width,
    height
}

define Rectangle.area() {
    return self.width * self.height
}

let rect = Rectangle{width: 10, height: 5}
print("Rectangle area:", rect.area())

// Functional programming
let numbers = [1, 2, 3, 4, 5]
let doubled = map(numbers, x => x * 2)
let evens = filter(numbers, x => x % 2 == 0)
let sum = reduce(numbers, (a, b) => a + b, 0)

// String interpolation
let message = "Hello, ${user.name}! Your area is ${area}."
print(message)

// Using built-in modules
print(math.floor(3.7))      // 3
print(json.stringify(user)) // JSON output
print(time.now())           // Current time
```

### What You Learn

- Constants prevent accidental changes
- Functions should return values, not just print
- Structs group related data
- Methods add behavior to data
- Lambda expressions enable functional patterns
- Built-in modules provide common functionality

### Not Yet Needed

- Type annotations
- Complex error handling
- Multiple modules
- Advanced patterns

---

## Level 3: Engineer

**Goal:** Write robust, maintainable, production-quality code.

At this level, you're learning *why* certain practices matter. Victoria provides type safety and structure.

### What You Use

```victoria
// Typed variables catch errors early
let count:int = 0
let name:string = "Victoria"
let prices:array = [9.99, 19.99, 29.99]

// Typed constants
const MAX_RETRIES:int = 3
const BASE_URL:string = "https://api.example.com"

// Typed function parameters and return types
define add(a:int, b:int) -> int {
    return a + b
}

define greet(name:string, age:int) -> string {
    return "Hello, ${name}! You are ${age} years old."
}

define calculateAverage(numbers:array) -> float {
    let sum = reduce(numbers, (a, b) => a + b, 0)
    return sum / len(numbers)
}

// The 'any' type for flexibility when needed
define printValue(value:any) -> void {
    print("Value:", value)
}

// Type errors are caught at runtime with helpful messages
// let x:int = "hello"  // ERROR: type mismatch

// Structured error handling
try {
    let result = riskyOperation()
    print("Success:", result)
} catch (error) {
    print("Error occurred:", error)
}

// Include statements for code organization
include "utils.vc"
include "models.vc"

// Well-structured struct with typed methods
struct User {
    name,
    email,
    age
}

define User.validate() -> bool {
    if (self.name == "" or self.email == "") {
        return false
    }
    if (self.age < 0 or self.age > 150) {
        return false
    }
    return true
}

define User.toJSON() -> string {
    return json.stringify({
        "name": self.name,
        "email": self.email,
        "age": self.age
    })
}

// Using everything together
define createUser(name:string, email:string, age:int) -> any {
    let user = User{name: name, email: email, age: age}
    if (!user.validate()) {
        return null
    }
    return user
}

let user = createUser("Alice", "alice@example.com", 30)
if (user != null) {
    print(user.toJSON())
}
```

### What You Learn

- Types catch bugs before they become problems
- Function signatures serve as documentation
- Error handling makes programs robust
- Code organization matters at scale
- The right constraints enable better code

### Victoria's Type System

| Type | Description | Example |
|------|-------------|---------|
| `int` | Integer numbers | `let x:int = 42` |
| `float` | Decimal numbers | `let pi:float = 3.14` |
| `string` | Text | `let s:string = "hello"` |
| `bool` | True/false | `let ok:bool = true` |
| `array` | List of values | `let arr:array = [1, 2, 3]` |
| `map` | Key-value pairs | `let m:map = {"a": 1}` |
| `any` | Any type (opt-out) | `let x:any = anything` |
| `void` | No return value | `-> void` |

---

## Mixing Levels

One of Victoria's strengths is that levels aren't walls—they're guidelines. You can mix styles in the same file:

```victoria
// Level 1: Quick script at the top
let data = [1, 2, 3, 4, 5]

// Level 2: Some structure
define process(items) {
    return map(items, x => x * 2)
}

// Level 3: Type safety where it matters
define validateAge(age:int) -> bool {
    return age >= 0 and age <= 150
}

// Use them together
let processed = process(data)
print(processed)

if (validateAge(25)) {
    print("Valid age")
}
```

This allows you to:
- Prototype quickly without types
- Add types as requirements solidify
- Enforce types in critical sections

---

## Level Progression Path

### Week 1-2: Explorer
- Write simple scripts
- Learn variables and printing
- Understand if/else and loops
- Create basic functions

### Week 3-4: Builder
- Use constants appropriately
- Return values from functions
- Create structs for data
- Explore built-in modules
- Try lambda expressions

### Month 2+: Engineer
- Add type annotations
- Handle errors properly
- Organize code into files
- Write self-documenting code

---

## The Victoria Promise

At every level, Victoria promises:
- **Clear error messages** that teach, not frustrate
- **Consistent syntax** that doesn't change between levels
- **No hidden gotchas** waiting to trip you up
- **Real capability** to build useful programs

Start where you're comfortable. Grow at your own pace. Victoria will be there at every step.

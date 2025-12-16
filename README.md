# Victoria

<p align="center">
  <img src="media/QueenVictoria.jpg" alt="Queen Victoria">
</p>

<p align="center">
  <strong>The programming language you wish you had learned first.</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Victoria-v0.1--alpha-blueviolet?style=flat-square">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square">
  <img src="https://img.shields.io/badge/build-passing-brightgreen?style=flat-square">
</p>

---

## What is Victoria?

Victoria is a **learning-first programming language** that combines the readability of Python with the structure of JavaScript and C-like style, wrapped in error messages that actually teach you something.

```victoria
// Simple and readable
let name = "World"
print("Hello, ${name}!")

// Grows with you
define greet(name:string) -> string {
    return "Hello, ${name}!"
}
```

## Who is Victoria For?

| You Are | Victoria Helps You |
|---------|-------------------|
| **A complete beginner** | Learn programming without cryptic errors |
| **An educator** | Teach with a language designed for teaching |
| **A hobbyist** | Prototype ideas quickly without boilerplate |
| **A developer** | Script tasks with minimal setup |

## Why Victoria Exists

Most programming languages were designed for experts, then adapted for beginners as an afterthought. Victoria takes the opposite approach.

**Every design decision asks:** *"How this would not confuse the learner?"*

The result:
- **Error messages that teach** — Not "SyntaxError: unexpected token", but explanations with suggestions.
- **Optional complexity** — Start simple, add types and structure when you're ready
- **Zero friction** — No build tools, no package managers, no configuration files

<details>
<summary><strong>See Victoria's error messages in action</strong></summary>

```
error[E0001]: type mismatch: cannot assign string to variable of type int
 --> example.vc:2:1
  |
2 | let x:int = "hello";
  | ^^^^ type mismatch
  |
  = note: Victoria supports optional static typing for type safety
  = help: either change the value to an int, or change the type annotation to :string
```

</details>

---

## Quick Start

### Installation

```bash
# Clone and build (requires Go 1.21+)
git clone https://github.com/theawakener0/VictoriaLang.git
cd VictoriaLang
go build -o victoria cmd/victoria/main.go
```

### Hello World

```bash
# Interactive REPL
./victoria

# Run a script
echo 'print("Hello, Victoria!")' > hello.vc
./victoria hello.vc
```

---

## Language at a Glance

<table>
<tr>
<td width="50%">

### Variables & Types
```victoria
// Dynamic (no types required)
let name = "Alice"
let age = 25

// Static (optional type safety)
let count:int = 0
let message:string = "Hello"

// Constants
const PI = 3.14159
const MAX_USERS:int = 100
```

</td>
<td width="50%">

### Functions
```victoria
// Simple function
define add(a, b) {
    return a + b
}

// Typed function
define greet(name:string) -> string {
    return "Hello, ${name}!"
}

// Lambda expressions
let double = x => x * 2
let sum = (a, b) => a + b
```

</td>
</tr>
<tr>
<td>

### Control Flow
```victoria
// Conditions
if (age >= 18) {
    print("Adult")
} else {
    print("Minor")
}

// Loops
for i in 0..5 {
    print(i)
}

for item in items {
    print(item)
}
```

</td>
<td>

### Data Structures
```victoria
// Arrays
let nums = [1, 2, 3, 4, 5]
let slice = nums[1:3]  // [2, 3]

// Hashes
let user = {
    "name": "Alice",
    "age": 30
}

// Structs
struct Point { x, y }
let p = new Point(10, 20)
```

</td>
</tr>
</table>

### Functional Programming
```victoria
let numbers = [1, 2, 3, 4, 5]

let doubled = map(numbers, x => x * 2)         // [2, 4, 6, 8, 10]
let evens = filter(numbers, x => x % 2 == 0)   // [2, 4]
let sum = reduce(numbers, (a, b) => a + b, 0)  // 15
```

### Built-in Modules
```victoria
// Math
print(math.floor(3.7))     // 3
print(math.random(1, 100)) // random 1-100

// JSON
let data = json.parse('{"name": "Victoria"}')
print(json.stringify(data))

// Time
print(time.now())
print(time.format(time.now(), "YYYY-MM-DD"))
```

---

## Documentation

| Document | Description |
|----------|-------------|
| [Language Reference](docs/LANGUAGE.md) | Complete syntax and features |
| [Philosophy](docs/PHILOSOPHY.md) | Why Victoria is designed this way |
| [Language Levels](docs/LEVELS.md) | Learn at your own pace |
| [Roadmap](docs/ROADMAP.md) | Where Victoria is heading |

---

## Examples

Check out the [`examples/`](examples/) directory for complete programs:

- `hello.vc` — Hello World
- `fib.vc` — Fibonacci sequence  
- `types_demo.vc` — Type system showcase
- `server.vc` — Simple HTTP server
- `structs.vc` — Object-oriented programming

---

## Contributing

Victoria is open source and welcomes contributions!

- **Found a bug?** [Open an issue](https://github.com/theawakener0/VictoriaLang/issues)
- **Have an idea?** Start a discussion
See the [Roadmap](docs/ROADMAP.md) for priority areas.

---

## License

MIT License — Use Victoria however you'd like.

---

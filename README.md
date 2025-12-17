# Victoria Programming Language

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

Victoria is a **learning-first programming language** designed specifically for:
- **Beginners** learning to code for the first time
- **Students** studying Data Structures and Algorithms (DSA)
- **Educators** teaching programming concepts
- **Competitive programmers** who want clean, readable code

It combines the best features of Python, JavaScript, Go, and C, without the confusion.

```victoria
// Clean and readable
let name = "World"
print("Hello, ${name}!")

// Grows with you — add types when ready
define factorial(n:int) -> int {
    if (n <= 1) { return 1 }
    return n * factorial(n - 1)
}

// DSA-ready constants
#make MOD 1000000007
#make INF 999999999
```

---

## Why Victoria Over Other Languages?

### The Problem with Existing Languages

| Language | Problem for Beginners |
|----------|----------------------|
| **Python** | Indentation errors are confusing; no optional typing until 3.10; slow for competitive programming |
| **JavaScript** | `undefined`, `null`, `NaN` chaos; `==` vs `===`; async confusion |
| **C/C++** | Pointer nightmares; segfaults; manual memory management; cryptic errors |
| **Java** | `public static void main(String[] args)` just to print "Hello World" |
| **Go** | Great but strict; error handling boilerplate; no generics (until recently) |

### Victoria's Solution

| Feature | Victoria | Python | JavaScript | C++ | Java |
|---------|----------|--------|------------|-----|------|
| **Readable syntax** | ✅ | ✅ | ⚠️ | ❌ | ⚠️ |
| **Optional static typing** | ✅ | ⚠️ (hints only) | ❌ | N/A (required) | N/A (required) |
| **Helpful error messages** | ✅ | ⚠️ | ❌ | ❌ | ⚠️ |
| **Zero configuration** | ✅ | ✅ | ⚠️ | ❌ | ❌ |
| **DSA-ready constants** | ✅ (`#make`) | ❌ | ❌ | ✅ (`#define`) | ❌ |
| **Enum support** | ✅ | ⚠️ (limited) | ❌ | ✅ | ✅ |
| **Character operations** | ✅ | ✅ | ⚠️ | ✅ | ✅ |

---

## Real-World Comparisons

### Example 1: Character Frequency (Classic DSA Problem)

<table>
<tr>
<th>Victoria</th>
<th>Python</th>
</tr>
<tr>
<td>

```victoria
define countFreq(s) {
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

print(countFreq("hello"))
// {h: 1, e: 1, l: 2, o: 1}
```

</td>
<td>

```python
def count_freq(s):
    freq = {}
    for c in s:
        if c not in freq:
            freq[c] = 1
        else:
            freq[c] += 1
    return freq

print(count_freq("hello"))
# {'h': 1, 'e': 1, 'l': 2, 'o': 1}
```

</td>
</tr>
</table>

**Victoria Advantage**: Same clarity as Python, but with C-style braces that make scope crystal clear for beginners.

---

### Example 2: Binary Search with Type Safety

<table>
<tr>
<th>Victoria</th>
<th>C++</th>
</tr>
<tr>
<td>

```victoria
define binarySearch(arr, target:int) -> int {
    let left = 0
    let right = len(arr) - 1
    
    while (left <= right) {
        let mid = left + (right - left) / 2
        if (arr[mid] == target) {
            return mid
        } else if (arr[mid] < target) {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return -1
}

let nums = [1, 3, 5, 7, 9, 11]
print(binarySearch(nums, 7))  // 3
```

</td>
<td>

```cpp
#include <vector>
#include <iostream>

int binarySearch(std::vector<int>& arr, int target) {
    int left = 0;
    int right = arr.size() - 1;
    
    while (left <= right) {
        int mid = left + (right - left) / 2;
        if (arr[mid] == target) {
            return mid;
        } else if (arr[mid] < target) {
            left = mid + 1;
        } else {
            right = mid - 1;
        }
    }
    return -1;
}

int main() {
    std::vector<int> nums = {1, 3, 5, 7, 9, 11};
    std::cout << binarySearch(nums, 7) << std::endl;
    return 0;
}
```

</td>
</tr>
</table>

**Victoria Advantage**: No includes, no `main()` boilerplate, no `std::` prefixes — just the algorithm.

---

### Example 3: Graph State Management with Enums

<table>
<tr>
<th>Victoria</th>
<th>Java</th>
</tr>
<tr>
<td>

```victoria
enum NodeState {
    UNVISITED,
    VISITING,
    VISITED
}

define hasCycle(graph, node, states) {
    if (states[node] == NodeState.VISITING) {
        return true  // Back edge = cycle
    }
    if (states[node] == NodeState.VISITED) {
        return false
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

</td>
<td>

```java
enum NodeState {
    UNVISITED, VISITING, VISITED
}

public class Graph {
    public boolean hasCycle(
        Map<Integer, List<Integer>> graph,
        int node,
        NodeState[] states
    ) {
        if (states[node] == NodeState.VISITING) {
            return true;
        }
        if (states[node] == NodeState.VISITED) {
            return false;
        }
        
        states[node] = NodeState.VISITING;
        
        for (int neighbor : graph.get(node)) {
            if (hasCycle(graph, neighbor, states)) {
                return true;
            }
        }
        
        states[node] = NodeState.VISITED;
        return false;
    }
}
```

</td>
</tr>
</table>

**Victoria Advantage**: Same enum power as Java, without the class ceremony.

---

### Example 4: Competitive Programming Setup

<table>
<tr>
<th>Victoria</th>
<th>C++</th>
</tr>
<tr>
<td>

```victoria
// Competitive programming template
#make MOD 1000000007
#make INF 999999999
#make MAXN 100005

// Power function with modulo
define power(base:int, exp:int) -> int {
    let result = 1
    base = base % MOD
    while (exp > 0) {
        if (exp % 2 == 1) {
            result = (result * base) % MOD
        }
        exp = exp / 2
        base = (base * base) % MOD
    }
    return result
}

print(power(2, 10))  // 1024
```

</td>
<td>

```cpp
#include <bits/stdc++.h>
using namespace std;

#define MOD 1000000007
#define INF 999999999
#define MAXN 100005

long long power(long long base, long long exp) {
    long long result = 1;
    base %= MOD;
    while (exp > 0) {
        if (exp & 1) {
            result = (result * base) % MOD;
        }
        exp >>= 1;
        base = (base * base) % MOD;
    }
    return result;
}

int main() {
    cout << power(2, 10) << endl;
    return 0;
}
```

</td>
</tr>
</table>

**Victoria Advantage**: `#make` works like `#define` but is cleaner. No `#include`, no `main()`.

---

## Error Messages That Teach

Victoria's error messages are designed to **teach**, not frustrate.

### Type Error
```
error: type mismatch: cannot assign string to variable of type int
   --> example.vc:3:14
    |
2 | let count:int = 0
3 | count = "hello"
    |         ^^^^^^^ expected int, got string
    |
  = note: Victoria uses optional static typing for safer code
  = help: either change the value to an integer, or change the type to :string
```

### Undefined Variable
```
error: identifier not found: userName
   --> example.vc:5:7
    |
5 | print(userName)
    |       ^^^^^^^^ not defined in this scope
    |
  = help: did you mean 'username'? (check capitalization)
  = note: variables must be declared with 'let' before use
```

### Constant Reassignment
```
error: cannot reassign constant variable: PI
   --> example.vc:4:1
    |
2 | const PI = 3.14159
3 | 
4 | PI = 3.0
    | ^^ cannot reassign constant
    |
  = note: constants declared with 'const' cannot be changed
  = help: use 'let' instead if you need to reassign this variable
```

### Array Index Out of Bounds
```
error: array index out of bounds: 10
   --> example.vc:3:1
    |
2 | let arr = [1, 2, 3]
3 | print(arr[10])
    |       ^^^^^^^ index 10 is out of bounds for array of length 3
    |
  = note: array indices start at 0 and go up to length-1
  = help: valid indices for this array are 0, 1, 2
```

---

## Quick Start

### Installation

```bash
# Clone and build (requires Go 1.21+)
git clone https://github.com/theawakener0/VictoriaLang.git
cd VictoriaLang/victoria
go build -o victoria cmd/victoria/main.go

# Or with one command
go install github.com/theawakener0/VictoriaLang/cmd/victoria@latest
```

### Hello World

```bash
# Interactive REPL
./victoria

# Run a script
echo 'print("Hello, Victoria!")' > hello.vc
./victoria hello.vc
```

### Your First DSA Program

Create `two_sum.vc`:
```victoria
// Classic Two Sum problem
define twoSum(nums, target:int) {
    let seen = {}
    for i in 0..len(nums) {
        let complement = target - nums[i]
        if (seen[complement] != null) {
            return [seen[complement], i]
        }
        seen[nums[i]] = i
    }
    return []
}

let nums = [2, 7, 11, 15]
print(twoSum(nums, 9))  // [0, 1]
```

Run it:
```bash
./victoria two_sum.vc
```

---

## Language Features

### Core Syntax

```victoria
// Variables (dynamic or typed)
let name = "Victoria"
let age:int = 1
const PI = 3.14159

// Functions
define greet(name:string) -> string {
    return "Hello, ${name}!"
}

// Lambda expressions
let double = x => x * 2
let add = (a, b) => a + b

// Control flow
if (age >= 18) {
    print("Adult")
} else {
    print("Minor")
}

// Loops
for i in 0..10 {
    print(i)
}

for item in items {
    print(item)
}
```

### DSA Features

```victoria
// Compile-time constants
#make MOD 1000000007
#make MAXN 100005

// Enums for state management
enum State { UNVISITED, VISITING, VISITED }

// Character operations
let ascii = ord('A')      // 65
let char = chr(65)        // "A"
let isNum = isDigit('5')  // true

// Array slicing
let arr = [1, 2, 3, 4, 5]
let first = arr[:2]       // [1, 2]
let last = arr[3:]        // [4, 5]
let mid = arr[1:4]        // [2, 3, 4]

// Hash maps with index assignment
let freq = {}
freq["a"] = 1
freq["a"] = freq["a"] + 1
```

### Functional Programming

```victoria
let nums = [1, 2, 3, 4, 5]

let doubled = map(nums, x => x * 2)       // [2, 4, 6, 8, 10]
let evens = filter(nums, x => x % 2 == 0) // [2, 4]
let sum = reduce(nums, (a, b) => a + b, 0) // 15
```

### Built-in Modules

```victoria
// Math
print(math.floor(3.7))      // 3
print(math.ceil(3.2))       // 4
print(math.sqrt(16))        // 4
print(math.random(1, 100))  // Random 1-100

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
| [Language Reference](docs/LANGUAGE.md) | Complete syntax, types, and features |
| [Philosophy](docs/PHILOSOPHY.md) | Why Victoria is designed this way |
| [Language Levels](docs/LEVELS.md) | Progressive learning path |
| [Roadmap](docs/ROADMAP.md) | Future plans and features |

---

## Example Programs

The [`examples/`](examples/) directory contains many ready-to-run programs:

| File | Description |
|------|-------------|
| `hello.vc` | Hello World |
| `fib.vc` | Fibonacci sequence |
| `dsa_demo.vc` | **DSA algorithms showcase** |
| `types_demo.vc` | Type system demo |
| `structs.vc` | Object-oriented programming |
| `server.vc` | Simple HTTP server |
| `try_catch.vc` | Error handling |

### Run the DSA Demo
```bash
./victoria examples/dsa_demo.vc
```

Output:
```
=== Algorithm: Character Frequency ===
Frequency of 'mississippi': {m: 1, i: 4, s: 4, p: 2}

=== Algorithm: Binary Search ===
binarySearch([1,3,5,7,9,11], 7): 3

=== Algorithm: Fibonacci (Memoization) ===
fib(30) = 832040

=== Algorithm: Sieve of Eratosthenes ===
Primes up to 50: [2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47]
```

---

## Contributing

Victoria is open source and welcomes contributions!

- **Found a bug?** [Open an issue](https://github.com/theawakener0/VictoriaLang/issues)
- **Have a feature idea?** Start a discussion
- **Want to contribute?** Check the [Roadmap](docs/ROADMAP.md)

### Priority Areas
- [ ] More built-in data structures (Stack, Queue, Heap)
- [ ] Standard library expansion
- [ ] Neovim & VS Code extension with syntax highlighting
- [ ] Package manager
- [ ] Compiler (currently interpreted only)

---

## License

MIT License — Use Victoria however you'd like.

---

<p align="center">
  <a href="docs/LANGUAGE.md">Learn More</a> •
  <a href="examples/">Examples</a> •
  <a href="https://github.com/theawakener0/VictoriaLang/issues">Report Bug</a>
</p>

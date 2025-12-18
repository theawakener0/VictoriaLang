# Victoria Roadmap

This document outlines the planned development trajectory for Victoria. Our goal is to evolve Victoria into a robust learning-first language while maintaining its core philosophy of simplicity and helpfulness.

## Current Status: v0.1 (Alpha)

Victoria is currently in active development. The core language is functional and can be used for learning, prototyping, and small projects.

### What Works Today

**Core Language**
- Variables (`let`) and constants (`const`)
- All basic data types (int, float, string, bool, array, hash)
- Control flow (if/else, for, while, switch)
- Functions with closures
- Structs with methods
- Lambda expressions
- String interpolation
- **Isolated Module System** - Modules run in their own scope to prevent name collisions

**Type System**
- Optional type annotations for variables
- Typed function parameters
- Return type annotations
- Runtime type checking
- Helpful type error messages

**Built-in Modules**
- `os` - Operating system interaction (now with `exec`, `platform`, `arch`)
- `math` - Mathematical operations
- `json` - JSON parsing/serialization
- `time` - Date/time handling
- `net` - HTTP client & TCP networking
- `path` - Path manipulation utilities

**Developer Experience**
- Rust-inspired error messages with suggestions
- REPL for interactive development
- File-based script execution
- `vic` Package Manager - Install libraries directly from GitHub

---

## Non-Goals

To stay focused, Victoria explicitly does **not** aim to:

- **Replace production languages** - Use Go, Rust, or Python for production systems
- **Maximize performance** - Clarity over speed
- **Support every paradigm** - Practical over pure
- **Have the largest ecosystem** - Quality over quantity

---

## How to Contribute

Victoria is open source and welcomes contributions! Here's how you can help:

### Code Contributions
1. Check the [GitHub Issues](https://github.com/theawakener0/VictoriaLang/issues) for open tasks
2. Look for "good first issue" labels for beginner-friendly tasks
3. Submit pull requests with tests

### Non-Code Contributions
- **Documentation** - Improve explanations and examples
- **Error Messages** - Suggest clearer, more helpful messages
- **Bug Reports** - Detailed reports help us fix issues faster
- **Feature Requests** - Share ideas for new features
- **Teaching** - Use Victoria and share your experience

### Priority Areas

Right now, we especially need help with:
1. **Error message improvements** - Making every error a learning moment
2. **Documentation** - More examples and tutorials
3. **Testing** - More test coverage
4. **VS Code extension** - Syntax highlighting and snippets

---

## Feedback

We'd love to hear from you! If you have:
- Feature suggestions
- Pain points to share
- Success stories

Please open an issue on GitHub or reach out to the maintainers.

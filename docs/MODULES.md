# Victoria Module System

Victoria uses a simple and powerful module system that supports both built-in modules and user-defined libraries.

## Using Modules

Use the `include` statement to load modules:

```victoria
include "os"
include "path"
include "mylib"
```

- **Built-in Modules**: `os`, `net`, `std`, `math`, `json`, `time`, `path`.
- **User Modules**: Victoria looks for `.vc` files in:
  1. The current directory.
  2. The `victoria_modules/` directory.
  3. `~/.victoria/modules/` (Global).

When a user module is included, it is evaluated in an **isolated environment**. All top-level variables and functions are exported as a Hash object named after the module.

Example `mylib.vc`:
```victoria
let version = "1.0.0"
let add = define(a: int, b: int) -> int { return a + b }
```

Example usage:
```victoria
include "mylib"
print(mylib.version)
print(mylib.add(1, 2))
```

## Package Manager (vic)

Victoria comes with a package manager called `vic`.

### Commands

- `vic init`: Initialize a new project with a `vic.json` file.
- `vic install <user>/<repo>`: Install a package from GitHub into `victoria_modules/`.
- `vic list`: List all installed packages.

### Creating and Publishing a Library

1. Create a new GitHub repository.
2. Add your Victoria code. Use `index.vc` or `main.vc` as the entry point.
3. Push your code to GitHub.
4. Others can now install your library using:
   ```bash
   vic install your-username/your-repo-name
   ```

### Example Library Structure

```
your-repo/
  ├── index.vc      (Entry point)
  ├── utils.vc
  └── README.md
```

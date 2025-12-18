package evaluator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"victoria/ast"
	"victoria/lexer"
	"victoria/object"
	"victoria/parser"
)

// moduleRegistry holds all available modules
var moduleRegistry = make(map[string]func() *object.Hash)

// createModule creates a hash object from a map of methods
func createModule(methods map[string]object.Object) *object.Hash {
	pairs := make(map[object.HashKey]object.HashPair)
	for name, method := range methods {
		key := &object.String{Value: name}
		pairs[key.HashKey()] = object.HashPair{Key: key, Value: method}
	}
	return &object.Hash{Pairs: pairs}
}

// createSocketObject creates a socket object for TCP connections
func createSocketObject(conn net.Conn) *object.Hash {
	methods := map[string]object.Object{
		"read": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				reader := bufio.NewReader(conn)
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						return &object.String{Value: ""}
					}
					return newError("failed to read from socket: %s", err.Error())
				}
				return &object.String{Value: strings.TrimSuffix(line, "\n")}
			},
		},
		"readAll": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				data, err := io.ReadAll(conn)
				if err != nil {
					return newError("failed to read from socket: %s", err.Error())
				}
				return &object.String{Value: string(data)}
			},
		},
		"write": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.STRING_OBJ {
					return newError("argument to `write` must be STRING")
				}
				_, err := conn.Write([]byte(args[0].(*object.String).Value))
				if err != nil {
					return newError("failed to write to socket: %s", err.Error())
				}
				return NULL
			},
		},
		"close": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				conn.Close()
				return NULL
			},
		},
	}
	return createModule(methods)
}

// RegisterBuiltinModules registers all built-in modules
func RegisterBuiltinModules() {
	// OS Module
	moduleRegistry["os"] = func() *object.Hash {
		osMethods := map[string]object.Object{
			"readFile": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `readFile` must be STRING, got %s", args[0].Type())
					}
					filename := args[0].(*object.String).Value
					content, err := os.ReadFile(filename)
					if err != nil {
						return newError("could not read file: %s", err.Error())
					}
					return &object.String{Value: string(content)}
				},
			},
			"writeFile": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("first argument to `writeFile` must be STRING")
					}
					if args[1].Type() != object.STRING_OBJ {
						return newError("second argument to `writeFile` must be STRING")
					}
					filename := args[0].(*object.String).Value
					content := args[1].(*object.String).Value
					err := os.WriteFile(filename, []byte(content), 0644)
					if err != nil {
						return newError("could not write file: %s", err.Error())
					}
					return TRUE
				},
			},
			"remove": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `remove` must be STRING")
					}
					filename := args[0].(*object.String).Value
					err := os.Remove(filename)
					if err != nil {
						return newError("could not remove file: %s", err.Error())
					}
					return TRUE
				},
			},
			"exists": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `exists` must be STRING")
					}
					filename := args[0].(*object.String).Value
					if _, err := os.Stat(filename); os.IsNotExist(err) {
						return FALSE
					}
					return TRUE
				},
			},
			"exit": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					code := 0
					if len(args) == 1 {
						if args[0].Type() != object.INTEGER_OBJ {
							return newError("argument to `exit` must be INTEGER")
						}
						code = int(args[0].(*object.Integer).Value)
					}
					os.Exit(code)
					return NULL
				},
			},
			"mkdir": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `mkdir` must be STRING")
					}
					dirname := args[0].(*object.String).Value
					err := os.MkdirAll(dirname, 0755)
					if err != nil {
						return newError("could not create directory: %s", err.Error())
					}
					return TRUE
				},
			},
			"readDir": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `readDir` must be STRING")
					}
					dirname := args[0].(*object.String).Value
					entries, err := os.ReadDir(dirname)
					if err != nil {
						return newError("could not read directory: %s", err.Error())
					}
					elements := make([]object.Object, len(entries))
					for i, entry := range entries {
						elements[i] = &object.String{Value: entry.Name()}
					}
					return &object.Array{Elements: elements}
				},
			},
			"stat": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `stat` must be STRING")
					}
					filename := args[0].(*object.String).Value
					info, err := os.Stat(filename)
					if err != nil {
						return newError("could not stat file: %s", err.Error())
					}
					// Return a hash with file info
					pairs := make(map[object.HashKey]object.HashPair)
					nameKey := &object.String{Value: "name"}
					pairs[nameKey.HashKey()] = object.HashPair{Key: nameKey, Value: &object.String{Value: info.Name()}}
					sizeKey := &object.String{Value: "size"}
					pairs[sizeKey.HashKey()] = object.HashPair{Key: sizeKey, Value: &object.Integer{Value: info.Size()}}
					isDir := FALSE
					if info.IsDir() {
						isDir = TRUE
					}
					isDirKey := &object.String{Value: "isDir"}
					pairs[isDirKey.HashKey()] = object.HashPair{Key: isDirKey, Value: isDir}
					modTimeKey := &object.String{Value: "modTime"}
					pairs[modTimeKey.HashKey()] = object.HashPair{Key: modTimeKey, Value: &object.Integer{Value: info.ModTime().Unix()}}
					return &object.Hash{Pairs: pairs}
				},
			},
			"rename": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
						return newError("arguments to `rename` must be STRING")
					}
					oldPath := args[0].(*object.String).Value
					newPath := args[1].(*object.String).Value
					err := os.Rename(oldPath, newPath)
					if err != nil {
						return newError("could not rename file: %s", err.Error())
					}
					return TRUE
				},
			},
			"getwd": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					wd, err := os.Getwd()
					if err != nil {
						return newError("could not get working directory: %s", err.Error())
					}
					return &object.String{Value: wd}
				},
			},
			"chdir": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `chdir` must be STRING")
					}
					dir := args[0].(*object.String).Value
					err := os.Chdir(dir)
					if err != nil {
						return newError("could not change directory: %s", err.Error())
					}
					return TRUE
				},
			},
			"env": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) == 0 {
						// Return all environment variables as hash
						envMap := make(map[object.HashKey]object.HashPair)
						for _, env := range os.Environ() {
							parts := strings.SplitN(env, "=", 2)
							if len(parts) == 2 {
								key := &object.String{Value: parts[0]}
								envMap[key.HashKey()] = object.HashPair{Key: key, Value: &object.String{Value: parts[1]}}
							}
						}
						return &object.Hash{Pairs: envMap}
					} else if len(args) == 1 {
						if args[0].Type() != object.STRING_OBJ {
							return newError("argument to `env` must be STRING")
						}
						return &object.String{Value: os.Getenv(args[0].(*object.String).Value)}
					} else if len(args) == 2 {
						if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
							return newError("arguments to `env` must be STRING")
						}
						os.Setenv(args[0].(*object.String).Value, args[1].(*object.String).Value)
						return TRUE
					}
					return newError("wrong number of arguments. got=%d, want=0, 1, or 2", len(args))
				},
			},
			"args": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					elements := make([]object.Object, len(os.Args))
					for i, arg := range os.Args {
						elements[i] = &object.String{Value: arg}
					}
					return &object.Array{Elements: elements}
				},
			},
			"exec": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 1 {
						return newError("wrong number of arguments. got=%d, want=at least 1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("first argument to `exec` must be STRING (command)")
					}
					name := args[0].(*object.String).Value
					cmdArgs := []string{}
					for i := 1; i < len(args); i++ {
						cmdArgs = append(cmdArgs, args[i].Inspect())
					}
					cmd := exec.Command(name, cmdArgs...)
					output, err := cmd.CombinedOutput()
					if err != nil {
						return newError("command failed: %s\nOutput: %s", err.Error(), string(output))
					}
					return &object.String{Value: string(output)}
				},
			},
			"platform": &object.String{Value: runtime.GOOS},
			"arch":     &object.String{Value: runtime.GOARCH},
			"pid":      &object.Integer{Value: int64(os.Getpid())},
			"hostname": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					name, err := os.Hostname()
					if err != nil {
						return newError("could not get hostname: %s", err.Error())
					}
					return &object.String{Value: name}
				},
			},
			"user": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					u, err := user.Current()
					if err != nil {
						return newError("could not get current user: %s", err.Error())
					}
					pairs := make(map[object.HashKey]object.HashPair)
					usernameKey := &object.String{Value: "username"}
					pairs[usernameKey.HashKey()] = object.HashPair{Key: usernameKey, Value: &object.String{Value: u.Username}}
					nameKey := &object.String{Value: "name"}
					pairs[nameKey.HashKey()] = object.HashPair{Key: nameKey, Value: &object.String{Value: u.Name}}
					homeKey := &object.String{Value: "home"}
					pairs[homeKey.HashKey()] = object.HashPair{Key: homeKey, Value: &object.String{Value: u.HomeDir}}
					uidKey := &object.String{Value: "uid"}
					pairs[uidKey.HashKey()] = object.HashPair{Key: uidKey, Value: &object.String{Value: u.Uid}}
					gidKey := &object.String{Value: "gid"}
					pairs[gidKey.HashKey()] = object.HashPair{Key: gidKey, Value: &object.String{Value: u.Gid}}
					return &object.Hash{Pairs: pairs}
				},
			},
			"tempDir": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					return &object.String{Value: os.TempDir()}
				},
			},
			"chmod": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.INTEGER_OBJ {
						return newError("arguments to `chmod` must be STRING (path) and INTEGER (mode)")
					}
					path := args[0].(*object.String).Value
					mode := args[1].(*object.Integer).Value
					err := os.Chmod(path, os.FileMode(mode))
					if err != nil {
						return newError("could not chmod: %s", err.Error())
					}
					return TRUE
				},
			},
			"kill": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `kill` must be INTEGER (pid)")
					}
					pid := int(args[0].(*object.Integer).Value)
					proc, err := os.FindProcess(pid)
					if err != nil {
						return newError("could not find process: %s", err.Error())
					}
					err = proc.Kill()
					if err != nil {
						return newError("could not kill process: %s", err.Error())
					}
					return TRUE
				},
			},
		}
		return createModule(osMethods)
	}

	// Net Module
	moduleRegistry["net"] = func() *object.Hash {
		netMethods := map[string]object.Object{
			"get": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `get` must be STRING")
					}
					url := args[0].(*object.String).Value
					resp, err := http.Get(url)
					if err != nil {
						return newError("HTTP GET failed: %s", err.Error())
					}
					defer resp.Body.Close()
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return newError("failed to read response: %s", err.Error())
					}
					// Return a hash with status, statusCode, and body
					pairs := make(map[object.HashKey]object.HashPair)
					statusKey := &object.String{Value: "status"}
					pairs[statusKey.HashKey()] = object.HashPair{Key: statusKey, Value: &object.String{Value: resp.Status}}
					statusCodeKey := &object.String{Value: "statusCode"}
					pairs[statusCodeKey.HashKey()] = object.HashPair{Key: statusCodeKey, Value: &object.Integer{Value: int64(resp.StatusCode)}}
					bodyKey := &object.String{Value: "body"}
					pairs[bodyKey.HashKey()] = object.HashPair{Key: bodyKey, Value: &object.String{Value: string(body)}}
					return &object.Hash{Pairs: pairs}
				},
			},
			"post": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 2 || len(args) > 3 {
						return newError("wrong number of arguments. got=%d, want=2 or 3", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("first argument to `post` must be STRING (URL)")
					}
					if args[1].Type() != object.STRING_OBJ {
						return newError("second argument to `post` must be STRING (body)")
					}
					url := args[0].(*object.String).Value
					body := args[1].(*object.String).Value
					contentType := "application/json"
					if len(args) == 3 {
						if args[2].Type() != object.STRING_OBJ {
							return newError("third argument to `post` must be STRING (content-type)")
						}
						contentType = args[2].(*object.String).Value
					}
					resp, err := http.Post(url, contentType, strings.NewReader(body))
					if err != nil {
						return newError("HTTP POST failed: %s", err.Error())
					}
					defer resp.Body.Close()
					respBody, err := io.ReadAll(resp.Body)
					if err != nil {
						return newError("failed to read response: %s", err.Error())
					}
					pairs := make(map[object.HashKey]object.HashPair)
					statusKey := &object.String{Value: "status"}
					pairs[statusKey.HashKey()] = object.HashPair{Key: statusKey, Value: &object.String{Value: resp.Status}}
					statusCodeKey := &object.String{Value: "statusCode"}
					pairs[statusCodeKey.HashKey()] = object.HashPair{Key: statusCodeKey, Value: &object.Integer{Value: int64(resp.StatusCode)}}
					bodyKey := &object.String{Value: "body"}
					pairs[bodyKey.HashKey()] = object.HashPair{Key: bodyKey, Value: &object.String{Value: string(respBody)}}
					return &object.Hash{Pairs: pairs}
				},
			},
			"head": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `head` must be STRING")
					}
					url := args[0].(*object.String).Value
					resp, err := http.Head(url)
					if err != nil {
						return newError("HTTP HEAD failed: %s", err.Error())
					}
					defer resp.Body.Close()
					pairs := make(map[object.HashKey]object.HashPair)
					statusKey := &object.String{Value: "status"}
					pairs[statusKey.HashKey()] = object.HashPair{Key: statusKey, Value: &object.String{Value: resp.Status}}
					statusCodeKey := &object.String{Value: "statusCode"}
					pairs[statusCodeKey.HashKey()] = object.HashPair{Key: statusCodeKey, Value: &object.Integer{Value: int64(resp.StatusCode)}}
					// Add headers
					headersKey := &object.String{Value: "headers"}
					headerPairs := make(map[object.HashKey]object.HashPair)
					for key, values := range resp.Header {
						hKey := &object.String{Value: key}
						headerPairs[hKey.HashKey()] = object.HashPair{Key: hKey, Value: &object.String{Value: strings.Join(values, ", ")}}
					}
					pairs[headersKey.HashKey()] = object.HashPair{Key: headersKey, Value: &object.Hash{Pairs: headerPairs}}
					return &object.Hash{Pairs: pairs}
				},
			},
			"delete": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `delete` must be STRING")
					}
					url := args[0].(*object.String).Value
					client := &http.Client{}
					req, err := http.NewRequest("DELETE", url, nil)
					if err != nil {
						return newError("failed to create request: %s", err.Error())
					}
					resp, err := client.Do(req)
					if err != nil {
						return newError("HTTP DELETE failed: %s", err.Error())
					}
					defer resp.Body.Close()
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return newError("failed to read response: %s", err.Error())
					}
					pairs := make(map[object.HashKey]object.HashPair)
					statusKey := &object.String{Value: "status"}
					pairs[statusKey.HashKey()] = object.HashPair{Key: statusKey, Value: &object.String{Value: resp.Status}}
					statusCodeKey := &object.String{Value: "statusCode"}
					pairs[statusCodeKey.HashKey()] = object.HashPair{Key: statusCodeKey, Value: &object.Integer{Value: int64(resp.StatusCode)}}
					bodyKey := &object.String{Value: "body"}
					pairs[bodyKey.HashKey()] = object.HashPair{Key: bodyKey, Value: &object.String{Value: string(body)}}
					return &object.Hash{Pairs: pairs}
				},
			},
			"put": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 2 || len(args) > 3 {
						return newError("wrong number of arguments. got=%d, want=2 or 3", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("first argument to `put` must be STRING (URL)")
					}
					if args[1].Type() != object.STRING_OBJ {
						return newError("second argument to `put` must be STRING (body)")
					}
					url := args[0].(*object.String).Value
					body := args[1].(*object.String).Value
					contentType := "application/json"
					if len(args) == 3 {
						if args[2].Type() != object.STRING_OBJ {
							return newError("third argument to `put` must be STRING (content-type)")
						}
						contentType = args[2].(*object.String).Value
					}
					client := &http.Client{}
					req, err := http.NewRequest("PUT", url, strings.NewReader(body))
					if err != nil {
						return newError("failed to create request: %s", err.Error())
					}
					req.Header.Set("Content-Type", contentType)
					resp, err := client.Do(req)
					if err != nil {
						return newError("HTTP PUT failed: %s", err.Error())
					}
					defer resp.Body.Close()
					respBody, err := io.ReadAll(resp.Body)
					if err != nil {
						return newError("failed to read response: %s", err.Error())
					}
					pairs := make(map[object.HashKey]object.HashPair)
					statusKey := &object.String{Value: "status"}
					pairs[statusKey.HashKey()] = object.HashPair{Key: statusKey, Value: &object.String{Value: resp.Status}}
					statusCodeKey := &object.String{Value: "statusCode"}
					pairs[statusCodeKey.HashKey()] = object.HashPair{Key: statusCodeKey, Value: &object.Integer{Value: int64(resp.StatusCode)}}
					bodyKey := &object.String{Value: "body"}
					pairs[bodyKey.HashKey()] = object.HashPair{Key: bodyKey, Value: &object.String{Value: string(respBody)}}
					return &object.Hash{Pairs: pairs}
				},
			},
			"parseQuery": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
						return newError("argument to `parseQuery` must be STRING")
					}
					query := args[0].(*object.String).Value
					values, err := url.ParseQuery(query)
					if err != nil {
						return newError("failed to parse query: %s", err.Error())
					}
					pairs := make(map[object.HashKey]object.HashPair)
					for key, vals := range values {
						k := &object.String{Value: key}
						v := &object.String{Value: strings.Join(vals, ",")}
						pairs[k.HashKey()] = object.HashPair{Key: k, Value: v}
					}
					return &object.Hash{Pairs: pairs}
				},
			},
			"lookupHost": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
						return newError("argument to `lookupHost` must be STRING")
					}
					host := args[0].(*object.String).Value
					addrs, err := net.LookupHost(host)
					if err != nil {
						return newError("failed to lookup host: %s", err.Error())
					}
					elements := make([]object.Object, len(addrs))
					for i, addr := range addrs {
						elements[i] = &object.String{Value: addr}
					}
					return &object.Array{Elements: elements}
				},
			},
			"interfaces": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					ifaces, err := net.Interfaces()
					if err != nil {
						return newError("failed to get interfaces: %s", err.Error())
					}
					elements := make([]object.Object, len(ifaces))
					for i, iface := range ifaces {
						pairs := make(map[object.HashKey]object.HashPair)
						nameKey := &object.String{Value: "name"}
						pairs[nameKey.HashKey()] = object.HashPair{Key: nameKey, Value: &object.String{Value: iface.Name}}
						indexKey := &object.String{Value: "index"}
						pairs[indexKey.HashKey()] = object.HashPair{Key: indexKey, Value: &object.Integer{Value: int64(iface.Index)}}
						mtuKey := &object.String{Value: "mtu"}
						pairs[mtuKey.HashKey()] = object.HashPair{Key: mtuKey, Value: &object.Integer{Value: int64(iface.MTU)}}
						macKey := &object.String{Value: "mac"}
						pairs[macKey.HashKey()] = object.HashPair{Key: macKey, Value: &object.String{Value: iface.HardwareAddr.String()}}
						elements[i] = &object.Hash{Pairs: pairs}
					}
					return &object.Array{Elements: elements}
				},
			},
			"dial": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("first argument to `dial` must be STRING (host)")
					}
					if args[1].Type() != object.INTEGER_OBJ {
						return newError("second argument to `dial` must be INTEGER (port)")
					}
					host := args[0].(*object.String).Value
					port := args[1].(*object.Integer).Value
					addr := fmt.Sprintf("%s:%d", host, port)
					conn, err := net.Dial("tcp", addr)
					if err != nil {
						return newError("failed to connect: %s", err.Error())
					}
					return createSocketObject(conn)
				},
			},
			"dialUdp": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("first argument to `dialUdp` must be STRING (host)")
					}
					if args[1].Type() != object.INTEGER_OBJ {
						return newError("second argument to `dialUdp` must be INTEGER (port)")
					}
					host := args[0].(*object.String).Value
					port := args[1].(*object.Integer).Value
					addr := fmt.Sprintf("%s:%d", host, port)
					conn, err := net.Dial("udp", addr)
					if err != nil {
						return newError("failed to connect UDP: %s", err.Error())
					}
					return createSocketObject(conn)
				},
			},
			"listenTcp": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.INTEGER_OBJ {
						return newError("first argument to `listenTcp` must be INTEGER (port)")
					}
					port := args[0].(*object.Integer).Value
					handler := args[1]

					// Create TCP listener
					listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
					if err != nil {
						return newError("failed to start TCP server: %s", err.Error())
					}

					fmt.Printf("TCP server listening on port %d\n", port)

					// Accept connections
					for {
						conn, err := listener.Accept()
						if err != nil {
							continue
						}

						// Create connection object
						connObj := createSocketObject(conn)

						// Call handler with connection object
						if fn, ok := handler.(*object.Function); ok {
							env := extendFunctionEnv(fn, []object.Object{connObj})
							Eval(fn.Body, env)
						} else if builtin, ok := handler.(*object.Builtin); ok {
							builtin.Fn(connObj)
						}
					}
				},
			},
			"listenUdp": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.INTEGER_OBJ {
						return newError("first argument to `listenUdp` must be INTEGER (port)")
					}
					port := args[0].(*object.Integer).Value
					handler := args[1]

					addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
					conn, err := net.ListenUDP("udp", addr)
					if err != nil {
						return newError("failed to start UDP server: %s", err.Error())
					}

					fmt.Printf("UDP server listening on port %d\n", port)

					buf := make([]byte, 4096)
					for {
						n, remoteAddr, err := conn.ReadFromUDP(buf)
						if err != nil {
							continue
						}

						data := string(buf[:n])
						remote := remoteAddr.String()

						pairs := make(map[object.HashKey]object.HashPair)
						dataKey := &object.String{Value: "data"}
						pairs[dataKey.HashKey()] = object.HashPair{Key: dataKey, Value: &object.String{Value: data}}
						remoteKey := &object.String{Value: "remote"}
						pairs[remoteKey.HashKey()] = object.HashPair{Key: remoteKey, Value: &object.String{Value: remote}}
						packetObj := &object.Hash{Pairs: pairs}

						if fn, ok := handler.(*object.Function); ok {
							env := extendFunctionEnv(fn, []object.Object{packetObj})
							Eval(fn.Body, env)
						} else if builtin, ok := handler.(*object.Builtin); ok {
							builtin.Fn(packetObj)
						}
					}
				},
			},
			"request": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 2 {
						return newError("wrong number of arguments. got=%d, want=at least 2", len(args))
					}
					method := strings.ToUpper(args[0].Inspect())
					urlStr := args[1].Inspect()
					var body io.Reader
					if len(args) > 2 && args[2].Type() == object.STRING_OBJ {
						body = strings.NewReader(args[2].(*object.String).Value)
					}

					req, err := http.NewRequest(method, urlStr, body)
					if err != nil {
						return newError("failed to create request: %s", err.Error())
					}

					if len(args) > 3 && args[3].Type() == object.HASH_OBJ {
						headers := args[3].(*object.Hash)
						for _, pair := range headers.Pairs {
							req.Header.Set(pair.Key.Inspect(), pair.Value.Inspect())
						}
					}

					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						return newError("HTTP request failed: %s", err.Error())
					}
					defer resp.Body.Close()

					respBody, _ := io.ReadAll(resp.Body)
					pairs := make(map[object.HashKey]object.HashPair)
					statusKey := &object.String{Value: "status"}
					pairs[statusKey.HashKey()] = object.HashPair{Key: statusKey, Value: &object.String{Value: resp.Status}}
					statusCodeKey := &object.String{Value: "statusCode"}
					pairs[statusCodeKey.HashKey()] = object.HashPair{Key: statusCodeKey, Value: &object.Integer{Value: int64(resp.StatusCode)}}
					bodyKey := &object.String{Value: "body"}
					pairs[bodyKey.HashKey()] = object.HashPair{Key: bodyKey, Value: &object.String{Value: string(respBody)}}

					headerPairs := make(map[object.HashKey]object.HashPair)
					for k, v := range resp.Header {
						hk := &object.String{Value: k}
						hv := &object.String{Value: strings.Join(v, ",")}
						headerPairs[hk.HashKey()] = object.HashPair{Key: hk, Value: hv}
					}
					headersKey := &object.String{Value: "headers"}
					pairs[headersKey.HashKey()] = object.HashPair{Key: headersKey, Value: &object.Hash{Pairs: headerPairs}}

					return &object.Hash{Pairs: pairs}
				},
			},
			"listen": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.INTEGER_OBJ {
						return newError("first argument to `listen` must be INTEGER (port)")
					}
					port := args[0].(*object.Integer).Value
					handler := args[1]

					http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						// Build request object
						reqPairs := make(map[object.HashKey]object.HashPair)

						methodKey := &object.String{Value: "method"}
						reqPairs[methodKey.HashKey()] = object.HashPair{Key: methodKey, Value: &object.String{Value: r.Method}}

						pathKey := &object.String{Value: "path"}
						reqPairs[pathKey.HashKey()] = object.HashPair{Key: pathKey, Value: &object.String{Value: r.URL.Path}}

						queryKey := &object.String{Value: "query"}
						reqPairs[queryKey.HashKey()] = object.HashPair{Key: queryKey, Value: &object.String{Value: r.URL.RawQuery}}

						// Headers
						headerPairs := make(map[object.HashKey]object.HashPair)
						for key, values := range r.Header {
							hKey := &object.String{Value: key}
							headerPairs[hKey.HashKey()] = object.HashPair{Key: hKey, Value: &object.String{Value: strings.Join(values, ", ")}}
						}
						headersKey := &object.String{Value: "headers"}
						reqPairs[headersKey.HashKey()] = object.HashPair{Key: headersKey, Value: &object.Hash{Pairs: headerPairs}}

						// Body
						body, _ := io.ReadAll(r.Body)
						bodyKey := &object.String{Value: "body"}
						reqPairs[bodyKey.HashKey()] = object.HashPair{Key: bodyKey, Value: &object.String{Value: string(body)}}

						reqObj := &object.Hash{Pairs: reqPairs}

						// Call handler
						var result object.Object
						if fn, ok := handler.(*object.Function); ok {
							env := extendFunctionEnv(fn, []object.Object{reqObj})
							result = Eval(fn.Body, env)
							result = unwrapReturnValue(result)
						} else if builtin, ok := handler.(*object.Builtin); ok {
							result = builtin.Fn(reqObj)
						}

						// Process result
						if result != nil {
							switch res := result.(type) {
							case *object.String:
								w.Header().Set("Content-Type", "text/plain")
								_, _ = w.Write([]byte(res.Value))
							case *object.Hash:
								// Check for status, headers, body
								for _, pair := range res.Pairs {
									if keyStr, ok := pair.Key.(*object.String); ok {
										switch keyStr.Value {
										case "status":
											if statusInt, ok := pair.Value.(*object.Integer); ok {
												w.WriteHeader(int(statusInt.Value))
											}
										case "headers":
											if headersHash, ok := pair.Value.(*object.Hash); ok {
												for _, hPair := range headersHash.Pairs {
													if hKey, ok := hPair.Key.(*object.String); ok {
														if hVal, ok := hPair.Value.(*object.String); ok {
															w.Header().Set(hKey.Value, hVal.Value)
														}
													}
												}
											}
										case "body":
											if bodyStr, ok := pair.Value.(*object.String); ok {
												_, _ = w.Write([]byte(bodyStr.Value))
											}
										}
									}
								}
							default:
								_, _ = w.Write([]byte(result.Inspect()))
							}
						}
					})

					fmt.Printf("HTTP server listening on port %d\n", port)
					err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
					if err != nil {
						return newError("failed to start server: %s", err.Error())
					}
					return NULL
				},
			},
			"serve": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					if args[0].Type() != object.INTEGER_OBJ {
						return newError("first argument to `serve` must be INTEGER (port)")
					}
					if args[1].Type() != object.HASH_OBJ {
						return newError("second argument to `serve` must be HASH (routes)")
					}
					port := args[0].(*object.Integer).Value
					routes := args[1].(*object.Hash)

					mux := http.NewServeMux()
					for _, pair := range routes.Pairs {
						path := pair.Key.Inspect()
						handler := pair.Value

						mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
							// Build request object
							reqPairs := make(map[object.HashKey]object.HashPair)
							methodKey := &object.String{Value: "method"}
							reqPairs[methodKey.HashKey()] = object.HashPair{Key: methodKey, Value: &object.String{Value: r.Method}}
							pathKey := &object.String{Value: "path"}
							reqPairs[pathKey.HashKey()] = object.HashPair{Key: pathKey, Value: &object.String{Value: r.URL.Path}}
							queryKey := &object.String{Value: "query"}
							reqPairs[queryKey.HashKey()] = object.HashPair{Key: queryKey, Value: &object.String{Value: r.URL.RawQuery}}

							body, _ := io.ReadAll(r.Body)
							bodyKey := &object.String{Value: "body"}
							reqPairs[bodyKey.HashKey()] = object.HashPair{Key: bodyKey, Value: &object.String{Value: string(body)}}

							reqObj := &object.Hash{Pairs: reqPairs}

							var result object.Object
							if fn, ok := handler.(*object.Function); ok {
								env := extendFunctionEnv(fn, []object.Object{reqObj})
								result = Eval(fn.Body, env)
								result = unwrapReturnValue(result)
							} else if builtin, ok := handler.(*object.Builtin); ok {
								result = builtin.Fn(reqObj)
							}

							if result != nil {
								switch res := result.(type) {
								case *object.String:
									w.Write([]byte(res.Value))
								case *object.Hash:
									for _, p := range res.Pairs {
										k := p.Key.Inspect()
										if k == "status" {
											if si, ok := p.Value.(*object.Integer); ok {
												w.WriteHeader(int(si.Value))
											}
										} else if k == "body" {
											w.Write([]byte(p.Value.Inspect()))
										} else if k == "headers" {
											if hh, ok := p.Value.(*object.Hash); ok {
												for _, hp := range hh.Pairs {
													w.Header().Set(hp.Key.Inspect(), hp.Value.Inspect())
												}
											}
										}
									}
								default:
									w.Write([]byte(result.Inspect()))
								}
							}
						})
					}

					fmt.Printf("HTTP server (mux) listening on port %d\n", port)
					err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
					if err != nil {
						return newError("failed to start server: %s", err.Error())
					}
					return NULL
				},
			},
		}
		return createModule(netMethods)
	}

	// Std Module
	moduleRegistry["std"] = func() *object.Hash {
		return createModule(map[string]object.Object{
			"version":  &object.String{Value: "1.0.0"},
			"first":    builtins["first"],
			"last":     builtins["last"],
			"rest":     builtins["rest"],
			"push":     builtins["push"],
			"pop":      builtins["pop"],
			"split":    builtins["split"],
			"join":     builtins["join"],
			"contains": builtins["contains"],
			"index":    builtins["index"],
			"upper":    builtins["upper"],
			"lower":    builtins["lower"],
			"keys":     builtins["keys"],
			"values":   builtins["values"],
		})
	}

	// Math Module
	moduleRegistry["math"] = func() *object.Hash {
		mathMethods := map[string]object.Object{
			"pi": &object.Float{Value: math.Pi},
			"e":  &object.Float{Value: math.E},
			"abs": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					switch arg := args[0].(type) {
					case *object.Integer:
						if arg.Value < 0 {
							return &object.Integer{Value: -arg.Value}
						}
						return arg
					case *object.Float:
						return &object.Float{Value: math.Abs(arg.Value)}
					default:
						return newError("argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())
					}
				},
			},
			"sin": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					val := getNumericValue(args[0])
					if val == nil {
						return newError("argument to `sin` must be FLOAT or INTEGER")
					}
					return &object.Float{Value: math.Sin(*val)}
				},
			},
			"cos": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					val := getNumericValue(args[0])
					if val == nil {
						return newError("argument to `cos` must be FLOAT or INTEGER")
					}
					return &object.Float{Value: math.Cos(*val)}
				},
			},
			"tan": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					val := getNumericValue(args[0])
					if val == nil {
						return newError("argument to `tan` must be FLOAT or INTEGER")
					}
					return &object.Float{Value: math.Tan(*val)}
				},
			},
			"sqrt": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					val := getNumericValue(args[0])
					if val == nil {
						return newError("argument to `sqrt` must be FLOAT or INTEGER")
					}
					return &object.Float{Value: math.Sqrt(*val)}
				},
			},
			"pow": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 2 {
						return newError("wrong number of arguments. got=%d, want=2", len(args))
					}
					x := getNumericValue(args[0])
					y := getNumericValue(args[1])
					if x == nil {
						return newError("argument 1 to `pow` must be FLOAT or INTEGER")
					}
					if y == nil {
						return newError("argument 2 to `pow` must be FLOAT or INTEGER")
					}
					return &object.Float{Value: math.Pow(*x, *y)}
				},
			},
			"floor": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					switch arg := args[0].(type) {
					case *object.Integer:
						return arg
					case *object.Float:
						return &object.Integer{Value: int64(math.Floor(arg.Value))}
					default:
						return newError("argument to `floor` must be INTEGER or FLOAT, got %s", args[0].Type())
					}
				},
			},
			"ceil": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					switch arg := args[0].(type) {
					case *object.Integer:
						return arg
					case *object.Float:
						return &object.Integer{Value: int64(math.Ceil(arg.Value))}
					default:
						return newError("argument to `ceil` must be INTEGER or FLOAT, got %s", args[0].Type())
					}
				},
			},
			"round": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					switch arg := args[0].(type) {
					case *object.Integer:
						return arg
					case *object.Float:
						return &object.Integer{Value: int64(math.Round(arg.Value))}
					default:
						return newError("argument to `round` must be INTEGER or FLOAT, got %s", args[0].Type())
					}
				},
			},
			"min": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 2 {
						return newError("wrong number of arguments. got=%d, want=at least 2", len(args))
					}
					minVal := math.MaxFloat64
					isAllIntegers := true
					for _, arg := range args {
						var val float64
						switch a := arg.(type) {
						case *object.Integer:
							val = float64(a.Value)
						case *object.Float:
							val = a.Value
							isAllIntegers = false
						default:
							return newError("arguments to `min` must be INTEGER or FLOAT, got %s", arg.Type())
						}
						if val < minVal {
							minVal = val
						}
					}
					if isAllIntegers {
						return &object.Integer{Value: int64(minVal)}
					}
					return &object.Float{Value: minVal}
				},
			},
			"max": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 2 {
						return newError("wrong number of arguments. got=%d, want=at least 2", len(args))
					}
					maxVal := -math.MaxFloat64
					isAllIntegers := true
					for _, arg := range args {
						var val float64
						switch a := arg.(type) {
						case *object.Integer:
							val = float64(a.Value)
						case *object.Float:
							val = a.Value
							isAllIntegers = false
						default:
							return newError("arguments to `max` must be INTEGER or FLOAT, got %s", arg.Type())
						}
						if val > maxVal {
							maxVal = val
						}
					}
					if isAllIntegers {
						return &object.Integer{Value: int64(maxVal)}
					}
					return &object.Float{Value: maxVal}
				},
			},
			"random": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) == 0 {
						return &object.Float{Value: rand.Float64()}
					} else if len(args) == 1 {
						if args[0].Type() != object.INTEGER_OBJ {
							return newError("argument to `random` must be INTEGER, got %s", args[0].Type())
						}
						n := args[0].(*object.Integer).Value
						if n <= 0 {
							return newError("argument to `random` must be positive")
						}
						return &object.Integer{Value: rand.Int63n(n)}
					} else if len(args) == 2 {
						if args[0].Type() != object.INTEGER_OBJ || args[1].Type() != object.INTEGER_OBJ {
							return newError("arguments to `random` must be INTEGER")
						}
						min := args[0].(*object.Integer).Value
						max := args[1].(*object.Integer).Value
						if max < min {
							return newError("max must be >= min in random(min, max)")
						}
						return &object.Integer{Value: min + rand.Int63n(max-min+1)}
					}
					return newError("wrong number of arguments. got=%d, want=0, 1, or 2", len(args))
				},
			},
			"log": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					val := getNumericValue(args[0])
					if val == nil {
						return newError("argument to `log` must be FLOAT or INTEGER")
					}
					return &object.Float{Value: math.Log(*val)}
				},
			},
			"log10": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					val := getNumericValue(args[0])
					if val == nil {
						return newError("argument to `log10` must be FLOAT or INTEGER")
					}
					return &object.Float{Value: math.Log10(*val)}
				},
			},
		}
		return createModule(mathMethods)
	}

	// JSON Module
	moduleRegistry["json"] = func() *object.Hash {
		jsonMethods := map[string]object.Object{
			"parse": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `json.parse` must be STRING, got %s", args[0].Type())
					}
					jsonStr := args[0].(*object.String).Value
					return parseJSON(jsonStr)
				},
			},
			"stringify": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 1 || len(args) > 2 {
						return newError("wrong number of arguments. got=%d, want=1 or 2", len(args))
					}
					indent := ""
					if len(args) == 2 {
						if args[1].Type() == object.INTEGER_OBJ {
							spaces := args[1].(*object.Integer).Value
							for i := int64(0); i < spaces; i++ {
								indent += " "
							}
						} else if args[1].Type() == object.STRING_OBJ {
							indent = args[1].(*object.String).Value
						}
					}
					return stringifyJSON(args[0], indent)
				},
			},
			"valid": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("wrong number of arguments. got=%d, want=1", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("argument to `json.valid` must be STRING, got %s", args[0].Type())
					}
					jsonStr := args[0].(*object.String).Value
					var js interface{}
					if err := json.Unmarshal([]byte(jsonStr), &js); err != nil {
						return FALSE
					}
					return TRUE
				},
			},
		}
		return createModule(jsonMethods)
	}

	// Time Module
	moduleRegistry["time"] = func() *object.Hash {
		timeMethods := map[string]object.Object{
			"now": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					return &object.Integer{Value: time.Now().Unix()}
				},
			},
			"nowMs": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
				},
			},
			"format": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 1 || len(args) > 2 {
						return newError("wrong number of arguments. got=%d, want=1 or 2", len(args))
					}
					var t time.Time
					if args[0].Type() == object.INTEGER_OBJ {
						t = time.Unix(args[0].(*object.Integer).Value, 0)
					} else {
						return newError("first argument to `time.format` must be INTEGER (unix timestamp)")
					}
					layout := "2006-01-02 15:04:05"
					if len(args) == 2 {
						if args[1].Type() != object.STRING_OBJ {
							return newError("second argument to `time.format` must be STRING (format)")
						}
						layout = convertTimeFormat(args[1].(*object.String).Value)
					}
					return &object.String{Value: t.Format(layout)}
				},
			},
			"parse": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 1 || len(args) > 2 {
						return newError("wrong number of arguments. got=%d, want=1 or 2", len(args))
					}
					if args[0].Type() != object.STRING_OBJ {
						return newError("first argument to `time.parse` must be STRING")
					}
					dateStr := args[0].(*object.String).Value
					layout := "2006-01-02 15:04:05"
					if len(args) == 2 {
						if args[1].Type() != object.STRING_OBJ {
							return newError("second argument to `time.parse` must be STRING (format)")
						}
						layout = convertTimeFormat(args[1].(*object.String).Value)
					}
					t, err := time.Parse(layout, dateStr)
					if err != nil {
						return newError("failed to parse time: %s", err.Error())
					}
					return &object.Integer{Value: t.Unix()}
				},
			},
			"year": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.year` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Year())}
				},
			},
			"month": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.month` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Month())}
				},
			},
			"day": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.day` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Day())}
				},
			},
			"hour": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.hour` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Hour())}
				},
			},
			"minute": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.minute` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Minute())}
				},
			},
			"second": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.second` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Second())}
				},
			},
			"weekday": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.weekday` must be INTEGER (unix timestamp)")
					}
					t := time.Unix(args[0].(*object.Integer).Value, 0)
					return &object.Integer{Value: int64(t.Weekday())}
				},
			},
			"sleep": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
						return newError("argument to `time.sleep` must be INTEGER (milliseconds)")
					}
					ms := args[0].(*object.Integer).Value
					time.Sleep(time.Duration(ms) * time.Millisecond)
					return NULL
				},
			},
			"date": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					var t time.Time
					if len(args) == 0 {
						t = time.Now()
					} else if len(args) == 1 && args[0].Type() == object.INTEGER_OBJ {
						t = time.Unix(args[0].(*object.Integer).Value, 0)
					} else {
						return newError("argument to `time.date` must be INTEGER (unix timestamp) or no arguments")
					}
					return &object.String{Value: t.Format("2006-01-02")}
				},
			},
			"time": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					var t time.Time
					if len(args) == 0 {
						t = time.Now()
					} else if len(args) == 1 && args[0].Type() == object.INTEGER_OBJ {
						t = time.Unix(args[0].(*object.Integer).Value, 0)
					} else {
						return newError("argument to `time.time` must be INTEGER (unix timestamp) or no arguments")
					}
					return &object.String{Value: t.Format("15:04:05")}
				},
			},
		}
		return createModule(timeMethods)
	}

	// Path Module
	moduleRegistry["path"] = func() *object.Hash {
		pathMethods := map[string]object.Object{
			"join": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					paths := make([]string, len(args))
					for i, arg := range args {
						if arg.Type() != object.STRING_OBJ {
							return newError("arguments to `path.join` must be STRING")
						}
						paths[i] = arg.(*object.String).Value
					}
					return &object.String{Value: filepath.Join(paths...)}
				},
			},
			"base": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
						return newError("argument to `path.base` must be STRING")
					}
					return &object.String{Value: filepath.Base(args[0].(*object.String).Value)}
				},
			},
			"dir": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
						return newError("argument to `path.dir` must be STRING")
					}
					return &object.String{Value: filepath.Dir(args[0].(*object.String).Value)}
				},
			},
			"ext": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
						return newError("argument to `path.ext` must be STRING")
					}
					return &object.String{Value: filepath.Ext(args[0].(*object.String).Value)}
				},
			},
			"abs": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
						return newError("argument to `path.abs` must be STRING")
					}
					abs, err := filepath.Abs(args[0].(*object.String).Value)
					if err != nil {
						return newError("could not get absolute path: %s", err.Error())
					}
					return &object.String{Value: abs}
				},
			},
			"isAbs": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
						return newError("argument to `path.isAbs` must be STRING")
					}
					if filepath.IsAbs(args[0].(*object.String).Value) {
						return TRUE
					}
					return FALSE
				},
			},
		}
		return createModule(pathMethods)
	}
}

// Helper function to get numeric value as float64 pointer
func getNumericValue(obj object.Object) *float64 {
	switch o := obj.(type) {
	case *object.Integer:
		val := float64(o.Value)
		return &val
	case *object.Float:
		return &o.Value
	default:
		return nil
	}
}

// parseJSON converts a JSON string to Victoria objects
func parseJSON(jsonStr string) object.Object {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return newError("failed to parse JSON: %s", err.Error())
	}
	return jsonToObject(data)
}

// jsonToObject converts a Go interface{} to Victoria object
func jsonToObject(data interface{}) object.Object {
	if data == nil {
		return NULL
	}
	switch v := data.(type) {
	case bool:
		if v {
			return TRUE
		}
		return FALSE
	case float64:
		if v == float64(int64(v)) {
			return &object.Integer{Value: int64(v)}
		}
		return &object.Float{Value: v}
	case string:
		return &object.String{Value: v}
	case []interface{}:
		elements := make([]object.Object, len(v))
		for i, elem := range v {
			elements[i] = jsonToObject(elem)
		}
		return &object.Array{Elements: elements}
	case map[string]interface{}:
		pairs := make(map[object.HashKey]object.HashPair)
		for key, val := range v {
			keyObj := &object.String{Value: key}
			valObj := jsonToObject(val)
			pairs[keyObj.HashKey()] = object.HashPair{Key: keyObj, Value: valObj}
		}
		return &object.Hash{Pairs: pairs}
	default:
		return newError("unsupported JSON type: %T", v)
	}
}

// stringifyJSON converts a Victoria object to JSON string
func stringifyJSON(obj object.Object, indent string) object.Object {
	goVal := objectToGo(obj)
	var jsonBytes []byte
	var err error
	if indent != "" {
		jsonBytes, err = json.MarshalIndent(goVal, "", indent)
	} else {
		jsonBytes, err = json.Marshal(goVal)
	}
	if err != nil {
		return newError("failed to stringify JSON: %s", err.Error())
	}
	return &object.String{Value: string(jsonBytes)}
}

// objectToGo converts a Victoria object to Go interface{}
func objectToGo(obj object.Object) interface{} {
	switch o := obj.(type) {
	case *object.Integer:
		return o.Value
	case *object.Float:
		return o.Value
	case *object.Boolean:
		return o.Value
	case *object.String:
		return o.Value
	case *object.Null:
		return nil
	case *object.Array:
		result := make([]interface{}, len(o.Elements))
		for i, elem := range o.Elements {
			result[i] = objectToGo(elem)
		}
		return result
	case *object.Hash:
		result := make(map[string]interface{})
		for _, pair := range o.Pairs {
			key := pair.Key.Inspect()
			result[key] = objectToGo(pair.Value)
		}
		return result
	default:
		return obj.Inspect()
	}
}

// convertTimeFormat converts common format tokens to Go's time format
func convertTimeFormat(format string) string {
	replacements := map[string]string{
		"YYYY": "2006",
		"YY":   "06",
		"MM":   "01",
		"DD":   "02",
		"HH":   "15",
		"hh":   "03",
		"mm":   "04",
		"ss":   "05",
		"SSS":  "000",
		"Z":    "-0700",
		"A":    "PM",
		"a":    "pm",
	}
	result := format
	for token, goFormat := range replacements {
		result = strings.Replace(result, token, goFormat, -1)
	}
	return result
}

func findModuleFile(name string) (string, bool) {
	// 1. Current directory
	path := name
	if !strings.HasSuffix(path, ".vc") {
		path += ".vc"
	}
	if _, err := os.Stat(path); err == nil {
		return path, true
	}

	// 2. victoria_modules directory
	// Check if it's a package name (e.g. "mylib")
	pkgPath := filepath.Join("victoria_modules", name)
	// Check for victoria_modules/name.vc
	if _, err := os.Stat(pkgPath + ".vc"); err == nil {
		return pkgPath + ".vc", true
	}
	// Check for victoria_modules/name/index.vc or victoria_modules/name/main.vc
	if _, err := os.Stat(filepath.Join(pkgPath, "index.vc")); err == nil {
		return filepath.Join(pkgPath, "index.vc"), true
	}
	if _, err := os.Stat(filepath.Join(pkgPath, "main.vc")); err == nil {
		return filepath.Join(pkgPath, "main.vc"), true
	}

	// 3. Global modules
	home, _ := os.UserHomeDir()
	if home != "" {
		globalPath := filepath.Join(home, ".victoria", "modules", name)
		if _, err := os.Stat(globalPath + ".vc"); err == nil {
			return globalPath + ".vc", true
		}
		if _, err := os.Stat(filepath.Join(globalPath, "index.vc")); err == nil {
			return filepath.Join(globalPath, "index.vc"), true
		}
	}

	return "", false
}

// evalIncludeStatement handles include statements for modules and files
func evalIncludeStatement(node *ast.IncludeStatement, env *object.Environment) object.Object {
	for _, moduleName := range node.Modules {
		if factory, ok := moduleRegistry[moduleName]; ok {
			module := factory()
			env.Set(moduleName, module)
		} else {
			// Try to load as file
			filename, found := findModuleFile(moduleName)
			if !found {
				return newError("module or file not found: %s", moduleName)
			}

			content, err := os.ReadFile(filename)
			if err != nil {
				return newError("could not read file: %s", err.Error())
			}

			l := lexer.New(string(content))
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				msg := fmt.Sprintf("parser errors in %s:\n", filename)
				for _, msgErr := range p.Errors() {
					msg += "\t" + msgErr + "\n"
				}
				return newError(msg)
			}

			// Create a new environment for the module to isolate it
			moduleEnv := object.NewEnvironment()

			result := Eval(program, moduleEnv)
			if isError(result) {
				return result
			}

			// Export the module's environment as a Hash
			moduleObj := moduleEnv.ToHash()

			// Use the base name of the module as the variable name
			baseName := filepath.Base(moduleName)
			baseName = strings.TrimSuffix(baseName, ".vc")

			env.Set(baseName, moduleObj)
		}
	}
	return NULL
}

// Blank identifier to use filepath package
var _ = filepath.Base

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		initProject()
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Usage: vic install <github-user>/<repo>")
			return
		}
		installPackage(os.Args[2])
	case "list":
		listPackages()
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
	}
}

func printHelp() {
	fmt.Println("Victoria Package Manager (vic)")
	fmt.Println("\nUsage:")
	fmt.Println("  vic init              Initialize a new Victoria project")
	fmt.Println("  vic install <pkg>     Install a package from GitHub (e.g., user/repo)")
	fmt.Println("  vic list              List installed packages")
	fmt.Println("  vic help              Show this help message")
}

func initProject() {
	if _, err := os.Stat("vic.json"); err == nil {
		fmt.Println("vic.json already exists")
		return
	}

	wd, _ := os.Getwd()
	name := filepath.Base(wd)

	config := Config{
		Name:         name,
		Version:      "1.0.0",
		Dependencies: make(map[string]string),
	}

	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile("vic.json", data, 0644)
	fmt.Println("Initialized vic.json")
}

func installPackage(pkg string) {
	// pkg is expected to be "user/repo"
	parts := strings.Split(pkg, "/")
	if len(parts) != 2 {
		fmt.Println("Invalid package format. Use user/repo")
		return
	}

	repoName := parts[1]
	targetDir := filepath.Join("victoria_modules", repoName)

	// Create victoria_modules if it doesn't exist
	os.MkdirAll("victoria_modules", 0755)

	fmt.Printf("Installing %s...\n", pkg)

	url := fmt.Sprintf("https://github.com/%s.git", pkg)
	cmd := exec.Command("git", "clone", "--depth", "1", url, targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "already exists") {
			fmt.Printf("Package %s already installed. Updating...\n", pkg)
			updateCmd := exec.Command("git", "-C", targetDir, "pull")
			updateCmd.Run()
		} else {
			fmt.Printf("Failed to install package: %s\n%s", err, string(output))
			return
		}
	}

	// Update vic.json
	updateConfig(pkg)
	fmt.Printf("Successfully installed %s\n", pkg)
}

func updateConfig(pkg string) {
	data, err := os.ReadFile("vic.json")
	if err != nil {
		return
	}

	var config Config
	json.Unmarshal(data, &config)

	if config.Dependencies == nil {
		config.Dependencies = make(map[string]string)
	}

	config.Dependencies[pkg] = "latest"

	newData, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile("vic.json", newData, 0644)
}

func listPackages() {
	if _, err := os.Stat("victoria_modules"); os.IsNotExist(err) {
		fmt.Println("No packages installed")
		return
	}

	entries, _ := os.ReadDir("victoria_modules")
	fmt.Println("Installed packages:")
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("  - %s\n", entry.Name())
		}
	}
}

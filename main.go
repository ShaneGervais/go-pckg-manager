package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		install_name := os.Args[2]
		if install_name == "" {
			fmt.Println("Please provide a package name to install")
			os.Exit(1)
		}
		install_package(install_name)
	case "remove":
		remove_name := os.Args[2]
		if remove_name == "" {
			fmt.Println("Please provide a package name to remove")
			os.Exit(1)
		}
		remove_package(remove_name)
	case "update":
		update_name := os.Args[2]
		if update_name == "" {
			fmt.Println("Please provide a package name to update")
			os.Exit(1)
		}
		update_package(update_name)
	case "search":
		search_name := os.Args[2]
		if search_name == "" {
			fmt.Println("Please provide a package to search for")
			os.Exit(1)
		}
		search_package(search_name)
	case "list":
		list_packages()
	default:
		fmt.Println("Invalid command.")
		os.Exit(1)
	}
}

func parseDependencies(output string) []string {
	dependencies := []string{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Dependency resolved:") {
			parts := strings.Split(line, " ")
			if len(parts) >= 4 {
				dependency := strings.TrimSpace(parts[3])
				dependencies = append(dependencies, dependency)
			}
		}
	}

	return dependencies
}

func install_package(name string) error {
	path_checker()

	fmt.Println("Installing:", name)

	install := exec.Command("dnf", "install", "--assumeno")
	output, err := install.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to check package dependencies: %v", err)
	}

	dependencies := parseDependencies(string(output))

	for _, i := range dependencies {
		err := install_package(i)
		if err != nil {
			return fmt.Errorf("failed to install dependency %s: %v", i, err)
		}
	}

	install = exec.Command("dnf", "install", name)
	install.Stdout = os.Stdout
	install.Stdin = os.Stdin

	err_run := install.Run()
	if err_run != nil {
		fmt.Printf("Error installing the package %s. Error: %v\n", name, err_run)
		os.Exit(1)
	} else {
		fmt.Println("Package", name, "has been successfully installed.")
	}
	return nil
}

func remove_package(name string) {
	path_checker()

	fmt.Println("Removing:", name)
	remove := exec.Command("dnf", "remove", name)
	remove.Stdout = os.Stdout
	remove.Stdin = os.Stdin

	err_run := remove.Run()
	if err_run != nil {
		fmt.Println("Error removing the package", name, "Error:", err_run)
	} else {
		fmt.Println("Package", name, "has been successfully removed")
	}
}

func update_package(name string) {
	path_checker()

	fmt.Println("Updating:", name)
	update := exec.Command("dnf", "update", name)
	update.Stdin = os.Stdin
	update.Stdout = os.Stdout

	err_run := update.Run()
	if err_run != nil {
		fmt.Println("Error updating package", name, "Error:", err_run)
	} else {
		fmt.Println("Package", name, "successfully updated")
	}
}

func list_packages() {
	path_checker()
	list := exec.Command("dnf", "list", "installed")
	list.Stdin = os.Stdin
	list.Stdout = os.Stdout

	err_run := list.Run()
	if err_run != nil {
		fmt.Println("Error listing installed packages. Error:", err_run)
	} else {
		fmt.Println("Listing installed packages")
	}
}

func search_package(name string) {
	path_checker()

	fmt.Println("Searching for:", name)

	search := exec.Command("dnf", "search", name)
	search.Stdin = os.Stdin
	search.Stdout = os.Stdout

	err_run := search.Run()
	if err_run != nil {
		fmt.Println("Error searching packages", name, "in installed packages. Error:", err_run)
	} else {
		fmt.Println("Found", name)
	}

}

func path_checker() {
	err_path := os.Setenv("PATH", "/usr/bin:/usr/sbin:/bin:/sbin")
	if err_path != nil {
		fmt.Println("Failed to set PATH environment variable:", err_path)
		os.Exit(1)
	}
}

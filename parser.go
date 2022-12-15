// THIS IS MY INITIAL PARSER FILE THE REAL IMPLEMENTATION WILL HAPPEN IN SYFT ITSELF
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Dependency represents a single dependency in the build.gradle file
type Dependency struct {
	Group   string
	Name    string
	Version string
}

// Plugin represents a single plugin in the build.gradle file
type Plugin struct {
	Id      string
	Version string
}

func main() {
	// Open the build.gradle file
	file, err := os.Open("build.gradle.kts")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Create slices to hold the dependencies and plugins
	dependencies := []Dependency{}
	plugins := []Plugin{}
	// Create a map to hold the variables
	variables := map[string]string{}

	// Keep track of whether we are in the dependencies or plugins section
	inDependenciesSection := false
	inPluginsSection := false

	// Loop over all lines in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Trim leading and trailing whitespace from the line
		line = strings.TrimSpace(line)

		// Check if the line starts with "dependencies {"
		if strings.HasPrefix(line, "dependencies {") {
			inDependenciesSection = true
			continue
		}

		// Check if the line starts with "plugins {"
		if strings.HasPrefix(line, "plugins {") {
			inPluginsSection = true
			continue
		}

		// Check if the line is "}"
		if line == "}" {
			inDependenciesSection = false
			inPluginsSection = false
			continue
		}

		// Check if we are in the plugins section
		if inPluginsSection {
			// Split the line on whitespace to extract the group, name, and version of the dependency
			fields := strings.Fields(line)
			// Check if the line contains at least 3 fields (group, version as a literal string, and version as the version number)
			if len(fields) >= 3 {
				start := strings.Index(fields[0], "(") + 1
				end := strings.Index(fields[0], ")")
				groupName := fields[0][start:end]
				groupName = strings.Trim(groupName, `"`)
				version := strings.Trim(fields[2], `"`)
				// Create a new Dependency struct and add it to the dependencies slice
				plugin := Plugin{Id: groupName, Version: version}
				plugins = append(plugins, plugin)
			}
		}

		// Check if we are in the dependencies section
		if inDependenciesSection {
			// Extract the group, name, and version from the function call
			start := strings.IndexFunc(line, func(r rune) bool {
				return r == '(' || r == ' '
			}) + 1
			// Extract the group, name, and version from the function call
			end := strings.IndexFunc(line, func(r rune) bool {
				return r == ')' || r == ' '
			})
			groupNameVersion := line[start:end]
			groupNameVersion = strings.Trim(groupNameVersion, "\"")
			parts := strings.Split(groupNameVersion, ":")
			// if we only have 2 sections the version is probably missing
			if len(parts) == 2 {
				// search for the version in the plugin section
				version := searchInPlugins(parts[0], plugins)
				// Create a new Dependency struct and add it to the dependencies slice
				dep := Dependency{Group: parts[0], Name: parts[1], Version: version}
				dependencies = append(dependencies, dep)
			}
			// we have a version directly specified
			if len(parts) == 3 {
				// Create a new Dependency struct and add it to the dependencies slice
				dep := Dependency{Group: parts[0], Name: parts[1], Version: parts[2]}
				dependencies = append(dependencies, dep)
			}
		}

		// Check if the line contains an assignment
		if strings.Contains(line, "=") {
			// Split the line on the "=" character to separate the key and value
			parts := strings.Split(line, "=")

			// Trim any leading and trailing whitespace from the key and value
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Add the key and value to the map
			variables[key] = value
		}

	}
	// Print the dependencies
	fmt.Println("Dependencies:")
	for _, dep := range dependencies {
		fmt.Printf("%s\n %s\n %s\n", dep.Group, dep.Name, dep.Version)
	}
	// Print the plugins
	fmt.Println("\nPlugins:")
	for _, plugin := range plugins {
		fmt.Printf("  %s:%s\n", plugin.Id, plugin.Version)
	}
	// Print the variables
	fmt.Println("Variables:")
	for key, value := range variables {
		fmt.Printf("%s = %s\n", key, value)
	}

}

func searchInPlugins(groupName string, plugins []Plugin) string {
	for _, v := range plugins {
		if v.Id == groupName {
			return v.Version
		}
	}
	return ""
}

// var regex = regexp.MustCompile(`\$\{([^\}]+)\}`)
// variable := regex.FindStringSubmatch(line)[1]

// // Use the strings.Split function to split the line on the "=" character and extract the value of the variable
// parts := strings.Split(line, "=")
// version := parts[len(parts)-1]

// // Trim any whitespace from the variable name and version
// variable = strings.TrimSpace(variable)
// version = strings.TrimSpace(version)

// // Replace the variable name with its value in the line
// line = strings.Replace(line, "${"+variable+"}", version, -1)

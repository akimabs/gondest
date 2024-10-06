package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var moduleName string

var generateCmd = &cobra.Command{
	Use:   "generate [type] [name]",
	Short: "Generate a new controller, service, or module",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		typ := args[0]
		name := args[1]

		// Automatically generate controller, service, and module if type is "module"
		if typ == "module" {
			moduleName = name // Set the module name for all files
			createDirectoryStructure(moduleName)
			createFileFromTemplate(name, "domains/"+moduleName+"/"+name+".controller.go", "controller.go.tpl")
			createFileFromTemplate(name, "domains/"+moduleName+"/"+name+".service.go", "service.go.tpl")
			createFileFromTemplate(name, "domains/"+moduleName+"/"+name+".module.go", "module.go.tpl")
			fmt.Printf("Module %s created with controller, service, and module files in domains.\n", name)
		} else {
			fmt.Printf("Key %s that's not included in the command \n", name)
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "init [app-name]",
	Short: "Initialize a new GoFiber app with gondest structure",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]

		// Create project directory
		err := os.MkdirAll(appName, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating project directory: %v\n", err)
			return
		}

		// Change directory to the newly created project
		err = os.Chdir(appName)
		if err != nil {
			fmt.Printf("Error changing to project directory: %v\n", err)
			return
		}

		// Initialize go.mod with GoFiber
		initGoMod(appName)

		createDirectoryStructure("domains")
		createFileFromTemplate(appName, "domains/app.controller.go", "default/controller.go.tpl")
		createFileFromTemplate(appName, "domains/app.service.go", "default/service.go.tpl")
		createFileFromTemplate(appName, "domains/app.module.go", "default/module.go.tpl")

		// Create main.go
		createFileFromTemplate(appName, "main.go", "default/main.go.tpl")

		fmt.Printf("Project %s initialized with GoFiber and default structure.\n", appName)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(initCmd)
	generateCmd.Flags().StringVarP(&moduleName, "module", "m", "", "Module name")
}

func initGoMod(appName string) {
	cmd := exec.Command("go", "mod", "init", appName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error initializing go.mod: %v\n", err)
		return
	}

	// // Add GoFiber dependency
	// cmd = exec.Command("go", "get", "github.com/gofiber/fiber/v2")
	// err = cmd.Run()
	// if err != nil {
	// 	fmt.Printf("Error adding GoFiber dependency: %v\n", err)
	// 	return
	// }
	// fmt.Println("go.mod initialized with GoFiber.")
}

// createFileFromTemplate creates a file from a template
func createFileFromTemplate(name, filename, tplName string) {
	// Get the directory of the executable
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	// Create the full path to the templates directory
	templatesDir := filepath.Join(filepath.Dir(execPath), "templates")

	// Parse the template file
	tplPath := filepath.Join(templatesDir, tplName)
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	data := map[string]string{
		"ModuleName":     cases.Title(language.English).String(moduleName),
		"ControllerName": cases.Title(language.English).String(name) + "Controller",
		"ServiceName":    cases.Title(language.English).String(name) + "Service",
	}

	if err := tpl.Execute(f, data); err != nil {
		fmt.Println("Error executing template:", err)
	}
}

// createDirectoryStructure creates the necessary directory structure for the module
func createDirectoryStructure(moduleName string) {
	// Create the base domains directory if it doesn't exist
	if err := os.MkdirAll(moduleName, os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
}

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var moduleName string
var templatePath = "/usr/local/share/gondest/templates"

var generateCmd = &cobra.Command{
	Use:   "generate [type] [name]",
	Short: "Generate a new controller, service, or module", // Deskripsi singkat
	Long: `Generate will create a new controller, service, or module based on the type.
For 'module', it automatically generates a controller, service, and module file.`,
	Example: `gondest generate module user 
gondest generate controller post`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		typ := args[0]
		name := args[1]

		// Automatically generate controller, service, and module if type is "module"
		if typ == "module" {
			moduleName = name // Set the module name for all files
			createFileFromTemplate(name, "domains/"+moduleName+"/"+name+".controller.go", "controller.go.tpl")
			createFileFromTemplate(name, "domains/"+moduleName+"/"+name+".service.go", "service.go.tpl")
			createFileFromTemplate(name, "domains/"+moduleName+"/"+name+".module.go", "module.go.tpl")
			updateMainGo(moduleName)
			fmt.Printf("Module %s created with controller, service, and module files in domains.\n", name)
		} else {
			fmt.Printf("Key %s that's not included in the command \n", name)
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "init [app-name]",
	Short: "Initialize a new GoFiber app with gondest structure",
	Long: `This command initializes a new GoFiber project with the standard structure.
It creates the base domains, services, and controller structure along with a main.go file.`,
	Example: `gondest init myApp`,
	Args:    cobra.ExactArgs(1),
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

		installDependency()

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

}

func installDependency() {
	// // Add GoFiber dependency
	cmd := exec.Command("go", "mod", "tidy")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error adding GoFiber dependency: %v\n", err)
		return
	}
	fmt.Println("go.mod initialized with GoFiber.")

}

// createFileFromTemplate creates a file from a template
func createFileFromTemplate(name, filename, tplName string) {
	// Ensure the directory structure exists
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, os.ModePerm) // Ensure the directory is created
	if err != nil {
		fmt.Println("Error creating directory structure:", err)
		return
	}

	// Parse the template file
	tplPath := filepath.Join(templatePath, tplName)
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
		"AppName":        name, // Pass appName to the template
		"ModuleName":     cases.Title(language.English).String(moduleName),
		"ControllerName": cases.Title(language.English).String(name) + "Controller",
		"ServiceName":    cases.Title(language.English).String(name) + "Service",
	}

	if err := tpl.Execute(f, data); err != nil {
		fmt.Println("Error executing template:", err)
	}
}

// updateMainGo updates main.go to add the new module import, module to fxApp, and controller to fx.Invoke
func updateMainGo(moduleName string) {
	// Convert the module name to a capitalized version for the import alias (e.g., Auth)
	moduleAlias := cases.Title(language.English).String(moduleName)

	// Define the module import and fxApp entry
	moduleImport := fmt.Sprintf(`%s "app_gondest/domains/%s"`, moduleAlias, moduleName)
	moduleEntry := fmt.Sprintf(`%s.Module,`, moduleAlias)
	controllerEntry := fmt.Sprintf("%sController *%s.%sController", moduleName, moduleAlias, moduleAlias)
	registerRouteLine := fmt.Sprintf("\t\t\t%sController.RegisterRoutes(app)", moduleName) // Ensure correct tabbing

	// Read the main.go file
	mainFile := "main.go"
	content, err := os.ReadFile(mainFile)
	if err != nil {
		fmt.Printf("Error reading main.go: %v\n", err)
		return
	}

	// Convert to string
	mainContent := string(content)

	// Check if the module is already added
	if strings.Contains(mainContent, moduleImport) {
		fmt.Printf("Module %s already exists in main.go\n", moduleName)
		return
	}

	// Add the import for the new module
	importSection := `import (`
	importIndex := strings.Index(mainContent, importSection)
	if importIndex != -1 {
		importIndex += len(importSection)
		mainContent = mainContent[:importIndex] + "\n\t" + moduleImport + mainContent[importIndex:]
	}

	// Add the module entry to fxApp
	fxNewSection := `fx.New(`
	fxNewIndex := strings.Index(mainContent, fxNewSection)
	if fxNewIndex != -1 {
		fxNewIndex += len(fxNewSection)
		mainContent = mainContent[:fxNewIndex] + "\n\t\t" + moduleEntry + mainContent[fxNewIndex:]
	}

	// Update the fx.Invoke section with the new controller and register route
	fxInvokeSection := `fx.Invoke(`
	fxInvokeIndex := strings.Index(mainContent, fxInvokeSection)
	if fxInvokeIndex != -1 {
		// Locate the function signature
		invokeEndIndex := strings.Index(mainContent[fxInvokeIndex:], `})`)
		invokeBlock := mainContent[fxInvokeIndex : fxInvokeIndex+invokeEndIndex]

		// Check if the controller is already included
		if !strings.Contains(invokeBlock, controllerEntry) {
			// Find the opening of the function signature
			funcIndex := strings.Index(invokeBlock, `func(`)
			if funcIndex != -1 {
				funcEndIndex := strings.Index(invokeBlock[funcIndex:], `)`)
				invokeBlock = invokeBlock[:funcIndex+funcEndIndex] + ", " + controllerEntry + invokeBlock[funcIndex+funcEndIndex:]
				mainContent = mainContent[:fxInvokeIndex] + invokeBlock + mainContent[fxInvokeIndex+invokeEndIndex:]
			}
		}

		// Add the RegisterRoutes call inside the fx.Invoke block
		if !strings.Contains(mainContent, registerRouteLine) {
			// Add after the last RegisterRoutes(app) call
			registerRoutesIndex := strings.LastIndex(mainContent[:fxInvokeIndex+invokeEndIndex], "RegisterRoutes(app)")
			if registerRoutesIndex != -1 {
				registerRoutesEnd := strings.Index(mainContent[registerRoutesIndex:], "\n")
				mainContent = mainContent[:registerRoutesIndex+registerRoutesEnd] + "\n" + registerRouteLine + mainContent[registerRoutesIndex+registerRoutesEnd:]
			}
		}
	}

	// Write the modified content back to main.go
	err = os.WriteFile(mainFile, []byte(mainContent), os.ModePerm)
	if err != nil {
		fmt.Printf("Error writing to main.go: %v\n", err)
		return
	}

	fmt.Printf("Module %s added to main.go\n", moduleName)
}

// createDirectoryStructure creates the necessary directory structure for the module
func createDirectoryStructure(moduleName string) {
	// Create the base domains directory if it doesn't exist
	if err := os.MkdirAll(moduleName, os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
}

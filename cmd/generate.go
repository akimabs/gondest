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
var dbType string
var templatePath = "/usr/local/share/gondest/templates"

var generateCmd = &cobra.Command{
	Use:   "generate [type] [name]",
	Short: "Generate a new controller, service, or module",
	Long: `Generate will create a new controller, service, or module based on the type.
For 'module', it automatically generates a controller, service, and module file.`,
	Example: `gondest generate module user 
gondest generate controller post`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		typ := args[0]
		name := args[1]

		// Automatically generate controller, service, and module if type is "module"
		switch typ {
		case "module":
			moduleName = name // Set the module name for all files
			createFileFromTemplate(name, "domains/"+moduleName+"/"+name+".controller.go", "module/controller.go.tpl", "")
			createFileFromTemplate(name, "domains/"+moduleName+"/"+name+".service.go", "module/service.go.tpl", "")
			createFileFromTemplate(name, "domains/"+moduleName+"/"+name+".module.go", "module/module.go.tpl", "")
			updateMainGo(moduleName)
			fmt.Printf("Module %s created with controller, service, model, validation, and module files in domains.\n", name)

		case "model":
			createFileFromTemplate(name, "models/"+name+".model.go", "model/model.go.tpl", "")
			fmt.Printf("Model %s created in models.\n", name)

		default:
			fmt.Printf("Key %s is not recognized.\n", typ)
		}

	},
}

var configCmd = &cobra.Command{
	Use:     "config [type]",
	Short:   "Config a new integration ex: db(mysql, postgres, sqlite, sqlserver), redis, etc",
	Long:    `Config will create a new db, redis, etc based on the type.`,
	Example: `gondest config db --postgres`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		typ := args[0]
		name := "default"

		// Check if the specific db flag is set and set the dbType accordingly
		if isMysql, _ := cmd.Flags().GetBool("mysql"); isMysql {
			dbType = "mysql"
		}
		if isPostgres, _ := cmd.Flags().GetBool("postgres"); isPostgres {
			dbType = "postgres"
		}
		if isSqlserver, _ := cmd.Flags().GetBool("sqlserver"); isSqlserver {
			dbType = "sqlserver"
		}

		// Check if dbType is valid
		if dbType == "" {
			fmt.Printf("Error: Unsupported or missing db type. Allowed flags are: --mysql, --postgres, --sqlite, --sqlserver.\n")
			return
		}

		// Switch based on type
		switch typ {
		case "db":
			createFileFromTemplate(name, "config/config.db.go", "config/db.go.tpl", dbType)
			createFileFromTemplate(name, "config/module.db.go", "config/module.go.tpl", dbType)
			updateENV()
			updateMainGoConfig("config")
			installDependency()
			fmt.Printf("Config created for %s database.\n", dbType)

		default:
			fmt.Printf("Type %s is not recognized.\n", typ)
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
		createFileFromTemplate(appName, "domains/app.controller.go", "default/controller.go.tpl", "")
		createFileFromTemplate(appName, "domains/app.service.go", "default/service.go.tpl", "")
		createFileFromTemplate(appName, "domains/app.module.go", "default/module.go.tpl", "")
		createFileFromTemplate(appName, "utils/response.go", "default/response.go.tpl", "")
		createFileFromTemplate(appName, "main.go", "default/main.go.tpl", "")
		createFileFromTemplate(appName, ".env", "default/.env.tpl", "")
		createFileFromTemplate(appName, ".env.example", "default/.env.example.tpl", "")
		createFileFromTemplate(appName, ".air.toml", "default/.air.toml.tpl", "")

		installDependency()

		fmt.Printf("Project %s initialized with GoFiber and default structure.\n", appName)
		fmt.Printf("Init your database value at environment variables")

	},
}

func init() {
	configCmd.Flags().Bool("mysql", false, "Use MySQL as the database")
	configCmd.Flags().Bool("postgres", false, "Use PostgreSQL as the database")
	configCmd.Flags().Bool("sqlserver", false, "Use SQL Server as the database")
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(initCmd)

}

func getGoModModuleName() string {
	goModFile := "go.mod"
	content, err := os.ReadFile(goModFile)
	if err != nil {
		fmt.Printf("Error reading go.mod: %v\n", err)
		return ""
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	fmt.Println("Module name not found in go.mod")
	return ""
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
func createFileFromTemplate(name, filename, tplName string, dbType string) {
	appName := getGoModModuleName() // Get the module name dynamically

	if dbType == "" {
		dbType = "postgres" // Default to postgres if not specified
	}

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
		"AppName":        appName,
		"ModelName":      cases.Title(language.English).String(name) + "Model",
		"ModuleName":     cases.Title(language.English).String(moduleName),
		"ControllerName": cases.Title(language.English).String(name) + "Controller",
		"ServiceName":    cases.Title(language.English).String(name) + "Service",
		"DatabaseDriver": dbType,
	}

	if err := tpl.Execute(f, data); err != nil {
		fmt.Println("Error executing template:", err)
	}
}

func updateENV() {
	envFiles := []string{".env", ".env.example"}

	for _, envFile := range envFiles {
		content, err := os.ReadFile(envFile)
		if err != nil {
			fmt.Printf("Error reading %s file: %v\n", envFile, err)
			return
		}

		envContent := string(content)

		// Check if DB environment variables are already present
		envVars := []string{
			"\n",
			"DB_HOST=",
			"DB_USER=",
			"DB_PASSWORD=",
			"DB_NAME=",
			"DB_PORT=",
		}

		updated := false
		insertIndex := strings.Index(envContent, "PORT=3000")
		if insertIndex != -1 {
			insertIndex += len("PORT=3000") // Move to the end of the PORT=3000 line
		}

		for _, envVar := range envVars {
			if !strings.Contains(envContent, envVar) {
				if insertIndex != -1 {
					// Insert new variables below PORT=3000
					envContent = envContent[:insertIndex] + envVar + "\n" + envContent[insertIndex:]
					insertIndex += len(envVar) + 1 // Adjust the index for the next variable
				} else {
					// If PORT=3000 not found, append at the end
					envContent += envVar + "\n"
				}
				updated = true
			}
		}

		// If any new environment variables were added, write back to the .env or .env.example file
		if updated {
			err = os.WriteFile(envFile, []byte(envContent), os.ModePerm)
			if err != nil {
				fmt.Printf("Error writing to %s file: %v\n", envFile, err)
				return
			}
			fmt.Printf("Updated %s with missing database variables.\n", envFile)
		} else {
			fmt.Printf("%s is already up-to-date.\n", envFile)
		}
	}
}

// updateMainGoDb updates main.go to add the new module import, module to fxApp, and controller to fx.Invoke (case config)
func updateMainGoConfig(moduleName string) {
	// Convert the module name to a capitalized version for the import alias (e.g., Auth)
	appName := getGoModModuleName() // Get the module name dynamically
	moduleAlias := moduleName

	// Define the module import and fxApp entry
	moduleEntry := fmt.Sprintf(`%s.Module,`, moduleAlias)
	moduleImport := fmt.Sprintf(`%s "%s/%s"`, moduleAlias, appName, moduleName)

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

	// Write the modified content back to main.go
	err = os.WriteFile(mainFile, []byte(mainContent), os.ModePerm)
	if err != nil {
		fmt.Printf("Error writing to main.go: %v\n", err)
		return
	}

	fmt.Printf("Module %s added to main.go\n", moduleName)
}

func updateMainGo(moduleName string) {
	appName := getGoModModuleName() // Get the module name dynamically

	moduleAlias := cases.Title(language.English).String(moduleName)

	moduleImport := fmt.Sprintf(`%s "%s/domains/%s"`, moduleAlias, appName, moduleName)
	moduleEntry := fmt.Sprintf(`%s.Module,`, moduleAlias)
	controllerEntry := fmt.Sprintf("%sController *%s.%sController", moduleName, moduleAlias, moduleAlias)
	registerRouteLine := fmt.Sprintf("\t\t\t%sController.RegisterRoutes(app)", moduleName)

	mainFile := "main.go"
	content, err := os.ReadFile(mainFile)
	if err != nil {
		fmt.Printf("Error reading main.go: %v\n", err)
		return
	}

	mainContent := string(content)

	if strings.Contains(mainContent, moduleImport) {
		fmt.Printf("Module %s already exists in main.go\n", moduleName)
		return
	}

	importSection := `import (`
	importIndex := strings.Index(mainContent, importSection)
	if importIndex != -1 {
		importIndex += len(importSection)
		mainContent = mainContent[:importIndex] + "\n\t" + moduleImport + mainContent[importIndex:]
	}

	fxNewSection := `fx.New(`
	fxNewIndex := strings.Index(mainContent, fxNewSection)
	if fxNewIndex != -1 {
		fxNewIndex += len(fxNewSection)
		mainContent = mainContent[:fxNewIndex] + "\n\t\t" + moduleEntry + mainContent[fxNewIndex:]
	}

	fxInvokeSection := `fx.Invoke(`
	fxInvokeIndex := strings.Index(mainContent, fxInvokeSection)
	if fxInvokeIndex != -1 {
		invokeEndIndex := strings.Index(mainContent[fxInvokeIndex:], `})`)
		invokeBlock := mainContent[fxInvokeIndex : fxInvokeIndex+invokeEndIndex]

		if !strings.Contains(invokeBlock, controllerEntry) {
			funcIndex := strings.Index(invokeBlock, `func(`)
			if funcIndex != -1 {
				funcEndIndex := strings.Index(invokeBlock[funcIndex:], `)`)
				invokeBlock = invokeBlock[:funcIndex+funcEndIndex] + ", " + controllerEntry + invokeBlock[funcIndex+funcEndIndex:]
				mainContent = mainContent[:fxInvokeIndex] + invokeBlock + mainContent[fxInvokeIndex+invokeEndIndex:]
			}
		}

		if !strings.Contains(mainContent, registerRouteLine) {
			registerRoutesIndex := strings.LastIndex(mainContent[:fxInvokeIndex+invokeEndIndex], "RegisterRoutes(app)")
			if registerRoutesIndex != -1 {
				registerRoutesEnd := strings.Index(mainContent[registerRoutesIndex:], "\n")
				mainContent = mainContent[:registerRoutesIndex+registerRoutesEnd] + "\n" + registerRouteLine + mainContent[registerRoutesIndex+registerRoutesEnd:]
			}
		}
	}

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

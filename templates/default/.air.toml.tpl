# Config file for Air

# Root directory for the project (where the main.go file is located)
root = "."
# Start the command that will be run when you start Air
cmd = "go run main.go"

# List of directories to watch for file changes
include_dir = ["."]
exclude_dir = ["vendor", "node_modules"]

# Files with these extensions will trigger a reload
include_ext = ["go", "tmpl", "html", "css", "env"]

# The path to your .env file, so Air can reload when it changes
load_dotenv = true

# For verbose logging of events during file watching and reloading
debug = false

# Graceful shutdown delay time (in seconds)
graceful_delay = 2

# Color mode for output: "on", "off", "auto"
color = "auto"
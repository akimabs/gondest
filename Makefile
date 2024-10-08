BINARY_NAME=gondest
INSTALL_PATH=/usr/local/bin
TEMPLATE_PATH=/usr/local/share/gondest/templates
TEMPLATE_SRC=templates

install:
	@if [ ! -d "$(TEMPLATE_SRC)" ]; then \
		echo "Template directory '$(TEMPLATE_SRC)' not found!"; \
		exit 1; \
	fi
	go build -o $(BINARY_NAME)
	@echo "Using sudo to create directories and copy files to system paths."
	sudo mkdir -p $(TEMPLATE_PATH)
	sudo cp -r $(TEMPLATE_SRC)/* $(TEMPLATE_PATH)
	sudo mv $(BINARY_NAME) $(INSTALL_PATH)
	@echo "Installation complete! Templates are located at $(TEMPLATE_PATH)."
	@echo "enjoy"
	@echo " ▗▄▄▖ ▗▄▖ ▗▖  ▗▖▗▄▄▄ ▗▄▄▄▖ ▗▄▄▖▗▄▄▄▖"
	@echo "▐▌   ▐▌ ▐▌▐▛▚▖▐▌▐▌  █▐▌   ▐▌     █  "
	@echo "▐▌▝▜▌▐▌ ▐▌▐▌ ▝▜▌▐▌  █▐▛▀▀▘ ▝▀▚▖  █  "
	@echo "▝▚▄▞▘▝▚▄▞▘▐▌  ▐▌▐▙▄▄▀▐▙▄▄▖▗▄▄▞▘  █  "


uninstall:
	@echo "Using sudo to remove files from system paths."
	sudo rm -rf $(TEMPLATE_PATH)
	sudo rm $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstallation complete!"

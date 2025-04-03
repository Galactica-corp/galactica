#? help: Get more info on make commands.
help:
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@cat $(MAKEFILE_LIST) | sed -n 's/^#?//p' | column -t -s ':' | sort | sed -e 's/^/ /'
.PHONY: help

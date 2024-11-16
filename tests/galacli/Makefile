uname=$(shell uname -s)
is_darwin :=$(filter Darwin,$(uname))
all: ## run tests
all: $(if $(is_darwin),network,) clean start_test

help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[[:alpha:][:punct:]]+:.*?## / {printf "  \033[32m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

clean:
	rm -rf node0* .pytest_cache __pycache__

.venv/bin/pytest .venv/bin/python:
	uv venv
	uv sync

start: ## just run galacli.py for debugging purpose
start: .venv/bin/python
	.venv/bin/python galacli.py

start_test: .venv/bin/pytest ## start tests
	.venv/bin/pytest -v .

ps: ## show running galacticad processes
	-@ps | grep galacticad | grep -v 'grep'

terminate: ps ## terminate running galacticad processes
	ps | grep galacticad | grep -v 'grep' | awk '{print $$1}'  | xargs -n1 kill -15

kill: ps ## kill running galacticad processes
	ps | grep galacticad | grep -v 'grep' | awk '{print $$1}'  | xargs -n1 sudo kill -9

ifeq ($(uname),Darwin)
network: /tmp/127.0.0.2 /tmp/127.0.0.3 /tmp/127.0.0.4 ## assume you are using macos for development
/tmp/127.0.0.%:
	sudo ifconfig lo0 alias $(shell basename $@) up && touch $@
endif

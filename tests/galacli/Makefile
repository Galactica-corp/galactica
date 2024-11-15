uname=$(shell uname -s)
is_darwin :=$(filter Darwin,$(uname))

all: $(if $(is_darwin),network,)  clean start

clean:
	rm -rf node0* .pytest_cache __pycache__

.venv/bin/python:
	uv venv
	uv sync

start: .venv/bin/python
	pytest -v .
# .venv/bin/python galacli.py
ps:
	-@ps | grep galacticad | grep -v 'grep'

terminate: ps
	ps | grep galacticad | grep -v 'grep' | awk '{print $$1}'  | xargs -n1 kill -15
kill: ps
	ps | grep galacticad | grep -v 'grep' | awk '{print $$1}'  | xargs -n1 sudo kill -9

ifeq ($(uname),Darwin)
network: /tmp/127.0.0.2 /tmp/127.0.0.3 /tmp/127.0.0.4
/tmp/127.0.0.%:
	sudo ifconfig lo0 alias $(shell basename $@) up && touch $@
endif

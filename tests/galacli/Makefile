uname=$(shell uname -s)
is_darwin :=$(filter Darwin,$(uname))

all: $(if $(is_darwin),network,) clean start

clean:
	rm -rf node0*

.venv/bin/python:
	uv venv
	uv sync

start: .venv/bin/python
	.venv/bin/python galacli.py

ifeq ($(uname),Darwin)
network: /tmp/127.0.0.2 /tmp/127.0.0.3 /tmp/127.0.0.4
/tmp/127.0.0.%:
	sudo ifconfig lo0 alias $(shell basename $@) up && touch $@
endif

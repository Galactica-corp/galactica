DOCKER := $(shell which docker)
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)

#
#  dish / Makefile
#

#
#  VARS
#

include .env.example
-include .env

PROJECT_NAME?=${APP_NAME}

DOCKER_DEV_IMAGE?=${PROJECT_NAME}-image
DOCKER_DEV_CONTAINER?=${PROJECT_NAME}-run

COMPOSE_FILE=deployments/docker-compose.yml

# define standard colors
# https://gist.github.com/rsperl/d2dfe88a520968fbc1f49db0a29345b9
ifneq (,$(findstring xterm,${TERM}))
	BLACK        := $(shell tput -Txterm setaf 0)
	RED          := $(shell tput -Txterm setaf 1)
	GREEN        := $(shell tput -Txterm setaf 2)
	YELLOW       := $(shell tput -Txterm setaf 3)
	LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
	PURPLE       := $(shell tput -Txterm setaf 5)
	BLUE         := $(shell tput -Txterm setaf 6)
	WHITE        := $(shell tput -Txterm setaf 7)
	RESET        := $(shell tput -Txterm sgr0)
else
	BLACK        := ""
	RED          := ""
	GREEN        := ""
	YELLOW       := ""
	LIGHTPURPLE  := ""
	PURPLE       := ""
	BLUE         := ""
	WHITE        := ""
	RESET        := ""
endif

export

#
#  FUNCTIONS
# 

define print_info
	@echo -e "\n>>> ${YELLOW}${1}${RESET}\n"
endef

define update_semver
	$(call print_info, Incrementing semver to ${1}...)
	@[ -f ".env" ] || cp .env.example .env
	@sed -i 's|APP_VERSION=.*|APP_VERSION=${1}|' .env
	@sed -i 's|APP_VERSION=.*|APP_VERSION=${1}|' .env.example
endef

#
#  TARGETS
#

all: info

.PHONY: build local_build logs major minor patch push run stop test version
info: 
	@echo -e "\n${GREEN} ${PROJECT_NAME} / Makefile ${RESET}\n"

	@echo -e "${YELLOW} make fmt     --- reformat the go source (gofmt) ${RESET}"
	@echo -e "${YELLOW} make test    --- run unit tests (go test) ${RESET}"
	@echo -e "${YELLOW} make build   --- build project (docker image) ${RESET}"
	@echo -e "${YELLOW} make run     --- run project ${RESET}"
	@echo -e "${YELLOW} make logs    --- fetch container's logs ${RESET}"
	@echo -e "${YELLOW} make stop    --- stop and purge project (only docker containers!) ${RESET}\n"

build:  
	@echo -e "\n${YELLOW} Building project (docker-compose build)... ${RESET}\n"
	@docker compose -f ${COMPOSE_FILE} build 

local_build: 
	@echo -e "\n${YELLOW} [local] Building project... ${RESET}\n"
	@go mod tidy
	@go build -tags dev -o dish cmd/dish/main.go

run:	build
	@echo -e "\n${YELLOW} Starting project (docker-compose up)... ${RESET}\n"
	@docker compose -f ${COMPOSE_FILE} up --force-recreate

logs:
	@echo -e "\n${YELLOW} Fetching container's logs (CTRL-C to exit)... ${RESET}\n"
	@docker logs ${DOCKER_DEV_CONTAINER} -f

stop:  
	@echo -e "\n${YELLOW} Stopping and purging project (docker-compose down)... ${RESET}\n"
	@docker compose -f ${COMPOSE_FILE} down

test:
	@echo -e "\n${YELLOW} [local] Running unit tests (go test)... ${RESET}\n"
	@go test ./...

push: 
	@git tag -fa 'v${APP_VERSION}' -m 'v${APP_VERSION}'
	@git push --follow-tags --set-upstream origin master


MAJOR := $(shell echo ${APP_VERSION} | cut -d. -f1)
MINOR := $(shell echo ${APP_VERSION} | cut -d. -f2)
PATCH := $(shell echo ${APP_VERSION} | cut -d. -f3)

major:
	$(eval APP_VERSION := $(shell echo $$(( ${MAJOR} + 1 )).0.0))
	$(call update_semver,${APP_VERSION})

minor:
	$(eval APP_VERSION := $(shell echo ${MAJOR}.$$(( ${MINOR} + 1 )).0))
	$(call update_semver,${APP_VERSION})

patch:
	$(eval APP_VERSION := $(shell echo ${MAJOR}.${MINOR}.$$(( ${PATCH} + 1 ))))
	$(call update_semver,${APP_VERSION})

version:
	$(call print_info, Current version: ${APP_VERSION}...)


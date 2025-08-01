# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DATABASE
# ==================================================================================== #
#
## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database="${ROSETTA_DB_CONN_STR}" up

## db/migrations/goto number=$1: target versiont to migrate to
.PHONY: db/migrations/goto
db/migrations/goto: confirm
	@echo 'Running down migrations...'
	migrate -path=./migrations -database="${ROSETTA_DB_CONN_STR}" goto ${number}

## db/migrations/down
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path=./migrations -database="${ROSETTA_DB_CONN_STR}" down


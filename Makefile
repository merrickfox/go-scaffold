.PHONY: default
default: local_db


.PHONY: local_db
local_db:
	-docker network create serv_network
	docker-compose -f docker-compose-local-db.yml up --build

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: .venv test.py test.go

.DEFAULT_GOAL := help

help:  ## Show this help
	@grep -E '^([a-zA-Z_.-]+):.*## ' $(MAKEFILE_LIST) | sort | awk -F ':.*## ' '{printf "%-20s %s\n", $$1, $$2}'

up:  ## Start all containers
	@docker compose up --remove-orphans -V -d

.PHONY: start
start: up  ## Start all containers but don't recreate

down: ## Stop and remove all containers
	@docker compose down

.PHONY: stop
stop: down  ## Stop but keep all containers

destroy:  ## Stop all containers and remove named volumes
	@docker compose down -v

logs:  ## Tail container logs
	@docker compose logs -f

ps:  ## List running containers
	@docker compose ps

kill:  ## Kill a random running NATS container
	@docker ps --format '{{.Names}}' | grep nats-nats | shuf -n 1 | xargs docker stop

kill_leader:  ## Specifically kill JetStream cluster leader
	@curl -s http://localhost:8222/jsz | jq -r .meta_cluster.leader | xargs -rI{} docker stop nats-{}-1

.ONESHELL:
account_info:  ## Connect to NATS and show account info
	@echo localhost:
	nats --user=a --password=a -s nats://localhost:4222 account info | grep -E '(Connected|Streams)'	
	docker compose ps -a --format '{{.Names}}' | grep nats-nats | while read c; do
		echo "\n$${c}:"
		ip=$$(docker inspect --format '{{range .NetworkSettings.Networks}} {{.IPAddress}}{{end}}' $${c} | grep -Eo '[0-9.]+')
		if [ -n "$${ip}" ]; then
			nats --user=a --password=a -s nats://$${ip}:4222 account info | grep -E '(Connected|Streams)'
		else
			echo "  server down"
		fi
	done
	true

.ONESHELL:
js_status:  ## Get JetStream status
	@curl http://localhost:8222/jsz
	for i in $$(seq 3); do echo nats-nats$${i}-1: && docker exec -it nats-nats$${i}-1 ls -l /data/jetstream/TK/streams/tk/msgs ; done
	true

.venv:
	@test -d .venv || python3 -m venv .venv
	@. .venv/bin/activate && pip install -qr requirements.txt

test.py: .venv  ## Run sample nats.py
	@. .venv/bin/activate && python3 test.py
	@true

test.go:  ## Run sample nats.go
	@go run test.go

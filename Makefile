run-backend:
	cd backend && \
	swag init --generalInfo cmd/main.go --output docs && \
	docker compose down && \
	docker compose up -d postgres && \
	docker compose up -d directus && \
	air

run-telegram-bot:
	cd telegram-bot && \
	docker compose up --build -d

stop:
	cd backend && \
	docker compose down \
	cd ../telegram-bot \
	docker compose down

clean:
	cd backend && \
	docker compose down --volumes --remove-orphans && \
	rm -rf tmp/* \
	cd ../telegram-bot \
	docker compose down --volumes --remove-orphans && \
	rm -rf tmp/*
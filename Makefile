dev:
	cd backend && \
	swag init --generalInfo cmd/main.go --output docs && \
	docker compose down && \
	docker compose up -d postgres && \
	docker compose up -d directus && \
	air

stop:
	cd backend && \
	docker compose down

clean:
	cd backend && \
	docker compose down --volumes --remove-orphans && \
	rm -rf tmp/*
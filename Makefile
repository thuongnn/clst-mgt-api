.PHONY: dev dev-down go proto

dev:
	docker-compose up -d

dev-down:
	docker-compose down

go:
	air


## Database
DB_URL=postgres://user:admin@localhost:5432/marketflow?sslmode=disable

## Создание новой миграции: make migrate-create name=название
migrate-create:
	@echo "Creating new migration: $(name)"
	migrate create -seq -ext=.sql -dir=./migrations $(name)

## Применить все миграции
migrate-up:
	migrate -path=./migrations -database "$(DB_URL)" up

## Применить N миграций: make migrate-upn n=2
migrate-upn:
	migrate -path=./migrations -database "$(DB_URL)" up $(n)

## Откатить одну миграцию
migrate-down1:
	migrate -path=./migrations -database "$(DB_URL)" down 1

## Откатить все миграции
migrate-down:
	migrate -path=./migrations -database "$(DB_URL)" down

## Посмотреть текущую версию миграций
migrate-version:
	migrate -path=./migrations -database "$(DB_URL)" version

## Пропустить миграции до определённой версии: make migrate-force v=2
migrate-force:
	migrate -path=./migrations -database "$(DB_URL)" force $(v)

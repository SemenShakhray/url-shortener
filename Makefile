include .env
LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@latest

migration-add:
	goose -dir ${MIGRATION_DIR} create ${MIGRATION_NAME} sql

migration-up:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} up -v

migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} down -v

migration-reset:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} reset -v

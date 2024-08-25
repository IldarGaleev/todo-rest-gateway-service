


ifeq ($(OS),Windows_NT)
RM_TOOL:=del
else
RM_TOOL:=rm
endif


.PHONY: cover gen_swag

cover:
	go test -short -count=1 -coverprofile=.\coverage.out ./...
	go tool cover -html=.\coverage.out
	@$(RM_TOOL) .\coverage.out

gen_swag:
	swag init -g ./cmd/todo/main.go
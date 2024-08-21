


ifeq ($(OS),Windows_NT)
RM_TOOL:=del
else
RM_TOOL:=rm
endif


.PHONY: cover

cover:
	go test -short -count=1 -coverprofile=.\coverage.out ./...
	go tool cover -html=.\coverage.out
	@$(RM_TOOL) .\coverage.out


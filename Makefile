BUILD_DIR = build
BUILD_NAME = repo-measurer

# Build the go program
build:
	go build -o $(BUILD_DIR)/$(BUILD_NAME) main.go

# Run the go program
run:
	go run main.go

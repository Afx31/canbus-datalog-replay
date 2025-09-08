BINARY_NAME=bin/canbus-datalog-replay
BINARY_DIR=./

# Compile executable
compile:
	go build -o ${BINARY_NAME} ${BINARY_DIR}

# Compile executable and runs
rebuild: compile
	./${BINARY_NAME}

# Runs pre-compiled executable
run:
	./${BINARY_NAME}

# Cleans and removes executable
purge:
	go clean
	rm ${BINARY_NAME}

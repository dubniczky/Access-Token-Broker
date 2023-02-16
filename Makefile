main = main.go

# Run the server
start::
	go run $(main) config.go

# Start the server in development mode (restart on file changes)
dev:: node_modules
	npx nodemon \
		--watch './**/*.go' \
		--signal SIGTERM \
		--exec 'make' start

# Install node modules
node_modules:
	npm install

main = ats
config = config.yaml
out_dir = bin

# Run the server
start::
	go run .

# Start the server in development mode (restart on file changes)
dev:: node_modules
	npx nodemon \
		--watch './**/*.go' \
		--signal SIGTERM \
		--exec 'make' start

# Build the project
build::
	mkdir -p $(out_dir)
	rm -rf $(out_dir)/*
	cp $(config) $(out_dir)
	go build -o $(out_dir)/$(main) .

# Install node modules
node_modules:
	npm install

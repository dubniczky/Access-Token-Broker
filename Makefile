main = main.go

# Run the server
start::
	go run $(main)

dev::
	npx nodemon \
		--watch './**/*.go' \
		--signal SIGTERM \
		--exec 'go' run $(main)

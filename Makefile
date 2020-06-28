# Set our default go compiler
GO := go
NAME := rpminfo

.DEFAULT: build

build:
	$(GO) build -o $(NAME)

test:
	$(GO) test

lint:
	$(GO) fmt

clean:
	rm -f $(NAME)

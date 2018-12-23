FUNCTION_ALIAS ?= prd
S3_BUCKET_NAME ?= ""
STACK_NAME ?= runtime-golang


.DEFAULT_GOAL := build
install:
	go get -t ./cmd/lambda
.PHONY: install



pre_build: install mocks
#	gometalinter.v2 ./...
#	go test -v ./...
.PHONY: pre_build

build: pre_build
	GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/lambda
	@zip -9 -r ./handler.zip bootstrap
.PHONY: build

deploy: build
	aws cloudformation package \
		--template-file cfn.yaml \
		--output-template-file cfn.out.yaml \
		--s3-bucket $(S3_BUCKET_NAME) \
		--s3-prefix cfn

	aws cloudformation deploy \
		--template-file cfn.out.yaml \
		--capabilities CAPABILITY_IAM \
		--stack-name $(STACK_NAME) \
        --parameter-overrides \
        	FunctionAlias=$(FUNCTION_ALIAS) \
        	FunctionName=$(STACK_NAME) \
        	DomainName=$(DOMAIN_NAME) \
        	SubjectAlternativeNames=$(SUBJECT_ALTERNATIVE_NAMES)
.PHONY: deploy
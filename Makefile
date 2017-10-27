OUTPUT_DIR=./_output
BINARY=${OUTPUT_DIR}/conntracker
ORGANIZATION=hyperpilot
IMAGE=k8sconntrack
TAG=test

build: clean
	go build -o ${BINARY} ./cmd

docker-test:
	docker build -t ${ORGANIZATION}/${IMAGE}:${TAG} .

docker:
	docker build -t ${ORGANIZATION}/${IMAGE}:test .

.PHONY: clean
clean:
	@: if [ -f ${OUTPUT_DIR} ]; then rm -rf ${OUTPUT_DIR};fi

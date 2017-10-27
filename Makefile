OUTPUT_DIR=./_output
BINARY=${OUTPUT_DIR}/conntracker
ORGANIZATION=hyperpilot
IMAGE=k8sconntrack
TAG=test

build: clean
	go build -o ${BINARY} ./cmd

docker:
	docker build -t ${ORGANIZATION}/${IMAGE}:latest .

generate-secret:
	kubectl create secret generic vmt-config --from-file ~/.kube/config

.PHONY: clean
clean:
	@: if [ -f ${OUTPUT_DIR} ]; then rm -rf ${OUTPUT_DIR};fi

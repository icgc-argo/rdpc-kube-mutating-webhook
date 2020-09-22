#
# Copyright (c) 2020. Ontario Institute for Cancer Research
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
#
.DEFAULT_GOAL := docker-image

K8S_NAMESPACE := webhook-demo
KUBECTL_EXE := /usr/bin/kubectl
K8S_CMD := $(KUBECTL_EXE) -n $(K8S_NAMESPACE)
DOCKER_ACCOUNT = rtisma1
IMAGE ?= $(DOCKER_ACCOUNT)/webhook-go-server:latest

# this is the legacy build that is also replicated in the Dockerfile
image/webhook-server: $(shell find . -name '*.go')
	./init-build.sh
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $@ ./cmd/webhook-server

docker-image:
	docker build -t $(IMAGE) ./

push-image: docker-image
	docker push $(IMAGE)

send-request:
	cd ./examples/ && ./run-server-request.sh 8080

start-docker:
	docker-compose up --build -d

stop-docker:
	docker-compose down -v

k8s-destroy:
	@./destroy.sh

k8s-deploy:
	@./deploy.sh

k8s-test:
	@$(K8S_CMD) apply -f ./examples/pod-with-defaults.yaml
	@$(K8S_CMD) get pods -oyaml pod-with-defaults

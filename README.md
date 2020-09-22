# Kubernetes Admission Controller Webhook Demo

## Tools
1. Take a left json and right json, generate a list of jsonPatches: https://json-patch-builder-online.github.io/
2. Take an input json and a list of patches, and generate what the output json would look like: https://json-schema-validator.herokuapp.com/jsonpatch.jsp

## Testing with docker
This will send a dummy admission review payload to the webhook server. This is useful for development of the webhook server

1. Start server with `make start-docker`
2. Send an example request located in ./examples/admission-review.example.json via `make send-request`
3. Destroy the docker environment with `make stop-docker`

## Building Docker Image
1. `make docker-image` : this builds the docker image locally
2. `make push-image` : this will build and push the docker image. 
                        Ensure you run `docker login` before running this command and change 
                        the `DOCKER_ACCOUNT` makefile variable to match your dockerhub account username.

## Testing in K8s
Assuming you have permissions to deploy pods, mutating webhook configurations, and create namespaces, 
the following will deploy the webhook service and test it in a k8s environment
1. `make k8s-deploy`: This will deploy the webhook service and mutatingwebhookconfiguration
2. `make k8s-test`: This will apply the resource ./examples/pod-with-defaults.yaml which is just a busybox echo command.
                    In addition, the pod description in yaml format will be dumped. 
                    You can observer the spec.volumes and spec.containers section for volume mounts and see emptyDirs there. 
                    This confirms the original ./examples/pod-with-defaults.yaml was mutated to have emptyDirs
3. `make k8s-destroy`: This will destroy the environment previously created




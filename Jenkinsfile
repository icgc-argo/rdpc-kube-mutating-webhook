@Library(value='jenkins-pipeline-library@master', changelog=false) _
pipelineRDPCRdpcKubeMutatingWebhook(
    buildImage: "node:12.6.0",
    dockerRegistry: "ghcr.io",
    dockerRepo: "icgc-argo/rdpc-kube-mutating-webhook",
    gitRepo: "icgc-argo/rdpc-kube-mutating-webhook",
    testCommand: "true",
    helmRelease: "kube-webhook"
)
# SPIFFE EXAMPLE

Simple Go script to communicate with the [SPIFFE](https://spiffe.io/) implementation: [Spire](https://spiffe.io/docs/latest/spire-about/spire-concepts/).

# How to run it

## Kubernetes

First, you may have a [Kubernetes](https://kubernetes.io/) cluster running, or run it locally using [Minikube](https://minikube.sigs.k8s.io/docs/).

`Obs`: If you are running this example in Minikube, you must start it running the following command in terminal:

```sh
minikube start \
    --extra-config=apiserver.service-account-signing-key-file=/var/lib/minikube/certs/sa.key \
    --extra-config=apiserver.service-account-key-file=/var/lib/minikube/certs/sa.pub \
    --extra-config=apiserver.service-account-issuer=api \
    --extra-config=apiserver.service-account-api-audiences=api,spire-server \
    --extra-config=apiserver.authorization-mode=Node,RBAC
```

And if you want to run Docker locally with Minikube, you must run:

```sh
eval $(minikube -p minikube docker-env)
```

## Spire deploy in K8s

After you have a K8s cluster running, you can deploy the Spire Server and Agent in it. For more info about each manifest and Spire in general, visit the [Spire Quickstart in K8s](https://spiffe.io/docs/latest/try/getting-started-k8s/) and the [Spire Documentation](https://spiffe.io/docs/latest/spire-about/spire-concepts/).

### Creating Spire Namespace
```sh
kubectl apply -f spire/spire-namespace.yaml
```

### Deploying the Spire Server
```sh
kubectl apply -f spire/server
```

### Deploying the Spire Agent
```sh
kubectl apply -f spire/agent
```


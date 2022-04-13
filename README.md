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
## Registering tha agent and the application

First, we need to register the Spire Agent to the server, running:

```sh
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s_sat:cluster:demo-cluster \
    -selector k8s_sat:agent_ns:spire \
    -selector k8s_sat:agent_sa:spire-agent \
    -node
```

Once the agent is registered, let's register the Workload.

```sh
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/app/server \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s:ns:app \
    -selector k8s:container-name:server
```

```sh
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/app/client \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s:ns:app \
    -selector k8s:container-name:client
```

## Running the application

After you run the Spire Server and Spire Agent and registered the Agent and the Workload, it's time to run the application.

First you must build the docker image by running:

```sh
docker build -t <image-name> .
```

into the `/application` directory.

After building the image, you must edit the `/application/k8s/application-deployment.yaml` to put the image name in the `containers.image` option.

Then, run the following commands:

```sh
kubectl apply -f application-namespace.yaml
kubectl apply -f k8s/application-deployment.yaml
```

After this, the application must be running.

## Testing if it is working

To verify if the application is working well, first lets check if the pod is Ready and Running.

```sh
kubectl get pods --namespace=app
```

It must return something like

```sh
NAME                   READY   STATUS             RESTARTS        AGE
app-57b7d6cc6d-2nkv8   0/1     CrashLoopBackOff   19 (3m7s ago)   75m
```

Now, let's see if the application was able to get tha SpiffeId. By running

```sh
kubectl logs <pod-name> --namespace=app
```

It must return something like:

```sh
Success fetching SVID
spiffe://example.org/ns/default/sa/default
```

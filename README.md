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
    -selector k8s:container-name:<server_image_name>
```

```sh
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/app/client \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s:ns:app \
    -selector k8s:container-name:<client_image_name>
```

## Running the application

After you run the Spire Server and Spire Agent and registered the Agent and the Workloads, it's time to run the application.

First you must enter in `/application` directory and create the application namespace by running:

```sh
kubectl apply -f application-namespace.yaml
```

`OBS.:` You must see that the default namespace of the application is `app`. If you want to change it, you must change the workload registration `k8s:ns` selector from `-selector k8s:ns:app` to `-selector k8s:ns:<your_custom_namespace>`.

After creating the namespace you must build the server image by running:

```sh
docker build -t <server_image_name> .
```

into the `/application/server` directory.

Then you must run

```sh
kubectl apply -f k8s/
```

to deploy you server application into Kubernetes. Notice that in the `server-deployment.yaml` you have the `imagePullPolicy` property, this property must be settled to `Never` if you are running the minikube and using the minikube docker env.
Notice that you must modify the `image` property too, and insert the image name that you've settled in `docker build` command.

Now, the server must be running. To check the deployment run:

```sh
kubectl get pods -n <application_namespace>
```

To run the client, you must build the its docker image too. So, into `/application/client` directory, run:

```sh
docker build -t <client_image_name> .
```

then edit the `client-deployment.yaml` file to set the `imagePullPolicy` property and the `image` property as needed.

After editing, run:

```sh
kubectl get pods -n <application_namespace>
```

If everything is running, it must return something like:

```sh
NAME                      READY   STATUS    RESTARTS   AGE
client-7d97f7958f-9mvn2   1/1     Running   0          4s
server-67f85cd649-rgsh2   1/1     Running   0          8s
```

## Testing if it is working

To verify if the application is working well, first lets check if the pods are Ready and Running.

```sh
kubectl get pods -n <application_namespace>
```

If everything is running, it must return something like:

```sh
NAME                      READY   STATUS    RESTARTS   AGE
client-7d97f7958f-9mvn2   1/1     Running   0          4s
server-67f85cd649-rgsh2   1/1     Running   0          8s
```

Now, let's see if the server and the client are communicating well.

First, let's see the client logs:

```sh
kubectl logs <clients_pod_name> -n <application_namespace>
```

If everything is running good, you must have something like:

```sh
Success fetching SVID
spiffe://example.org/ns/app/client
2022/04/18 20:57:23 Server says: "Hello client!\n"
2022/04/18 20:57:28 Server says: "Hello client!\n"
2022/04/18 20:57:33 Server says: "Hello client!\n"
2022/04/18 20:57:38 Server says: "Hello client!\n"
```

Let's see the server logs now:

Run:

```sh
kubectl logs <server_pod_name> -n <application_namespace>
```

If everything is running good, you must have something like:

```sh
Success fetching SVID
spiffe://example.org/ns/app/server
2022/04/18 20:57:23 Client says: "Hello server!\n"
2022/04/18 20:57:28 Client says: "Hello server!\n"
2022/04/18 20:57:33 Client says: "Hello server!\n"
2022/04/18 20:57:38 Client says: "Hello server!\n"
```

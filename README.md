# Demo gRPC server/client on K8s with Istio Load balance

## Prerequisites

1. Acces to k8s cluster
2. Istio installed

## Deploy

```bash
make compile
make build_client
make build_server
```

Optionally push the built images

```
docker push -t $GRPC_GO_CLIENT .
docker push -t $GRPC_GO_SERVER .
```

Deploy services to Kubernetes
```
kubectl apply -f deployment.yml
```

Wait for pods to come up. Check the logs.

```
kubectl get po
```
```
NAME                      READY   STATUS            RESTARTS   AGE
client-6675b67b45-9x8hp   0/2     PodInitializing   0          5s
server-5795c976fc-4l9kn   0/2     PodInitializing   0          5s
server-5795c976fc-5f2r9   0/2     PodInitializing   0          5s
server-5795c976fc-7bthn   2/2     Running           0          5s
```
```
kubectl logs -f client-6675b67b45-9x8hp
```
```
2021/10/27 11:06:59 Connecting to server on grpc-server:8000
2021/10/27 11:07:04 Unary response from server: Hello Primoz from server 172.17.0.8
2021/10/27 11:07:04 Server stream response: Hello Primoz from server 172.17.0.6
2021/10/27 11:07:04 Server stream response: Hello Primoz from server 172.17.0.6
2021/10/27 11:07:04 Server stream response: Hello Primoz from server 172.17.0.6
2021/10/27 11:07:04 Server stream response: Hello Primoz from server 172.17.0.6
2021/10/27 11:07:04 Bistream server response: Hello Primoz from server 172.17.0.7
2021/10/27 11:07:04 Bistream server response: Hello Marko from server 172.17.0.7
2021/10/27 11:07:04 Bistream server response: Hello Matej from server 172.17.0.7
2021/10/27 11:07:09 Unary response from server: Hello Primoz from server 172.17.0.8
...
```
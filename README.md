# Custom Resource Definition in Kubernetes  + Handling CRD with Client-go's Dynamic Package

### Create CRD 
```shell
$ kubectl apply -f crd.yaml
```

### Now run the client-go code using:
```shell
$ go run main.go
```
## Necessary Command
```shell
$ kubectl get crds #get the crds
$ kubectl get songs #get the cr(song)
$ kubectl get songs.music.sportshead.dev my-favourite-song -n default -o yaml #Detailed information of existing cr
```

## Delete CRD
```shell
$ kubectl delete -f crd.yaml
```
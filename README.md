# apps-scanner

```shell
$ kustomize build https://github.com/weaveworks-gitops-poc/wego-sockshop/apps/sockshop/environments/dev | kubectl apply -f -
```

And build and run...

```shell
$ go build ./cmd/scanner
$ ./scanner
2022/01/19 21:18:39 Starting to scan for applications
2022/01/19 21:18:40 found 14 pods
2022/01/19 21:18:40       applications sockshop
2022/01/19 21:18:40       envs dev
2022/01/19 21:18:40       services cart,cart-db,catalog,catalog-db,frontend,orders,orders-db,payment,queue-master,rabbitmq,session-db,shipping,user,user-db
```

# Installation from Flux

```shell
$ flux install --components source-controller,kustomize-controller
```

```yaml
apiVersion: source.toolkit.fluxcd.io/v1beta1
kind: GitRepository
metadata:
  name: sockshop-repo
  namespace: default
spec:
  interval: 15m
  url: https://github.com/weaveworks-gitops-poc/wego-sockshop.git
```

```yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1beta2
kind: Kustomization
metadata:
  name: sockshop-dev
  namespace: default
spec:
  interval: 5m
  path: ./apps/sockshop/environments
  prune: true
  sourceRef:
    kind: GitRepository
    name: sockshop-repo
```

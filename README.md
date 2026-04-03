Старт `minicube`
```bash
minikube start
```

Команды
```bash
kubectl apply -f base/00-namespace.yaml
kubectl apply -f base/01-resource-quota.yaml
kubectl apply -f base/02-limit-range.yaml
```
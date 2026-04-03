## Старт `minicube`
```bash
minikube start
```

## Команды для создания неймспейса и лимитов
```bash
kubectl apply -f base/00-namespace.yaml
kubectl apply -f base/01-resource-quota.yaml
kubectl apply -f base/02-limit-range.yaml
```
Проверка статуса 
```bash
kubectl describe quota compute-resources -n k8s-lifecycle-prod
```

## Создать файлы с настройками проекта и зашифрованными данными
```bash
kubectl apply -f configs/01-app-configmap.yaml
kubectl apply -f configs/02-app-secret.yaml
```

Получение описания проекта
```bash
kubectl get configmap app-settings -n k8s-lifecycle-prod -o yaml
```

Получение зашифрованных данных с проекта
```bash
kubectl get secret app-credentials -n k8s-lifecycle-prod -o jsonpath='{.data.DB_USER}' | base64 --decode; echo ""
```

Проверить что запустилось прямо сейчас
```bash
kubectl get all -n k8s-lifecycle-prod
```

Получить переменные окружения
```bash
kubectl exec -it deployment/lifecycle-app -n k8s-lifecycle-prod -- printenv | grep APP_
```
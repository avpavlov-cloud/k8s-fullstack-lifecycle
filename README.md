## k8s-fullstack-lifecycle
Проект по созданию полной инфраструктуры приложения в Kubernetes: от базовых лимитов и баз данных до собственного Kubernetes Operator на языке Go.
## 🚀 Быстрый старт## 1. Подготовка кластера

minikube start# (Опционально) Добавление второго узла для проверки DaemonSet
minikube node add

## 2. Фундамент (Namespace, Quotas, Limits)
Создаем изолированную среду и устанавливаем правила игры.

kubectl apply -f base/# Проверка квоты (почему может не хватать ресурсов на 5 реплик)
kubectl describe quota compute-resources -n k8s-lifecycle-prod

## 3. Конфигурация и Секреты

kubectl apply -f configs/# Проверка расшифровки секрета
kubectl get secret app-credentials -n k8s-lifecycle-prod -o jsonpath='{.data.DB_USER}' | base64 --decode; echo ""

## 4. Развертывание слоев (Stateless & Stateful)

# Приложение (Deployment)
kubectl apply -f app-stateless/# База данных (StatefulSet + PV + StorageClass)
kubectl apply -f app-stateful/

## 5. Системные сервисы и Задачи (DaemonSet, Jobs)

# Сборщик логов на каждом узле
kubectl apply -f infrastructure/# Бэкап и автоматическая чистка (CronJob)
kubectl apply -f jobs/

## 6. Расширение Kubernetes (CRD + Operator)
Самый продвинутый этап: обучаем Kubernetes нашему ресурсу MySite.
Регистрация нового типа ресурсов:

kubectl apply -f custom-res/01-crd.yaml

Запуск "мозга" (Оператора) на Go:
Убедитесь, что вы находитесь в папке operator/

go mod tidy
go run main.go

Создание своего сайта через оператор:

# В новом терминале
kubectl apply -f custom-res/02-sample-site.yaml

## 🔍 Полезные команды для отладки

* Посмотреть всё сразу: kubectl get all -n k8s-lifecycle-prod
* Проверить работу оператора: kubectl get ms -n k8s-lifecycle-prod (сокращенно от mysites)
* Логи приложения: kubectl logs -l app=lifecycle-app -n k8s-lifecycle-prod
* Статус масштабирования: kubectl describe deployment marketing-site-deploy -n k8s-lifecycle-prod


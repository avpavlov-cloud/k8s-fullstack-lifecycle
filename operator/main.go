package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1" // Для работы с Deployment
	corev1 "k8s.io/api/core/v1" // Для работы с Pods и Containers
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// --- СТРУКТУРЫ ДАННЫХ ---
type MySite struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MySiteSpec `json:"spec"`
}

type MySiteSpec struct {
	Image    string `json:"image"`
	Replicas int32  `json:"replicas"`
}

func (in *MySite) DeepCopyObject() runtime.Object {
	out := &MySite{}
	*out = *in
	return out
}

type MySiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MySite `json:"items"`
}

func (in *MySiteList) DeepCopyObject() runtime.Object {
	out := &MySiteList{}
	*out = *in
	if in.Items != nil {
		out.Items = make([]MySite, len(in.Items))
		copy(out.Items, in.Items)
	}
	return out
}

// --- ЛОГИКА ОПЕРАТОРА ---
type MySiteReconciler struct {
	client.Client
}

func (r *MySiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// 1. Получаем сам объект MySite из базы
	var mySite MySite
	if err := r.Get(ctx, req.NamespacedName, &mySite); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	fmt.Printf("\n[ОПЕРАТОР] Обработка %s. Образ: %s, Реплик: %d\n", mySite.Name, mySite.Spec.Image, mySite.Spec.Replicas)

	// 2. Описываем желаемый Deployment
	deploymentName := mySite.Name + "-deploy"
	foundDeployment := &appsv1.Deployment{}

	err := r.Get(ctx, client.ObjectKey{Name: deploymentName, Namespace: mySite.Namespace}, foundDeployment)

	if err != nil && errors.IsNotFound(err) {
		// 3. ЕСЛИ ДЕПЛОЙМЕНТА НЕТ — СОЗДАЕМ ЕГО
		fmt.Printf("[ОПЕРАТОР] Создаю новый Deployment: %s\n", deploymentName)

		newDeploy := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: mySite.Namespace,
				// Привязываем деплоймент к MySite (если удалить MySite, удалится и деплоймент)
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(&mySite, schema.GroupVersionKind{
						Group:   "stable.example.com",
						Version: "v1",
						Kind:    "MySite",
					}),
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &mySite.Spec.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": deploymentName},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": deploymentName}},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "web",
							Image: mySite.Spec.Image,
						}},
					},
				},
			},
		}

		if err := r.Create(ctx, newDeploy); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func main() {
	ctrl.SetLogger(zap.New())
	scheme := runtime.NewScheme()
	gv := schema.GroupVersion{Group: "stable.example.com", Version: "v1"}
	scheme.AddKnownTypes(gv, &MySite{}, &MySiteList{})
	metav1.AddToGroupVersion(scheme, gv)
	// Добавляем стандартные типы (Deployment), чтобы оператор их понимал
	_ = appsv1.AddToScheme(scheme)

	mgr, _ := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{Scheme: scheme})
	_ = ctrl.NewControllerManagedBy(mgr).For(&MySite{}).Complete(&MySiteReconciler{Client: mgr.GetClient()})
	_ = mgr.Start(ctrl.SetupSignalHandler())
}

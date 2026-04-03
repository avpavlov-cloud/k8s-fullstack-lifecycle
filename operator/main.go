package main

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// 1. Структура одиночного объекта
type MySite struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec struct {
		Image    string `json:"image"`
		Replicas int    `json:"replicas"`
	} `json:"spec"`
}

func (in *MySite) DeepCopyObject() runtime.Object {
	out := &MySite{}
	*out = *in
	return out
}

// 2. СТРУКТУРА СПИСКА (Этого не хватало!)
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

// 3. Логика контроллера
type MySiteReconciler struct {
	client.Client
}

func (r *MySiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	fmt.Printf("\n[ОПЕРАТОР] >>> Вижу объект MySite: %s в пространстве %s\n", req.Name, req.Namespace)
	return ctrl.Result{}, nil
}

func main() {
	ctrl.SetLogger(zap.New())
	fmt.Println(">>> Запуск оператора k8s-fullstack-lifecycle...")

	scheme := runtime.NewScheme()
	gv := schema.GroupVersion{Group: "stable.example.com", Version: "v1"}
	
	// Регистрируем и объект, и СПИСОК объектов
	scheme.AddKnownTypes(gv, &MySite{}, &MySiteList{})
	metav1.AddToGroupVersion(scheme, gv)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		os.Exit(1)
	}

	err = ctrl.NewControllerManagedBy(mgr).
		For(&MySite{}). 
		Complete(&MySiteReconciler{Client: mgr.GetClient()})
	
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		os.Exit(1)
	}
}

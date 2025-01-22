/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	// 导入 Kubernetes 核心 API 包
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "github.com/fuxiansen/i-operator/api/v1"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.crd.fuxiansen.com,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.crd.fuxiansen.com,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.crd.fuxiansen.com,resources=applications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.4/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	fmt.Println("-------------------------test-------------------------")
	logger := log.FromContext(ctx)
	//获取Application api 实例对象
	application := &corev1.Application{}
	err := r.Get(ctx, req.NamespacedName, application)
	if err != nil {
		logger.Error(err, "unable to fetch Application")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// 获取当前api对象关联运行的pod
	podlist := &v1.PodList{} // 使用 Kubernetes 核心 API 中的 PodList 类型
	err = r.List(ctx, podlist, client.InNamespace(req.Namespace), client.MatchingLabels{"app": application.Name})
	if err != nil {
		logger.Error(err, "unable to fetch Pod")
		return ctrl.Result{}, err
	}

	currentPodNum := len(podlist.Items)

	if currentPodNum < int(application.Spec.Replicas) {
		// 当前运行的Pod少于期望状态Pod的数量，创建pod
		for i := currentPodNum; i < int(application.Spec.Replicas); i++ {
			pod := &v1.Pod{ // 使用 Kubernetes 核心 API 中的 Pod 类型
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%d", application.Name, i),
					Namespace: req.Namespace,
					Labels: map[string]string{
						"app": application.Name,
					},
				},
				Spec: application.Spec.Template.Spec,
			}
			// 建立关联
			if err := ctrl.SetControllerReference(application, pod, r.Scheme); err != nil {
				logger.Error(err, "unable to set ownerReference")
				return ctrl.Result{}, err
			}
			// 创建pod
			if err := r.Create(ctx, pod); err != nil {
				logger.Error(err, "unable to create Pod")
				return ctrl.Result{}, err
			}
		}

	} else if currentPodNum > int(application.Spec.Replicas) {
		// 当前运行的Pod大于期望状态Pod的数量，删除pod
		deletePodNum := podlist.Items[:currentPodNum-int(application.Spec.Replicas)]
		for _, pod := range deletePodNum {
			if err := r.Delete(ctx, &pod); err != nil {
				logger.Error(err, "unable to delete Pod")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Application{}).
		Named("application").
		Complete(r)
}

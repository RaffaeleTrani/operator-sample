/*


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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	webappv1 "operator-example/api/v1"
)

// OperatorReconciler reconciles a Operator object
type OperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.my.domain,resources=operators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.my.domain,resources=operators/status,verbs=get;update;patch

func (r *OperatorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("operator", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *OperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := ctrl.NewControllerManagedBy(mgr).For(&webappv1.Operator{}).Build(r)

	if err == nill && c != nil {
		errWatch := c.Watch(
			&source.Kind{Type: &appsv1.Pod{}),
			handler.EnqueueRequestsFromMapFunc(func(a client.Object) []reconcile.Request {
        			return []reconcile.Request{
            				{NamespacedName: types.NamespacedName{
                			Name: a.GetName(),
                			Namespace: "default",
            				}},
			}),
			predicate.Funcs{
				CreateFunc: func(e event.CreateEvent) {
					
		if errWatch != nil {
			fmt.Println("watch non settato")
		} else {
			fmt.Println("watch settato")
		}
	}
	return err
}

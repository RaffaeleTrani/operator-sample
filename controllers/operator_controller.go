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
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"github.com/vishvananda/netlink"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	batchv1 "operator-example/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	ctx := context.Background()
	log := r.Log.WithValues("operator", req.NamespacedName)

	p := &corev1.Pod{}
	if err := r.Get(ctx, req.NamespacedName, p); err != nil {
		log.Error(err, "unable to fetch pod")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	//ottengo la lista di link (in quale namespace???)
	linkList, _ := netlink.LinkList()
	vethLink := linkList[len(linkList)]
	brExists := false
	var brLink netlink.Link
	//cerco se esiste già il bridge
	for _, link := range linkList {
		if link.Attrs().Name == "br1" {
			brExists = true
			brLink = link
			break
		}
	}
	//il bridge non esiste --> lo creo e gli attacco la veth --> NOTA: per il momento gli attacco l'ultima veth creata
	if brExists != true {
		la := netlink.NewLinkAttrs()
		la.Name = "foo"
		mybridge := &netlink.Bridge{LinkAttrs: la}
		err := netlink.LinkAdd(mybridge)
		if err != nil  {
			fmt.Printf("could not add %s: %v\n", la.Name, err)
		}
		netlink.LinkSetMaster(vethLink, mybridge)
	} else {
		//il bridge esiste, attacco semplicemente la veth pair al bridge
		netlink.LinkSetMaster(vethLink, brLink)
	}

	return ctrl.Result{}, nil
}

func (r *OperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := ctrl.NewControllerManagedBy(mgr).For(&batchv1.Operator{}).Build(r)

	if err == nil && c != nil {
		errWatch := c.Watch(
			&source.Kind{Type: &corev1.Pod{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {

					return []reconcile.Request{
						{NamespacedName: types.NamespacedName{
							Name:      a.Meta.GetName(),
							Namespace: "default",
						}},
					}
				}),
			},
			predicate.Funcs{
				CreateFunc: func(e event.CreateEvent) bool {
					annotations := e.Meta.GetAnnotations()
					for key, element := range annotations {
						if key == "k8s.v1.cni.cncf.io/network" {
							fmt.Println("Multus annotation: ", element)
							return true
						}
					}
					return false
				},
			},
		)
		if errWatch != nil {
			fmt.Println("setted watch")
		} else {
			fmt.Println("watch not setted")
		}
	}
	return err
}

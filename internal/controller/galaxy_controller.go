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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	v1 "kubesphere.io/galaxy/api/v1"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GalaxyReconciler reconciles a Galaxy object
type GalaxyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=astronomy.galaxy.kubesphere.io,resources=galaxies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=astronomy.galaxy.kubesphere.io,resources=galaxies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=astronomy.galaxy.kubesphere.io,resources=galaxies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Galaxy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *GalaxyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	planet := &v1.Planet{}
	if err := r.Client.Get(ctx, req.NamespacedName, planet); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	galaxy := &v1.Galaxy{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name:      planet.GetGalaxy(),
		Namespace: req.Namespace,
	}, galaxy); err != nil {
		return ctrl.Result{}, err
	}

	for _, item := range galaxy.Spec.Planets {
		if item.Name != planet.Name {
			continue
		}
		planet.Spec = v2.PlanetSpec{
			Name:       item.Name,
			DiameterKm: item.DiameterKm,
			HasLife:    item.HasLife,
		}
	}

	if err := r.Update(ctx, planet); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GalaxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&astronomyv1.Galaxy{}).
		Named("galaxy").
		Complete(r)
}

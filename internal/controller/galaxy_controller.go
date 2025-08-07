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
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	v1 "kubesphere.io/galaxy/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// GalaxyReconciler reconciles a Galaxy object
// It generate planets corresponding to galaxy when galaxy is created or updated
type GalaxyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *GalaxyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	galaxy := &v1.Galaxy{}
	if err := r.Client.Get(ctx, req.NamespacedName, galaxy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for _, galaxyPlanetSpec := range galaxy.Spec.Planets {
		planetName := fmt.Sprintf("%s-%s", galaxy.Name, galaxyPlanetSpec.Name)
		desired := &v1.Planet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      planetName,
				Namespace: galaxy.Namespace,
			},
			Spec: galaxyPlanetSpec.ToPlanetSpec(),
		}

		if err := controllerutil.SetControllerReference(galaxy, desired, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		current := &v1.Planet{}
		err := r.Client.Get(ctx, client.ObjectKey{Name: planetName, Namespace: galaxy.Namespace}, current)
		if errors.IsNotFound(err) {
			if err := r.Client.Create(ctx, desired); err != nil && !errors.IsAlreadyExists(err) {
				return ctrl.Result{}, err
			}
			continue
		}
		if err != nil {
			return ctrl.Result{}, err
		}

		needsUpdate := false

		if !reflect.DeepEqual(current.Spec, desired.Spec) {
			current.Spec = desired.Spec
			needsUpdate = true
		}

		// Find old version(not garbage-collected correctly), and mark reference
		if !metav1.IsControlledBy(current, galaxy) {
			if err := controllerutil.SetControllerReference(galaxy, current, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}
			needsUpdate = true
		}

		if needsUpdate {
			if err := r.Client.Update(ctx, current); err != nil {
				return ctrl.Result{}, err
			}
			log.V(1).Info("Updated Planet to match Galaxy spec", "planet", planetName)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GalaxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Galaxy{}).
		Named("galaxy").
		Complete(r)
}

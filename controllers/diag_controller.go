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
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	srev1 "diag/api/v1"
	core "k8s.io/api/core/v1"
)

// DiagReconciler reconciles a Diag object
type DiagReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

// +kubebuilder:rbac:groups=sre.example.com,resources=diags,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sre.example.com,resources=diags/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;watch;list

func (r *DiagReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("diag", req.NamespacedName)

	var Diag srev1.Diag
	if err := r.Get(ctx, req.NamespacedName, &Diag); err != nil {
		log.Error(err, "unable to fetch Server")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, ignoreNotFound(err)
	}

	pods := &core.PodList{}
	if err := r.List(ctx, pods); err != nil {
		log.Error(err, "unable to list pods")
		return ctrl.Result{}, err
	}

	shootDiagStatus := make([]string, 0)

	// probeAPIServer
	if Diag.Spec.ProbeAPIServer {
		probeAPIServer := "probeAPIServer NOK"
		if len(pods.Items) > 0 {
			probeAPIServer = "probeAPIServer OK"
		} else {
			probeAPIServer = "probeAPIServer NOK"
		}
		shootDiagStatus = append(shootDiagStatus, probeAPIServer)
	}

	Diag.Status.ShootDiagStatus = shootDiagStatus

	if err := r.Status().Update(ctx, &Diag); err != nil {
		log.Error(err, "unable to update Diag status")
		return ctrl.Result{}, err
	}

	log.V(1).Info("Updated shootDiagStatus")
	return ctrl.Result{}, nil
}

func (r *DiagReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&srev1.Diag{}).
		Owns(&core.Pod{}).
		Complete(r)
}

package controller

import (
	"github.com/go-logr/logr"
	"github.com/goodhosts/hostsfile"

	"k8s.io/apimachinery/pkg/labels"
	ctrl "sigs.k8s.io/controller-runtime"

	routelisters "github.com/openshift/client-go/route/listers/route/v1"
)

type HostRouteSync struct {
	Address   string
	Lister    routelisters.RouteLister
	Log       logr.Logger
	HostsFile *hostsfile.Hosts
}

func (r *HostRouteSync) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("Starting reconcile")
	routes, err := r.Lister.List(labels.Everything())
	if err != nil {
		r.Log.Error(err, "Failed to list routes")
		return ctrl.Result{}, err
	}
	if err = r.HostsFile.Load(); err != nil {
		r.Log.Error(err, "Failed to load hosts file")
		return ctrl.Result{}, err
	}
	modified := false
	for _, route := range routes {
		if r.HostsFile.HasHostname(route.Spec.Host) {
			if r.HostsFile.Has(r.Address, route.Spec.Host) {
				continue
			}
			if err = r.HostsFile.RemoveByHostname(route.Spec.Host); err != nil {
				r.Log.Error(err, "Failed to remove host", "host", route.Spec.Host)
				return ctrl.Result{}, err
			}
		}
		if err = r.HostsFile.Add(r.Address, route.Spec.Host); err != nil {
			r.Log.Error(err, "Failed to add host", "host", route.Spec.Host)
			return ctrl.Result{}, err
		}
		modified = true
		r.Log.Info("Added host", "host", route.Spec.Host)
	}
	if !modified {
		return ctrl.Result{}, nil
	}
	if err = r.HostsFile.Flush(); err != nil {
		r.Log.Error(err, "Failed to flush hosts file")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

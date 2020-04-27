package controller

import (
	"time"

	"github.com/goodhosts/hostsfile"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"

	routeclient "github.com/openshift/client-go/route/clientset/versioned"
	routeinformers "github.com/openshift/client-go/route/informers/externalversions"

	"github.com/csrwng/host-router/pkg/config"
)

func Setup(cfg *config.HostRouter) error {
	client, err := routeclient.NewForConfig(cfg.Config())
	if err != nil {
		return err
	}
	informerFactory := routeinformers.NewSharedInformerFactory(client, 2*time.Hour)
	cfg.Manager().Add(manager.RunnableFunc(func(stopCh <-chan struct{}) error {
		informerFactory.Start(stopCh)
		return nil
	}))
	routes := informerFactory.Route().V1().Routes()
	hostsFile, err := hostsfile.NewCustomHosts(cfg.HostsFile)
	if err != nil {
		return err
	}
	reconciler := &HostRouteSync{
		Lister:    routes.Lister(),
		Log:       cfg.Logger().WithName("route sync"),
		HostsFile: &hostsFile,
		Address:   cfg.RouterAddress,
	}
	c, err := controller.New("host-route-sync", cfg.Manager(), controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}
	if err := c.Watch(&source.Informer{Informer: routes.Informer()}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}
	return nil
}

package config

import (
	"os"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	ctrl "sigs.k8s.io/controller-runtime"

	routev1 "github.com/openshift/api/route/v1"
)

type SetupFunc func(*HostRouter) error

type HostRouter struct {
	manager ctrl.Manager
	config  *rest.Config
	logger  logr.Logger
	scheme  *runtime.Scheme

	Kubeconfig    string
	RouterAddress string
	HostsFile     string
	SetupFuncs    []SetupFunc
}

func (r *HostRouter) Manager() ctrl.Manager {
	if r.manager != nil {
		return r.manager
	}
	var err error
	r.manager, err = ctrl.NewManager(r.Config(), ctrl.Options{
		Scheme:             r.Scheme(),
		MetricsBindAddress: "0", // disable metrics serving
	})
	if err != nil {
		r.Fatal(err, "failed to create controller manager")
	}
	return r.manager
}

func (r *HostRouter) Start() error {
	if err := r.ensureHostsFile(); err != nil {
		return err
	}
	stopCh := make(chan struct{})
	for _, f := range r.SetupFuncs {
		if err := f(r); err != nil {
			return err
		}
	}
	return r.Manager().Start(stopCh)
}

func (r *HostRouter) Config() *rest.Config {
	if r.config != nil {
		return r.config
	}
	var err error
	r.config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: r.Kubeconfig},
		&clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		r.Fatal(err, "cannot get the cluster's rest config")
	}
	return r.config
}

func (r *HostRouter) Fatal(err error, msg string) {
	r.Logger().Error(err, msg)
	os.Exit(1)
}

func (r *HostRouter) Logger() logr.Logger {
	if r.logger != nil {
		return r.logger
	}
	r.logger = ctrl.Log.WithName("control-plane-operator")
	return r.logger
}

func (r *HostRouter) Scheme() *runtime.Scheme {
	if r.scheme == nil {
		r.scheme = runtime.NewScheme()
		kubescheme.AddToScheme(r.scheme)
		routev1.Install(r.scheme)
	}
	return r.scheme
}

func (r *HostRouter) ensureHostsFile() error {
	_, err := os.Stat(r.HostsFile)
	if err == nil || !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(r.HostsFile)
	if err != nil {
		return err
	}
	return f.Close()
}

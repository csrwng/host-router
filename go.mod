module github.com/csrwng/host-router

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/goodhosts/hostsfile v0.0.2
	github.com/openshift/api v0.0.0-20200424083944-0422dc17083e
	github.com/openshift/client-go v0.0.0-20200422192633-6f6c07fc2a70
	github.com/spf13/cobra v1.0.0
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	sigs.k8s.io/controller-runtime v0.4.0
)

replace (
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20190923180330-3b6373338c9b
	k8s.io/apimachinery => github.com/openshift/kubernetes-apimachinery v0.0.0-20190926190123-4ba2b154755f
	k8s.io/client-go => github.com/openshift/kubernetes-client-go v0.0.0-20190926190130-2917f17b9089
)

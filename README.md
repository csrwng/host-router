# host-router
OpenShift Route to Hosts file synchronizer

Takes a hosts file path, IP address of router, and KUBECONFIG.
Continually manages hosts file based on routes present in OpenShift cluster.

Environment Variables:
- KUBECONFIG: Configuration to communicate with OpenShift cluster
- HOSTSFILE: File containing hosts to manage
- ROUTERIP: External IP of OpenShift router

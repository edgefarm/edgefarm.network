package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"k8s.io/client-go/util/homedir"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/edgefarm/edgefarm.network/cmd/network-dependencies-webhook/pkg/validate"
)

type ServerParameters struct {
	port     int    // webhook server port
	certFile string // path to the x509 certificate for https
	keyFile  string // path to the x509 private key matching `CertFile`
}

var parameters ServerParameters
var config *rest.Config

func main() {
	useKubeConfig := os.Getenv("USE_KUBECONFIG")
	kubeConfigFilePath := os.Getenv("KUBECONFIG")

	flag.IntVar(&parameters.port, "port", 8443, "Webhook server port.")
	flag.StringVar(&parameters.certFile, "tlsCertFile", "/etc/webhook/certs/tls.crt", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&parameters.keyFile, "tlsKeyFile", "/etc/webhook/certs/tls.key", "File containing the x509 private key to --tlsCertFile.")
	flag.Parse()
	if len(useKubeConfig) == 0 {
		// default to service account in cluster token
		c, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		config = c
	} else {
		//load from a kube config
		var kubeconfig string

		if kubeConfigFilePath == "" {
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = filepath.Join(home, ".kube", "config")
			}
		} else {
			kubeconfig = kubeConfigFilePath
		}

		c, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		config = c
	}
	http.HandleFunc("/validating-network-edgefarm-io-v1alpha1-pods", func(w http.ResponseWriter, r *http.Request) {
		validate.Pods(config, w, r)
	})
	http.HandleFunc("/validating-network-edgefarm-io-v1alpha1-accounts", func(w http.ResponseWriter, r *http.Request) {
		validate.Accounts(config, w, r)
	})
	http.HandleFunc("/validating-network-edgefarm-io-v1alpha1-users", func(w http.ResponseWriter, r *http.Request) {
		validate.Users(config, w, r)
	})
	http.HandleFunc("/validating-network-edgefarm-io-v1alpha1-providerconfigs", func(w http.ResponseWriter, r *http.Request) {
		validate.ProviderConfigs(config, w, r)
	})
	http.HandleFunc("/validating-network-edgefarm-io-v1alpha1-streams", func(w http.ResponseWriter, r *http.Request) {
		validate.Streams(config, w, r)
	})
	err := http.ListenAndServeTLS(":"+strconv.Itoa(parameters.port), parameters.certFile, parameters.keyFile, nil)
	if err != nil {
		panic(err)
	}
}

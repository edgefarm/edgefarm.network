package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type DomainsReponse struct {
	Domains []Domain `json:"domains"`
}

type Domain struct {
	Namespace  string  `json:"namespace"`
	Network    string  `json:"network"`
	Subnetwork string  `json:"subnetwork"`
	Node       string  `json:"node"`
	Domain     string  `json:"domain"`
	Stream     *string `json:"stream,omitempty"`
}

type Path struct {
	Namespace  string  `json:"namespace"`
	Network    string  `json:"network"`
	Subnetwork string  `json:"subnetwork"`
	Stream     *string `json:"stream,omitempty"`
}

func getPath(path []string) (*Path, error) {
	if len(path) < 8 {
		return nil, fmt.Errorf("invalid path")
	}
	if path[1] != "domains" {
		return nil, fmt.Errorf("invalid path")
	}
	if path[2] != "namespace" {
		return nil, fmt.Errorf("invalid path")
	}
	if path[4] != "network" {
		return nil, fmt.Errorf("invalid path")
	}
	if path[6] != "subnetwork" {
		return nil, fmt.Errorf("invalid path")
	}
	p := &Path{
		Namespace:  path[3],
		Network:    path[5],
		Subnetwork: path[7],
	}
	if len(path) == 10 {
		if path[8] != "stream" {
			return nil, fmt.Errorf("invalid path")
		}
		p.Stream = &path[9]
	}
	return p, nil
}

func subNetworkDomains(ctx context.Context, c *kubernetes.Clientset, path *Path) ([]Domain, error) {
	domains := []Domain{}
	labelSelector := "network.edgefarm.io/type=leaf"
	if path.Network != "*" {
		labelSelector += "," + fmt.Sprintf("name.network.edgefarm.io/%s=", path.Network)
	}
	if path.Subnetwork != "*" {
		labelSelector += "," + fmt.Sprintf("subnetwork.network.edgefarm.io/%s=", path.Subnetwork)
	}
	podList, err := c.CoreV1().Pods(path.Namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}

	for _, pod := range podList.Items {
		net := ""
		sub := ""
		for l := range pod.Labels {
			if strings.HasPrefix(l, "name.network.edgefarm.io/") {
				net = strings.TrimPrefix(l, "name.network.edgefarm.io/")
			}
			if strings.HasPrefix(l, "subnetwork.network.edgefarm.io/") {
				sub = strings.TrimPrefix(l, "subnetwork.network.edgefarm.io/")
			}
		}

		domains = append(domains, Domain{
			Namespace:  path.Namespace,
			Network:    net,
			Subnetwork: sub,
			Node:       pod.Spec.NodeName,
			Domain:     fmt.Sprintf("%s-%s-%s", net, sub, pod.Spec.NodeName),
		})

	}
	return domains, nil
}

func streamDomains(ctx context.Context, client *dynamic.DynamicClient, path *Path) ([]Domain, error) {
	gvr := schema.GroupVersionResource{
		Group:    "streams.network.edgefarm.io",
		Version:  "v1alpha1",
		Resource: "xstreams",
	}
	labelSelector := "streams.network.edgefarm.io/type=Standard"

	if path.Network != "*" {
		labelSelector += fmt.Sprintf(",streams.network.edgefarm.io/network=%s", path.Network)
	}
	if path.Namespace != "*" {
		labelSelector += fmt.Sprintf(",streams.network.edgefarm.io/namespace=%s", path.Namespace)
	}
	if path.Subnetwork != "*" {
		labelSelector += fmt.Sprintf(",streams.network.edgefarm.io/subnetwork=%s", path.Subnetwork)
	}
	if path.Stream != nil {
		labelSelector += fmt.Sprintf(",streams.network.edgefarm.io/stream=%s", *path.Stream)
	}

	streams, err := client.Resource(gvr).List(context.Background(), v1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}
	list := []Domain{}
	for _, stream := range streams.Items {
		namespace := stream.Object["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["streams.network.edgefarm.io/namespace"].(string)
		network := stream.Object["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["streams.network.edgefarm.io/network"].(string)
		subnetwork := stream.Object["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["streams.network.edgefarm.io/subnetwork"].(string)
		node := stream.Object["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["streams.network.edgefarm.io/node"].(string)
		domain := stream.Object["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["streams.network.edgefarm.io/domain"].(string)
		list = append(list, Domain{
			Namespace:  namespace,
			Network:    network,
			Subnetwork: subnetwork,
			Node:       node,
			Domain:     domain,
			Stream:     path.Stream,
		})
	}
	return list, nil
}

func handleDomainsEndpoint(ctx context.Context, clients *Clients, w http.ResponseWriter, r *http.Request) {
	// Get the namespace, network, and subnetwork parameters from the path.
	s := strings.Split(r.URL.Path, "/")
	path, err := getPath(s)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid URL: %s", err.Error()), http.StatusBadRequest)
		return
	}
	var domains []Domain
	if path.Stream != nil {
		// Get the domains for the specified stream.
		domains, err = streamDomains(ctx, clients.Dynamic, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		domains, err = subNetworkDomains(ctx, clients.Clientset, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Encode the response object as JSON
	jsonResp, err := json.Marshal(DomainsReponse{Domains: domains})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

type Clients struct {
	// Kubernetes clientset.
	Clientset *kubernetes.Clientset
	// Dynamic client.
	Dynamic *dynamic.DynamicClient
}

func main() {
	// Create a Kubernetes client with an incluster config.
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("Error creating incluster config: %v\n", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating clientset: %v\n", err)
		os.Exit(1)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("error creating dynamic client: %v\n", err)
		os.Exit(1)
	}

	clients := &Clients{
		Clientset: clientset,
		Dynamic:   dynamicClient,
	}
	// Define the HTTP server address and timeout values.
	serverAddr := ":9090"
	readTimeout := 5 * time.Second
	writeTimeout := 10 * time.Second

	// Define the HTTP handler for the /domains endpoint.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Define a context with a timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Get the namespace, network, and subnetwork parameters from the path.
		s := strings.Split(r.URL.Path, "/")
		log.Println("Got request for: ", s)
		switch s[1] {
		case "domains":
			handleDomainsEndpoint(ctx, clients, w, r)
		default:
			http.Error(w, "Invalid URL", http.StatusBadRequest)
		}
	})

	// Create a new HTTP server with the specified address and timeouts.
	server := &http.Server{
		Addr:         serverAddr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Start the HTTP server in a goroutine.
	go func() {
		fmt.Printf("Starting server at %s\n", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Wait for a SIGINT or SIGTERM signal to shut down the server gracefully.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)
	<-sigCh
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
	}
}

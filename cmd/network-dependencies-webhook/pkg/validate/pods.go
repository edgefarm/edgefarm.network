package validate

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/edgefarm/edgefarm.network/internal/common"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

// GetFromSuffixLabelKey returns the value of a label with a given prefix
// Note, that the prefix must be unique otherwise the first match is returned
func GetFromSuffixLabelKey(labels map[string]string, prefix string) (string, error) {
	for k := range labels {
		if strings.Contains(k, prefix) {
			network := ""
			format := prefix + "%s"
			n, err := fmt.Sscanf(k, format, &network)
			if err != nil || n != 1 {
				return "", fmt.Errorf("can't parse label %s", k)
			}
			return network, nil
		}
	}
	return "", fmt.Errorf("label %s not found", prefix)
}

func Pods(config *rest.Config, w http.ResponseWriter, r *http.Request) {
	scheme := runtime.NewScheme()
	codecFactory := serializer.NewCodecFactory(scheme)
	deserializer := codecFactory.UniversalDeserializer()
	admissionReview, err := admissionReviewFromRequest(r, deserializer)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't retrieve admission review from request: %v", err))
		return
	}

	if admissionReview.Request.Operation != v1beta1.Delete {
		writeAllowedResponse(w, admissionReview)
		return
	}

	name := admissionReview.Request.Name
	namespace := admissionReview.Request.Namespace
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't create client: %v", err))
		return
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't get pods: %v", err))
		return
	}
	pod.Labels[common.MarkedForDeletionLabelKey] = common.MarkedForDeletionLabelValue
	pod, err = clientset.CoreV1().Pods(namespace).Update(context.Background(), pod, metav1.UpdateOptions{})
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't update pod: %v", err))
		return
	}

	dynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't create dynamic client: %v", err))
		return
	}
	networkName, err := GetFromSuffixLabelKey(pod.Labels, podNetworkLabelKeyPrefix)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't get network name: %v", err))
		return
	}
	subnetworkname, err := GetFromSuffixLabelKey(pod.Labels, podSubnetworkLabelKeyPrefix)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't get subnetwork name: %v", err))
		return
	}

	streams, err := GetStreams(context.Background(), dynamic, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s,%s=%s", streamNetworkLabelKey, networkName, streamSubnetworkLabelKey, subnetworkname),
	})
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't list streams: %v", err))
		return
	}

	if len(streams.Items) > 0 {
		writeDeniedResponse(w, admissionReview, "cannot delete pod %s: streams still present", pod.Name)
		return
	}

	writeAllowedResponse(w, admissionReview)
}

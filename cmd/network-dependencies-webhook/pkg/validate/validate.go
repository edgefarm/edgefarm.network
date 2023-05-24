package validate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/edgefarm/edgefarm.network/internal/common"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

var streamResource = schema.GroupVersionResource{
	Group:    "streams.network.edgefarm.io",
	Version:  "v1alpha1",
	Resource: "xstreams",
}

func ListStreams(ctx context.Context, client dynamic.Interface, options metav1.ListOptions) ([]unstructured.Unstructured, error) {
	list, err := client.Resource(streamResource).List(ctx, options)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func Validate(config *rest.Config, w http.ResponseWriter, r *http.Request) {
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
	hostname := pod.Spec.NodeName
	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("streams.network.edgefarm.io/node=%s", hostname),
	}
	dynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't create dynamic client: %v", err))
		return
	}
	items, err := ListStreams(context.Background(), dynamic, listOptions)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't list streams: %v", err))
		return
	}

	if len(items) > 0 {
		writeDeniedResponse(w, admissionReview, "cannot delete pod %s: streams still present", pod.Name)
		return
	}

	writeAllowedResponse(w, admissionReview)
}

func admissionReviewFromRequest(r *http.Request, deserializer runtime.Decoder) (*v1beta1.AdmissionReview, error) {
	var body []byte
	if r.Body != nil {
		defer r.Body.Close()
		body = readAll(r.Body)
	}
	admissionReview := &v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, admissionReview); err != nil {
		return nil, err
	}
	return admissionReview, nil
}

func writeAllowedResponse(w http.ResponseWriter, admissionReview *v1beta1.AdmissionReview) {
	admissionReview.Response = &v1beta1.AdmissionResponse{
		UID:     admissionReview.Request.UID,
		Allowed: true,
	}
	writeResponse(w, admissionReview)
}

func writeDeniedResponse(w http.ResponseWriter, admissionReview *v1beta1.AdmissionReview, reason string, args ...interface{}) {
	admissionReview.Response = &v1beta1.AdmissionResponse{
		UID:     admissionReview.Request.UID,
		Allowed: false,
		Result: &metav1.Status{
			Message: fmt.Sprintf(reason, args...),
		},
	}
	writeResponse(w, admissionReview)
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func writeResponse(w http.ResponseWriter, admissionReview *v1beta1.AdmissionReview) {
	w.Header().Set("Content-Type", "application/json")
	responseBytes, err := json.Marshal(admissionReview)
	if err != nil {
		panic(fmt.Sprintf("can't encode response: %v", err))
	}
	w.Write(responseBytes)
}

func readAll(reader io.Reader) []byte {
	var buf []byte
	for {
		chunk := make([]byte, 1024)
		n, err := reader.Read(chunk)
		if err != nil && err != io.EOF {
			panic(fmt.Sprintf("can't read request body: %v", err))
		}
		if n == 0 {
			break
		}
		buf = append(buf, chunk[:n]...)
	}
	return buf
}

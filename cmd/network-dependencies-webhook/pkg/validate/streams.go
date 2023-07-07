package validate

import (
	"context"
	"fmt"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

func Streams(config *rest.Config, w http.ResponseWriter, r *http.Request) {
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

	dynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't create dynamic client: %v", err))
		return
	}
	stream, err := GetStreamByName(context.Background(), dynamic, name)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't get stream: %v", err))
		return
	}
	uid := string(stream.ObjectMeta.UID)
	consumers, err := GetConsumersByDependsOnUidLabel(context.Background(), dynamic, uid)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't get consumer: %v", err))
		return
	}
	if len(consumers.Items) > 0 {
		writeDeniedResponse(w, admissionReview, "cannot delete stream %s: consumers still present", name)
		return
	}
	writeAllowedResponse(w, admissionReview)
}

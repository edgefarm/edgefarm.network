package validate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

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

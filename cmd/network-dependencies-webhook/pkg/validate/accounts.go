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

// Accounts handles the deletion of accounts
// accounts can only be deleted if no user depends on it
func Accounts(config *rest.Config, w http.ResponseWriter, r *http.Request) {
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

	account, err := GetAccountByName(context.Background(), dynamic, name)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't get account: %v", err))
		return
	}

	uid := account.ObjectMeta.UID

	users, err := GetUsersByDependsOnUidLabel(context.Background(), dynamic, string(uid))
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't get user: %v", err))
		return
	}
	if len(users.Items) > 0 {
		writeDeniedResponse(w, admissionReview, "cannot delete account %s: users still present", name)
		return
	}
	writeAllowedResponse(w, admissionReview)
}

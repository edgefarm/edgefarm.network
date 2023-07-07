package validate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	userv1 "github.com/edgefarm/provider-natssecrets/apis/user/v1alpha1"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

func Users(config *rest.Config, w http.ResponseWriter, r *http.Request) {
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

	user, err := GetUserByName(context.Background(), dynamic, name)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't get user: %v", err))
		return
	}
	if val, ok := user.ObjectMeta.Annotations["crossplane.io/composition-resource-name"]; ok {
		if val == systemUserSuffix {
			SystemUser(w, admissionReview, dynamic, user)
			return
		} else if val == systemAccountUserSuffix {
			SysAccountUser(w, admissionReview, dynamic, user)
			return
		}
	}
	// normal user always allowed
	writeAllowedResponse(w, admissionReview)
}

// SystemUser handles the deletion of system users
// system users are allowed to be deleted if there is no ProviderConfig with the UID referenced in the DependsOnUid label
func SystemUser(w http.ResponseWriter, admissionReview *v1beta1.AdmissionReview, dynamic *dynamic.DynamicClient, user *userv1.User) {
	dependsOnUid, err := GetUidFromDependsOnUidLabel(user.ObjectMeta.Labels)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("can't get dependsOnUid from user: %v", err))
		return
	}
	providerConfig, err := GetProviderConfigByUid(context.Background(), dynamic, dependsOnUid)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeAllowedResponse(w, admissionReview)
		} else {
			writeErrorResponse(w, err)
		}
		return
	}
	if providerConfig != nil {
		writeDeniedResponse(w, admissionReview, "cannot delete system user %s: provider config %s still present", user.ObjectMeta.Name, providerConfig.ObjectMeta.Name)
		return
	}
}

// SysAccountUser handles the deletion of sys-account users
// sys-account users are allowed to be deleted if there is no Account that belongs to the user
func SysAccountUser(w http.ResponseWriter, admissionReview *v1beta1.AdmissionReview, dynamic *dynamic.DynamicClient, user *userv1.User) {
	name := user.ObjectMeta.Name
	accountName := strings.TrimSuffix(name, "-"+systemAccountUserSuffix)
	account, err := GetAccountByName(context.Background(), dynamic, accountName)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeAllowedResponse(w, admissionReview)
		} else {
			writeErrorResponse(w, err)
		}
	}
	if account != nil {
		writeDeniedResponse(w, admissionReview, "cannot delete sys-account-user %s: account %s still present", user.ObjectMeta.Name, accountName)
		return
	}
}

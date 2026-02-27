package admission

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func HandleValidate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	var review admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &review); err != nil {
		http.Error(w, "cannot decode admission review", http.StatusBadRequest)
		return
	}

	req := review.Request
	if req == nil {
		http.Error(w, "no admission request", http.StatusBadRequest)
		return
	}

	log.Printf("Admission request UID=%s Namespace=%s Name=%s Operation=%s",
		req.UID, req.Namespace, req.Name, req.Operation)

	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		http.Error(w, "cannot decode pod", http.StatusBadRequest)
		return
	}

	result := ValidatePod(&pod)

	response := admissionv1.AdmissionResponse{
		UID:     req.UID,
		Allowed: result.Allowed,
		Result: &metav1.Status{
			Message: result.Message,
		},
	}

	respReview := admissionv1.AdmissionReview{
		TypeMeta: review.TypeMeta,
		Response: &response,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respReview)
}

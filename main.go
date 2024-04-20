package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var userRights map[string]map[string]bool

func readRightsFromFile(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read rights file: %v", err)
	}
	userRights = make(map[string]map[string]bool)
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.Split(line, ":")
		if len(parts) == 4 {
			user := parts[0]
			right := parts[1] + ":" + parts[2] + ":" + parts[3] // Combines verb, resource, and namespace
			if _, ok := userRights[user]; !ok {
				userRights[user] = make(map[string]bool)
			}
			userRights[user][right] = true
		}
	}
}

func handleSAR(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Received a SubjectAccessReview request")

	var sar authv1.SubjectAccessReview
	if err := json.NewDecoder(req.Body).Decode(&sar); err != nil {
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	allowed := false
	if rights, ok := userRights[sar.Spec.User]; ok {
		rightKey := sar.Spec.ResourceAttributes.Verb + ":" + sar.Spec.ResourceAttributes.Resource + ":" + sar.Spec.ResourceAttributes.Namespace
		if _, allowed = rights[rightKey]; allowed {
			allowed = true
		}
	}

	response := authv1.SubjectAccessReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "authorization.k8s.io/v1beta1",
			Kind:       "SubjectAccessReview",
		},
		Status: authv1.SubjectAccessReviewStatus{
			Allowed: allowed,
		},
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

func main() {
	readRightsFromFile("rights.txt")
	http.HandleFunc("/authorize", handleSAR)
	log.Println("Listening on TCP port 8888...")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}

package admission

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

type ValidationResult struct {
	Allowed bool
	Message string
}

func ValidatePod(pod *corev1.Pod, cfg *PolicyConfig) ValidationResult {
	for _, container := range pod.Spec.Containers {
		// Rule 1: Disallow :latest tag or missing tag
		if cfg.Policies.DenyLatestTag && isLatestOrNoTag(container.Image) {
			return ValidationResult{
				Allowed: false,
				Message: fmt.Sprintf("container %s uses disallowed image tag (latest or no tag)", container.Name),
			}
		}

		// Rule 2: Require CPU & Memory requests and limits
		if cfg.Policies.RequireResources && !hasRequiredResources(container) {
			return ValidationResult{
				Allowed: false,
				Message: fmt.Sprintf("container %s must define cpu and memory requests and limits", container.Name),
			}
		}
	}

	return ValidationResult{
		Allowed: true,
		Message: "pod is valid",
	}
}

func isLatestOrNoTag(image string) bool {
	parts := strings.Split(image, ":")
	if len(parts) == 1 {
		return true // no tag -> implicit latest
	}

	tag := parts[len(parts)-1]
	return tag == "latest"
}

func hasRequiredResources(container corev1.Container) bool {
	req := container.Resources.Requests
	lim := container.Resources.Limits

	if req == nil || lim == nil {
		return false
	}

	if req.Cpu().IsZero() || req.Memory().IsZero() {
		return false
	}

	if lim.Cpu().IsZero() || lim.Memory().IsZero() {
		return false
	}

	return true
}

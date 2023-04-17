package v1alpha1

type ConfigSpec struct {
	WebhookURL string `json:"webhookURL"`
}

type Config struct {
	APIVersion string     `json:"apiVersion"`
	Kind       string     `json:"kind"`
	Spec       ConfigSpec `json:"spec"`
}

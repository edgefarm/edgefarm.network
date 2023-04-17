package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"log2webhook/apis/v1alpha1"

	fnv1alpha1 "github.com/crossplane/crossplane/apis/apiextensions/fn/io/v1alpha1"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

type Datasource struct {
	Composite fnv1alpha1.ObservedComposite `json:"composite"`

	Resources map[string]fnv1alpha1.ObservedResource `json:"resources,omitempty"`
}

func main() {
	// Read the function IO from stdin.
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		failFatal(fnv1alpha1.FunctionIO{}, errors.Wrap(err, "cannot read stdin"))
		return
	}

	j, err := yaml.YAMLToJSON(stdin)
	if err != nil {
		panic(err)
	}

	// Unmarshal the function IO.
	in := fnv1alpha1.FunctionIO{}
	if err = yaml.Unmarshal([]byte(strings.TrimSpace(string(stdin))), &in); err != nil {
		failFatal(fnv1alpha1.FunctionIO{}, errors.Wrap(err, "cannot unmarshal as FunctionIO"))
		return
	}

	out := *(in.DeepCopy())

	// Build the datasource.
	// We simply pass everything under "observed" of the input as the datasource (i.e. similar to values.yaml in helm).
	// For convenience, we convert the observed resources to a map keyed by their name so that they can be referenced
	// in the templates by their name.
	ds := Datasource{}
	ds.Composite = in.Observed.Composite
	ds.Resources = make(map[string]fnv1alpha1.ObservedResource, len(in.Observed.Resources))
	for _, r := range in.Observed.Resources {
		ds.Resources[r.Name] = r
	}

	// Parse the function configuration.
	cfg := v1alpha1.Config{}
	if in.Config != nil {
		if err = yaml.Unmarshal(in.Config.Raw, &cfg); err != nil {
			failFatal(out, errors.Wrap(err, "cannot unmarshal as Config"))
			return
		}
	}

	err = sendRequestBin(cfg.Spec.WebhookURL, bytes.NewReader(j))
	if err != nil {
		failFatal(fnv1alpha1.FunctionIO{}, errors.Wrap(err, "cannot send to RequestBin"))
		return
	}

	// Marshal and write the output.
	b, err := yaml.Marshal(out)
	if err != nil {
		failFatal(out, errors.Wrap(err, "cannot marshal output"))
		return
	}

	fmt.Println(string(b))
}

func failFatal(io fnv1alpha1.FunctionIO, err error) {
	io.Results = append(io.Results, fnv1alpha1.Result{
		Severity: fnv1alpha1.SeverityFatal,
		Message:  err.Error(),
	})
	b, _ := yaml.Marshal(io)
	fmt.Println(string(b))
}

func sendRequestBin(URI string, body io.Reader) error {
	req, err := http.NewRequest("POST", URI, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

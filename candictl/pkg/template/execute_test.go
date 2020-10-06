package template

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestExecuteTemplate(t *testing.T) {
	var data map[string]interface{}

	err := yaml.Unmarshal([]byte(`
nodeIP: "127.0.0.1"
clusterConfiguration:
  kubernetesVersion: "1.15"
  clusterType: "Cloud"
  serviceSubnetCIDR: "127.0.0.1/24"
  podSubnetCIDR: "127.0.0.1/24"
  clusterDomain: "%s.example.com"
extraArgs: {}
`), &data)
	if err != nil {
		t.Errorf("Loading templates error: %v", err)
	}

	_, err = RenderTemplate("/deckhouse/candi/control-plane-kubeadm/", data)
	if err != nil {
		t.Errorf("Rendering templates error: %v", err)
	}
}

func TestExecuteTemplate_DefineAndInclude(t *testing.T) {
	var data map[string]interface{}

	err := yaml.Unmarshal([]byte(`
nodeIP: "127.0.0.1"
`), &data)
	if err != nil {
		t.Errorf("Loading templates error: %v", err)
	}

	rendered, err := RenderTemplate("testdata/execute", data)
	if err != nil {
		t.Errorf("Rendering templates error: %v", err)
	}
	if len(rendered) == 0 {
		t.Errorf("Should render a template, got 0 rendered templates")
	}
	content := rendered[0].Content.String()
	if !strings.Contains(content, "DEFINE GOT 127.0.0.1") {
		t.Errorf("Define and include should work in templates, got '%s'", content)
	}
}
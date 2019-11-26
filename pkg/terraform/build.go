package terraform

import (
	"fmt"

	"github.com/deislabs/porter/pkg/exec/builder"
	"gopkg.in/yaml.v2"
)

const terraformClientVersion = "0.11.11"
const dockerfileLines = `ENV TERRAFORM_VERSION=%s
RUN apt-get update && apt-get install -y wget unzip && \
wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin
`

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig `yaml:"config"`
}

// MixinConfig represents configuration that can be set on the terraform mixin in porter.yaml, e.g.
// mixins:
// - terraform:
//		TODO
type MixinConfig struct {
	BackendConfig map[string]string `yaml:"backendConfig"`
}

// Build installs the terraform cli and runs terraform init,
// with backend config if supplied.
func (m *Mixin) Build() error {
	var input BuildInput
	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(m.Out, dockerfileLines, terraformClientVersion)

	var initCmd = "RUN terraform init"
	var backendConfig = input.Config.BackendConfig
	if len(backendConfig) > 0 {
		initCmd = initCmd + " -backend=true"
		for _, k := range sortKeys(backendConfig) {
			initCmd = initCmd + fmt.Sprintf(" -backend-config=%s=%s", k, backendConfig[k])
		}
		initCmd = initCmd + " -reconfigure"
	}
	fmt.Fprint(m.Out, initCmd)

	return nil
}

package terraform

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMixin_Build(t *testing.T) {
	const buildOutput = `ENV TERRAFORM_VERSION=0.11.11
RUN apt-get update && apt-get install -y wget unzip && \
wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin
RUN terraform init`

	t.Run("build with config", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-with-config.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)

		err = m.Build()
		require.NoError(t, err, "build failed")

		wantOutput := buildOutput + ` -backend=true` +
			` -backend-config=access_key=myaccesskey` +
			` -backend-config=container_name=mycontainer` +
			` -backend-config=key=my.tfstate` +
			` -backend-config=storage_account_name=mystorageaccount` +
			` -reconfigure`
		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
	})

	t.Run("build without config", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-without-config.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)

		err = m.Build()
		require.NoError(t, err, "build failed")

		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, buildOutput, gotOutput)
	})
}

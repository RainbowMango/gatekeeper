/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/open-policy-agent/gatekeeper/v3/apis"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var cfg *rest.Config

func TestMain(m *testing.M) {
	t := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "..", "vendor", "github.com", "open-policy-agent", "frameworks", "constraint", "deploy", "crds.yaml"),
			filepath.Join("..", "..", "..", "..", "config", "crd", "bases"),
		},
		ErrorIfCRDPathMissing: true,
	}
	///TODO(ritazh): remove when vap is GAed in k/k
	args := t.ControlPlane.GetAPIServer().Configure()
	args.Append("runtime-config", "api/all=true")
	args.Append("feature-gates", "ValidatingAdmissionPolicy=true")

	if err := apis.AddToScheme(scheme.Scheme); err != nil {
		log.Fatal(err)
	}

	var err error
	if cfg, err = t.Start(); err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	if err = t.Stop(); err != nil {
		log.Printf("error while trying to stop server: %v", err)
	}
	os.Exit(code)
}

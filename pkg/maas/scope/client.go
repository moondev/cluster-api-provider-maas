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

package scope

import (
	"os"

	"github.com/canonical/gomaasclient/client"

)

// NewMaasClient creates a new MaaS client for a given session
// TODO (looking up on Env really the besT? though it is kind of what EC2 does

func NewMaasClient(_ *ClusterScope) *client.Client {

	maasEndpoint := os.Getenv("MAAS_ENDPOINT")
	if maasEndpoint == "" {
		panic("missing env MAAS_ENDPOINT; e.g: MAAS_ENDPOINT=http://10.11.130.11:5240/MAAS")
	}

	maasAPIKey := os.Getenv("MAAS_API_KEY")
	if maasAPIKey == "" {
		panic("missing env MAAS_API_KEY; e.g: MAAS_API_KEY=x:y:z>")
	}

	maasClient, err := client.GetClient(maasEndpoint, maasAPIKey, "2.0")
	if err != nil {
		return nil
	}
	return maasClient
}

// NewSpectroMaasClient creates a new Spectro MaaS client for a given session
func NewSpectroMaasClient(_ *ClusterScope) maasclient.ClientSetInterface {


	maasEndpoint := os.Getenv("MAAS_ENDPOINT")
	if maasEndpoint == "" {
		panic("missing env MAAS_ENDPOINT; e.g: MAAS_ENDPOINT=http://10.11.130.11:5240/MAAS")
	}

	maasAPIKey := os.Getenv("MAAS_API_KEY")
	if maasAPIKey == "" {
		panic("missing env MAAS_API_KEY; e.g: MAAS_API_KEY=x:y:z>")
	}

	maasClient := client.NewAuthenticatedClientSet(maasEndpoint, maasAPIKey)
	return maasClient
}

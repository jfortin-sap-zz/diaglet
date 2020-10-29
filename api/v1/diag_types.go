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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DiagSpec defines the desired state of Diag
type DiagSpec struct {
	// +kubebuilder:validation:Enum=true;false
	ProbeAPIServer bool `json:"probeAPIServer,omitempty"`
	// +kubebuilder:validation:Enum=true;false
	ProbeShootLastOperation bool `json:"probeShootLastOperation,omitempty"`
	// +kubebuilder:validation:Enum=true;false
	ProbeShootConditions bool `json:"probeShootConditions,omitempty"`
	// +kubebuilder:validation:Enum=true;false
	ProbeWorkerNodes bool `json:"probeWorkerNodes,omitempty"`
	// +kubebuilder:validation:Enum=true;false
	ProbeControlPlane bool `json:"probeControlPlane,omitempty"`
}

// DiagStatus defines the observed state of Diag
type DiagStatus struct {
	ShootDiagStatus []string `json:"shootDiagStatus,omitempty"`
}

// +kubebuilder:object:root=true

// Diag is the Schema for the diags API
// +kubebuilder:printcolumn:name="DiagStatus",type=string,JSONPath=`.status.shootDiagStatus`
// +kubebuilder:subresource:status
type Diag struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DiagSpec   `json:"spec,omitempty"`
	Status DiagStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DiagList contains a list of Diag
type DiagList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Diag `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Diag{}, &DiagList{})
}

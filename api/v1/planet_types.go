/*
Copyright 2025.

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

const PlanetGalaxyLabel = "kubesphere.io/galaxy"

// PlanetSpec defines the desired state of Planet
type PlanetSpec struct {
	// Name of the planet
	Name string `json:"name,omitempty"`
	// Diameter in kilometers
	DiameterKm float64 `json:"diameterKm,omitempty"`
	// Does the planet support life?
	HasLife bool `json:"hasLife,omitempty"`
	// List of moons
	Moons []string `json:"moons,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,path=planets,shortName=planet
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Planet is the Schema for the planets API
// Label: kubesphere.io/galaxy - references the GalaxySpec
type Planet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PlanetSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// PlanetList contains a list of Planet
type PlanetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Planet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Planet{}, &PlanetList{})
}

func (in *Planet) GetGalaxy() string {
	if in.Labels == nil {
		return ""
	}
	return in.Labels[PlanetGalaxyLabel]
}

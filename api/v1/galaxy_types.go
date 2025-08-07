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

type GalaxyPlanetSpec struct {
	Name       string   `json:"name,omitempty"`
	DiameterKm float64  `json:"diameterKm,omitempty"`
	HasLife    bool     `json:"hasLife,omitempty"`
	Moons      []string `json:"moons,omitempty"`
}

type GalaxySpec struct {
	Name    string             `json:"name,omitempty"`
	Planets []GalaxyPlanetSpec `json:"planets,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,path=galaxies,shortName=galaxy
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Galaxy is the Schema for the galaxies API
type Galaxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GalaxySpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// GalaxyList contains a list of Galaxy
type GalaxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Galaxy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Galaxy{}, &GalaxyList{})
}

func (in *GalaxyPlanetSpec) ToPlanetSpec() PlanetSpec {
	spec := PlanetSpec{
		Name:       in.Name,
		DiameterKm: in.DiameterKm,
		HasLife:    in.HasLife,
		Moons:      in.Moons,
	}
	return spec
}

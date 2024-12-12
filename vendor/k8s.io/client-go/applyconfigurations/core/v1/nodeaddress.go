/*
Copyright The Kubernetes Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	corev1 "k8s.io/api/core/v1"
)

// NodeAddressApplyConfiguration represents a declarative configuration of the NodeAddress type for use
// with apply.
type NodeAddressApplyConfiguration struct {
	Type    *corev1.NodeAddressType `json:"type,omitempty"`
	Address *string                 `json:"address,omitempty"`
}

// NodeAddressApplyConfiguration constructs a declarative configuration of the NodeAddress type for use with
// apply.
func NodeAddress() *NodeAddressApplyConfiguration {
	return &NodeAddressApplyConfiguration{}
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *NodeAddressApplyConfiguration) WithType(value corev1.NodeAddressType) *NodeAddressApplyConfiguration {
	b.Type = &value
	return b
}

// WithAddress sets the Address field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Address field is set to the value of the last call.
func (b *NodeAddressApplyConfiguration) WithAddress(value string) *NodeAddressApplyConfiguration {
	b.Address = &value
	return b
}

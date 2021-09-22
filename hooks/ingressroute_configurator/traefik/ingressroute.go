package traefik

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressRouteSpec is a specification for a IngressRouteSpec resource.
type IngressRouteSpec struct {
	Routes      []Route  `json:"routes"`
	EntryPoints []string `json:"entryPoints,omitempty"`
}

// Route contains the set of routes.
type Route struct {
	Match string `json:"match"`
	// +kubebuilder:validation:Enum=Rule
	Kind        string          `json:"kind"`
	Priority    int             `json:"priority,omitempty"`
	Services    []Service       `json:"services,omitempty"`
	Middlewares []MiddlewareRef `json:"middlewares,omitempty"`
}



// LoadBalancerSpec can reference either a Kubernetes Service object (a load-balancer of servers),
// or a TraefikService object (a traefik load-balancer of services).
type LoadBalancerSpec struct {
	// Name is a reference to a Kubernetes Service object (for a load-balancer of servers),
	// or to a TraefikService object (service load-balancer, mirroring, etc).
	// The differentiation between the two is specified in the Kind field.
	Name string `json:"name"`
	// +kubebuilder:validation:Enum=Service;TraefikService
	Kind      string          `json:"kind,omitempty"`
	Namespace string          `json:"namespace,omitempty"`

	// Port and all the fields below are related to a servers load-balancer,
	// and therefore should only be specified when Name references a Kubernetes Service.

	Port               int32        `json:"port,omitempty"`
	Scheme             string                      `json:"scheme,omitempty"`
	Strategy           string                      `json:"strategy,omitempty"`
	PassHostHeader     *bool                       `json:"passHostHeader,omitempty"`
	ServersTransport   string                      `json:"serversTransport,omitempty"`

	// Weight should only be specified when Name references a TraefikService object
	// (and to be precise, one that embeds a Weighted Round Robin).
	Weight *int `json:"weight,omitempty"`
}

// Service defines an upstream to proxy traffic.
type Service struct {
	LoadBalancerSpec  `json:",inline"`
}

// MiddlewareRef is a ref to the Middleware resources.
type MiddlewareRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:storageversion

// IngressRoute is an Ingress CRD specification.
type IngressRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec IngressRouteSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressRouteList is a list of IngressRoutes.
type IngressRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []IngressRoute `json:"items"`
}

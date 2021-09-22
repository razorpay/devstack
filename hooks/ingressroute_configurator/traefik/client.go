package traefik

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

// GetIngressRoute gets the ingressroute for the provided arguments
func GetIngressRoute(clientset *kubernetes.Clientset,namespace string , resourceName string) (*IngressRoute,error){
	var ingressRoute IngressRoute
	path := fmt.Sprintf(INGRESS_ABS_PATH,namespace)
	data, err := clientset.RESTClient().
				 Get().
				 AbsPath(path).
				 Name(resourceName).
				 DoRaw(context.TODO())

	if err != nil {
		return &IngressRoute{},err
	}
	jerr := json.Unmarshal(data, &ingressRoute)
	if err != nil {
		return &IngressRoute{}, jerr
	}
	return &ingressRoute,nil
}

// PatchIngressRoute patches the ingress route with the provided arguments
func PatchIngressRoute(clientset *kubernetes.Clientset,namespace string , resourceName string , ingressroute IngressRoute) (*IngressRoute,error) {
	var updatedIR IngressRoute
	data , err := json.Marshal(&ingressroute)
	if err != nil {
		return &IngressRoute{},err
	}
	path := fmt.Sprintf(INGRESS_ABS_PATH,namespace)

	upsertData , err := clientset.RESTClient().
		Patch(types.MergePatchType).
		AbsPath(path).
		Body(data).
		Name(resourceName).
		DoRaw(context.TODO())
	if err != nil  {
		return &IngressRoute{},err
	}
	err = json.Unmarshal(upsertData, &updatedIR)
	if err != nil {
		return &IngressRoute{}, err
	}
	return &updatedIR,nil
}


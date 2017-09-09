package converter

import (
	"reflect"

	"github.com/projectcalico/libcalico-go/lib/api"
	"github.com/projectcalico/libcalico-go/lib/backend/k8s"
	backendConverter "github.com/projectcalico/libcalico-go/lib/converter"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type policyConverter struct {
}

// NewPolicyConverter Constructor for policyConverter
func NewPolicyConverter() Converter {
	return &policyConverter{}
}

func (p *policyConverter) Convert(k8sObj interface{}) (interface{}, error) {
	if reflect.TypeOf(k8sObj) != reflect.TypeOf(&v1beta1.NetworkPolicy{}) {
		log.Fatalf("can not convert object %#v to calico policy. Object is not of type *v1beta1.NetworkPolicy", k8sObj)
	}

	np := k8sObj.(*v1beta1.NetworkPolicy)

	var policyConverter k8s.Converter
	kvpair, err := policyConverter.NetworkPolicyToPolicy(np)
	if err != nil {
		return nil, err
	}

	var backendConverter backendConverter.PolicyConverter
	policy, err := backendConverter.ConvertKVPairToAPI(kvpair)
	if err != nil {
		return nil, err
	}
	calicoPolicy := policy.(*api.Policy)
	return *calicoPolicy, err
}

// GetKey returns name of Policy as its key.  For Policies created by this controller
// and backed by NetworkPolicy objects, the name is of the format
// `knp.default.namespace.name`.
func (p *policyConverter) GetKey(obj interface{}) string {
	if reflect.TypeOf(obj) != reflect.TypeOf(api.Policy{}) {
		log.Fatalf("can not construct key for object %#v. Object is not of type api.Policy", obj)
	}
	policy := obj.(api.Policy)
	return policy.Metadata.Name
}

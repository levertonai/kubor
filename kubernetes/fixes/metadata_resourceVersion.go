package fixes

import (
	"fmt"
	"github.com/echocat/kubor/common"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"reflect"
)

func init() {
	registerUpdateFix(updateFixForResourceVersionIsAbsent)
}

func updateFixForResourceVersionIsAbsent(original unstructured.Unstructured, target *unstructured.Unstructured) error {
	if !groupVersionKindMatches(&original, target) {
		return nil
	}

	resourceVersion := common.GetObjectPathValue(original.Object, "metadata", "resourceVersion")
	if resourceVersion == nil {
		return nil
	}
	sResourceVersion, ok := resourceVersion.(string)
	if !ok || sResourceVersion == "" {
		return nil
	}

	metadata, ok := target.Object["metadata"]
	if !ok {
		metadata = map[string]interface{}{}
		target.Object["metadata"] = metadata
	}
	mMetadata, ok := metadata.(map[string]interface{})
	if !ok {
		return fmt.Errorf("'metadata' property of target does already exists but is not of type map[string]interface{} it is %v", reflect.TypeOf(metadata))
	}

	mMetadata["resourceVersion"] = sResourceVersion

	return nil
}
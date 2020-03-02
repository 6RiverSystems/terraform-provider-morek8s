package morek8s

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/streaming"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

func getNamespace(obj metav1.Object) string {
	ns := obj.GetNamespace()
	if ns == "" {
		return "default"
	}
	return ns
}

func buildID(obj metav1.Object) string {
	return fmt.Sprintf("%s/%s", getNamespace(obj), obj.GetName())
}

func nameChanged(old, new, meta interface{}) bool {
	// ignore errors. If any exist they should be catched by validator function.
	oldU, _ := expandResourceFromStr(old.(string))
	newU, _ := expandResourceFromStr(new.(string))
	return buildID(&newU) != buildID(&oldU)
}

func decode(in string) (unstructured.Unstructured, error) {
	yamlDecoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	reader := json.YAMLFramer.NewFrameReader(ioutil.NopCloser(bytes.NewReader([]byte(in))))
	d := streaming.NewDecoder(reader, yamlDecoder)

	var objects []runtime.Object

	for {
		obj, _, err := d.Decode(nil, nil)
		if err != nil {
			if err == io.EOF {
				break
			}
			return unstructured.Unstructured{}, err
		}
		objects = append(objects, obj)
	}

	if len(objects) == 0 {
		return unstructured.Unstructured{}, nil
	}

	if cnt := len(objects); cnt > 1 {
		err := fmt.Errorf("unexpected number of resources: %d, expected 1", cnt)
		return unstructured.Unstructured{}, err
	}

	return *objects[0].(*unstructured.Unstructured), nil
}

func expandResourceFromStr(in string) (unstructured.Unstructured, error) {
	if len(in) == 0 {
		return unstructured.Unstructured{}, nil
	}
	u, err := decode(in)
	if err != nil {
		return unstructured.Unstructured{}, err
	}
	return u, nil
}

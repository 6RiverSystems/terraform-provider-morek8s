package morek8s

import (
	"fmt"
)

func validateK8sFile(v interface{}, key string) (ws []string, es []error) {
	s, ok := v.(string)
	if !ok {
		es = []error{fmt.Errorf("%s: must be a non-nil string", key)}
		return
	}

	_, err := expandResourceFromStr(s)
	if err != nil {
		es = []error{fmt.Errorf("%s: must be valid k8s JSON or YAML", key), err}
		return
	}
	return
}

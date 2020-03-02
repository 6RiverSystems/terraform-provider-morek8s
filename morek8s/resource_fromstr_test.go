package morek8s

import (
	"context"
	"fmt"
	"testing"

	"github.com/6RiverSystems/terraform-provider-helpers/kubernetes"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAccResourceFromStr_basic(t *testing.T) {
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMoreK8sFromStrDestroy,
		Steps: []resource.TestStep{
			{
				// Create secret and check it exists
				Config: testAccMoreK8sFromStrConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMoreK8sFromStrExists("morek8s_from_str.test"),
				),
			},
		},
	})
}

func testAccCheckMoreK8sFromStrDestroy(s *terraform.State) error {
	cl := testAccProvider.Meta().(client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "morek8s_from_str" {
			continue
		}
		namespace, name, err := kubernetes.IDParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		data, ok := rs.Primary.Attributes["data"]
		if !ok {
			return fmt.Errorf("Not found data attribute")
		}
		u, err := expandResourceFromStr(data)
		if err != nil {
			return err
		}

		found := unstructured.Unstructured{}
		found.SetGroupVersionKind(u.GetObjectKind().GroupVersionKind())
		err = cl.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &found)

		if err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			return err
		}

		return fmt.Errorf("resource %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckMoreK8sFromStrExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		cl := testAccProvider.Meta().(client.Client)

		namespace, name, err := kubernetes.IDParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		data, ok := rs.Primary.Attributes["data"]
		if !ok {
			return fmt.Errorf("Not found data attribute")
		}
		u, err := expandResourceFromStr(data)
		if err != nil {
			return err
		}

		found := unstructured.Unstructured{}
		found.SetGroupVersionKind(u.GetObjectKind().GroupVersionKind())
		err = cl.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &found)
		return err
	}
}

func testAccMoreK8sFromStrConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "morek8s_from_str" "test" {
	data =<<EOF
{
		"apiVersion": "v1",
		"data": {
				"password": "MWYyZDFlMmU2N2Rm",
				"username": "YWRtaW4="
		},
		"kind": "Secret",
		"metadata": {
				"name": "%s",
				"namespace": "default"
		},
		"type": "Opaque"
}
EOF
}
`, name)
}

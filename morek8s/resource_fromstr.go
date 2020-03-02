package morek8s

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func resourceFromStr() *schema.Resource {
	return &schema.Resource{
		Create: resourceFromStrCreate,
		Read:   resourceFromStrRead,
		Update: resourceFromStrUpdate,
		Delete: resourceFromStrDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
			Delete: schema.DefaultTimeout(10 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"data": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Resource content",
				Required:     true,
				ValidateFunc: validateK8sFile,
			},
		},

		CustomizeDiff: customdiff.All(
			// Set ForceNew to true if namespace or name is changed
			customdiff.ForceNewIfChange("data", namespacedNameChanged),
		),
	}
}

func resourceFromStrCreate(d *schema.ResourceData, m interface{}) error {
	data := d.Get("data").(string)
	u, err := expandResourceFromStr(data)
	if err != nil {
		return err
	}

	cl := m.(client.Client)

	log.Printf("[INFO] Creating new k8s resource: %#v", u)

	if err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err := cl.Create(context.TODO(), &u)
		if err == nil {
			return nil
		}
		e := fmt.Errorf("failed to create k8s resource %#v, %s", u, err)
		return resource.RetryableError(e)
	}); err != nil {
		return err
	}

	log.Printf("[INFO] Submitted new k8s resource: %#v", u)

	d.SetId(buildID(&u))

	return resourceFromStrRead(d, m)
}

func resourceFromStrRead(d *schema.ResourceData, m interface{}) error {
	data := d.Get("data").(string)
	u, err := expandResourceFromStr(data)
	if err != nil {
		return err
	}

	key, err := idToKey(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading k8s resource %s", d.Id())

	found := unstructured.Unstructured{}
	found.SetGroupVersionKind(u.GetObjectKind().GroupVersionKind())
	cl := m.(client.Client)
	err = cl.Get(context.TODO(), key, &found)

	if err != nil {
		// If the resource does not exist, inform Terraform. We want to immediately
		// return here to prevent further processing
		if errors.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}

	log.Printf("[INFO] Received k8s resource: %#v", found)
	return nil
}

func resourceFromStrUpdate(d *schema.ResourceData, m interface{}) error {
	data := d.Get("data").(string)
	u, err := expandResourceFromStr(data)
	if err != nil {
		return err
	}

	key, err := idToKey(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating k8s resource %s", d.Id())

	found := unstructured.Unstructured{}
	found.SetGroupVersionKind(u.GetObjectKind().GroupVersionKind())
	cl := m.(client.Client)
	err = cl.Get(context.TODO(), key, &found)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}

	// update resource version
	u.SetResourceVersion(found.GetResourceVersion())
	err = cl.Update(context.TODO(), &u)
	if err != nil {
		return err
	}

	return resourceFromStrRead(d, m)
}

func resourceFromStrDelete(d *schema.ResourceData, m interface{}) error {
	data := d.Get("data").(string)
	u, err := expandResourceFromStr(data)
	if err != nil {
		return err
	}

	key, err := idToKey(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting k8s resource %s", d.Id())

	found := unstructured.Unstructured{}
	found.SetGroupVersionKind(u.GetObjectKind().GroupVersionKind())
	cl := m.(client.Client)
	err = cl.Get(context.TODO(), key, &found)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	err = cl.Delete(context.TODO(), &found)

	if err != nil {
		return err
	}

	if err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		err := cl.Get(context.TODO(), key, &found)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		e := fmt.Errorf("k8s resource %s still exists", d.Id())
		return resource.RetryableError(e)
	}); err != nil {
		return err
	}

	log.Printf("[INFO] k8s resource deleted %s", d.Id())
	return nil
}

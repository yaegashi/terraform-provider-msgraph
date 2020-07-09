package msgraph

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	P "github.com/yaegashi/msgraph.go/ptr"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func resourceApplicationPasswordResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationPasswordCreate,
		Read:   resourceApplicationPasswordRead,
		Delete: resourceApplicationPasswordDelete,
		Schema: map[string]*schema.Schema{
			"application_id":  {Type: schema.TypeString, Required: true, ForceNew: true},
			"display_name":    {Type: schema.TypeString, Required: true, ForceNew: true},
			"start_date_time": {Type: schema.TypeString, ValidateFunc: validation.IsRFC3339Time, Optional: true, Computed: true, ForceNew: true},
			"end_date_time":   {Type: schema.TypeString, ValidateFunc: validation.IsRFC3339Time, Optional: true, Computed: true, ForceNew: true},
			"secret_text":     {Type: schema.TypeString, Computed: true, Sensitive: true},
		},
	}
}

type resourceApplicationPassword struct {
	graph    *graph
	resource *schema.ResourceData
}

func newResourceApplicationPassword(d *schema.ResourceData, meta interface{}) *resourceApplicationPassword {
	return &resourceApplicationPassword{
		graph:    newGraph(meta),
		resource: d,
	}
}

func resourceApplicationPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	return newResourceApplicationPassword(d, meta).create()
}

func resourceApplicationPasswordRead(d *schema.ResourceData, meta interface{}) error {
	return newResourceApplicationPassword(d, meta).read()
}

func resourceApplicationPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	return newResourceApplicationPassword(d, meta).delete()
}

func (r *resourceApplicationPassword) graphSet(pc *msgraph.PasswordCredential) {
	r.resource.Set("display_name", pc.DisplayName)
	r.resource.Set("start_date_time", pc.StartDateTime.Format(time.RFC3339))
	r.resource.Set("end_date_time", pc.EndDateTime.Format(time.RFC3339))
	// pc.SecretText appears only when creating.  We should preserve the old value when it's nil.
	if pc.SecretText != nil {
		r.resource.Set("secret_text", pc.SecretText)
	}
}

func (r *resourceApplicationPassword) graphGet() *msgraph.PasswordCredential {
	pc := &msgraph.PasswordCredential{}
	if val, ok := r.resource.GetOkExists("display_name"); ok {
		pc.DisplayName = P.CastString(val)
	}
	if val, ok := r.resource.GetOkExists("start_date_time"); ok {
		t, _ := time.Parse(time.RFC3339, val.(string))
		pc.StartDateTime = &t
	}
	if val, ok := r.resource.GetOkExists("end_date_time"); ok {
		t, _ := time.Parse(time.RFC3339, val.(string))
		pc.EndDateTime = &t
	}
	return pc
}

func (r *resourceApplicationPassword) create() error {
	applicationID := r.resource.Get("application_id").(string)
	newPC := r.graphGet()
	pc, err := r.graph.applicationAddPassword(
		applicationID,
		&msgraph.ApplicationAddPasswordRequestParameter{PasswordCredential: newPC},
	)
	if err != nil {
		return err
	}
	r.resource.SetId(string(*pc.KeyID))
	r.graphSet(pc)
	return nil
}

func (r *resourceApplicationPassword) read() error {
	id := r.resource.Id()
	applicationID := r.resource.Get("application_id").(string)
	application, err := r.graph.applicationRead(applicationID)
	if err != nil {
		return err
	}
	for i := range application.PasswordCredentials {
		if string(*application.PasswordCredentials[i].KeyID) != id {
			continue
		}
		r.graphSet(&application.PasswordCredentials[i])
		return nil
	}
	r.resource.SetId("")
	return nil
}

func (r *resourceApplicationPassword) delete() error {
	id := msgraph.UUID(r.resource.Id())
	applicationID := r.resource.Get("application_id").(string)
	err := r.graph.applicationRemovePassword(
		applicationID,
		&msgraph.ApplicationRemovePasswordRequestParameter{KeyID: &id},
	)
	if err != nil {
		return err
	}
	r.resource.SetId("")
	return nil
}

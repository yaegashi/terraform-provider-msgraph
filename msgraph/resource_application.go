package msgraph

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	P "github.com/yaegashi/msgraph.go/ptr"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func resourceApplicationResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_id":                       {Type: schema.TypeString, Computed: true},
			"display_name":                 {Type: schema.TypeString, Required: true},
			"sign_in_audience":             {Type: schema.TypeString, Optional: true, ValidateFunc: validation.StringInSlice([]string{"AzureADMyOrg", "AzureADMultipleOrgs", "AzureADandPersonalMicrosoftAccount"}, false)},
			"identifier_uris":              {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"home_page_url":                {Type: schema.TypeString, Optional: true},
			"logout_url":                   {Type: schema.TypeString, Optional: true},
			"redirect_uris":                {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"enable_id_token_issuance":     {Type: schema.TypeBool, Optional: true},
			"enable_access_token_issuance": {Type: schema.TypeBool, Optional: true},
			"api": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accept_mapped_claims":      {Type: schema.TypeBool, Optional: true, Default: false},
						"known_client_applications": {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsUUID}},
						"oauth2_permission_scope": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"admin_consent_description":  {Type: schema.TypeString, Required: true},
									"admin_consent_display_name": {Type: schema.TypeString, Required: true},
									"id":                         {Type: schema.TypeString, Computed: true, Optional: true, ValidateFunc: validation.IsUUID},
									"is_enabled":                 {Type: schema.TypeBool, Optional: true, Default: true},
									"origin":                     {Type: schema.TypeString, Optional: true},
									"type":                       {Type: schema.TypeString, Required: true, ValidateFunc: validation.StringInSlice([]string{"User", "Admin"}, false)},
									"user_consent_description":   {Type: schema.TypeString, Optional: true},
									"user_consent_display_name":  {Type: schema.TypeString, Optional: true},
									"value":                      {Type: schema.TypeString, Required: true},
								},
							},
						},
						"pre_authorized_applications": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_id":                   {Type: schema.TypeString, Required: true, ValidateFunc: validation.IsUUID},
									"delegated_permission_ids": {Type: schema.TypeSet, Required: true, MinItems: 1, Elem: &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsUUID}},
								},
							},
						},
						"requested_access_token_version": {Type: schema.TypeInt, Optional: true, Default: 1, ValidateFunc: validation.IntInSlice([]int{1, 2})},
					},
				},
			},
			"app_role": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":                   {Type: schema.TypeString, Computed: true, Optional: true, ValidateFunc: validation.IsUUID},
						"allowed_member_types": {Type: schema.TypeSet, Required: true, MinItems: 1, Elem: &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringInSlice([]string{"User", "Application"}, false)}},
						"description":          {Type: schema.TypeString, Required: true},
						"display_name":         {Type: schema.TypeString, Required: true},
						"is_enabled":           {Type: schema.TypeBool, Optional: true, Default: true},
						"value":                {Type: schema.TypeString, Required: true},
					},
				},
			},
			"required_resource_access": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_app_id": {Type: schema.TypeString, Required: true, ValidateFunc: validation.IsUUID},
						"resource_access": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id":   {Type: schema.TypeString, Required: true, ValidateFunc: validation.IsUUID},
									"type": {Type: schema.TypeString, Required: true, ValidateFunc: validation.StringInSlice([]string{"Scope", "Role"}, false)},
								},
							},
						},
					},
				},
			},
		},
	}
}

type resourceApplication struct {
	graph    *graph
	resource *schema.ResourceData
}

func newResourceApplication(d *schema.ResourceData, meta interface{}) *resourceApplication {
	return &resourceApplication{
		graph:    newGraph(meta),
		resource: d,
	}
}

func resourceApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	return newResourceApplication(d, meta).create()
}

func resourceApplicationRead(d *schema.ResourceData, meta interface{}) error {
	return newResourceApplication(d, meta).read()
}

func resourceApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	return newResourceApplication(d, meta).update()
}

func resourceApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	return newResourceApplication(d, meta).delete()
}

func (r *resourceApplication) graphSet(application *msgraph.Application) {
	r.resource.Set("app_id", application.AppID)
	r.resource.Set("display_name", application.DisplayName)
	r.resource.Set("sign_in_audience", application.SignInAudience)
	r.resource.Set("identifier_uris", application.IdentifierUris)
	r.resource.Set("home_page_url", application.Web.HomePageURL)
	r.resource.Set("logout_url", application.Web.LogoutURL)
	r.resource.Set("redirect_uris", application.Web.RedirectUris)
	r.resource.Set("enable_id_token_issuance", application.Web.ImplicitGrantSettings.EnableIDTokenIssuance)
	r.resource.Set("enable_access_token_issuance", application.Web.ImplicitGrantSettings.EnableAccessTokenIssuance)
	if application.API != nil {
		vPermissionScopes := make([]map[string]interface{}, len(application.API.Oauth2PermissionScopes))
		for i, permissionScope := range application.API.Oauth2PermissionScopes {
			vPermissionScopes[i] = map[string]interface{}{
				"admin_consent_description":  permissionScope.AdminConsentDescription,
				"admin_consent_display_name": permissionScope.AdminConsentDisplayName,
				"id":                         permissionScope.ID,
				"is_enabled":                 permissionScope.IsEnabled,
				"origin":                     permissionScope.Origin,
				"type":                       permissionScope.Type,
				"user_consent_description":   permissionScope.UserConsentDescription,
				"user_consent_display_name":  permissionScope.UserConsentDisplayName,
				"value":                      permissionScope.Value,
			}
		}
		vPreAuthorizedApplications := make([]map[string]interface{}, len(application.API.PreAuthorizedApplications))
		for i, preAuthorizedApplication := range application.API.PreAuthorizedApplications {
			vPreAuthorizedApplications[i] = map[string]interface{}{
				"app_id":                   preAuthorizedApplication.AppID,
				"delegated_permission_ids": preAuthorizedApplication.DelegatedPermissionIDs,
			}
		}
		vAPI := map[string]interface{}{
			"accept_mapped_claims":           application.API.AcceptMappedClaims,
			"known_client_applications":      application.API.KnownClientApplications,
			"oauth2_permission_scope":        vPermissionScopes,
			"pre_authorized_applications":    vPreAuthorizedApplications,
			"requested_access_token_version": application.API.RequestedAccessTokenVersion,
		}
		r.resource.Set("api", []map[string]interface{}{vAPI})
	} else {
		r.resource.Set("api", []map[string]interface{}{})
	}
	vAppRoles := make([]map[string]interface{}, len(application.AppRoles))
	for i, appRole := range application.AppRoles {
		vAppRoles[i] = map[string]interface{}{
			"id":                   appRole.ID,
			"allowed_member_types": appRole.AllowedMemberTypes,
			"description":          appRole.Description,
			"display_name":         appRole.DisplayName,
			"is_enabled":           appRole.IsEnabled,
			"value":                appRole.Value,
		}
	}
	r.resource.Set("app_role", vAppRoles)
	vRequiredResourceAccesses := make([]map[string]interface{}, len(application.RequiredResourceAccess))
	for i1, requiredResourceAccess := range application.RequiredResourceAccess {
		vResourceAccesses := make([]map[string]interface{}, len(requiredResourceAccess.ResourceAccess))
		for i2, resourceAccess := range requiredResourceAccess.ResourceAccess {
			vResourceAccesses[i2] = map[string]interface{}{
				"id":   resourceAccess.ID,
				"type": resourceAccess.Type,
			}
		}
		vRequiredResourceAccesses[i1] = map[string]interface{}{
			"resource_app_id": requiredResourceAccess.ResourceAppID,
			"resource_access": vResourceAccesses,
		}
	}
	r.resource.Set("required_resource_access", vRequiredResourceAccesses)
}

func (r *resourceApplication) graphGet() *msgraph.Application {
	application := &msgraph.Application{Web: &msgraph.WebApplication{ImplicitGrantSettings: &msgraph.ImplicitGrantSettings{}}}
	if val, ok := r.resource.GetOkExists("app_id"); ok {
		application.AppID = P.CastString(val)
	}
	if val, ok := r.resource.GetOkExists("display_name"); ok {
		application.DisplayName = P.CastString(val)
	}
	if val, ok := r.resource.GetOkExists("sign_in_audience"); ok {
		application.SignInAudience = P.CastString(val)
	}
	for _, val := range r.resource.Get("identifier_uris").([]interface{}) {
		application.IdentifierUris = append(application.IdentifierUris, val.(string))
	}
	if val, ok := r.resource.GetOkExists("home_page_url"); ok {
		application.Web.HomePageURL = P.CastString(val)
	}
	if val, ok := r.resource.GetOkExists("logout_url"); ok {
		application.Web.LogoutURL = P.CastString(val)
	}
	for _, val := range r.resource.Get("redirect_uris").([]interface{}) {
		application.Web.RedirectUris = append(application.Web.RedirectUris, val.(string))
	}
	if val, ok := r.resource.GetOkExists("enable_id_token_issuance"); ok {
		application.Web.ImplicitGrantSettings.EnableIDTokenIssuance = P.CastBool(val)
	}
	if val, ok := r.resource.GetOkExists("enable_access_token_issuance"); ok {
		application.Web.ImplicitGrantSettings.EnableAccessTokenIssuance = P.CastBool(val)
	}
	if v0, ok := r.resource.GetOkExists("api"); ok {
		vAPIs := v0.(*schema.Set).List()
		if len(vAPIs) > 0 {
			vAPI := vAPIs[0].(map[string]interface{})
			application.API = &msgraph.APIApplication{
				AcceptMappedClaims:          P.CastBool(vAPI["accept_mapped_claims"]),
				RequestedAccessTokenVersion: P.CastInt(vAPI["requested_access_token_version"]),
			}
			vKnownClientApplications := vAPI["known_client_applications"].(*schema.Set).List()
			if len(vKnownClientApplications) > 0 {
				application.API.KnownClientApplications = make([]msgraph.UUID, len(vKnownClientApplications))
				for i1, v1 := range vKnownClientApplications {
					application.API.KnownClientApplications[i1] = msgraph.UUID(v1.(string))
				}
			} else {
				application.API.KnownClientApplications = nil
				application.API.SetAdditionalData("knownClientApplications", []struct{}{})
			}
			if len(application.API.KnownClientApplications) == 0 {
				application.API.KnownClientApplications = nil
				application.API.SetAdditionalData("knownClientApplications", []struct{}{})
			}
			vPermissionScopes := vAPI["oauth2_permission_scope"].(*schema.Set).List()
			if len(vPermissionScopes) > 0 {
				application.API.Oauth2PermissionScopes = make([]msgraph.PermissionScope, len(vPermissionScopes))
				for i1, v1 := range vPermissionScopes {
					vPermissionScope := v1.(map[string]interface{})
					id := ""
					if vID, ok := vPermissionScope["id"]; ok {
						id = vID.(string)
					}
					if id == "" {
						id = uuid.New().String()
					}
					application.API.Oauth2PermissionScopes[i1] = msgraph.PermissionScope{
						AdminConsentDescription: P.CastString(vPermissionScope["admin_consent_description"]),
						AdminConsentDisplayName: P.CastString(vPermissionScope["admin_consent_display_name"]),
						ID:                      (*msgraph.UUID)(&id),
						IsEnabled:               P.CastBool(vPermissionScope["is_enabled"]),
						Type:                    P.CastString(vPermissionScope["type"]),
						UserConsentDescription:  P.CastString(vPermissionScope["user_consent_description"]),
						UserConsentDisplayName:  P.CastString(vPermissionScope["user_consent_display_name"]),
						Value:                   P.CastString(vPermissionScope["value"]),
						// Origin is for internal use?
						// Origin: P.CastString(vPermissionScope["origin"]),
					}
				}
			} else {
				application.API.Oauth2PermissionScopes = nil
				application.API.SetAdditionalData("oauth2PermissionScopes", []struct{}{})
			}
			vPreAuthorizedApplications := vAPI["pre_authorized_applications"].(*schema.Set).List()
			if len(vPreAuthorizedApplications) > 0 {
				application.API.PreAuthorizedApplications = make([]msgraph.PreAuthorizedApplication, len(vPreAuthorizedApplications))
				for i1, v1 := range vPreAuthorizedApplications {
					vPreAuthorizedApplication := v1.(map[string]interface{})
					vDelegatedPermissionIDs := vPreAuthorizedApplication["delegated_permission_ids"].(*schema.Set).List()
					application.API.PreAuthorizedApplications[i1] = msgraph.PreAuthorizedApplication{
						AppID:                  P.CastString(vPreAuthorizedApplication["app_id"]),
						DelegatedPermissionIDs: make([]string, len(vDelegatedPermissionIDs)),
					}
					for i2, v2 := range vDelegatedPermissionIDs {
						application.API.PreAuthorizedApplications[i1].DelegatedPermissionIDs[i2] = v2.(string)
					}
				}
			} else {
				application.API.PreAuthorizedApplications = nil
				application.API.SetAdditionalData("preAuthorizedApplications", []struct{}{})
			}
		}
	}
	if application.API == nil {
		application.SetAdditionalData("api", nil)
	}
	if v0, ok := r.resource.GetOkExists("app_role"); ok {
		vAppRoles := v0.(*schema.Set).List()
		application.AppRoles = make([]msgraph.AppRole, len(vAppRoles))
		for i1, v1 := range vAppRoles {
			vAppRole := v1.(map[string]interface{})
			id := ""
			if vID, ok := vAppRole["id"]; ok {
				id = vID.(string)
			}
			if id == "" {
				id = uuid.New().String()
			}
			vAllowedMemberTypes := vAppRole["allowed_member_types"].(*schema.Set).List()
			application.AppRoles[i1] = msgraph.AppRole{
				ID:                 (*msgraph.UUID)(&id),
				AllowedMemberTypes: make([]string, len(vAllowedMemberTypes)),
				Description:        P.CastString(vAppRole["description"]),
				DisplayName:        P.CastString(vAppRole["display_name"]),
				IsEnabled:          P.CastBool(vAppRole["is_enabled"]),
				Value:              P.CastString(vAppRole["value"]),
			}
			for i2, v2 := range vAllowedMemberTypes {
				application.AppRoles[i1].AllowedMemberTypes[i2] = v2.(string)
			}
		}
	}
	if len(application.AppRoles) == 0 {
		application.AppRoles = nil
		application.SetAdditionalData("appRoles", []struct{}{})
	}
	if v0, ok := r.resource.GetOkExists("required_resource_access"); ok {
		vRequiredResourceAccesses := v0.(*schema.Set).List()
		application.RequiredResourceAccess = make([]msgraph.RequiredResourceAccess, len(vRequiredResourceAccesses))
		for i1, v1 := range vRequiredResourceAccesses {
			vRequiredResourceAccess := v1.(map[string]interface{})
			vResourceAccesses := vRequiredResourceAccess["resource_access"].(*schema.Set).List()
			application.RequiredResourceAccess[i1] = msgraph.RequiredResourceAccess{
				ResourceAppID:  P.CastString(vRequiredResourceAccess["resource_app_id"]),
				ResourceAccess: make([]msgraph.ResourceAccess, len(vResourceAccesses)),
			}
			for i2, v2 := range vResourceAccesses {
				vResourceAccess := v2.(map[string]interface{})
				application.RequiredResourceAccess[i1].ResourceAccess[i2] = msgraph.ResourceAccess{
					ID:   (*msgraph.UUID)(P.CastString(vResourceAccess["id"])),
					Type: P.CastString(vResourceAccess["type"]),
				}
			}
		}
	}
	if len(application.RequiredResourceAccess) == 0 {
		application.RequiredResourceAccess = nil
		application.SetAdditionalData("requiredResourceAccess", []struct{}{})
	}
	return application
}

func (r *resourceApplication) create() error {
	newApplication := r.graphGet()
	application, err := r.graph.applicationCreate(newApplication)
	if err != nil {
		return err
	}
	r.resource.SetId(*application.ID)
	r.graphSet(application)
	return nil
}

func (r *resourceApplication) read() error {
	application, err := r.graph.applicationRead(r.resource.Id())
	if err != nil {
		r.resource.SetId("")
		if errRes, ok := err.(*msgraph.ErrorResponse); ok {
			if errRes.StatusCode() == http.StatusNotFound {
				return nil
			}
		}
		return err
	}
	r.graphSet(application)
	return nil
}

func (r *resourceApplication) update() error {
	newApplication := r.graphGet()
	application, err := r.graph.applicationUpdate(r.resource.Id(), newApplication)
	if err != nil {
		return err
	}
	r.graphSet(application)
	return nil
}

func (r *resourceApplication) delete() error {
	r.graph.applicationDelete(r.resource.Id())
	r.resource.SetId("")
	return nil
}

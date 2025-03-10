package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NewTeamTemplateResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: teamTemplateCreate,
		ReadContext:   teamTemplateRead,
		UpdateContext: teamTemplateUpdate,
		DeleteContext: teamTemplateDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Identifier of the template.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the template.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the template.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"template_data": {
				Description: "The template data in JSON format.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					jsonStr := val.(string)
					var js json.RawMessage
					if err := json.Unmarshal([]byte(jsonStr), &js); err != nil {
						errs = append(errs, fmt.Errorf("%q contains invalid JSON: %s", key, err))
					}
					return
				},
			},
			"team_id": {
				Description: "The team ID the template belongs to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Description: "The template type (defaults to 'issue').",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "issue",
			},
			"sort_order": {
				Description: "The sort order of the template.",
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func teamTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*LinearClient)

	name := d.Get("name").(string)
	templateData := d.Get("template_data").(string)
	teamID := d.Get("team_id").(string)
	templateType := d.Get("type").(string)

	// Prepare variables for GraphQL mutation
	variables := map[string]interface{}{
		"templateCreateInput": map[string]interface{}{
			"name":         name,
			"templateData": templateData,
			"teamId":       teamID,
			"type":         templateType,
		},
	}

	// Add optional fields if provided
	if description, ok := d.GetOk("description"); ok {
		variables["templateCreateInput"].(map[string]interface{})["description"] = description.(string)
	}

	if sortOrder, ok := d.GetOk("sort_order"); ok {
		variables["templateCreateInput"].(map[string]interface{})["sortOrder"] = sortOrder.(float64)
	}

	// Create the GraphQL mutation
	mutation := `
		mutation TemplateCreate($templateCreateInput: TemplateCreateInput!) {
			templateCreate(input: $templateCreateInput) {
				success
				template {
					id
					sortOrder
				}
			}
		}
	`

	// Execute the mutation
	response := make(map[string]interface{})
	if err := client.gqlClient.Execute(mutation, variables, &response); err != nil {
		return diag.FromErr(fmt.Errorf("error creating template: %v", err))
	}

	// Extract the template ID from the response
	templateCreate, ok := response["templateCreate"].(map[string]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf("unexpected response format"))
	}

	success, ok := templateCreate["success"].(bool)
	if !ok || !success {
		return diag.FromErr(fmt.Errorf("template creation failed"))
	}

	template, ok := templateCreate["template"].(map[string]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf("template data missing in response"))
	}

	templateID, ok := template["id"].(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("template ID missing in response"))
	}

	d.SetId(templateID)

	// If sortOrder was computed, save it to state
	if sortOrder, ok := template["sortOrder"].(float64); ok {
		d.Set("sort_order", sortOrder)
	}

	return teamTemplateRead(ctx, d, meta)
}

func teamTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*LinearClient)

	// Prepare variables for GraphQL query
	variables := map[string]interface{}{
		"templateId": d.Id(),
	}

	// Create the GraphQL query
	query := `
		query Template($templateId: String!) {
			template(id: $templateId) {
				id
				name
				description
				templateData
				team {
					id
				}
				type
				sortOrder
			}
		}
	`

	// Execute the query
	response := make(map[string]interface{})
	if err := client.gqlClient.Execute(query, variables, &response); err != nil {
		return diag.FromErr(fmt.Errorf("error reading template: %v", err))
	}

	// Extract the template data from the response
	templateData, ok := response["template"].(map[string]interface{})
	if !ok {
		return diag.Errorf("template not found or invalid response format")
	}

	// Set template data to resource state
	if err := d.Set("name", templateData["name"]); err != nil {
		return diag.FromErr(err)
	}

	if description, ok := templateData["description"]; ok && description != nil {
		if err := d.Set("description", description); err != nil {
			return diag.FromErr(err)
		}
	}

	if templateDataJson, ok := templateData["templateData"]; ok && templateDataJson != nil {
		if err := d.Set("template_data", templateDataJson); err != nil {
			return diag.FromErr(err)
		}
	}

	if team, ok := templateData["team"].(map[string]interface{}); ok {
		if err := d.Set("team_id", team["id"]); err != nil {
			return diag.FromErr(err)
		}
	}

	if templateType, ok := templateData["type"]; ok && templateType != nil {
		if err := d.Set("type", templateType); err != nil {
			return diag.FromErr(err)
		}
	}

	if sortOrder, ok := templateData["sortOrder"]; ok && sortOrder != nil {
		if err := d.Set("sort_order", sortOrder); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func teamTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*LinearClient)

	// Prepare variables for GraphQL mutation
	variables := map[string]interface{}{
		"templateUpdateInput": map[string]interface{}{
			"id": d.Id(),
		},
	}

	// Add fields that have changed
	if d.HasChange("name") {
		variables["templateUpdateInput"].(map[string]interface{})["name"] = d.Get("name").(string)
	}

	if d.HasChange("description") {
		variables["templateUpdateInput"].(map[string]interface{})["description"] = d.Get("description").(string)
	}

	if d.HasChange("template_data") {
		variables["templateUpdateInput"].(map[string]interface{})["templateData"] = d.Get("template_data").(string)
	}

	if d.HasChange("type") {
		variables["templateUpdateInput"].(map[string]interface{})["type"] = d.Get("type").(string)
	}

	if d.HasChange("sort_order") {
		variables["templateUpdateInput"].(map[string]interface{})["sortOrder"] = d.Get("sort_order").(float64)
	}

	// Create the GraphQL mutation
	mutation := `
		mutation TemplateUpdate($templateUpdateInput: TemplateUpdateInput!) {
			templateUpdate(input: $templateUpdateInput) {
				success
			}
		}
	`

	// Execute the mutation
	response := make(map[string]interface{})
	if err := client.gqlClient.Execute(mutation, variables, &response); err != nil {
		return diag.FromErr(fmt.Errorf("error updating template: %v", err))
	}

	// Check if update was successful
	templateUpdate, ok := response["templateUpdate"].(map[string]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf("unexpected response format"))
	}

	success, ok := templateUpdate["success"].(bool)
	if !ok || !success {
		return diag.FromErr(fmt.Errorf("template update failed"))
	}

	return teamTemplateRead(ctx, d, meta)
}

func teamTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*LinearClient)

	// Prepare variables for GraphQL mutation
	variables := map[string]interface{}{
		"templateDeleteInput": map[string]interface{}{
			"id": d.Id(),
		},
	}

	// Create the GraphQL mutation
	mutation := `
		mutation TemplateDelete($templateDeleteInput: TemplateDeleteInput!) {
			templateDelete(input: $templateDeleteInput) {
				success
			}
		}
	`

	// Execute the mutation
	response := make(map[string]interface{})
	if err := client.gqlClient.Execute(mutation, variables, &response); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting template: %v", err))
	}

	// Check if deletion was successful
	templateDelete, ok := response["templateDelete"].(map[string]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf("unexpected response format"))
	}

	success, ok := templateDelete["success"].(bool)
	if !ok || !success {
		return diag.FromErr(fmt.Errorf("template deletion failed"))
	}

	// Clear the ID
	d.SetId("")

	return nil
} 
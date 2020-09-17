package squadcast

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-provider-squadcast/types"
)

var serviceRes types.ServiceRes

const servicePath string = "/services"

func resourceSquadcastService() *schema.Resource {
	return &schema.Resource{
		Create: resourceSquadcastServiceCreate,
		Read:   resourceSquadcastServiceRead,
		Update: resourceSquadcastServiceUpdate,
		Delete: resourceSquadcastServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sid": {
				Type:        schema.TypeString,
				Description: "Unique service data ID",
				Computed:    true,
				Required:    false,
				// ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				// ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Service description",
				Default:     "Service created via Terraform provider",
				Optional:    true,
				// ForceNew: true,
			},
			"escalation_policy_id": {
				Type:        schema.TypeString,
				Description: "Escalation policy id to be associated with the service",
				Required:    true,
				// ForceNew:    true,
			},
			"email_prefix": {
				Type:        schema.TypeString,
				Description: `Email prefix for the service`,
				Required:    true,
				// ForceNew:    true,
			},
		},
	}
}

func resourceSquadcastServiceCreate(resourceData *schema.ResourceData, configMetaData interface{}) error {
	var squadcastConfig = configMetaData.(Config)

	if squadcastConfig.AccessToken != "" {
		log.Printf("[INFO] Access token is not set")
	}

	var serviceName = resourceData.Get("name").(string)
	var serviceDescription = resourceData.Get("description").(string)
	var escalationPolicyId = resourceData.Get("escalation_policy_id").(string)
	var emailPrefix = resourceData.Get("email_prefix").(string)

	log.Printf("[INFO] Creating new service: %s", serviceName)

	reqBody, err := json.Marshal(map[string]string{
		"name":                 serviceName,
		"description":          serviceDescription,
		"escalation_policy_id": escalationPolicyId,
		"email_prefix":         emailPrefix,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, getAPIFullURL(servicePath), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", squadcastConfig.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	responseData, err := ioutil.ReadAll(res.Body)

	if res.StatusCode > 299 {
		return errors.New(string(responseData))
	}

	json.Unmarshal(responseData, &serviceRes)

	resourceData.Set("name", serviceRes.Data.Name)
	resourceData.Set("sid", serviceRes.Data.Id)
	resourceData.SetId(serviceRes.Data.Name)

	log.Printf("[INFO] Successfully created service: %s", serviceName)

	return nil
}

func resourceSquadcastServiceRead(resourceData *schema.ResourceData, configMetaData interface{}) error {
	var serviceName = resourceData.Get("name").(string)
	var squadcastConfig = configMetaData.(Config)

	if squadcastConfig.AccessToken != "" {
		log.Printf("[INFO] Access token is not set")
	}

	reqBody, err := json.Marshal(map[string]string{})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, getAPIFullURL(servicePath)+"?name="+serviceName, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", squadcastConfig.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	responseData, err := ioutil.ReadAll(res.Body)

	if res.StatusCode > 299 {
		return errors.New(string(responseData))
	}

	json.Unmarshal(responseData, &serviceRes)

	return nil
}

func resourceSquadcastServiceUpdate(resourceData *schema.ResourceData, configMetaData interface{}) error {
	var squadcastConfig = configMetaData.(Config)

	if squadcastConfig.AccessToken != "" {
		log.Printf("[INFO] Access token is not set")
	}

	var serviceName = resourceData.Get("name").(string)
	var serviceDescription = resourceData.Get("description").(string)
	var escalationPolicyId = resourceData.Get("escalation_policy_id").(string)
	var emailPrefix = resourceData.Get("email_prefix").(string)
	var serviceID = resourceData.Get("sid").(string)

	log.Printf("[INFO] Updating service: %s", serviceName)

	reqBody, err := json.Marshal(map[string]string{
		"name":                 serviceName,
		"description":          serviceDescription,
		"escalation_policy_id": escalationPolicyId,
		"email_prefix":         emailPrefix,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, getAPIFullURL(servicePath)+"/"+serviceID, bytes.NewBuffer(reqBody)) // serviceID
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", squadcastConfig.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	responseData, err := ioutil.ReadAll(res.Body)

	if res.StatusCode > 299 {
		return errors.New(string(responseData))
	}

	json.Unmarshal(responseData, &serviceRes)
	log.Printf("[INFO] Successfully updated service: %s", serviceName)
	return nil
}

func resourceSquadcastServiceDelete(resourceData *schema.ResourceData, configMetaData interface{}) error {

	var squadcastConfig = configMetaData.(Config)

	if squadcastConfig.AccessToken != "" {
		log.Printf("[INFO] Access token is not set")
	}

	var serviceID = resourceData.Get("sid").(string)

	log.Printf("[INFO] Deleting service: %s", serviceID)

	reqBody, err := json.Marshal(map[string]string{})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, getAPIFullURL(servicePath)+"/"+serviceID, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", squadcastConfig.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	responseData, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(responseData, &serviceRes)

	log.Printf("[INFO] Successfully deleted service: %s", serviceID)
	return nil
}

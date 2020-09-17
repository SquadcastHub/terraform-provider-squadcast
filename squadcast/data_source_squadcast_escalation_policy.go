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

const escalationPolicyPath string = "/escalation-policies"

var escalationPolicyRes types.EscalationPolicyRes

func dataSourceSquadcastEscalationPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSquadcastEscalationPolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				// Sensitive:   true,
				Description: "Escalation policy name",
			},
			"id": {
				Type:        schema.TypeString,
				Description: "ObjectId of the escalation policy",
				Optional:    true,
			},
		},
	}
}

func dataSourceSquadcastEscalationPolicyRead(resourceData *schema.ResourceData, configMetaData interface{}) error {
	var escalationPolicyName = resourceData.Get("name").(string)
	var squadcastConfig = configMetaData.(Config)

	if squadcastConfig.AccessToken != "" {
		log.Printf("[INFO] Access token is not set")
	}

	reqBody, err := json.Marshal(map[string]string{})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, getAPIFullURL(escalationPolicyPath)+"?name="+escalationPolicyName, bytes.NewBuffer(reqBody))
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

	json.Unmarshal(responseData, &escalationPolicyRes)
	resourceData.Set("id", escalationPolicyRes.Data[0].Id)
	resourceData.SetId(escalationPolicyRes.Data[0].Name)

	return nil
}

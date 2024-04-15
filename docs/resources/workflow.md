---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "squadcast_workflow Resource - terraform-provider-squadcast"
subcategory: ""
description: |-
  
---

# squadcast_workflow (Resource)



## Example Usage

```terraform
resource "squadcast_workflow" "example_workflow" {
   title = "test workflow"
   description = "Test workflow description"
   owner_id = "63bfabae865e9c93cd31756e"
   enabled = true
   trigger = "incident_triggered"
   # TODO: Needs to accomodate to the new structure
   filters {
      fields {
         value = "P1"
      }
      type = "priority_is"
   }
   entity_owner {
      type = "user" 
      id = "63209531af0f36245bfac82f"
   }
   tags {
      key = "tagKey"
      value = "tagValue"
      color = "#000000"
   }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `entity_owner` (Block List, Min: 1) The entity owner of the workflow (see [below for nested schema](#nestedblock--entity_owner))
- `owner_id` (String) The ID of the user who owns the workflow
- `title` (String) The title of the workflow
- `trigger` (String) The trigger for the workflow

### Optional

- `description` (String) The description of the workflow
- `enabled` (Boolean) Whether the workflow is enabled or not
- `filters` (Block List, Max: 1) The filters to be applied on the workflow (see [below for nested schema](#nestedblock--filters))
- `tags` (Block List) The tags to be applied on the workflow (see [below for nested schema](#nestedblock--tags))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--entity_owner"></a>
### Nested Schema for `entity_owner`

Required:

- `type` (String)

Read-Only:

- `id` (String) The ID of this resource.


<a id="nestedblock--filters"></a>
### Nested Schema for `filters`

Required:

- `condition` (String) Condition to be applied on the filters (and / or)

Optional:

- `filters` (Block List) (see [below for nested schema](#nestedblock--filters--filters))

<a id="nestedblock--filters--filters"></a>
### Nested Schema for `filters.filters`

Optional:

- `condition` (String) Condition to be applied on the filters (and / or)
- `filters` (Block List) (see [below for nested schema](#nestedblock--filters--filters--filters))
- `type` (String)
- `value` (String)

<a id="nestedblock--filters--filters--filters"></a>
### Nested Schema for `filters.filters.filters`

Optional:

- `key` (String)
- `type` (String)
- `value` (String)




<a id="nestedblock--tags"></a>
### Nested Schema for `tags`

Required:

- `color` (String)
- `key` (String)
- `value` (String)
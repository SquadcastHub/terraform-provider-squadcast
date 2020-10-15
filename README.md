Terraform Provider
==================

- Website: https://www.terraform.io
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------
-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

Using the provider
----------------------
```hcl
terraform {
  required_providers {
    squadcast = {
      source  = "SquadcastHub/squadcast"
    }
  }
}

provider "squadcast" {
  squadcast_token = "YOUR-SQUADCAST-TOKEN"
}
```

Developing the Provider
---------------------------
TODO

Acceptance test prerequisites
-----------------------------
`make testacc`

### Squadcast personal refresh token
You will need to create a [personal refresh token](https://app.squadcast.com) 
Once the token has been created, it must be exported in your environment as `squadcast_token`.


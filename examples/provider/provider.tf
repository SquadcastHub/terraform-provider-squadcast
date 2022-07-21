provider "squadcast" {
  # Hard-coding credentials into any Terraform configuration is not recommended
  # refresh_token and region can also be passed via environment variables (SQUADCAST_REFRESH_TOKEN and SQUADCAST_REGION)
  refresh_token = "YOUR-SQUADCAST-TOKEN"
  region        = "us"
}

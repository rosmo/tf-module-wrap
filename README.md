# Terraform module wrapper

Read any Terraform module and wraps all its parameters into an object variable and also generates
a passthrough module source.

## Usage

```sh
$ tfmodulewrap 
  -add-defaults
        Add default values
  -ignore-vars string
        Variables to ignore
  -module-path string
        Path containing the module
  -module-var string
        Variable for the module configuration
```

## Example

Using it on [Cloud Foundation Fabric](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric)'s [`net-address`](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/master/modules/net-address) module:

```sh
$ tfmodulewrap modules/net-address -add-defaults=true
````

Result:

```terraform
variable "net_address" {
  type = object({
    # Map of internal addresses used for Private Service Access.
    psa_addresses = optional(map(object({
      address       = string
      network       = string
      prefix_length = number
      description   = optional(string, "Terraform managed.")
      name          = optional(string)
      })), {
    })
    # Map of internal addresses used for Private Service Connect.
    psc_addresses = optional(map(object({
      address          = optional(string)
      description      = optional(string, "Terraform managed.")
      name             = optional(string)
      network          = optional(string)
      region           = optional(string)
      subnet_self_link = optional(string)
      service_attachment = optional(object({ # so we can safely check if service_attachemnt != null in for_each
        psc_service_attachment_link = string
        global_access               = optional(bool)
      }))
      })), {
    })
    # Map of external addresses, keyed by name.
    external_addresses = optional(map(object({
      region      = string
      description = optional(string, "Terraform managed.")
      ipv6 = optional(object({
        endpoint_type = string
      }))
      labels     = optional(map(string), {})
      name       = optional(string)
      subnetwork = optional(string) # for IPv6
      tier       = optional(string)
      })), {
    })
    # List of global addresses to create.
    global_addresses = optional(map(object({
      description = optional(string, "Terraform managed.")
      ipv6        = optional(map(string)) # To be left empty for ipv6
      name        = optional(string)
      })), {
    })
    # Map of internal addresses to create, keyed by name.
    internal_addresses = optional(map(object({
      region      = string
      subnetwork  = string
      address     = optional(string)
      description = optional(string, "Terraform managed.")
      ipv6        = optional(map(string)) # To be left empty for ipv6
      labels      = optional(map(string))
      name        = optional(string)
      purpose     = optional(string)
      })), {
    })
    # Map of internal addresses used for HPA VPN over Cloud Interconnect.
    ipsec_interconnect_addresses = optional(map(object({
      region        = string
      address       = string
      network       = string
      description   = optional(string, "Terraform managed.")
      name          = optional(string)
      prefix_length = number
      })), {
    })
    # PSC network attachments, names as keys.
    network_attachments = optional(map(object({
      subnet_self_link      = string
      automatic_connection  = optional(bool, false)
      description           = optional(string, "Terraform-managed.")
      producer_accept_lists = optional(list(string))
      producer_reject_lists = optional(list(string))
      })), {
    })
    # Project where the addresses will be created.
    project_id = string
  })
}

module "net_address" {
  source = "/Users/taneli/demos/cloud-foundation-fabric/modules/net-address"

  external_addresses           = var.net_address.external_addresses
  global_addresses             = var.net_address.global_addresses
  internal_addresses           = var.net_address.internal_addresses
  ipsec_interconnect_addresses = var.net_address.ipsec_interconnect_addresses
  network_attachments          = var.net_address.network_attachments
  project_id                   = var.net_address.project_id
  psa_addresses                = var.net_address.psa_addresses
  psc_addresses                = var.net_address.psc_addresses
}
```
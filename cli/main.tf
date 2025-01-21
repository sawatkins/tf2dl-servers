module "tf2_server_us" {
  source = "./modules/tf2_server_module"

  region               = var.region_us
  ami                  = var.ami_us
  instance_type        = var.instance_type
  key_name             = var.key_name
  security_group_id    = var.security_group_id_us
  iam_instance_profile = var.iam_instance_profile
  elastic_ip           = var.eip_us
  name                 = "tf2_server_us"
}

module "tf2_server_eu" {
  source = "./modules/tf2_server_module"

  region               = var.region_eu
  ami                  = var.ami_eu
  instance_type        = var.instance_type
  key_name             = var.key_name
  security_group_id    = var.security_group_id_eu
  iam_instance_profile = var.iam_instance_profile
  elastic_ip           = var.eip_eu
  name                 = "tf2_server_eu"
}


output "tf2_server_us_id" {
  value = module.tf2_server_us.instance_id
}

output "tf2_server_us_public_ip" {
  value = module.tf2_server_us.instance_public_ip
}

output "tf2_server_us_public_dns" {
  value = module.tf2_server_us.instance_public_dns
}


output "tf2_server_eu_id" {
  value = module.tf2_server_eu.instance_id
}

output "tf2_server_eu_public_ip" {
  value = module.tf2_server_eu.instance_public_ip
}

output "tf2_server_eu_public_dns" {
  value = module.tf2_server_eu.instance_public_dns
}

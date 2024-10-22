
# module "tf2_server_jump_01" {
#   source = "./modules/tf2_server_module"

#   region               = var.region
#   instance_type        = var.instance_type
#   key_name             = var.key_name
#   security_group_id    = var.security_group_id
#   iam_instance_profile = var.iam_instance_profile
# }

module "tf2_server_red" {
  source = "./modules/tf2_server_module"

  region               = var.region
  instance_type        = var.instance_type
  key_name             = var.key_name
  security_group_id    = var.security_group_id
  iam_instance_profile = var.iam_instance_profile
  elastic_ip           = var.eip_red
  name                 = "tf2_server_red"
}


// output variables
# output "tf2_server_jump_01_id" {
#   value = module.tf2_server_jump_01.instance_id
# }

# output "tf2_server_jump_01_public_ip" {
#   value = module.tf2_server_jump_01.instance_public_ip
# }

# output "tf2_server_jump_01_public_dns" {
#   value = module.tf2_server_jump_01.instance_public_dns
# }

output "tf2_server_red_id" {
  value = module.tf2_server_red.instance_id
}

output "tf2_server_red_public_ip" {
  value = module.tf2_server_red.instance_public_ip
}

output "tf2_server_red_public_dns" {
  value = module.tf2_server_red.instance_public_dns
}

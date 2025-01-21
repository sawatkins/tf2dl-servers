variable "region_us" {
  type = string
  default = "us-west-1"
}

variable "region_eu" {
  type = string
  default = "eu-central-1"
}

variable "ami_us" {
  type = string
  description = "AMI of base tf2 server in us"
}

variable "ami_eu" {
  type = string
  description = "AMI of base tf2 server in eu"
}

variable "instance_type" {
  type = string
  default = "t3a.micro"
}

variable "key_name" {
  type = string
  description = "The ssh key pair to connect to the ec2 instances"
}

variable "security_group_id" {
  type = string
  description = "AWS security group for the ec2 instances"
}

variable "iam_instance_profile" {
  type = string
  description = "AWS IAM instance profile for the ec2 instances"
}

variable "eip_us" {
  type = string
  description = "Elastic ip for the us tf2 server"
}

variable "eip_eu" {
  type = string
  description = "Elastic ip for the eu tf2 server"
}



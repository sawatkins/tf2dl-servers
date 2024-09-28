variable "region" {
  type = string
  default = "us-west-1"
}

variable "instance_type" {
  type = string
  default = "t2.micro"
}

variable "key_name" {
  type = string
}

variable "security_group_id" {
  type = string
}

variable "iam_instance_profile" {
  type = string
}

provider "aws" {
  region = var.region
}

resource "aws_instance" "tf2_server" {
  ami           = "ami-0fda60cefceeaa4d3" // Amazon Linux 2 AMI
  instance_type = var.instance_type
  key_name      = var.key_name

  vpc_security_group_ids = [var.security_group_id]
  iam_instance_profile   = var.iam_instance_profile

  root_block_device {
    volume_size = 20
    volume_type = "gp3"
  }

  user_data = <<-EOF
              #!/bin/bash
              sudo yum update -y
              sudo yum install -y docker htop
              sudo systemctl enable docker
              sudo systemctl start docker
              sudo usermod -aG docker ec2-user
              newgrp docker
              python3 -m ensurepip --upgrade
              mkdir -p /home/ec2-user/maps
              EOF

  tags = {
    Name = "upfast-tf2-server"
  }
}

// output variables
output "instance_id" {
  value = aws_instance.tf2_server.id
}

output "instance_public_ip" {
  value = aws_instance.tf2_server.public_ip
}

output "instance_public_dns" {
  value = aws_instance.tf2_server.public_dns
}

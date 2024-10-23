provider "aws" {
  region = var.region
}

resource "aws_instance" "tf2_server" {
  ami           = "ami-00fd4a75e141e98d5" # debian ami (prev: "ami-0fda60cefceeaa4d3" amzn linux 2)
  instance_type = var.instance_type
  key_name      = var.key_name

  vpc_security_group_ids = [var.security_group_id]
  iam_instance_profile   = var.iam_instance_profile

  root_block_device {
    volume_size = 15
    volume_type = "gp3"
  }

  user_data = <<-EOF
              #!/bin/bash
              sudo apt-get update
              sudo apt-get upgrade -y
              sudo apt-get install -y htop
              python3 -m ensurepip --upgrade
              EOF

  tags = {
    Name = var.name
  }
}

resource "aws_eip_association" "tf2_server_eip_assoc" {
  instance_id   = aws_instance.tf2_server.id
  allocation_id = var.elastic_ip
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

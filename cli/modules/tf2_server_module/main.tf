provider "aws" {
  region = var.region
}

resource "aws_instance" "tf2_server" {
  ami           = var.ami # base tf2 server image i previously created
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
              sudo apt-get update
              sudo apt-get upgrade -y
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

provider "aws" {
  region = "us-west-1"
}

resource "aws_instance" "tf2_server" {
  ami           = "ami-0fda60cefceeaa4d3" // Amazon Linux 2 AMI
  instance_type = "t2.micro"
  key_name      = "x1-wsl-ubuntu22.04"

  vpc_security_group_ids = ["sg-00eb70ad5a661788f"]
  iam_instance_profile   = "upfast-read-s3"  

  root_block_device {
    volume_size = 20
    volume_type = "gp3"
  }

  user_data = <<-EOF
              #!/bin/bash
              sudo yum update -y
              sudo yum install -y docker htop
              sudo systemctl start docker
              sudo systemctl enable docker
              python3 -m ensurepip --upgrade
              pip3 install boto3 botocore
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
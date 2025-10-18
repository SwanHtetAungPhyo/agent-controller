output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "vpc_cidr" {
  description = "CIDR block of the VPC"
  value       = aws_vpc.main.cidr_block
}

output "public_subnet_ids" {
  description = "IDs of public subnets"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "IDs of private subnets"
  value       = aws_subnet.private[*].id
}

output "nat_gateway_id" {
  description = "ID of NAT Gateway"
  value       = aws_nat_gateway.main.id
}

output "web_sg_id" {
  description = "ID of web security group"
  value       = aws_security_group.web.id
}

output "app_sg_id" {
  description = "ID of application security group"
  value       = aws_security_group.app.id
}

output "web_instance_id" {
  description = "ID of web server instance"
  value       = aws_instance.web.id
}

output "web_instance_public_ip" {
  description = "Public IP of web server"
  value       = aws_instance.web.public_ip
}

output "web_instance_private_ip" {
  description = "Private IP of web server"
  value       = aws_instance.web.private_ip
}

output "app_instance_id" {
  description = "ID of application server instance"
  value       = aws_instance.app.id
}

output "app_instance_private_ip" {
  description = "Private IP of application server"
  value       = aws_instance.app.private_ip
}

output "web_server_url" {
  description = "Web server URL"
  value       = "http://${aws_instance.web.public_ip}"
}

output "ssh_web_server" {
  description = "SSH command for web server"
  value       = "ssh -i ~/.ssh/${var.key_pair_name}.pem ec2-user@${aws_instance.web.public_ip}"
}

output "ssh_app_server" {
  description = "SSH command for app server (via bastion)"
  value       = "ssh -i ~/.ssh/${var.key_pair_name}.pem -J ec2-user@${aws_instance.web.public_ip} ec2-user@${aws_instance.app.private_ip}"
}
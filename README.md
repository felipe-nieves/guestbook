# Guestbook 

## TODO:

EVERYTHING

## 💡 Idea

1. Develop a Go application with input validation and DynamoDB integration.
2. Use Terraform to define and create the necessary AWS infrastructure:
    * DynamoDB table
    * ECR repository
    * ECS cluster, task definition, and service
    * ALB and target group
3. Containerize the application, push the image to ECR, and deploy it to ECS.
4. Test it for high-availability 
5. Create some sort of dashboard for it using Grafana
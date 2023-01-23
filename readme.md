# Health Checker

This application allows you to upload a file into S3 as your input and checks the health status of the URLs every 10
minutes. Application listens on port 8080.

## Variables

The following variables is required when building the image to push to ECR:

| Name           | Description                                                                                                  |
|----------------|--------------------------------------------------------------------------------------------------------------|
| `AWS_ROLE_ARN` | This is the role to assume for the CICD pipeline to have permissions to push to ECR. This is a secret value. |
| `PROJECT_NAME` | The project name as specified inside your infrastructure code                                                |


#!/bin/bash

# Update the package index
sudo apt-get update

# Install necessary packages for Docker
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common

# Add Docker’s official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

# Add the Docker repository
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

# Update the package index again
sudo apt-get update

# Install Docker
sudo apt-get install -y docker-ce

# Start Docker
sudo systemctl start docker

# Enable Docker to start on boot
sudo systemctl enable docker

# Install Docker Compose
DOCKER_COMPOSE_VERSION=1.29.2
sudo curl -L "https://github.com/docker/compose/releases/download/$DOCKER_COMPOSE_VERSION/docker-compose-
$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify Docker installation
sudo docker run hello-world

# Create Docker volume for persistent data
sudo docker volume create safeline_data

# Create a docker-compose.yml file
cat <<EOL > docker-compose.yml
version: '3.8'

services:
  app:
    image: your_app_image
    ports:
      - "80:80"
    volumes:
      - safeline_data:/data
    environment:
      - ENV_VAR1=value1
      - ENV_VAR2=value2

  metrics:
    image: metrics_image
    ports:
      - "3000:3000"

  audit:
    image: audit_image
    ports:
      - "4000:4000"

  ssl:
    container_name: ssl_service
    image: nginx
    ports:
      - "443:443"
    volumes:
      - ./ssl_cert:/etc/nginx/ssl

volumes:
  safeline_data:

EOL

# Start the Docker services
sudo docker-compose up -d
# Management Microservice for SafeLine

## Overview
The Management Microservice is a part of the SafeLine architecture designed to handle various management tasks and workflows. This service interacts with other microservices to provide a coherent management interface.

## Architecture
- **Microservice Architecture**: The system is designed based on microservice principles to promote scalability and maintainability.
- **Service Interactions**: The Management Microservice communicates with the following services:
  - User Service: Handles user management and authentication.
  - Notification Service: Manages notifications and alerts.
  - Inventory Service: Keeps track of inventory items and stock levels.

## Setup
### Prerequisites
- Install [Node.js](https://nodejs.org/) (Recommended version: 14.x or 16.x)
- Install [Docker](https://www.docker.com/) for containerized deployments.

### Environment Variables
Create a `.env` file in the root directory to configure your environment variables:
```plaintext
DATABASE_URL=your_database_url
REDIS_URL=your_redis_url
JWT_SECRET=your_jwt_secret
```

### Running Locally
1. Clone the repository:
   ```bash
   git clone https://github.com/clee699/SafeLine.git
   cd SafeLine/management
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the service:
   ```bash
   npm start
   ```

### Docker Setup
To run the Management Microservice in a Docker container, use the following command:
```bash
docker build -t safeline-management .
docker run -d -p 3000:3000 --env-file .env safeline-management
```

## Usage
### API Endpoints
- **GET /api/v1/users**: Fetch a list of users.
- **POST /api/v1/users**: Create a new user.
- **GET /api/v1/notifications**: Retrieve notifications for a user.

### Example Requests
#### Get Users
```bash
curl -X GET http://localhost:3000/api/v1/users
```
#### Create User
```bash
curl -X POST http://localhost:3000/api/v1/users \ 
-H 'Content-Type: application/json' \ 
-d '{"name":"John Doe", "email":"john.doe@example.com"}'
```

## Error Handling
### Common Errors
- **400 Bad Request**: Returned when the request parameters are invalid.
- **401 Unauthorized**: Returned when authentication is required but not provided or failed.
- **404 Not Found**: Indicates that the requested resource does not exist.

### Error Response Format
Each error response will be structured as follows:
```json
{
  "status": "error",
  "error": {
    "code": "error_code",
    "message": "Detailed error message"
  }
}
```

## Conclusion
The Management Microservice provides essential functionalities for managing users and notifications within the SafeLine ecosystem. Follow the above documentation to set up, use, and troubleshoot the service effectively.
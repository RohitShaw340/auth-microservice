# auth-microservice

## Setup

### Prerequisites

- Go 1.16 or higher
- Docker (optional, for containerization)
- Make sure you have `GOPATH` set up correctly

### Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/auth-microservice.git
   cd auth-microservice
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Run the service:
   ```sh
   go run main.go
   ```

### Environment Variables

Create a `.env` file in the root directory and add the following variables:

```env
MYSQL_TEST_DSN="username:password@/"
MONGODB_TEST_DSN="MongoDB URL"
```

### Docker

To run the service in a Docker container:

1. Build the Docker image:

   ```sh
   docker build -t auth-microservice .
   ```

2. Run the Docker container:
   ```sh
   docker run -d -p 8080:8080 --env-file .env auth-microservice
   ```

### Testing

Run tests using:

```sh
go test ./...
```

### Service Description

Here's a detailed Service Description Document for your authentication microservice, which provides essential information about its functionality, architecture, API endpoints, and operational details.

---

# **Service Description Document for Authentication Microservice**

### **1. Overview**

The Authentication Microservice provides secure identity management and user authentication functionalities for the platform. It manages consumer creation, user registration, and login, ensuring that users can securely access their respective accounts within a multi-tenant environment.

### **2. Service Architecture**

- **Technology Stack**:

  - **Programming Language**: Go (Golang)
  - **Database**: MongoDB (for consumer metadata) and MySQL (for user data for each consumer)
  - **Protocol**: gRPC for inter-service communication
  - **Frameworks & Libraries**: `go-grpc`, `mongo-driver`, `go-sql-driver/mysql`

- **Components**:
  - **Consumer Database (MongoDB)**: Stores consumer metadata, including name, phone, email, primary key fields, and data schemas.
  - **User Database (MySQL)**: Stores user information for each consumer, based on a dynamically created schema for each consumer's database.
  - **gRPC Server**: Manages requests related to consumer creation, user registration, and authentication.

---

### **3. Service Functionalities**

The microservice provides the following core functionalities:

1. **Consumer Management**:

   - `GenerateConsumerID`: Registers a new consumer and assigns them a unique ID.
   - `GetConsumerID`: Retrieves a consumer’s ID based on their email.

2. **User Management**:
   - `Signup`: Allows a user to register within the consumer's environment.
   - `Login`: Authenticates a user’s credentials and provides access if valid.

---

### **4. API Specifications**

#### **1. GenerateConsumerID**

- **Purpose**: Registers a new consumer.
- **Endpoint**: `GenerateConsumerID`
- **Request Fields**:
  - `name` (string): Name of the consumer
  - `phone` (string): Consumer’s contact number
  - `email` (string): Consumer’s email address
  - `schema` (map[string]string): The schema defining the fields and types for the consumer’s users
  - `primary_key_field` (string): The primary key field for identifying users within the consumer
- **Response Fields**:
  - `consumer_id` (string): Unique identifier assigned to the consumer
  - `message` (string): Success message

#### **2. GetConsumerID**

- **Purpose**: Retrieves a consumer ID based on their email.
- **Endpoint**: `GetConsumerID`
- **Request Fields**:
  - `email` (string): The email of the consumer
- **Response Fields**:
  - `consumer_id` (string): The unique identifier of the consumer
  - `message` (string): Success message or error message if not found

#### **3. Signup**

- **Purpose**: Registers a new user under a specific consumer.
- **Endpoint**: `Signup`
- **Request Fields**:
  - `consumer_id` (string): Unique identifier for the consumer
  - `user_data` (map[string]string): User data in key-value pairs (e.g., username, password)
  - `primary_key_field` (string): Field for identifying the user (e.g., username)
- **Response Fields**:
  - `message` (string): Success message or error if registration fails

#### **4. Login**

- **Purpose**: Authenticates a user within a consumer environment.
- **Endpoint**: `Login`
- **Request Fields**:
  - `consumer_id` (string): Consumer’s unique identifier
  - `primary_key_field` (string): The primary field to identify the user
  - `primary_key_value` (string): The value of the primary key (e.g., username)
  - `password` (string): User’s password
- **Response Fields**:
  - `user_details` (map[string]string): Authenticated user’s data
  - `message` (string): Success message or error if authentication fails

---

### **5. Database Schema**

#### **MongoDB (Consumer Metadata)**

- **Database Name**: `auth_service`
- **Collection**: `consumers`
- **Fields**:
  - `name` (string): Consumer's name
  - `phone` (string): Consumer's phone number
  - `email` (string): Consumer's email
  - `user_schema` (map): Schema structure for the user table in MySQL
  - `primary_key_field` (string): Primary key field for users
  - `_id` (ObjectId): Auto-generated unique identifier

#### **MySQL (User Data for Each Consumer)**

- **Database Naming Convention**: `consumer_<consumer_id>`
- **Table**: `users`
- **Dynamic Fields**: Defined by the `user_schema` provided during consumer registration

---

### **6. Error Handling and Status Codes**

- **gRPC Status Codes**:

  - `OK`: Request was successful.
  - `InvalidArgument`: One or more request fields were missing or invalid.
  - `NotFound`: Resource (consumer or user) was not found.
  - `AlreadyExists`: Resource already exists (e.g., consumer with the same email).
  - `Internal`: General server error or unexpected condition.

- **Error Responses**:
  - **GenerateConsumerID**: Returns an error if the consumer insertion fails in MongoDB.
  - **GetConsumerID**: Returns an error if the consumer is not found.
  - **Signup**: Returns an error if user registration fails.
  - **Login**: Returns an error if credentials are invalid.

---

### **7. Security Measures**

- **Data Encryption**: Passwords are stored in hashed format.
- **gRPC Security**: Supports TLS encryption for secure data transfer.
- **User Role Management**: The microservice relies on primary keys and unique identifiers to restrict access.

---

### **8. Deployment Details**

- **Deployment Environment**: Docker containers for ease of scalability
- **Microservice Scaling**: Horizontal scaling via Kubernetes for high availability
- **Logging**: Log entries for all major events, stored in a centralized logging service

---

### **9. Service Dependencies**

- **MongoDB**: Manages consumer metadata and facilitates multi-tenant support.
- **MySQL**: Provides relational storage of user data for each consumer, ensuring consumer data isolation.
- **gRPC**: Inter-service communication, ensuring low-latency communication between microservices.

---

### **10. Testing and Quality Assurance**

- **Unit Testing**: Covers each gRPC handler method.
- **Integration Testing**: Ensures compatibility between MongoDB, MySQL, and gRPC server.
- **Continuous Integration (CI)**: Automated testing suite triggered on new deployments.
- **Quality Assurance**: Manual testing for API endpoint validations and edge cases.

---

### **11. Future Enhancements**

1. **OAuth and SSO Support**: Extend to support third-party authentication providers.
2. **Rate Limiting**: Limit requests to prevent abuse and ensure fair use.
3. **Enhanced Logging and Monitoring**: Integrate with systems like Prometheus and Grafana for real-time analytics.

---

### **12. Glossary**

- **Consumer**: Represents a tenant (organization or company) registered in the system.
- **User**: An end user of a consumer (e.g., employees or customers).
- **Primary Key Field**: Unique identifier used to identify users in a consumer's database.
- **Schema**: A custom data structure specified by each consumer for user data.

---

This document provides a foundational understanding of the Authentication Microservice's functionality, enabling further enhancements or integration with other microservices within the broader system architecture.

### License

This project is licensed under the MIT License.

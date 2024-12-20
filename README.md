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

The Authentication Microservice provides secure identity management and user authentication functionalities for the platform. It manages client creation, user registration, and login, ensuring that users can securely access their respective accounts within a multi-tenant environment.

### **2. Service Architecture**

- **Technology Stack**:

  - **Programming Language**: Go (Golang)
  - **Database**: MongoDB (for client metadata) and MySQL (for user data for each client)
  - **Protocol**: gRPC for inter-service communication
  - **Frameworks & Libraries**: `go-grpc`, `mongo-driver`, `go-sql-driver/mysql`

- **Components**:
  - **Client Database (MongoDB)**: Stores client metadata, including name, phone, email, primary key fields, and data schemas.
  - **User Database (MySQL)**: Stores user information for each client, based on a dynamically created schema for each client's database.
  - **gRPC Server**: Manages requests related to client creation, user registration, and authentication.

---

### **3. Service Functionalities**

The microservice provides the following core functionalities:

1. **Client Management**:

   - `GenerateClientID`: Registers a new client and assigns them a unique ID.
   - `GetClientID`: Retrieves a client’s ID based on their email.

2. **User Management**:
   - `Signup`: Allows a user to register within the client's environment.
   - `Login`: Authenticates a user’s credentials and provides access if valid.

---

### **4. API Specifications**

#### **1. GenerateClientID**

- **Purpose**: Registers a new client.
- **Endpoint**: `GenerateClientID`
- **Request Fields**:
  - `name` (string): Name of the client
  - `phone` (string): Client’s contact number
  - `email` (string): Client’s email address
  - `schema` (map[string]string): The schema defining the fields and types for the client’s users
  - `primary_key_field` (string): The primary key field for identifying users within the client
- **Response Fields**:
  - `client_id` (string): Unique identifier assigned to the client
  - `message` (string): Success message

#### **2. GetClientID**

- **Purpose**: Retrieves a client ID based on their email.
- **Endpoint**: `GetClientID`
- **Request Fields**:
  - `email` (string): The email of the client
- **Response Fields**:
  - `client_id` (string): The unique identifier of the client
  - `message` (string): Success message or error message if not found

#### **3. Signup**

- **Purpose**: Registers a new user under a specific client.
- **Endpoint**: `Signup`
- **Request Fields**:
  - `client_id` (string): Unique identifier for the client
  - `user_data` (map[string]string): User data in key-value pairs (e.g., username, password)
  - `primary_key_field` (string): Field for identifying the user (e.g., username)
- **Response Fields**:
  - `message` (string): Success message or error if registration fails

#### **4. Login**

- **Purpose**: Authenticates a user within a client environment.
- **Endpoint**: `Login`
- **Request Fields**:
  - `client_id` (string): Client’s unique identifier
  - `primary_key_field` (string): The primary field to identify the user
  - `primary_key_value` (string): The value of the primary key (e.g., username)
  - `password` (string): User’s password
- **Response Fields**:
  - `user_details` (map[string]string): Authenticated user’s data
  - `message` (string): Success message or error if authentication fails

---

### **5. Database Schema**

#### **MongoDB (Client Metadata)**

- **Database Name**: `auth_service`
- **Collection**: `clients`
- **Fields**:
  - `name` (string): Client's name
  - `phone` (string): Client's phone number
  - `email` (string): Client's email
  - `user_schema` (map): Schema structure for the user table in MySQL
  - `primary_key_field` (string): Primary key field for users
  - `_id` (ObjectId): Auto-generated unique identifier

#### **MySQL (User Data for Each Client)**

- **Database Naming Convention**: `client_<client_id>`
- **Table**: `users`
- **Dynamic Fields**: Defined by the `user_schema` provided during client registration

---

### **6. Error Handling and Status Codes**

- **gRPC Status Codes**:

  - `OK`: Request was successful.
  - `InvalidArgument`: One or more request fields were missing or invalid.
  - `NotFound`: Resource (client or user) was not found.
  - `AlreadyExists`: Resource already exists (e.g., client with the same email).
  - `Internal`: General server error or unexpected condition.

- **Error Responses**:
  - **GenerateClientID**: Returns an error if the client insertion fails in MongoDB.
  - **GetClientID**: Returns an error if the client is not found.
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

- **MongoDB**: Manages client metadata and facilitates multi-tenant support.
- **MySQL**: Provides relational storage of user data for each client, ensuring client data isolation.
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

- **Client**: Represents a tenant (organization or company) registered in the system.
- **User**: An end user of a client (e.g., employees or customers).
- **Primary Key Field**: Unique identifier used to identify users in a client's database.
- **Schema**: A custom data structure specified by each client for user data.

---

This document provides a foundational understanding of the Authentication Microservice's functionality, enabling further enhancements or integration with other microservices within the broader system architecture.

### License

This project is licensed under the MIT License.

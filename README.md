# Service Catalog API

## Description
This is a service catalog API built with Go. It allows users to view, search, paginate, and sort services within an organization.

## Features
- List services with support for filtering, sorting, pagination
- Get details of a specific service
- Search for the specific service by name or description
- Create a new service
- Update an existing service
- Delete a service

## Design Considerations and Assumptions
- The API is designed to be as simple as possible for the given requirements.
- MongoDB is used as the database for its flexibility and scalability.

## Setup and Installation
1. Clone the repository
    ```sh
    git clone <repository_url>
    cd service-catalog
    ```

2. Install dependencies
    ```sh
    go mod tidy
    ```

3. Setup Mongo
   Mac (mongo already installed)
   ```sh
    brew services start mongodb-community
    mongsh
   ```

4. Configure env variable:
    Change the mongo URI in .env file where your instance is running and port where you want to run the app (deafult to port 8080)
    ```
    MONGODB_URI=mongodb://localhost:27017
    PORT=8080
    ```

5. Run the application
    ```sh
    go run cmd/server/main.go
    ```

4. The API will be available at `http://localhost:8080` by default

## API Endpoints

- `GET /services`: List all services
- `GET /services/:id`: Get details of a specific service
- `POST /services`: Create a new service
- `PUT /service/:id`: Update an existing service
- `DELETE /service/:id`: Delete a service

Query Params supported in List all services endpoint:
- page
- pageSize 
  Example: http://localhost:8080/services?page=2&pageSize=5
- sortField (created_at (default), name)
- sortOrder (1 for ascending order (default), -1 for descending)
  Example: http://localhost:8080/services?page=1&pageSize=5&sortField=name&sortOrder=-1
           http://localhost:8080/services?page=1&pageSize=10&sortField=created_at&sortOrder=1
- search 
  Example: http://localhost:8080/services?page=1&pageSize=10&sortField=name&sortOrder=1&search=touch 


## Example Requests:
- List Services with Pagination, Sorting, and Search
```sh
curl -X GET "http://localhost:8080/services?page=1&pageSize=10&sortField=name&sortOrder=1&search=example" \
-H "Authorization: Bearer <your_jwt_token>"
```

- Create a New Service
```sh
curl -X POST http://localhost:8080/services \
-H "Content-Type: application/json" \
-d '{
    "name": "Contact Us",
    "description": "Get in touch with us...",
    "versions": ["1.0", "1.1", "1.2"]
}'
```

- Update a service
```sh
curl -X PUT http://localhost:8080/services/<service_id> \
-H "Content-Type: application/json" \
-d '{
    "name": "Updated Service Name",
    "description": "updated description...",
    "versions": ["2.0", "2.1"]
}'
```

## Running Tests
```sh
go test ./internal/handlers
```

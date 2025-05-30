# YourVibes Go API

## Overview
YourVibes Go API is a backend service built using a **Domain-Driven Design (DDD) Monolithic Structure**, designed to power a social media platform with features like posting, commenting, messaging, and more. It integrates with various services and clients to provide a seamless user experience.

---

## DDD Monolithic Structure

### Structure Overview
This project adopts a **Domain-Driven Design (DDD)** approach within a monolithic architecture, inspired by the [go-ddd](https://github.com/sklinkert/go-ddd) repository. The DDD structure organizes the codebase into layers to ensure separation of concerns, maintainability, and scalability.

The DDD layers are structured as follows (as depicted in the architecture diagram):

- **Domain Layer**: The core of the application, containing business logic, entities, and domain services. This layer is independent of external frameworks and infrastructure.
- **Application Layer**: Orchestrates the use cases and coordinates the domain logic with external systems (e.g., APIs, databases).
- **Interface Layer**: Handles external communication, such as HTTP APIs (using Gin-Gonic) and gRPC endpoints.
- **Infrastructure Layer**: Manages external dependencies like databases (PostgreSQL), caching (Redis), message queues (RabbitMQ), and external services (AI censor service).

![Domain driven design](https://github.com/poin4003/images/blob/master/yourvibes_ddd_structure.png?raw=true)

### Purpose of DDD Structure
- **Separation of Concerns**: Each layer has a distinct responsibility, making the codebase easier to maintain and test.
- **Business Focus**: The domain layer ensures that business rules are central to the application, independent of infrastructure concerns.
- **Scalability**: While the application is monolithic, the layered structure allows for easier refactoring into microservices if needed.
- **Maintainability**: Changes in one layer (e.g., switching databases) have minimal impact on other layers.

---

## YourVibes Ecosystem

The YourVibes ecosystem consists of several components working together to deliver a complete social media platform:

- **Go API (Gin-Gonic)**: The core backend service, handling API requests, business logic, and integrations.
- **Redis**: Used for caching to improve performance for frequently accessed data.
- **RabbitMQ**: Facilitates asynchronous communication, including pushing notifications and interacting with the Python-based AI service.
- **gRPC**: Enables communication between the Go API and the Python server for AI censoring.
- **PostgreSQL**: The primary database for persistent data storage.
- **AI Service**: A Python-based service ([yourvibes_ai_service](https://github.com/poin4003/yourvibes_ai_service.git)) for content moderation (e.g., censoring sensitive content in posts and comments).
- **Clients**:
   - **Mobile App** ([yourvibes_app_V2](https://github.com/Thanh-Phuog/yourvibes_app_V2.git)): Built with React Native for mobile app users.
   - **Web App** ([yourvibes-web-client-v2](https://github.com/Trunks-Pham/yourvibes-web-client-v2.git)): Built with React for web app users.
   - **CMS for Admin** ([yourvibes-web-cms-v2](https://github.com/Trunks-Pham/yourvibes-web-cms-v2.git)): Built with React for admin management.

The ecosystem architecture is illustrated below:

![Ecosystem Architecture](https://github.com/poin4003/images/blob/master/yourvibes_architect_design.png?raw=true)

---

## Database Structure

The database schema, stored in PostgreSQL, is designed to support the core functionalities of the platform. Below is the Entity-Relationship Diagram (ERD):

![Database ERD](https://github.com/poin4003/images/blob/master/yourvibes_database.png?raw=true)

Key tables include:
- **users**: Stores user information (e.g., name, email, password, role).
- **posts**: Manages user posts with privacy settings (public, friend_only, private).
- **comments**: Supports infinite-layer comments on posts.
- **conversations** and **messages**: Handles messaging between users.
- **notifications**: Manages user notifications (e.g., new posts, comments).
- **friend_requests**, **friends**: Manages friendships and friend requests.
- **advertises**, **new_feeds**: Supports advertising and newsfeed features.
- **reports**, **bills**: Handles user reports and advertisement payments.

---

## Main Go Libraries and Frameworks

This project leverages several Go libraries and frameworks to build a robust backend:

- **Gin-Gonic**: A high-performance HTTP web framework for building RESTful APIs. It handles routing, middleware, and request processing for the YourVibes API.
- **Gorm**: An ORM library for managing database operations with PostgreSQL, providing a simple and efficient way to interact with the database.

For a full list of dependencies, check the `go.mod` file.

---

## Features

### User Functions
- **Post a Post**: Posts are pushed to friends' newsfeeds, with notifications. Supports privacy settings (public, friend_only, private) and AI censoring to block sensitive content.
- **Like, Share, Comment**: Supports liking and sharing posts, infinite-layer comments, and liking comments. AI censors sensitive comments (replaced with *).
- **Notifications**: Managed via socket notifications and a notification dashboard.
- **Friend Management**: Send friend requests, add/unfriend, get friend list, friend suggestions, and birthday reminders.
- **Profile Management**: Edit avatar, cover photo, and personal info with privacy settings (public, friend_only, private).
- **Newsfeed & Trending**: Get personal posts, friend posts, ads, and featured posts. Trending posts based on interactions (10 likes, 5 comments, 10 clicks, 10 views in 7 days).
- **Advertising**: Users can promote posts as ads (visible to all) for 33,000 VND/day (max 30 days), with a 6-hour push limit per user. Cleared after payment expires.
- **Featured Posts**: Pushed for 7 days if interaction thresholds are met (10 likes, 5 comments, 10 clicks, 10 views), with a 6-hour limit. Cleared after 7 days.
- **Messaging**: Socket-based messaging with 1:1 or group conversations. Roles include owner (can kick members, delete conversations) and member.
- **Authentication**: Login, signup, and Google login support.

### Admin Functions
- **Revenue Management**: View system revenue.
- **Report Handling**: Manage reports on posts, users, and comments. Actions include blocking users (email notification, temporary post/comment block), blocking posts/comments (with notifications), and re-opening if needed.
- **Transaction History**: View all advertisement payment transactions.
- **Super Admin**: Manages admin accounts (create, block).

### Cron Jobs
- Cleans expired ads and featured posts from newsfeeds.
- Pushes posts to friends' newsfeeds and manages ad/feature post limits.

---

## Setup Instructions

### Prerequisites
- Go (latest version)
- PostgreSQL, Redis, RabbitMQ
- YourVibes AI Service ([yourvibes_ai_service](https://github.com/poin4003/yourvibes_ai_service.git))

### Steps

1. **Clone the Repository**
   ```bash
   git clone https://github.com/poin4003/yourVibes_GoApi.git
   cd yourVibes_GoApi
   ```

2. **Create Configuration File**
   - Create a `config` folder in the root directory.
   - Add `dev.yaml`, `prod.yaml`, or `cloud.yaml` with the following template (fill in private data):
     ```yaml
     server:
       port: 8080
       mode: "dev"
       server_endpoint: "http://localhost"
     postgresql:
       host: localhost
       port: 5432
       username: postgres
       password: 
       dbname: yourvibes_db
       max_idle_conns: 10
       max_open_conns: 100
       conn_max_lifetime: 3600
       ssl_mode: disable
     logger:
       log_level: debug
       file_log_name: "./storages/logs/dev.001.log"
       max_size: 500
       max_backup: 3
       max_age: 10
       compress: true
     media:
       folder: "./storages/media"
     redis:
       host: localhost
       port: 6379
       password:
       database: 0
     authentication:
       jwtSecretKey: 
       jwtAdminSecretKey: 
     mail_service:
       smtp_host: smtp.gmail.com
       smtp_port: 587
       smtp_username: 
       smtp_password: ""
     momo:
       partner_code: "MOMO"
       access_key: "F8BBA842ECF85"
       secret_key: "K951B6PE1waDMi640xX08PD3vg6EkVlz"
       redirect_url: "http://localhost:8080/v1/2024/bill/"
       ipn_url: "http://localhost:8080/v1/2024/bill/"
       endpoint_host: "test-payment.momo.vn"
       endpoint_path: "/v2/gateway/api/create"
     google:
       google_tokens_url: "https://oauth2.googleapis.com/token"
       secret_id: ""
       web_client_id: ""
       android_client_id: ""
       ios_client_id: ""
     rabbitmq:
       url: "amqp://guest:guest@localhost:5672/"
       username: "guest"
       password: "guest"
       vhost: "/"
       connection_timeout: 10
       max_reconnect_attempts: 5
     moderate_service:
       health_url: "http://localhost:5000/health"
     comment_censor_grpc_conn:
       port: 50051
       host: localhost
     ```

3. **Download Dependencies**
   ```bash
   go mod tidy
   ```

4. **Migrate Database**
   - Start PostgreSQL.
   - Run:
     ```bash
     make migrate CONFIG_FILE=dev
     ```
     (Use `CONFIG_FILE=prod` or `CONFIG_FILE=cloud` for respective environments.)

5. **Start Services**
   - Start Redis, PostgreSQL, and RabbitMQ.
   - Start the YourVibes AI Service.

6. **Run the Server**
   - For development (HTTP):
     ```bash
     make dev
     ```
   - For production (HTTPS, requires certificates at `/etc/ssl/certs/fullchain.pem` and `/etc/ssl/certs/privkey.pem`):
     ```bash
     make prod
     ```
   - For cloud environment:
     ```bash
     make cloud
     ```

7. **Access the API**
   - Open `http://localhost:8080/swagger/index.html` to use the API.
   - If Swagger fails due to CORS, update the `@host` to `localhost:8080` and `@schema` to `http` in `main.go`.

### Makefile Commands
- read `Makefile` for more information 

---

**Note**: Ensure all private configuration values (e.g., passwords, keys) are filled in the YAML file before running the application.

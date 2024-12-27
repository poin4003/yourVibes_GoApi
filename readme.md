# Social Network API

## Overview
This is a simple social network API built using Go and the Gin framework. The API provides basic functionalities for managing users, posts, and interactions.

---

## Features
- User authentication and profile management
- Create, read, update, and delete posts
- Like and comment on posts
- Retrieve user activity and news feed

---

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/poin4003/yourVibes_GoApi.git
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run main.go
   ```

The server will run at `http://localhost:8080` by default.

---

## Notes
- Make sure to configure the database connection in the `config` folder.
- API routes and documentation are defined in the `routes` folder.


# Go Link Vault
![Build Status](https://github.com/DebroyeAntoine/GoGoLinkVault/actions/workflows/go_unit_tests.yml/badge.svg)

Go Link Vault is a link management application that allows users to create, organize, update, delete, and share links with tags and categories. The project uses **Go (Gin)** for the backend and **React (TypeScript, Tailwind CSS, Redux)** for the frontend.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
  - [Backend](#backend)
  - [Frontend](#frontend)
- [Configuration](#configuration)
  - [Backend](#backend-configuration)
  - [Frontend](#frontend-configuration)
- [Running the Application](#running-the-application)
- [API](#api)
  - [Endpoints](#endpoints)
- [Tests](#tests)
  - [Backend Tests](#backend-tests)
  - [Frontend Tests](#frontend-tests)
- [Contributing](#contributing)
- [License](#license)

---

## Prerequisites

Make sure the following tools are installed on your machine:

- **Go**: Version 1.18+ (for the backend)
- **Node.js** and **npm**: Version 16+ (for the frontend)
- **Docker** (optional, for database management and development environment)

---

## Installation

### Backend

1. Clone the backend repository:

```bash
git clone https://github.com/your-username/go-link-vault.git
cd go-link-vault
```

2. Install Go dependencies:

```bash
go mod tidy
```

3. Configure the database and environment variables:
   - Create a `.env` file at the root of the project with the following:

```env
DB_USER=username
DB_PASSWORD=password
DB_NAME=go_link_vault
JWT_SECRET_KEY=secretkey
```

4. Run the server:

```bash
go run main.go
```

The server will be available at `http://localhost:8080`.

### Frontend

1. Clone the frontend repository (if separate):

```bash
git clone https://github.com/your-username/go-link-vault-frontend.git
cd go-link-vault-frontend
```

2. Install frontend dependencies:

```bash
npm install
```

3. Start the React development server:

```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`.

---

## Configuration

### Backend Configuration

1. Make sure the database is ready (PostgreSQL or other) and tables are created.

2. Modify the environment variables in the `.env` file to connect to your database.

### Frontend Configuration

1. If you're using a different backend or port for the API, modify the API URL in `src/api/index.ts`.

---

## Running the Application

### Running Backend and Frontend Together

1. Run the backend in one terminal:

```bash
go run main.go
```

2. Run the frontend in another terminal:

```bash
npm run dev
```

3. Access the application in your browser at [http://localhost:3000](http://localhost:3000).

---

## API

### Endpoints

#### Authentication

- **POST /auth/login**  
  Allows the user to log in and get a JWT token.
  - **Parameters**: 
    - `email`: User's email
    - `password`: User's password
  - **Response**: A JSON object containing the JWT token.

#### Links

- **GET /links**  
  Retrieve all links for the logged-in user.
  - **Response**: A list of links associated with the user.

- **POST /links**  
  Create a new link for the logged-in user.
  - **Parameters**: 
    - `url`: The link URL
    - `title`: The title of the link
    - `tags`: List of tags associated with the link
  - **Response**: The created link.

- **PUT /links/{id}**  
  Update an existing link.
  - **Parameters**: 
    - `id`: Link ID
    - `url`: New URL of the link
    - `title`: New title of the link
    - `tags`: New tags
  - **Response**: The updated link.

- **DELETE /links/{id}**  
  Delete a link.
  - **Parameters**: 
    - `id`: The link ID to delete
  - **Response**: Confirmation of the deletion.

---

## Tests

### Backend Tests

1. Install test dependencies for the backend:

```bash
go get github.com/stretchr/testify
```

2. Run the backend tests:

```bash
go test ./...
```

3. Test handler logic like link creation, authentication, etc.

### Frontend Tests

1. Run the frontend tests with Jest and React Testing Library:

```bash
npm run test
```

2. Test components and Redux logic with unit tests.

---

## Contributing

1. Fork this repository
2. Create a branch for your feature or bugfix.
3. Submit a Pull Request with a detailed description of your changes.

---

## License

Distributed under the MIT License. See the [LICENSE](./LICENSE) file for more information.

---

### 📝 Explanation

- This **README.md** covers installation, configuration, API usage, and tests for a full-stack project with a Go backend (Gin) and a React frontend with Redux.
- Feel free to customize further based on your specific needs, such as detailing entity structures, future features, or adding additional setup steps.

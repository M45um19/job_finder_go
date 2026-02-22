JobFinder API
JobFinder API is a high-performance, concurrent job portal backend built with Go (Standard Library) and PostgreSQL. It follows a clean architecture pattern and is designed to handle multiple job seekers and employers simultaneously using Go's powerful concurrency primitives.

This is a Free & Open Source API. Developers are welcome to clone, use, and integrate this backend into their own frontend applications.

Key Features:
Framework-less Core: Built primarily using Go's net/http and the lightweight Chi router for professional routing.

Role-Based Access Control (RBAC): Distinct permissions for Seekers (applying for jobs) and Employers (posting jobs).

Concurrency at Scale: Utilizes Goroutines for background tasks like notification logging and application processing.

Secure Auth: JWT-based authentication with custom middleware for route protection.

Clean Architecture: Organized into Handlers, Services, and Repositories for maximum maintainability.

Tech Stack:
Language: Go (1.20+)
Router: go-chi/chi
Database: PostgreSQL
Driver: jackc/pgx
Authentication: JWT (JSON Web Tokens)
Environment Management: joho/godotenv

Project Structure:

jobstream-api/
├── cmd/api/            # Entry point
├── internal/
│   ├── handlers/       # HTTP Request Handlers
│   ├── services/       # Business Logic Layer
│   ├── repository/     # Database Interaction (SQL)
│   ├── models/         # DB Entities & Structs
│   ├── middleware/     # Auth & Role Guards
│   └── router/         # Chi Route Definitions
├── migrations/         # SQL Schema files
└── pkg/                # Helper packages (Logger etc.)

Database Design: 
The system is built on a relational PostgreSQL schema designed for integrity and performance.
Users: Stores credentials, roles (seeker, employer), and profile info.
Jobs: Managed by Employers; contains job details and location types.
Applications: Links Seekers to Jobs with a unique constraint to prevent duplicate applications.

API Endpoints:

Authentication (Public):

Method | Endpoint | Description
POST | /api/v1/auth/register | Create a new account (Seeker/Employer)
POST | /api/v1/auth/login | Login and receive a JWT Token

Job Listings (Public & Protected)

Method | Endpoint | Description | Access
GET | /api/v1/jobs | List all jobs (supports filters) | Public
GET | /api/v1/jobs/{id} | Get specific job details | Public
POST | /api/v1/jobs | Post a new job | Employer
PUT | /api/v1/jobs/{id} | Update an existing job,Employer (Owner)
DELETE | /api/v1/jobs/{id} | Delete a job post | Employer (Owner)

Job Applications (Protected)

Method | Endpoint | Description | Access
POST | /api/v1/jobs/{id}/apply | Apply for a job (Concurrent processing) | Seeker
GET | /api/v1/applications | View all your submitted applications | Seeker
GET | /api/v1/jobs/{id}/applicants | View list of people who applied | Employer


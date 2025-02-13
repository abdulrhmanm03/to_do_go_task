# To-Do List API with Authentication and Database

## Overview

This project is a REST API for managing a to-do list built with Go, Gin, GORM, PostgreSQL, JWT authentication, and Docker Compose.

## Features

- **User Authentication:** Register and log in with JWT-based authentication.
- **Task Management:** Create, read, update, and delete tasks.
- **Database:** Uses PostgreSQL with GORM for ORM.
- **Dockerized:** Easily run the API and PostgreSQL using Docker Compose.
- **Testing:** Basic unit tests using Testify.

## Prerequisites

- Docker & Docker Compose installed on your machine.

## Setup Instructions

1. **Clone the repository:**
   ```bash
   git clone https://github.com/abdulrhmanm03/to_do_go_task.git
   cd to_do_go_task
   ```
2. **To run the container:**
   ```bash
   docker-compose up
   ```

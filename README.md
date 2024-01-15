>75

## AMS - Backend

Welcome to >75, an open-source platform attendance management application designed to help students stay on top of their attendance requirements! This is the backend for the app.

## Installation

To run this project locally, ensure you have Go and PostgreSQL installed on your machine.

1. **Go Installation:**
   - Follow the official [Go installation guide](https://golang.org/doc/install) to install Go on your system.

2. **PostgreSQL Installation:**
   - Download and install PostgreSQL from the [official website](https://www.postgresql.org/download/).
   - Set up a PostgreSQL database and note down the connection details.

## Setup

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/your-username/your-project.git
   cd your-project

2. **Setup env variables**
   copy paste the .env.example to .env file and change the details.
    ```bash
   cp .env.example .env
3. **Install Dependencies**
   ```bash
   go mod tidy

4. **Run the File**
    ```bash
   go run main.go

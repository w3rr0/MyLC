# MyLC for IAESTE

My Local Committee is a comprehensive application
designed to streamline the organization of stands and other events
by simplifying the signup process for members and organizers.
The system automates account creation based on @iaeste.pl email addresses
and tracks attendance efficiently.
Most importantly, the application is fully responsive,
ensuring full compatibility with both mobile and desktop devices.

This guide provides instructions on how to set up the environment
and run the application locally using Docker.
If you wish to build a complete system for your committee
based on this project - which I highly encourage — please find the resources below.

## 📋 Prerequisites
Before you begin, ensure you have the following tools installed on your machine.
Thanks to Docker, you do not need to install Go or PostgreSQL locally.

- `Docker desktop`: Includes the Docker Engine and Docker Compose.

- `Git`: For cloning the repository.

## 🚀 Getting Started

Follow these steps to get the application up and running.

#### 1. Clone the Repository


Fork this repository to your GitHub account, then clone it locally:

```bash
  # Using SSH
  git clone git@github.com:<YOUR_NICK>/IAESTE_stands_server.git

  # Or using HTTPS
  git clone https://github.com/<YOUR_NICK>/IAESTE_stands_server.git
```

Note: The main branch always contains the latest stable version of the code.
If a new version is released, update your local copy using the Sync fork button
on your GitHub repository.

#### 2. Build and Run via Docker

Navigate to the project directory and start the application using Docker Compose.
This command will build the API image and pull the necessary database image.

```bash
# Build containrest
docker compose build
# Then start with logs
docker compose up
# Or without (recommended unless you want to contribute to development)
docker compose up -d
```

- The API will be available at: http://localhost:8080
- The Database will be exposed on port: 5432


To stop the application:
```bash
docker compose down
```

You do not need to manually create the database or import SQL files.

The Docker setup is configured to automatically initialize the PostgreSQL database
using the `schema.sql` file upon the first launch.
The database `mylc` (and user postgres) will be created and populated
with the necessary tables automatically.

## 🛠️ Troubleshooting

- Setup Assistance: If you need help with the setup, feel free to DM me.
- Reporting Bugs: If you find a bug, please open an Issue in this repository.
Ideally, include a description of the error and steps to reproduce it.
Before submitting, please check if the issue has already been reported
or fixed in the latest version of the code.

## 📄 Note

This project is in an early stage of development
and is subject to significant changes. If you are involved in its development,
I recommend checking this repository regularly for updates.
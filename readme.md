# Setu Project App

## 🚀 Getting Started
This project runs inside a Docker container with a MySQL database.

### Prerequisites
Make sure you have the following installed:
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## 📦 Setup & Run the Application

### 1️⃣ Clone the Repository
```sh
git clone https://github.com/MetsysEht/setuProject
cd setuProject
```
### 2️⃣ Build and Run with Docker Compose
Before running the application, update the .env file with the necessary configuration:

```sh
docker-compose up -d --build
```
- `-d`: Runs the containers in detached mode.
- `--build`: Rebuilds the images if needed.

### 3️⃣ Verify Running Containers
```sh
docker ps
```
This should show both `setu` and `mysql-db` containers running.

### 4️⃣ Access the Application
- Https server is running on:  
  **http://localhost:8081**
- RPC server is running on:  
  **http://localhost:8080**
- MySQL will be accessible on **localhost:3307**

## 🛠 Environment Variables
Modify the `docker-compose.yml` or use `.env` file for database settings:
```
DB_HOST=db
DB_USER=root
DB_PASSWORD=password
DB_NAME=setu
```

## 🛑 Stopping and Cleaning Up
To stop the application and remove containers:
```sh
docker-compose down
```
To remove all Docker images and volumes:
```sh
docker system prune -a
```

## 📜 License
This project is licensed under the MIT License.

---
Feel free to contribute or raise issues if you encounter any problems! 🚀


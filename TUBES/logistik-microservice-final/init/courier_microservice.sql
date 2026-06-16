CREATE DATABASE IF NOT EXISTS courier_microservice;
USE courier_microservice;

CREATE TABLE IF NOT EXISTS deliveries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resi VARCHAR(100) UNIQUE,
    courier_id INT,
    status VARCHAR(50),
    assigned_zone VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
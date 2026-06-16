CREATE DATABASE IF NOT EXISTS gudang_microservice;
USE gudang_microservice;

CREATE TABLE IF NOT EXISTS packages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resi VARCHAR(100) UNIQUE,
    user_id INT,
    nama_barang VARCHAR(255),
    berat INT,
    status VARCHAR(50),
    warehouse_zone VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sorted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE outbox_events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_type VARCHAR(100),
    payload TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP NULL
);
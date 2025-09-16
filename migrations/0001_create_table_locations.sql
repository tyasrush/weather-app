CREATE TABLE IF NOT EXISTS locations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,       
    region VARCHAR(255),
    country VARCHAR(100),            
    latitude DECIMAL(9,6) DEFAULT 0.0,          
    longitude DECIMAL(9,6) DEFAULT 0.0,        
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_modified_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS weathers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    location_id BIGINT NOT NULL,
    temperature DECIMAL(5,2),          
    humidity INT,                      
    wind_speed DECIMAL(5,2),           
    condition_status VARCHAR(100),            
    condition_icon_url VARCHAR(255),            
    forecast_time DATETIME NOT NULL,   
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE
);

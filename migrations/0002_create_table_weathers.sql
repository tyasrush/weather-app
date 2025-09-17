CREATE TABLE IF NOT EXISTS weathers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    location_id BIGINT NOT NULL,
    temperature_celcius DECIMAL(5,2),          
    temperature_fahrenheit DECIMAL(5,2),          
    humidity INT,                      
    wind_speed DECIMAL(5,2),           
    condition_status VARCHAR(100),            
    condition_icon_url VARCHAR(255),            
    forecast_time TIMESTAMP NOT NULL,   
    forecast_type ENUM('day', 'hour') NOT NULL DEFAULT 'hour',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_modified_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE
);

ALTER TABLE weathers ADD UNIQUE KEY unique_location_forecast (location_id, forecast_time, forecast_type);


DELIMITER $$

CREATE TRIGGER trigger_weather_last_modified_at
BEFORE UPDATE ON weathers
FOR EACH ROW
BEGIN
    SET NEW.last_modified_at = NOW();
END$$

DELIMITER ;



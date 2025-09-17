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

CREATE INDEX idx_location_deleted_at ON locations(deleted_at);

DELIMITER $$

CREATE TRIGGER trigger_locations_last_modified_at
BEFORE UPDATE ON locations
FOR EACH ROW
BEGIN
    SET NEW.last_modified_at = NOW();
END$$

DELIMITER ;

-- seed data to make testing more efficient. not recommended implementation, but good to go for now 
INSERT INTO locations(name, region, country, latitude, longitude) VALUES('Jakarta', 'DKI Jakarta', 'Indonesia', -6.2088, 106.8456);
INSERT INTO locations(name, region, country, latitude, longitude) VALUES('Bandung', 'West Java', 'Indonesia', -6.9175, 107.6191);
INSERT INTO locations(name, region, country, latitude, longitude) VALUES('Surabaya', 'East Java', 'Indonesia', -7.2575, 112.7521);

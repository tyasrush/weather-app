CREATE TABLE IF NOT EXISTS orders (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  customer_id VARCHAR(36) NOT NULL,
  amount DECIMAL(10,2) NOT NULL,
  status ENUM('PENDING', 'PAID', 'CANCELLED') NOT NULL,
  created_at DATETIME NOT NULL
);

INSERT INTO orders (customer_id, amount, status, created_at) VALUES
('john_doe',       120.50, 'PENDING',   '2025-09-01 10:15:00'),
('jane_doe',       300.00, 'PAID',      '2025-09-02 11:30:00'),
('harry_potter',    99.99, 'CANCELLED', '2025-09-03 09:45:00'),
('hermione_granger',450.25,'PAID',      '2025-09-04 14:20:00'),
('ron_weasley',    200.75, 'PENDING',   '2025-09-05 16:10:00'),
('tony_stark',      75.00, 'PAID',      '2025-09-06 12:00:00'),
('bruce_wayne',    600.10, 'CANCELLED', '2025-09-07 08:40:00'),
('clark_kent',     150.35, 'PENDING',   '2025-09-08 18:25:00'),
('diana_prince',   325.00, 'PAID',      '2025-09-09 13:50:00'),
('steve_rogers',   500.90, 'PENDING',   '2025-09-10 17:15:00'),
('natasha_romanoff',250.00,'CANCELLED', '2025-09-11 15:30:00'),
('bruce_banner',   135.60, 'PAID',      '2025-09-12 09:05:00'),
('peter_parker',   420.40, 'PENDING',   '2025-09-13 11:45:00'),
('wade_wilson',     85.20, 'PAID',      '2025-09-14 19:00:00'),
('logan_howlett',  710.75, 'CANCELLED', '2025-09-15 20:10:00'),
('charles_xavier', 180.50, 'PAID',      '2025-09-16 10:30:00'),
('erik_lensherr',   95.00, 'PENDING',   '2025-09-17 08:55:00'),
('stephen_strange',310.80, 'PAID',      '2025-09-18 14:40:00'),
('scott_summers',  275.25, 'CANCELLED', '2025-09-19 16:20:00'),
('jean_grey',      520.99, 'PENDING',   '2025-09-20 12:15:00');


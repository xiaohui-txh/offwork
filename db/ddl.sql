CREATE DATABASE offwork DEFAULT CHARSET utf8mb4;

USE offwork;

CREATE TABLE offwork_checkin (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  lng DOUBLE NOT NULL COMMENT '经度',
  lat DOUBLE NOT NULL COMMENT '纬度',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '打卡时间',
  KEY idx_created_lng_lat (created_at, lng, lat),
  KEY idx_lng_lat (lng, lat)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='下班打卡记录表';

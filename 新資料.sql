CREATE TABLE appointments (
    appointment_id    INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    department_id     INT UNSIGNED NOT NULL,
    doctor_id         INT UNSIGNED NOT NULL,
    patient_id        VARCHAR(20) NOT NULL,
    appointment_time  DATETIME NOT NULL,
    status            VARCHAR(20) DEFAULT 'booked',
    service_type      VARCHAR(50),
    checkin_time      DATETIME DEFAULT NULL,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- 建議加上外鍵約束（若相關表存在）
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (doctor_id) REFERENCES doctors(id)
    -- 若有病人資料表，也可加：FOREIGN KEY (patient_id) REFERENCES patients(id)
);

CREATE TABLE departments (
    id                INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    dept_name         VARCHAR(100) NOT NULL,
    dept_description  TEXT,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE doctors (
    doctor_id    INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    dept_id      INT UNSIGNED NOT NULL,
    doctor_name  VARCHAR(100) NOT NULL,
    doctor_info  TEXT,
    password     VARCHAR(255) NOT NULL,
    FOREIGN KEY (dept_id) REFERENCES departments(id) ON DELETE CASCADE
);

CREATE TABLE doctor_leaves (
    leave_id              INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    doctor_id             INT UNSIGNED NOT NULL,
    start_time            DATETIME NOT NULL,
    end_time              DATETIME NOT NULL,
    substitute_doctor_id  INT UNSIGNED DEFAULT NULL,
    FOREIGN KEY (doctor_id) REFERENCES doctors(doctor_id) ON DELETE CASCADE,
    FOREIGN KEY (substitute_doctor_id) REFERENCES doctors(doctor_id) ON DELETE SET NULL
);

CREATE TABLE feedbacks (
    id               BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    appointment_id   INT UNSIGNED NOT NULL,
    feedback_rating  INT NOT NULL,
    patient_comment  TEXT,
    created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at       DATETIME DEFAULT NULL,

    FOREIGN KEY (appointment_id) REFERENCES appointments(appointment_id) ON DELETE CASCADE
);

CREATE TABLE managers (
    manager_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE patients (
    patient_id VARCHAR(255) PRIMARY KEY,
    patient_name VARCHAR(255) NOT NULL,
    patient_gender VARCHAR(50),
    patient_birth VARCHAR(50),
    patient_phone VARCHAR(50),
    password VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    emergency_name VARCHAR(255),
    emergency_phone VARCHAR(50),
    emergency_relation VARCHAR(50),
    drug_allergy JSON,
    food_allergy JSON,
    medical_history JSON
);

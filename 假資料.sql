-- 部門表
DROP TABLE IF EXISTS Department;
CREATE TABLE Department (
  Dept_ID INT PRIMARY KEY,
  Dept_Name VARCHAR(20),
  Dept_Description TEXT
);

INSERT INTO Department (Dept_ID, Dept_Name, Dept_Description) VALUES
(1, '外科', '外科主要負責以手術方式治療各類器官系統的疾病，包括消化道、內分泌系統、乳房、血管與腹壁等部位。外科醫師透過手術處置腫瘤、感染、外傷與結構異常等問題，並與其他科別合作提供整合性照護。'),
(2, '內科', '針對人體內部器官系統疾病進行非手術性治療的專業科別，涵蓋範圍廣泛，為病患進行初步診斷、藥物治療、慢性病管理與轉診協調。內科也是大多數患者的首診科別。'),
(3, '耳鼻喉科', '耳鼻喉科負責診療耳、鼻、咽喉及頭頸部相關疾病，並涵蓋聲音、吞嚥、平衡、睡眠呼吸功能障礙等子領域，服務對象涵蓋成人與兒童族群。');

-- 醫師表
DROP TABLE IF EXISTS Doctor;
CREATE TABLE Doctor (
  Doctor_ID CHAR(5) PRIMARY KEY,
  Dept_ID INT,
  Doctor_Name VARCHAR(10) NOT NULL,
  Doctor_Hire_Date DATE NOT NULL,
  Doctor_Gender ENUM('男', '女') NOT NULL,
  Doctor_Education VARCHAR(100),
  Doctor_Experience TEXT,
  Doctor_Phone VARCHAR(15),
  Password VARCHAR(20) NOT NULL,
  Doctor_Expertise TEXT
  FOREIGN KEY (Dept_ID) REFERENCES Departments(Dept_ID)
);

INSERT INTO Doctor (Doctor_ID, Doctor_Name, Doctor_Hire_Date, Doctor_Gender, Doctor_Education, Doctor_Experience, Doctor_Phone, Password, Doctor_Expertise) VALUES
('10001', 1, '陳建民', '2018-03-15', '男', '國立台灣大學醫學系', '台大醫院、榮總住院醫師', '0912345678', 'Ab123456', '大腸直腸癌手術、甲狀腺與乳房腫瘤切除、腹腔鏡微創手術'),
('10002', 1, '林芳瑜', '2020-07-22', '女', '國立陽明交通大學醫學系', '亞東醫院主治醫師', '0922333444', 'Cd789012', '肛門痔瘡手術、體表腫瘤/囊腫切除'),
('20001', 2, '張育誠', '2015-11-01', '男', '高雄醫學大學醫學系', '馬偕醫院', '0933444555', '123Ak4567', '糖尿病、慢性腎臟病'),
('20002', 2, '周怡君', '2019-02-28', '女', '中山醫學大學醫學系', '彰基、澄清醫院', '0911222333', 'QnB234560', '腸胃疾病、呼吸道疾病'),
('30001', 3, '李冠廷', '2017-05-18', '男', '中國醫藥大學醫學系', '中醫附醫', '0988777666', 'Hdd284901', '過敏性鼻炎、耳鳴、眩暈與中耳炎診治'),
('30002', 3, '吳宛蓉', '2021-09-05', '女', '馬偕醫學院醫學系', '馬偕耳鼻喉科診所', '0977555333', '5672jKKLQ', '鼻過敏與長期鼻塞治療、聽力檢查與耳道異物處理');

-- 病人表
DROP TABLE IF EXISTS Patient;
CREATE TABLE Patient (
  Patient_ID CHAR(10) PRIMARY KEY,
  Patient_Name VARCHAR(20) NOT NULL,
  Patient_Gender ENUM('男', '女') NOT NULL,
  Patient_Birth DATE NOT NULL,
  Patient_Phone CHAR(10) NOT NULL UNIQUE,
  Password VARCHAR(20) NOT NULL,
  drug_allergy TEXT,
  food_allergy TEXT,
  medical_history TEXT
);

INSERT INTO Patient (Patient_ID, Patient_Name, Patient_Gender, Patient_Birth, Patient_Phone, Password, drug_allergy, food_allergy, medical_history) VALUES
('A123456789', '王小明', '男', '1990-05-12', '0912345678', 'Password123', '青黴素', '花生', '高血壓'),
('B234567890', '李美麗', '女', '1985-09-23', '0987654321', 'passw0rd30', '無', '海鮮', '右手撕裂傷'),
('T174892018', '陳大同', '男', '1978-02-03', '0933555666', 'secure123Pwd', '阿司匹林', '無', '無'),
('S274850138', '周佳怡', '女', '2000-12-31', '0966888999', '12A345678', '無', '無', '無'),
('J160315637', '黃志強', '男', '1982-03-05', '0911222333', 'abC123','青黴素, 阿莫西林','花生, 蛋','高血壓,鼻竇炎 ');

-- 看診預約表
DROP TABLE IF EXISTS Appointment;
CREATE TABLE Appointment (
  Appointment_ID INT AUTO_INCREMENT PRIMARY KEY,
  Dept_ID INT NOT NULL,
  Doctor_ID CHAR(5) NOT NULL,
  Patient_ID CHAR(10) NOT NULL,
  Appointment_Time DATETIME NOT NULL,
  Status ENUM('booked', 'completed', 'cancelled', 'no_show') DEFAULT 'booked',
  FOREIGN KEY (Patient_ID) REFERENCES Patients(Patient_ID),
  FOREIGN KEY (Doctor_ID) REFERENCES Doctors(Doctor_ID),
  FOREIGN KEY (Dept_ID) REFERENCES Departments(Dept_ID)
);

INSERT INTO Appointment (Dept_ID, Doctor_ID, Patient_ID, Appointment_Time, Status) VALUES
(2, '10001', 'A123456789', '2025-06-10 09:00:00', 'booked'),
(2, '10001', 'B234567890', '2025-06-11 10:30:00', 'completed'),
(1, '20001', 'T174892018', '2025-06-12 14:00:00', 'cancelled'),
(2, '10002', 'S274850138', '2025-06-13 10:00:00', 'no_show'),
(3, '30002', 'J160315637', '2025-06-14 16:00:00', 'booked');

-- 回饋表（以後端生成 Feedback_ID 為前提）
DROP TABLE IF EXISTS Feedback;
CREATE TABLE Feedback (
  Feedback_ID VARCHAR(10) PRIMARY KEY,
  Appointment_ID INT UNIQUE,
  Feedback_Rating INT CHECK (Feedback_Rating BETWEEN 1 AND 5),
  Patient_Comment VARCHAR(200),
  FOREIGN KEY (Appointment_ID) REFERENCES Appointments(Appointment_ID)
);

INSERT INTO Feedback (Feedback_ID, Appointment_ID, Feedback_Rating, Patient_Comment) VALUES
('A1', 1, 5, '醫生很親切，講解詳細。'),
('A2', 2, 4, '等候時間稍長，但服務不錯。'),
('A3', 3, 2, '看診太匆忙，沒講清楚。'),
('A4', 4, 1, '臨時取消也沒通知，太糟糕了。'),
('A5', 5, 3, '普通，沒有特別好或壞的感覺。');

-- 管理員表
DROP TABLE IF EXISTS Manager;
CREATE TABLE Manager (
  Manager_ID CHAR(2) PRIMARY KEY,
  Doctor_ID CHAR(5),
  Appointment_ID INT,
  Generated_Time DATETIME,
  Login_History DATETIME,
  Report_Type VARCHAR(100),
  FOREIGN KEY (Doctor_ID) REFERENCES Doctors(Doctor_ID),
  FOREIGN KEY (Appointment_ID) REFERENCES Appointments(Appointment_ID)
);

INSERT INTO Manager (Manager_ID, Doctor_ID, Appointment_ID, Generated_Time, Login_History, Report_Type) VALUES
('M1', '10001', 1, '2025-06-06 08:00:00', '2025-06-06 08:00:00', '每日報告'),
('M2', '10002', 2, '2025-06-07 08:00:00', '2025-06-07 08:00:00', '統計分析');

-- 報告表（改名避免保留字衝突）
DROP TABLE IF EXISTS ManagerReport;
CREATE TABLE ManagerReport (
  Report_ID INT PRIMARY KEY AUTO_INCREMENT,
  Manager_ID CHAR(2),
  Patient_Noshow_Rate FLOAT,
  Doctor_Utilization_Rate FLOAT,
  Daily_Report TEXT,
  Total_Appointments INT,
  FOREIGN KEY (Manager_ID) REFERENCES Managers(Manager_ID)
);

INSERT INTO ManagerReport (Manager_ID, Patient_Noshow_Rate, Doctor_Utilization_Rate, Daily_Report, Total_Appointments) VALUES
('M1', 0.1, 0.85, '2025-06-06 report content', 18),
('M2', 0.05, 0.92, '2025-06-07 report content', 20);

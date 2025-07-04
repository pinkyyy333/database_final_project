# Database Final Project 診所預約系統
## 課程：NYCU 113-2 資料庫管理
### 組員：111704024葉昱欣、111704032曾鏡恩、111704035蕭佳蘋、111704050林珮瑜
**分工**
<br/>葉昱欣、曾鏡恩：前端+ERD圖+簡報
<br/>林珮瑜：後端
<br/>蕭佳蘋：前後端整合

**系統架構**
<br/>前端：使用react框架模組化html+css+js
<br/>後端：使用golang搭配MySQL

## 🧩 系統功能說明

### 👤 病患功能

**帳號管理**：可註冊（使用身分證字號）、登入與編輯個人資料。

**預約功能**：
- 可瀏覽各科別與醫師（包含醫師簡介）。
- 查詢可預約時段。
- 預約或取消預約。
- 預約流程：**選擇科別 → 醫師 → 日期 → 時段**

**看診紀錄查詢**：
- 可查看即將到來與過去的預約紀錄。

**回饋與評分**：
- 看診後可對醫師或整體經驗進行評分與回饋。

---

### 🦥 醫師功能

**排班與預約管理**：可檢視自身排班與病患名單。

**看診處理**：可更新看診狀態（報到、完成、未出席等）。

**病患紀錄與回饋**：
- 可查看病患過往的就診紀錄與回饋（平均評分、留言、簡單統計等）。

---

### 🧑‍💼 管理者功能

**醫師與排班管理**：
- 新增、編輯、刪除醫師資料。
- 設定可預約時段。

**預約監控與流程管理**：
- 監控所有預約紀錄。
- 管理衝突、取消與未出席情況。

**報表與分析**：
- 每日/每週產出報表，如：預約總數、缺席率、醫師利用率等。

---

## 🏥 基本系統要求

- 至少支援三個科別（如皮膚科、小兒科、內科），每科至少2位醫師。
- 系統需具備以下能力：
  - 資料驗證（例如不能輸入錯誤資料）
  - 衝突偵測（如時段重疊）
  - 重要操作須有確認機制（如刪除預約）

---

## 🌟 進階挑戰功能—Bonus

- **預約提醒**：透過 Email、Line、簡訊提醒病患。
- **自動未出席偵測**：若過時未報到，自動標記為未出席。
- **醫師替代管理**：臨時請假時可快速指派代診醫師並通知病患。
- **即時看診進度**：顯示現場候診進度與預估等待時間（需動態更新）。
- **非看診服務**：例如疫苗採種預約。

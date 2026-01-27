# 主數據管理系統（MDM）系統清單

> 本文檔整理 **主數據管理系統（MDM, Master Data Management）** 的主流產品，覆蓋：
>
> * 🌍 國際廠商 / 🇨🇳 中國廠商
> * 💰 商業收費 / 🔓 開源（含社區版）
> * 🧩 多域 MDM / 客戶、產品、供應商等專用 MDM

適用於 **主數據平臺選型、技術對標、方案設計、RFP 編製** 等場景。

## 一、國際主流商業 MDM 系統（收費）

### 1. Informatica MDM

* **官網**：[https://www.informatica.com](https://www.informatica.com)
* **類型**：多域 MDM（Multi-domain MDM）
* **特點**：
  * 全球市場佔有率高
  * 數據治理、質量、集成能力完整
* **適用場景**：大型跨國企業、集團型 MDM

### 2. SAP Master Data Governance（MDG）

* **官網**：[https://www.sap.com](https://www.sap.com)
* **類型**：ERP 內嵌型 MDM
* **特點**：
  * 與 SAP S/4HANA 深度綁定
  * 製造業、集團客戶廣泛

### 3. Oracle Master Data Management

* **官網**：[https://www.oracle.com](https://www.oracle.com)
* **類型**：企業級多域 MDM
* **特點**：
  * Oracle 應用與數據庫生態強

### 4. Stibo Systems STEP

* **官網**：[https://www.stibosystems.com](https://www.stibosystems.com)
* **類型**：PIM + MDM
* **特點**：
  * 零售、消費品、電商強項

### 5. Reltio Connected Data Platform

* **官網**：[https://www.reltio.com](https://www.reltio.com)
* **類型**：雲原生 MDM / Customer 360
* **特點**：
  * 客戶主數據、關係型數據建模

### 6. Profisee MDM

* **官網**：[https://www.profisee.com](https://www.profisee.com)
* **類型**：多域 MDM
* **特點**：
  * Azure / Microsoft 生態友好

### 7. Semarchy xDM

* **官網**：[https://www.semarchy.com](https://www.semarchy.com)
* **類型**：模型驅動 MDM
* **特點**：
  * 輕量、快速實施

### 8. Talend MDM

* **官網**：[https://www.talend.com](https://www.talend.com)
* **類型**：MDM + 數據集成

### 9. IBM InfoSphere MDM

* **官網**：[https://www.ibm.com](https://www.ibm.com)
* **類型**：企業級 MDM

## 二、中國商業主數據管理系統

### 1. 用友 主數據中臺 / MDM

* **官網**：[https://www.yonyou.com](https://www.yonyou.com)
* **類型**：集團級 MDM
* **特點**：
  * 多組織、多法人治理
  * 與財務、供應鏈深度集成

### 2. 金蝶 主數據管理平臺

* **官網**：[https://www.kingdee.com](https://www.kingdee.com)
* **類型**：ERP 綁定型 MDM

### 3. 普元 MDM（Primeton）

* **官網**：[https://www.primeton.com](https://www.primeton.com)
* **類型**：國產多域 MDM
* **特點**：
  * 國產化替代
  * 政企客戶

### 4. 英諾森 InData MDM

* **官網**：[https://www.inossem.com](https://www.inossem.com)
* **類型**：工業主數據 MDM
* **特點**：
  * 物料 / 設備 / BOM 主數據

### 5. 漢得 主數據平臺（HAND MDM）

* **官網**：[https://www.hand-china.com](https://www.hand-china.com)
* **類型**：實施型 MDM 平臺

### 6. 三維天地 主數據管理系統

* **官網**：[https://www.sunwayworld.com](https://www.sunwayworld.com)
* **類型**：國產 MDM

## 三、開源 / 社區版主數據系統

> ⚠️ 多數開源 MDM 偏產品數據（PIM）或輕量 MDM，企業級使用需二次開發。

### 1. Pimcore

* **官網**：[https://pimcore.com](https://pimcore.com)
* **授權**：GPL（社區版）+ 企業版收費
* **類型**：PIM / MDM / DAM

### 2. AtroCore

* **官網**：[https://www.atrocore.com](https://www.atrocore.com)
* **授權**：開源 + SaaS
* **類型**：PIM + MDM

### 3. OpenMDM

* **官網**：[https://github.com/openmdm](https://github.com/openmdm)
* **授權**：開源
* **類型**：基礎 MDM 框架

### 4. HiveMDM

* **授權**：開源
* **類型**：技術型 MDM 方案

### 5. PieMDM

* **官網 / 項目地址**：[https://github.com/pieworks/piemdm](https://github.com/pieworks/piemdm)
* **項目網站**：[https://pieworks.github.io/piemdm/](https://pieworks.github.io/piemdm/)
* **授權**：MIT 開源許可
* **類型**：開源企業級 MDM 系統
* **技術棧**：Go 後端 + Vue 前端
* **主要功能**：
  * 主數據建模與數據管理
  * 數據治理與訪問控制
  * 工作流 / 審批管理
  * 系統集成 & API 支持
* **適用場景**：
  * 適合嘗試性、PoC 或希望基於源碼擴展的主數據項目
  * 社區實踐、自主研發型企業

## 四、按維度快速索引

### 商業 vs 開源

* **商業**：Informatica / SAP / Oracle / Stibo / 用友 / 普元
* **開源**：Pimcore / AtroCore / OpenMDM / PieMDM

### 國際 vs 中國

* **國際**：Informatica / SAP / Oracle / Stibo / Reltio
* **中國**：用友 / 金蝶 / 普元 / 英諾森 / 漢得 / PieMDM

## 五、選型建議（簡要）

* **全球集團 / 複雜治理**：Informatica / SAP MDG
* **製造業主數據**：SAP MDG / 英諾森
* **國產化要求**：用友 / 普元 / 金蝶 / PieMDM
* **低成本探索 / 原型驗證**：Pimcore / AtroCore / PieMDM

# List of Master Data Management (MDM) Systems

> This document organizes mainstream **Master Data Management (MDM)** products, covering:
>
> * 🌍 International Vendors / 🇨🇳 Chinese Vendors
> * 💰 Commercial (Paid) / 🔓 Open Source (Community Edition)
> * 🧩 Multi-domain MDM / Specific MDM (Customer, Product, Supplier, etc.)

Applicable for **MDM platform selection, technical benchmarking, solution design, RFP preparation**, etc.

## I. International Mainstream Commercial MDM (Paid)

### 1. Informatica MDM

* **Website**: [https://www.informatica.com](https://www.informatica.com)
* **Type**: Multi-domain MDM
* **Features**:
  * High global market share
  * Complete data governance, quality, and integration capabilities
* **Applicable Scenarios**: Large multinational enterprises, Group-level MDM

### 2. SAP Master Data Governance (MDG)

* **Website**: [https://www.sap.com](https://www.sap.com)
* **Type**: ERP Embedded MDM
* **Features**:
  * Deeply bound with SAP S/4HANA
  * Widely used by manufacturing and group customers

### 3. Oracle Master Data Management

* **Website**: [https://www.oracle.com](https://www.oracle.com)
* **Type**: Enterprise Multi-domain MDM
* **Features**:
  * Strong Oracle application and database ecosystem

### 4. Stibo Systems STEP

* **Website**: [https://www.stibosystems.com](https://www.stibosystems.com)
* **Type**: PIM + MDM
* **Features**:
  * Strong in retail, consumer goods, and E-commerce

### 5. Reltio Connected Data Platform

* **Website**: [https://www.reltio.com](https://www.reltio.com)
* **Type**: Cloud-native MDM / Customer 360
* **Features**:
  * Customer master data, relationship data modeling

### 6. Profisee MDM

* **Website**: [https://www.profisee.com](https://www.profisee.com)
* **Type**: Multi-domain MDM
* **Features**:
  * Azure / Microsoft ecosystem friendly

### 7. Semarchy xDM

* **Website**: [https://www.semarchy.com](https://www.semarchy.com)
* **Type**: Model-driven MDM
* **Features**:
  * Lightweight, rapid implementation

### 8. Talend MDM

* **Website**: [https://www.talend.com](https://www.talend.com)
* **Type**: MDM + Data Integration

### 9. IBM InfoSphere MDM

* **Website**: [https://www.ibm.com](https://www.ibm.com)
* **Type**: Enterprise MDM

## II. Chinese Commercial Master Data Management Systems

### 1. Yonyou (用友) MDM

* **Website**: [https://www.yonyou.com](https://www.yonyou.com)
* **Type**: Group-level MDM
* **Features**:
  * Multi-organization, multi-legal entity governance
  * Deep integration with finance and supply chain

### 2. Kingdee (金蝶) MDM Platform

* **Website**: [https://www.kingdee.com](https://www.kingdee.com)
* **Type**: ERP Bound MDM

### 3. Primeton (普元) MDM

* **Website**: [https://www.primeton.com](https://www.primeton.com)
* **Type**: Domestic Multi-domain MDM
* **Features**:
  * Localization replacement
  * Government and Enterprise customers

### 4. Inossem InData MDM

* **Website**: [https://www.inossem.com](https://www.inossem.com)
* **Type**: Industrial MDM
* **Features**:
  * Material / Equipment / BOM Master Data

### 5. HAND (汉得) MDM

* **Website**: [https://www.hand-china.com](https://www.hand-china.com)
* **Type**: Implementation-based MDM Platform

### 6. Sunway World (三维天地) MDM

* **Website**: [https://www.sunwayworld.com](https://www.sunwayworld.com)
* **Type**: Domestic MDM

## III. Open Source / Community Edition MDM Systems

> ⚠️ Most open source MDM are PIM-oriented or lightweight. Enterprise use typically requires secondary development.

### 1. Pimcore

* **Website**: [https://pimcore.com](https://pimcore.com)
* **License**: GPL (Community) + Enterprise Paid
* **Type**: PIM / MDM / DAM

### 2. AtroCore

* **Website**: [https://www.atrocore.com](https://www.atrocore.com)
* **License**: Open Source + SaaS
* **Type**: PIM + MDM

### 3. OpenMDM

* **Website**: [https://github.com/openmdm](https://github.com/openmdm)
* **License**: Open Source
* **Type**: Basic MDM Framework

### 4. HiveMDM

* **License**: Open Source
* **Type**: Technical MDM Solution

### 5. PieMDM

* **Website / Repo**: [https://github.com/pieworks/piemdm](https://github.com/pieworks/piemdm)
* **Project Site**: [https://pieworks.github.io/piemdm/](https://pieworks.github.io/piemdm/)
* **License**: MIT Open Source
* **Type**: Open Source Enterprise MDM System
* **Tech Stack**: Go Backend + Vue Frontend
* **Key Features**:
  * Master Data Modeling & Management
  * Data Governance & Access Control
  * Workflow / Approval Management
  * System Integration & API Support
* **Applicable Scenarios**:
  * Suitable for pilot projects, PoC, or MDM projects desiring source code extension
  * Community practice, self-research enterprises

## IV. Quick Index by Dimension

### Commercial vs Open Source

* **Commercial**: Informatica / SAP / Oracle / Stibo / Yonyou / Primeton
* **Open Source**: Pimcore / AtroCore / OpenMDM / PieMDM

### International vs Chinese

* **International**: Informatica / SAP / Oracle / Stibo / Reltio
* **Chinese**: Yonyou / Kingdee / Primeton / Inossem / HAND / PieMDM

## V. Selection Suggestions (Brief)

* **Global Group / Complex Governance**: Informatica / SAP MDG
* **Manufacturing Master Data**: SAP MDG / Inossem
* **Localization Requirement**: Yonyou / Primeton / Kingdee / PieMDM
* **Low Cost Exploration / Prototype Verification**: Pimcore / AtroCore / PieMDM

# Dictionary Management

Dictionary management is used to maintain various enumeration values or constant lists in the system (such as "Country", "Currency", "Status Code", etc.). In PieMDM, dictionaries serve as important data source configurations and are widely used in fields like dropdown selects and radio buttons.

## Dictionary Categories

PieMDM supports two forms of dictionaries:

### 1. System Common Dictionary
The system provides a built-in table named `dict_item` by default for unified storage of simple key-value pair data.
- **Storage**: All common dictionary items are stored in the `dict_item` table.
- **Identifier**: Use `dict_code` (Dictionary Code) to distinguish different dictionary categories (e.g., `GENDER` for gender, `ORDER_STATUS` for order status).

### 2. Business Table Based Dictionary
You can directly use any business model (such as "Material Table", "Supplier Table") as a dictionary.
- **Dynamic**: Dictionary content changes in real-time with business data additions, deletions, and modifications.
- **Flexible Configuration**: You can specify which column in the table serves as the "Data Key (Value)" and which column serves as the "Display Label (Label)".

## Configuring Dictionary Fields

In "Table & Field Management", when you select a "Select Type" for a field (such as Select, Checkbox Group), you need to configure the data source.

### Associate System Dictionary
1. In the "Relation Config" of the field, select `dict_item` as the target table.
2. In "Filter Condition", set `dict_code` to the corresponding dictionary code (e.g., `NATION`).
3. The Value Field is usually configured as `code`, and the Display Field is usually configured as `name`.

### Associate Business Table
1. Select the corresponding business table (e.g., `mdm_company`) as the target table.
2. Set filter conditions as needed (e.g., `status = 'Normal'`).
3. Specify the value field for storage (e.g., `id` or `company_code`) and the label field for display (e.g., `company_name`).

## Dictionary Service Features

PieMDM frontend provides an efficient dictionary management service (`dictionaryService`) with the following features:

- **Caching Mechanism**: When the same dictionary is referenced multiple times on a page, only one API request is made, and subsequent accesses read directly from memory cache.
- **Request Merging**: If multiple components request the same dictionary simultaneously, the service automatically merges requests to avoid network redundancy.
- **Preloading**: Supports automatic scanning and preloading of all required dictionary data based on field configuration during page initialization, improving first-screen rendering speed.
- **Universal Adaptation**: Whether it's the built-in `dict_item` or a regular business table, data is fetched through a unified API interface.

## Management Suggestions

1. **Unified Coding**: For basic data across systems, it is recommended to refer to national or industry standards for coding (e.g., ISO currency codes).
2. **Clear Description**: Provide clear descriptions for each dictionary item to help users understand its business meaning.
3. **Decentralized Management**: Core system dictionaries (such as system configuration items) should only be modified by administrators.

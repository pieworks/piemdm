<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ entityName }} {{ $t('Update') }}</div>

    <div class="card-body mt-2">
      <div id="create_wrapper">
        <AppForm v-model:dataInfo="dataInfo" :fields="tableFields" :dictionarys="dictionarys"
          :readonlyFields="readonlyFields" @fetch-data="fetchData" />

        <div class="form-group row mt-3 mb-3">
          <div class="col-sm-auto">
            <button type="button" class="btn btn-outline-primary btn-sm me-2" @click="submitForm(dataInfo)">
              {{ $t('Submit') }}
            </button>
            <button type="button" class="btn btn-outline-secondary btn-sm me-2" @click="goIndex">
              {{ $t('Cancel') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {
  findEntity,
  getEntityList,
  updateEntity,
} from '@/api/entity';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import AppForm from './Form.vue';
import { getTableList } from '@/api/table';
import { AppModal } from '@/components/Modal/modal';

const router = useRouter();
const { t } = useI18n();
const dataInfo = ref({});
const originalData = ref({}); // Store original data for change detection
const tableFields = ref([]);
const params = ref({});
const dictionarys = ref([]);
const entityName = ref('');
const readonlyFields = ref([]);

onMounted(async () => {
  params.value = router.currentRoute.value.query;
  await getEntityName();
  await getEntityInfo(params.value.table_code);
});

const displayMode = ref('List');

// Get table name
const getEntityName = async () => {
  const where = {
    status: 'Normal',
  };
  if (params.value.table_code) {
    where.code = params.value.table_code;
  }
  const res = await getTableList(where);
  if (res && res.data && res.data.length > 0) {
    entityName.value = res.data[0].Name;
    displayMode.value = res.data[0].DisplayMode || res.data[0].display_mode || 'List';
  }
};

const getEntityInfo = async table_code => {
  const res = await findEntity({
    table_code: table_code,
    id: params.value.id,
  });
  if (res.data) {
    // Format date and datetime fields
    const formattedData = { ...res.data.info };
    res.data.tableFields.forEach(field => {
      const fieldValue = formattedData[field.Code];

      // Skip empty values
      if (!fieldValue || typeof fieldValue !== 'string') {
        return;
      }

      try {
        // Date field: ensure format is yyyy-MM-dd
        if (field.FieldType === 'date' || field.Type === 'Date') {
          let dateStr = fieldValue;

          // Remove timezone info +08:00 or Z
          dateStr = dateStr.replace(/[+\-]\d{2}:\d{2}$/, '').replace(/Z$/, '');

          // If contains time part, only take date
          if (dateStr.includes('T')) {
            formattedData[field.Code] = dateStr.split('T')[0];
          } else if (dateStr.includes(' ')) {
            formattedData[field.Code] = dateStr.split(' ')[0];
          } else {
            formattedData[field.Code] = dateStr;
          }
        }
        // DateTime field: ensure format is yyyy-MM-dd HH:mm:ss
        else if (field.FieldType === 'datetime' || field.Type === 'DateTime') {
          let datetime = fieldValue.trim();

          // Remove timezone info +08:00 or Z
          datetime = datetime.replace(/[+\-]\d{2}:\d{2}$/, '').replace(/Z$/, '');

          // Remove milliseconds .000
          if (datetime.includes('.')) {
            datetime = datetime.split('.')[0];
          }

          // Replace T with space (VueDatePicker model-type uses space format)
          if (datetime.includes('T')) {
            datetime = datetime.replace('T', ' ');
          }

          // Ensure correct format: yyyy-MM-dd HH:mm:ss
          if (datetime.match(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/)) {
            formattedData[field.Code] = datetime;
          } else {
            console.warn(`Invalid datetime format for ${field.Code}:`, fieldValue, '→', datetime);
          }
        }
        // Attachment, checkbox group, multi-select fields: deserialize JSON array
        else if (field.FieldType === 'attachment' ||
          field.FieldType === 'checkboxgroup' ||
          field.FieldType === 'multiselect') {
          if (typeof fieldValue === 'string' && fieldValue.startsWith('[')) {
            try {
              formattedData[field.Code] = JSON.parse(fieldValue);
            } catch (e) {
              console.warn(`Failed to parse array field ${field.Code}:`, fieldValue);
              // If parsing fails, initialize as empty array
              formattedData[field.Code] = [];
            }
          } else if (!fieldValue) {
            // If value is empty, initialize as empty array
            formattedData[field.Code] = [];
          }
        }
      } catch (error) {
        console.error(`Error formatting field ${field.Code}:`, error, fieldValue);
      }
    });

    dataInfo.value = formattedData;
    // Save deep copy of original data for change detection
    originalData.value = JSON.parse(JSON.stringify(formattedData));

    const dicts = {};

    // Use Promise.all to wait for all async operations to complete (consistent with create.vue)
    const promises = res.data.tableFields.map(async field => {
      // Get relation config
      const relationTable = field.Options?.relation?.target;
      const relationFilter = field.Options?.relation?.filter;

      // Load options for all fields with relation config
      // Exclude empty strings and undefined
      if (relationTable && relationTable.trim() !== '') {
        try {
          const queryParams = {
            table_code: relationTable,
            pageSize: 1000, // Load enough options to match initial value
            ...relationFilter, // Apply filter conditions
          };

          const res2 = await getEntityList(queryParams);

          if (res2 && res2.data) {
            dicts['dict-' + field.Code] = res2.data;
          }
        } catch (error) {
          console.error('Error loading options for', field.Code, ':', error);
        }
      }
    });

    // Wait for all options to load
    await Promise.all(promises);

    // Set options first, then set tableFields, ensure options are ready when form renders
    dictionarys.value = dicts;
    // Filter out system fields, only keep business fields
    tableFields.value = res.data.tableFields.filter(field => {
      if (field.IsSystem) return false;
      // If tree table, filter out level and path
      if (displayMode.value === 'Tree') {
        const code = (field.Code || field.code || '').toLowerCase();
        if (code === 'level' || code === 'path') {
          return false;
        }
      }
      return true;
    });

    // Get table config, identify relation field (ItemField) and set as readonly
    try {
      const tableInfo = await getTableList({
        code: table_code,
        status: 'Normal',
      });
      if (tableInfo && tableInfo.data && tableInfo.data.length > 0) {
        const itemField = tableInfo.data[0].ItemField;
        if (itemField) {
          readonlyFields.value = [itemField];
        }
      }
    } catch (error) {
      console.error('Failed to get table config:', error);
    }
  }
};

const fetchData = async (search, loading, field) => {
  const relationTable = field.Options?.relation?.target;

  // If no relation table config, return directly
  if (!relationTable) {
    return;
  }

  if (loading) loading(true);
  try {
    const relationCode = field.Options?.relation?.valueField;
    const relationName = field.Options?.relation?.labelField;

    // Read filter conditions from relation.filter
    const filterConditions = field.Options?.relation?.filter || {};

    // Build query params
    const queryParams = {
      table_code: relationTable,
      pageSize: 100,
      // Add all filter conditions
      ...filterConditions
    };

    // Only add search condition when search string is not empty and has valueField/labelField config
    if (search && search.length > 0 && relationCode && relationName) {
      queryParams['(' + relationCode + ' like ? or ' + relationName + ' like ? )'] =
        '%' + search + '%,%' + search + '%';
    }

    const res = await getEntityList(queryParams);

    if (res && res.data) {
      // Use reactive update, avoid directly clearing array
      dictionarys.value = {
        ...dictionarys.value,
        ['dict-' + field.Code]: res.data
      };
    }
  } catch (error) {
    console.error('Error fetching options:', error);
  } finally {
    if (loading) loading(false);
  }
};

const submitForm = async data => {
  try {
    // Required field validation
    const errors = {};
    // Required can be boolean true or string "Yes"
    const requiredFields = tableFields.value.filter(field =>
      field.Required === true || field.Required === 'Yes'
    );

    for (const field of requiredFields) {
      const value = data[field.Code];
      // Check if empty value
      const isEmpty = value === undefined ||
        value === null ||
        value === '' ||
        (Array.isArray(value) && value.length === 0);

      if (isEmpty) {
        errors[field.Code] = t('{fieldName} is required', { fieldName: field.Name });
      }
    }

    // Hard-coded validation for reason field
    if (!data.reason || data.reason.trim() === '') {
      errors['reason'] = t('Modification reason is required');
    }

    // If there are validation errors, show message and return
    if (Object.keys(errors).length > 0) {
      const errorList = Object.values(errors).join('<br>');
      await AppModal.alert({
        title: t('Validation Failed'),
        okTitle: t('OK'),
        bodyHtml: true,
        bodyContent: `<p>${t('Please fill in the following required fields:')}</p><p style="color: red;">${errorList}</p>`
      });
      return;
    }

    // Filter out autocode fields to prevent modification
    const filteredData = { ...data };
    tableFields.value.forEach(field => {
      if (field.FieldType === 'autocode') {
        delete filteredData[field.Code];
      }
    });

    // Filter out system fields (backend will set automatically)
    delete filteredData.created_at;
    delete filteredData.updated_at;
    delete filteredData.created_by;
    delete filteredData.updated_by;
    delete filteredData.deleted_at;

    const res = await updateEntity({
      table_code: params.value.table_code,
      ...filteredData,
    });

    // Only show success Modal on successful response
    if (res && res.status === 200) {
      const result = await AppModal.alert({
        title: 'Success',
        bodyContent: 'Your submission has been approved。',
      });
      if (!result) {
        return;
      }
      goIndex();
    }
  } catch (error) {
    // Error already handled and displayed in axios interceptor
    console.error('Update entity failed:', error);
  }
};

function goIndex() {
  if (params.value.redirect_url) {
    const redirectUrl = decodeURIComponent(params.value.redirect_url);
    router.push(redirectUrl);
  } else {
    router.push('/entity/index?table_code=' + params.value.table_code);
  }
}
</script>

<style scoped></style>

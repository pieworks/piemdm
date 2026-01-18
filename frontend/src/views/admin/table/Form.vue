<template>
  <form name="entityForm" id="entityForm" method="post" enctype="multipart/form-data">
    <div class="col-12 col-sm-12">
      <div class="row">
        <div class="col-sm-12">

          <div class="col-sm-12">
            <div class="form-group row">
              <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.code }]">
                {{ $t('Code') }}:
              </legend>
              <div class="col-sm-auto">
                <input type="text" class="form-control form-control-sm" v-model="code" v-bind="codeAttrs" name="code"
                  :placeholder="$t('Code')" maxlength="64" size="64" :disabled="!!props.dataInfo?.ID" />
                <div v-if="errors.Code" class="text-danger small mt-1">
                  {{ errors.Code }}
                </div>
              </div>
            </div>
          </div>
          <div class="col-sm-12">
            <div class="form-group row">
              <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.name }]">
                {{ $t('Name') }}:
              </legend>
              <div class="col-sm-auto">
                <input type="text" class="form-control form-control-sm" v-model="name" v-bind="nameAttrs" name="name"
                  :placeholder="$t('Name')" maxlength="64" size="64" />
                <div v-if="errors.Name" class="text-danger small mt-1">
                  {{ errors.Name }}
                </div>
              </div>
            </div>
          </div>
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.displaymode }]">
              {{ $t('DisplayMode') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect v-model="displayMode" v-bind="displayModeAttrs" name="displayMode"
                :reduce="option => option.value" :placeholder="$t('Please Select')" :options="tableContentOptions">
              </VSelect>
              <div v-if="errors.DisplayMode" class="text-danger small mt-1">
                {{ errors.DisplayMode }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.tabletype }]">
              {{ $t('TableType') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect v-model="tableType" v-bind="tableTypeAttrs" name="tableType" :reduce="option => option.value"
                :placeholder="$t('Please Select')" :options="tableRelationOptions"></VSelect>

              <div v-if="errors.TableType" class="text-danger small mt-1">
                {{ errors.TableType }}
              </div>
            </div>
          </div>
        </div>
        <!-- Show Tree Checkbox (only for Entity type) -->
        <div class="col-sm-12" v-if="tableType === 'Entity'">
          <div class="form-group row">
            <legend class="col-form-label col-sm-2">
              {{ $t('ShowTree') }}:
            </legend>
            <div class="col-sm-10">
              <div class="form-check">
                <input class="form-check-input" type="checkbox" id="showTree" v-model="showTree" />
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12" v-if="tableType === 'Item' || (tableType === 'Entity' && showTree)">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.parenttable }]">
              {{ $t('ParentTable') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect v-model="parentTable" v-bind="parentTableAttrs" name="ParentTable"
                :reduce="option => option.value" :placeholder="$t('Please Select')" :options="parentTableOptions"
                :loading="loadingParentTables">
                <template v-slot:option="option">
                  {{ option.value }} {{ option.label.split(' ')[1] }}
                </template>
                <template v-slot:selected-option="option">
                  {{ option.value }} {{ option.label.split(' ')[1] }}
                </template>
              </VSelect>
              <div v-if="errors.ParentTable" class="text-danger small mt-1">
                {{ errors.ParentTable }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12" v-if="tableType === 'Item' || (tableType === 'Entity' && showTree)">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.parentfield }]">
              {{ $t('ParentField') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect v-model="parentField" v-bind="parentFieldAttrs" name="ParentField"
                :reduce="option => option.value" :placeholder="$t('Please Select')" :options="parentFieldOptions"
                :loading="loadingParentFields" :disabled="!parentTable">
                <template v-slot:option="option">
                  {{ option.value }} {{ option.label.split(' ')[1] }}
                </template>
                <template v-slot:selected-option="option">
                  {{ option.value }} {{ option.label.split(' ')[1] }}
                </template>
              </VSelect>
              <div v-if="errors.ParentField" class="text-danger small mt-1">
                {{ errors.ParentField }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12" v-if="tableType === 'Item' || (tableType === 'Entity' && showTree)">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.selffield }]">
              {{ $t('SelfField') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect v-model="selfField" v-bind="selfFieldAttrs" name="SelfField" :reduce="option => option.value"
                :placeholder="$t('Please Select')" :options="selfFieldOptions" :loading="loadingSelfFields"
                :disabled="!code">
                <template v-slot:option="option">
                  {{ option.value }} {{ option.label.split(' ')[1] }}
                </template>
                <template v-slot:selected-option="option">
                  {{ option.value }} {{ option.label.split(' ')[1] }}
                </template>
              </VSelect>
              <div v-if="errors.SelfField" class="text-danger small mt-1">
                {{ errors.SelfField }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.sort }]">
              {{ $t('Sort') }}:
            </legend>
            <div class="col-sm-auto">
              <input type="number" class="form-control form-control-sm" v-model="sort" v-bind="sortAttrs" name="sort"
                :placeholder="$t('Sort')" min="1" />
              <div v-if="errors.Sort" class="text-danger small mt-1">
                {{ errors.Sort }}
              </div>
            </div>
          </div>
        </div>

        <div class="col-sm-12 mb-1">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.description }]">
              {{ $t('Description') }}:
            </legend>
            <div class="col-sm-6">
              <textarea class="form-control form-control-sm" v-model="description" v-bind="descriptionAttrs"
                name="description" :placeholder="$t('Description')" maxlength="255" rows="3"></textarea>
              <div v-if="errors.Description" class="text-danger small mt-1">
                {{ errors.Description }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.status }]">
              {{ $t('Status') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect v-model="status" v-bind="statusAttrs" name="status" :reduce="option => option.value"
                :options="statusOptions"></VSelect>
              <div v-if="errors.Status" class="text-danger small mt-1">
                {{ errors.Status }}
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>
    <div class="form-group row my-3">
      <div class="col-sm-auto">
        <input v-if="props.dataInfo && props.dataInfo.ID" type="hidden" name="id" :value="props.dataInfo.ID" />
        <button type="button" class="btn btn-outline-primary btn-sm me-2" @click="onSubmit()">
          {{ $t('Submit') }}
        </button>
        <button type="button" class="btn btn-outline-secondary btn-sm" @click="$emit('goIndex')">
          {{ $t('Cancel') }}
        </button>
      </div>
    </div>
  </form>
</template>

<script setup>
import { useFormOptions } from '@/composables/useFormOptions';
import { useForm } from 'vee-validate';
import { computed, watch, ref, onMounted } from 'vue';
import VSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';
import * as yup from 'yup';
import { getTableList } from '@/api/table';
import { getTableFieldList } from '@/api/table_field';

const { statusOptions, tableContentOptions, tableRelationOptions } = useFormOptions();

const parentTableOptions = ref([]);
const parentFieldOptions = ref([]);
const selfFieldOptions = ref([]);
const loadingParentTables = ref(false);
const loadingParentFields = ref(false);
const loadingSelfFields = ref(false);
const showTree = ref(false);

// Define props to receive dataInfo from parent component
const props = defineProps({
  dataInfo: {
    type: Object,
    default: () => ({}),
  },
});


const validationSchema = yup.object({
  Name: yup.string().required(),
  Code: yup.string().required(),
  DisplayMode: yup.string().required(),
  TableType: yup.string().required(),
  ParentTable: yup.string(),
  ParentField: yup.string(),
  SelfField: yup.string(),
  Status: yup.string().required(),
  Sort: yup.number(),
  Description: yup.string().max(255),
});


const requiredFields = computed(() => {
  const requiredMap = {};
  Object.keys(validationSchema.fields).forEach(key => {
    const field = validationSchema.fields[key];
    requiredMap[key.toLowerCase()] = field.tests.some(test => test.OPTIONS?.name === 'required');
  });
  return requiredMap;
});


const { values, errors, defineField, handleSubmit, setValues } = useForm({
  validationSchema,
  initialValues: props.dataInfo,
});


const [code, codeAttrs] = defineField('Code');
const [name, nameAttrs] = defineField('Name');
const [displayMode, displayModeAttrs] = defineField('DisplayMode');
const [sort, sortAttrs] = defineField('Sort');
const [tableType, tableTypeAttrs] = defineField('TableType');
const [parentTable, parentTableAttrs] = defineField('ParentTable');
const [parentField, parentFieldAttrs] = defineField('ParentField');
const [selfField, selfFieldAttrs] = defineField('SelfField');
const [status, statusAttrs] = defineField('Status');
const [description, descriptionAttrs] = defineField('Description');

const emit = defineEmits(['submitForm', 'goIndex']);

const onSubmit = handleSubmit(values => {
  emit('submitForm', values);
});


watch(
  () => props.dataInfo,
  newDataInfo => {
    if (newDataInfo && Object.keys(newDataInfo).length > 0) {
      setValues(newDataInfo);
      // Auto-check showTree if ParentTable is set for Entity type
      if (newDataInfo.TableType === 'Entity' && newDataInfo.ParentTable) {
        showTree.value = true;
      }
    }
  },
  { immediate: true, deep: true }
);

const loadTableOptions = async (filterFn, targetRef, loadingRef) => {
  loadingRef.value = true;
  try {
    const res = await getTableList({ table_type: 'Entity', pageSize: -1 });
    if (res?.data) {
      const tables = filterFn ? res.data.filter(filterFn) : res.data;
      targetRef.value = tables.map(table => ({
        label: `${table.Code} ${table.Name}`,
        value: table.Code
      }));
    }
  } catch (error) {
    console.error('Failed to load table options:', error);
  } finally {
    loadingRef.value = false;
  }
};

const loadParentTables = () => {
  // For Entity type with showTree, only show Tree tables
  // For Item type, show all Entity tables
  const filterFn = (tableType.value === 'Entity' && showTree.value)
    ? (table => table.DisplayMode === 'Tree' || table.display_mode === 'Tree')
    : null;

  loadTableOptions(
    filterFn,
    parentTableOptions,
    loadingParentTables
  );
};


const loadParentFields = async (tableCode) => {
  if (!tableCode) {
    parentFieldOptions.value = [];
    return;
  }
  loadingParentFields.value = true;
  try {
    const res = await getTableFieldList({ table_code: tableCode, pageSize: -1 });
    if (res && res.data) {
      parentFieldOptions.value = res.data.map(field => ({
        label: `${field.Code} ${field.Name}`,
        value: field.Code
      }));
    }
  } catch (error) {
    console.error('Failed to load parent fields:', error);
  } finally {
    loadingParentFields.value = false;
  }
};


const loadSelfFields = async (tableCode) => {
  if (!tableCode) {
    selfFieldOptions.value = [];
    return;
  }
  loadingSelfFields.value = true;
  try {
    const res = await getTableFieldList({ table_code: tableCode, pageSize: -1 });
    if (res && res.data) {
      selfFieldOptions.value = res.data.map(field => ({
        label: `${field.Code} ${field.Name}`,
        value: field.Code
      }));
    }
  } catch (error) {
    console.error('Failed to load self fields:', error);
  } finally {
    loadingSelfFields.value = false;
  }
};


watch(() => parentTable.value, (newValue) => {
  if (newValue) {
    loadParentFields(newValue);
  } else {
    parentFieldOptions.value = [];
    parentField.value = '';
  }
});


// Watch code changes to load self fields
watch(() => code.value, (newValue) => {
  if (newValue) {
    loadSelfFields(newValue);
  } else {
    selfFieldOptions.value = [];
    selfField.value = '';
  }
});


onMounted(() => {
  loadParentTables();

  if (props.dataInfo?.ParentTable) {
    loadParentFields(props.dataInfo.ParentTable);
    // Auto-check showTree if ParentTable is set for Entity type
    if (tableType.value === 'Entity') {
      showTree.value = true;
    }
  }

  // Load self fields if code exists
  if (props.dataInfo?.Code) {
    loadSelfFields(props.dataInfo.Code);
  }
});
</script>

<style scoped>
.required:after {
  content: ' *';
  color: #dc3545;
  font-weight: bold;
}
</style>

<template>
  <form name="fieldForm" id="fieldForm" @submit.prevent="onSubmit">
    <h5 class="mb-3">基本信息</h5>

    <!-- 字段类型 -->
    <div class="form-group row mb-3">
      <label class="col-form-label col-sm-2 required">字段类型:</label>
      <div class="col-sm-4">
        <select v-if="!dataInfo?.ID" v-model="formData.fieldType" class="form-select form-select-sm" required>
          <option value="">请选择字段类型</option>
          <optgroup v-for="group in fieldTypeGroups" :key="group.name" :label="group.label">
            <option v-for="type in group.types" :key="type" :value="type">
              {{ fieldTypePresets[type].label }}
            </option>
          </optgroup>
        </select>
        <input v-else type="text" class="form-control form-control-sm" :value="selectedTypePreset?.label" disabled />
      </div>
    </div>

    <!-- 字段标识 -->
    <div class="form-group row mb-2">
      <label class="col-form-label col-sm-2 required">字段标识:</label>
      <div class="col-sm-auto">
        <input type="text" class="form-control form-control-sm" v-model="formData.code" placeholder="字段标识"
          maxlength="64" size="64" required :disabled="!!dataInfo?.ID" />
        <div v-if="errors.code" class="text-danger small mt-1">
          {{ errors.code }}
        </div>
      </div>
    </div>

    <!-- 字段名称 -->
    <div class="form-group row mb-2">
      <label class="col-form-label col-sm-2 required">字段名称:</label>
      <div class="col-sm-auto">
        <input type="text" class="form-control form-control-sm" v-model="formData.name" placeholder="字段名称"
          maxlength="128" size="64" required />
        <div v-if="errors.name" class="text-danger small mt-1">
          {{ errors.name }}
        </div>
      </div>
    </div>

    <!-- 排序 -->
    <div class="form-group row mb-2">
      <label class="col-form-label col-sm-2">排序:</label>
      <div class="col-sm-2">
        <input type="number" class="form-control form-control-sm" v-model.number="formData.sort" placeholder="排序"
          min="0" />
      </div>
    </div>

    <!-- 分组名称 -->
    <div class="form-group row mb-2">
      <label class="col-form-label col-sm-2">分组名称:</label>
      <div class="col-sm-4">
        <input type="text" class="form-control form-control-sm" v-model="formData.groupName"
          placeholder="例如: 基本信息, 高级配置" maxlength="64" />
      </div>
    </div>

    <!-- 描述 -->
    <div class="form-group row mb-3">
      <label class="col-form-label col-sm-2">描述:</label>
      <div class="col-sm-6">
        <textarea class="form-control form-control-sm" v-model="formData.description" placeholder="描述" rows="3"
          maxlength="256"></textarea>
      </div>
    </div>

    <!-- 条件显示:关联配置 -->
    <div v-if="needsRelationConfig" class="relation-config mb-3">
      <h6 class="text-muted mb-2">关联配置</h6>
      <div class="form-group row mb-2">
        <label class="col-form-label col-sm-2 required">关联表:</label>
        <div class="col-sm-4">
          <v-select v-model="relationConfig.target" :options="tables" :reduce="table => table.Code" label="Name"
            :filterable="true" :searchable="true" placeholder="搜索并选择关联表..." class="form-select-sm">
            <template v-slot:option="option">
              {{ option.Code }} {{ option.Name }}
            </template>
            <template v-slot:selected-option="option">
              {{ option.Code }} {{ option.Name }}
            </template>
          </v-select>
        </div>
      </div>
      <div class="form-group row mb-2">
        <label class="col-form-label col-sm-2">存储字段:</label>
        <div class="col-sm-3">
          <v-select v-model="relationConfig.valueField" :options="targetTableFields" :reduce="field => field.code"
            label="name" :filterable="true" :filter-by="filterFieldOptions" placeholder="选择存储字段...">
            <template v-slot:option="option">
              {{ option.code }} {{ option.name }}
            </template>
            <template v-slot:selected-option="option">
              {{ option.code }} {{ option.name }}
            </template>
          </v-select>
        </div>
        <small class="col-sm-6 text-muted">保存到数据库的字段名</small>
      </div>
      <div class="form-group row mb-2">
        <label class="col-form-label col-sm-2">显示字段:</label>
        <div class="col-sm-3">
          <v-select v-model="relationConfig.labelField" :options="targetTableFields" :reduce="field => field.code"
            label="name" :filterable="true" :filter-by="filterFieldOptions" placeholder="选择显示字段...">
            <template v-slot:option="option">
              {{ option.code }} {{ option.name }}
            </template>
            <template v-slot:selected-option="option">
              {{ option.code }} {{ option.name }}
            </template>
          </v-select>
        </div>
        <small class="col-sm-6 text-muted">下拉列表显示的字段名</small>
      </div>
      <div class="form-group row mb-2">
        <label class="col-form-label col-sm-2">过滤条件:</label>
        <div class="col-sm-6">
          <div class="row g-2">
            <div class="col-sm-6">
              <v-select v-model="relationFilter.field" :options="targetTableFields" :reduce="field => field.code"
                label="name" :filterable="true" :filter-by="filterFieldOptions" placeholder="选择字段..."
                @option:selected="loadFilterFieldOptions">
                <template v-slot:option="option">
                  {{ option.code }} {{ option.name }}
                </template>
                <template v-slot:selected-option="option">
                  {{ option.code }} {{ option.name }}
                </template>
              </v-select>
            </div>
            <div class="col-sm-6">
              <!-- 字典表:一次加载全部数据,前端过滤搜索 -->
              <v-select v-if="relationConfig.target === 'dict_item' && relationFilter.field === 'dict_code'"
                v-model="relationFilter.value" :options="relationFilter.options" :reduce="opt => opt.code" label="name"
                :filterable="true" :filter-by="filterFieldOptions" placeholder="搜索字典分类..."
                @open="onOpenDictionaryClass">
                <template v-slot:option="option">
                  {{ option.code }} {{ option.name }}
                </template>
                <template v-slot:selected-option="option">
                  {{ option.code }} {{ option.name }}
                </template>
              </v-select>
              <!-- 其他字段:如果有选项显示 v-select,否则显示输入框 -->
              <v-select v-else-if="relationFilter.options && relationFilter.options.length > 0"
                v-model="relationFilter.value" :options="relationFilter.options" :reduce="opt => opt.code" label="name"
                :filterable="true" :filter-by="filterFieldOptions" placeholder="选择或搜索过滤值...">
                <template v-slot:option="option">
                  {{ option.code }} {{ option.name }}
                </template>
                <template v-slot:selected-option="option">
                  {{ option.code }} {{ option.name }}
                </template>
              </v-select>
              <input v-else type="text" v-model="relationFilter.value" class="form-control form-control-sm"
                placeholder="输入过滤值" />
            </div>
          </div>
        </div>
        <small class="col-sm-4 text-muted">
          根据指定字段过滤关联表数据
        </small>
      </div>
    </div>

    <!-- 条件显示:附件配置 -->
    <div v-if="needsAttachmentConfig" class="attachment-config mb-3">
      <h6 class="text-muted mb-2">附件配置</h6>

      <div class="form-group row mb-2">
        <label class="col-form-label col-sm-2">多文件上传:</label>
        <div class="col-sm-3">
          <select v-model="attachmentConfig.multiple" class="form-select form-select-sm">
            <option :value="false">单文件</option>
            <option :value="true">多文件</option>
          </select>
        </div>
        <small class="col-sm-6 text-muted">是否允许上传多个文件</small>
      </div>

      <div class="form-group row mb-2">
        <label class="col-form-label col-sm-2">允许类型:</label>
        <div class="col-sm-6">
          <input type="text" class="form-control form-control-sm" v-model="attachmentConfig.acceptInput"
            placeholder="如: image/*, .pdf, .doc, .docx" />
        </div>
        <small class="col-sm-3 text-muted">多个类型用逗号分隔</small>
      </div>

      <div class="form-group row mb-2">
        <label class="col-form-label col-sm-2">最大文件大小:</label>
        <div class="col-sm-3">
          <div class="input-group input-group-sm">
            <input type="number" class="form-control form-control-sm" v-model.number="attachmentConfig.maxSizeMB"
              placeholder="10" min="1" max="100" />
            <span class="input-group-text">MB</span>
          </div>
        </div>
        <small class="col-sm-6 text-muted">单个文件最大大小 (1-100 MB)</small>
      </div>
    </div>

    <!-- 条件显示:自动编码配置 -->
    <div v-if="needsAutocodeConfig" class="autocode-config mb-3">
      <h6 class="text-muted mb-2">自动编码配置</h6>

      <!-- 模式列表 -->
      <div v-for="(pattern, index) in autocodePatterns" :key="index" class="pattern-item mb-2 rounded">
        <div class="row g-2 align-items-center">
          <div class="col-sm-2">
            <select v-model="pattern.type" class="form-select form-select-sm">
              <option value="string">固定字符串</option>
              <option value="date">日期</option>
              <option value="field">数据字段</option>
              <option value="integer">序列号</option>
            </select>
          </div>

          <!-- string 配置 -->
          <div v-if="pattern.type === 'string'" class="col-sm-6">
            <input v-model="pattern.options.value" class="form-control form-control-sm" placeholder="如: C, PO-, INV" />
          </div>

          <!-- date 配置 -->
          <div v-if="pattern.type === 'date'" class="col-sm-6">
            <input v-model="pattern.options.format" class="form-control form-control-sm"
              placeholder="如: YYYYMMDD, YYMM, YYYY-MM-DD" />
          </div>

          <!-- field 配置 -->
          <div v-if="pattern.type === 'field'" class="col-sm-6">
            <select v-model="pattern.options.fieldCode" class="form-select form-select-sm">
              <option value="">请选择字段，值自动转大写,仅支持integer/string类型</option>
              <option v-for="field in currentTableBusinessFields" :key="field.code" :value="field.code">
                {{ field.name }} ({{ field.code }})
              </option>
            </select>
          </div>

          <!-- integer 配置 -->
          <div v-if="pattern.type === 'integer'" class="col-sm-6">
            <div class="row g-2">
              <div class="col">
                <input v-model.number="pattern.options.digits" type="number" class="form-control form-control-sm"
                  placeholder="位数(默认5)" min="1" max="10" />
              </div>
              <div class="col">
                <input v-model.number="pattern.options.start" type="number" class="form-control form-control-sm"
                  placeholder="起始值(默认1)" min="0" />
              </div>
              <div class="col">
                <select v-model="pattern.options.cycle" class="form-select form-select-sm">
                  <option value="none">不重置</option>
                  <option value="daily">每日重置</option>
                  <option value="monthly">每月重置</option>
                  <option value="yearly">每年重置</option>
                </select>
              </div>
            </div>
          </div>

          <div class="col-sm-auto">
            <button @click="removePattern(index)" class="btn btn-sm btn-outline-danger" type="button">
              <i class="bi bi-trash"></i>
            </button>
          </div>
        </div>
      </div>

      <button @click="addPattern" class="btn btn-sm btn-outline-primary" type="button">
        <i class="bi bi-plus-circle me-1"></i>添加模式
      </button>

      <!-- 预览 -->
      <div v-if="previewCode" class="mt-3 p-2 bg-light border rounded">
        <small class="text-muted">预览示例: </small>
        <code class="text-primary">{{ previewCode }}</code>
      </div>
    </div>

    <!-- 条件显示：日期时间配置 -->
    <div v-if="needsDateTimeConfig" class="datetime-config mb-3">
      <h6 class="text-muted mb-2">日期时间配置</h6>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="showTime" v-model="datetimeConfig.showTime" />
        <label class="form-check-label" for="showTime">
          显示时间选择器
        </label>
      </div>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="defaultToCurrentTime"
          v-model="datetimeConfig.defaultToCurrentTime" />
        <label class="form-check-label" for="defaultToCurrentTime">
          创建时默认当前时间
        </label>
      </div>
    </div>

    <!-- 字段配置选项 -->
    <h6 class="text-muted mb-2 mt-3">字段配置</h6>

    <!-- 是否必填 -->
    <div class="form-group row mb-2">
      <label class="col-form-label col-sm-2">是否必填:</label>
      <div class="col-sm-auto">
        <select v-model="formData.required" class="form-select form-select-sm">
          <option value="No">否</option>
          <option value="Yes">是</option>
        </select>
      </div>
    </div>

    <!-- 是否列表显示 -->
    <div class="form-group row mb-2">
      <label class="col-form-label col-sm-2">列表显示:</label>
      <div class="col-sm-auto">
        <select v-model="formData.isShow" class="form-select form-select-sm">
          <option value="Yes">是</option>
          <option value="No">否</option>
        </select>
      </div>
    </div>

    <!-- 是否作为筛选条件 -->
    <div class="form-group row mb-2">
      <label class="col-form-label col-sm-2">筛选条件:</label>
      <div class="col-sm-auto">
        <select v-model="formData.isFilter" class="form-select form-select-sm">
          <option value="No">否</option>
          <option value="Yes">是</option>
        </select>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="form-actions mt-4 mb-3">
      <button type="button" class="btn btn-primary btn-sm" @click="onSubmit">
        提交
      </button>
      <button type="button" class="btn btn-outline-secondary btn-sm ms-2" @click="$emit('goIndex')">
        取消
      </button>
    </div>
  </form>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import { getFieldPreset, fieldTypeGroups, fieldTypePresets } from '@/config/fieldTypePresets';
import { findTableList } from '@/api/table';
import { getTableFields, getTableOptions } from '@/api/table_field';
import { getEntityList } from '@/api/entity';
import { AppModal } from '@/components/Modal/modal';
import vSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';

const props = defineProps({
  dataInfo: {
    type: Object,
    default: () => ({}),
  },
});

const emit = defineEmits(['submitForm', 'goIndex']);

// 生成随机字段标识
const generateFieldCode = () => {
  const randomStr = Math.random().toString(36).substring(2, 8); // 生成6位随机字符
  return `f${randomStr}`;
};

// 创建默认表单数据工厂函数
const createDefaultFormData = (tableCode = '') => ({
  fieldType: '',
  code: generateFieldCode(), // 自动生成默认字段标识
  name: '',
  description: '',
  sort: 0,
  tableCode,
  required: 'No',
  isShow: 'Yes',
  tableCode,
  required: 'No',
  isShow: 'Yes',
  isFilter: 'No',
  groupName: '', // 分组名称
});

// 基本信息
const formData = ref(createDefaultFormData());

// 标志:是否正在从数据库加载数据(编辑模式)
// 用于防止 watch 覆盖从数据库加载的值
const isLoadingFromData = ref(false);

// 关联配置 - 默认关联到字典表
const relationConfig = ref({
  target: 'dict_item',  // 默认关联到字典表
  valueField: 'code',    // 存储字段 = 关联键
  labelField: 'name',
});

// 日期时间配置
const datetimeConfig = ref({
  showTime: false,
  defaultToCurrentTime: false,
});

// 自动编码配置
const autocodePatterns = ref([]);

// 附件配置
const attachmentConfig = ref({
  multiple: false,
  acceptInput: '',  // 用于输入的字符串
  maxSizeMB: 10,    // MB 单位
});

// 关联过滤条件 (单个对象)
const relationFilter = ref({
  field: '',
  value: '',
  options: []
});

// 目标表字段列表
const targetTableFields = ref([]);

// 当前表的业务字段列表(用于自动编码字段选择)
const currentTableFields = ref([]);

// 表列表
const tables = ref([]);

// 错误信息
const errors = ref({});

// 获取选中的字段类型预设
const selectedTypePreset = computed(() => {
  return formData.value.fieldType ? getFieldPreset(formData.value.fieldType) : null;
});

// 当前表的业务字段(排除系统字段)
const currentTableBusinessFields = computed(() => {
  return currentTableFields.value.filter(field => {
    // 确保字段有 code 和 name
    if (!field || !field.code || !field.name) {
      return false;
    }

    // 排除系统字段(使用 is_system 标记)
    if (field.is_system === true) {
      return false;
    }

    // 只保留字符串和整数类型字段
    // field_type: text, integer
    // 排除: textarea, number, autocode, date, datetime, select, relation 等
    const allowedFieldTypes = ['text', 'integer'];
    if (field.field_type && !allowedFieldTypes.includes(field.field_type)) {
      return false;
    }

    return true;
  });
});

// 是否需要关联配置
const needsRelationConfig = computed(() => {
  return ['belongsto', 'hasmany', 'manytomany', 'select', 'multiselect', 'radio', 'checkboxgroup'].includes(formData.value.fieldType);
});

// 是否需要日期时间配置
const needsDateTimeConfig = computed(() => {
  return ['date', 'datetime'].includes(formData.value.fieldType);
});

// 是否需要自动编码配置
const needsAutocodeConfig = computed(() => {
  return formData.value.fieldType === 'autocode';
});

// 是否需要附件配置
const needsAttachmentConfig = computed(() => {
  return formData.value.fieldType === 'attachment';
});

// 自定义字段过滤函数 - 支持同时搜索 name 和 code
// vue-select filter-by 函数签名: (option, label, search) => boolean
const filterFieldOptions = (option, label, search) => {
  if (!search) return true;
  const searchLower = search.toLowerCase();
  return (
    option.name?.toLowerCase().includes(searchLower) ||
    option.code?.toLowerCase().includes(searchLower)
  );
};

// 预览编码
const previewCode = computed(() => {
  if (!needsAutocodeConfig.value || autocodePatterns.value.length === 0) {
    return '';
  }

  let preview = '';
  const now = new Date();

  autocodePatterns.value.forEach(pattern => {
    if (pattern.type === 'string') {
      preview += pattern.options.value || '';
    } else if (pattern.type === 'date') {
      const format = pattern.options.format || 'YYYYMMDD';
      // 简单的日期格式化示例
      const formatted = format
        .replace('YYYY', now.getFullYear())
        .replace('YY', String(now.getFullYear()).slice(-2))
        .replace('MM', String(now.getMonth() + 1).padStart(2, '0'))
        .replace('DD', String(now.getDate()).padStart(2, '0'));
      preview += formatted;
    } else if (pattern.type === 'field') {
      // 显示字段代码的示例值
      const fieldCode = pattern.options.fieldCode;
      if (fieldCode) {
        // 查找字段名称
        const field = currentTableBusinessFields.value.find(f => f.code === fieldCode);
        preview += field ? `{${field.name}}` : `{${fieldCode}}`;
      } else {
        preview += '{字段}';
      }
    } else if (pattern.type === 'integer') {
      const digits = pattern.options.digits || 5;
      preview += '1'.padStart(digits, '0');
    }
  });

  return preview || '(未配置)';
});

// 添加模式
const addPattern = () => {
  autocodePatterns.value.push({
    type: 'string',
    options: {
      value: '',
      format: '',
      digits: null,  // 改为 null,要求用户必填
      start: null,   // 改为 null,要求用户必填
      cycle: 'none'
    }
  });
};

// 删除模式
const removePattern = (index) => {
  autocodePatterns.value.splice(index, 1);
};

// 加载目标表字段列表
const loadTargetTableFields = async (tableCode) => {
  if (!tableCode) {
    targetTableFields.value = [];
    return;
  }

  try {
    const response = await getTableFields({ table_code: tableCode });
    console.log('loadTargetTableFields response:', response);
    // axios 返回的是 response.status 和 response.data
    targetTableFields.value = response.data || [];
    console.log('targetTableFields loaded:', targetTableFields.value);
  } catch (error) {
    console.error('Failed to load target table fields:', error);
    targetTableFields.value = [];
  }
};

// 加载当前表的字段列表(用于自动编码字段选择)
const loadCurrentTableFields = async () => {
  if (!formData.value.tableCode) {
    currentTableFields.value = [];
    return;
  }

  try {
    const response = await getTableFields({ table_code: formData.value.tableCode });
    currentTableFields.value = response.data || [];
    console.log('currentTableFields loaded:', currentTableFields.value);
  } catch (error) {
    console.error('Failed to load current table fields:', error);
    currentTableFields.value = [];
  }
};

// 加载字段的可选值
// 1. 如果关联表是 dict_item 且过滤字段是 dict_code,直接从 t_dict 表加载
// 2. 否则如果字段有 relation 配置,从关联表加载
const loadFilterFieldOptions = async () => {
  if (!relationFilter.value.field) {
    relationFilter.value.options = [];
    return;
  }

  // 特殊处理:字典表的 dict 字段 - 加载全部数据用于前端过滤搜索
  if (relationConfig.value.target === 'dict_item' && relationFilter.value.field === 'dict') {
    try {
      // 从 t_dict 表加载全部选项
      const response = await getEntityList({ table_code: 't_dict', pageSize: 500 });
      // 将返回数据转换为 {code, name} 格式
      relationFilter.value.options = (response.data || []).map(item => ({
        code: item.code,
        name: item.name
      }));
    } catch (error) {
      console.error('Failed to load t_dict options:', error);
      relationFilter.value.options = [];
    }
    return;
  }

  // 查找字段配置
  const fieldConfig = targetTableFields.value.find(f => f.code === relationFilter.value.field);
  if (!fieldConfig || !fieldConfig.options) {
    relationFilter.value.options = [];
    return;
  }

  // 检查是否有 relation 配置
  const options = typeof fieldConfig.options === 'string'
    ? JSON.parse(fieldConfig.options)
    : fieldConfig.options;

  if (options.relation && options.relation.target) {
    try {
      const response = await getTableOptions(options.relation.target, options.relation.filter);
      relationFilter.value.options = response.data || [];
    } catch (error) {
      console.error('Failed to load field options:', error);
      relationFilter.value.options = [];
    }
  }
};

// 打开字典分类下拉时加载数据
const onOpenDictionaryClass = async () => {
  // 如果选项已经加载过,直接返回
  if (relationFilter.value.options.length > 0) {
    return;
  }

  try {
    const response = await getEntityList({
      table_code: 't_dict',
      pageSize: 500
    });
    // 将返回数据转换为 {code, name} 格式
    relationFilter.value.options = (response.data || []).map(item => ({
      code: item.code,
      name: item.name
    }));
  } catch (error) {
    console.error('Failed to load t_dict:', error);
    relationFilter.value.options = [];
  }
};

// 重置表单数据函数
const resetFormData = () => {
  const currentTableCode = formData.value.tableCode;
  const currentFieldType = formData.value.fieldType;

  // 重置表单数据,保留 tableCode 和 fieldType
  formData.value = createDefaultFormData(currentTableCode);
  formData.value.fieldType = currentFieldType;

  // 重置关联配置
  relationConfig.value = {
    target: '',
    valueField: 'code',
    labelField: 'name',
  };

  // 重置日期时间配置
  datetimeConfig.value = {
    showTime: false,
    defaultToCurrentTime: false,
  };

  // 重置自动编码配置
  autocodePatterns.value = [];
};

// 监听字段类型变化
watch(() => formData.value.fieldType, (newType, oldType) => {
  // 只在切换类型时重置(不包括初始加载和编辑模式)
  if (oldType && newType !== oldType && !props.dataInfo?.ID) {
    resetFormData();

    // 应用预设配置
    const preset = getFieldPreset(newType);
    if (preset) {
      console.log('Pending field type preset:', preset);
    }
  }
});

// 监听关联目标表变化,加载字段列表
watch(() => relationConfig.value.target, (newTarget) => {
  if (newTarget) {
    loadTargetTableFields(newTarget);

    // 如果正在从数据库加载数据,不要覆盖已加载的 valueField/labelField
    if (isLoadingFromData.value) {
      return;
    }

    // 如果新选择的是字典表,设置默认值
    if (newTarget === 'dict_item') {
      // 字典表:使用默认值 code/name 和默认过滤条件 dict
      relationConfig.value.valueField = 'code';
      relationConfig.value.labelField = 'name';
      relationFilter.value = {
        field: 'dict',
        value: '',
        options: []
      };
      // 加载 t_dict 的可选值
      loadTargetTableFields(newTarget).then(() => {
        loadFilterFieldOptions();
      });
    } else {
      // 非字典表:清空字段配置,让用户选择
      relationConfig.value.valueField = '';
      relationConfig.value.labelField = '';
      relationFilter.value = {
        field: '',
        value: '',
        options: []
      };
    }
  } else {
    targetTableFields.value = [];
    // 重置单个 filter 对象
    relationFilter.value = {
      field: '',
      value: '',
      options: []
    };
  }
});

// 加载表单数据
const loadFormData = () => {
  if (props.dataInfo && props.dataInfo.ID) {
    // 编辑模式
    formData.value = {
      fieldType: props.dataInfo.FieldType || '',
      code: props.dataInfo.Code || '',
      name: props.dataInfo.Name || '',
      description: props.dataInfo.Description || '',
      sort: props.dataInfo.Sort || 0,
      tableCode: props.dataInfo.TableCode || '',
      required: props.dataInfo.Required || 'No',
      isShow: props.dataInfo.IsShow || 'Yes',
      required: props.dataInfo.Required || 'No',
      isShow: props.dataInfo.IsShow || 'Yes',
      isFilter: props.dataInfo.IsFilter || 'No',
      groupName: props.dataInfo.GroupName || '',
    };

    // 解析 Options
    if (props.dataInfo.Options) {
      const options = typeof props.dataInfo.Options === 'string'
        ? JSON.parse(props.dataInfo.Options)
        : props.dataInfo.Options;

      console.log('Parsed options:', options);

      // 关联配置
      if (options.relation) {
        console.log('Loading relation config:', options.relation);

        // 设置标志,防止 watch 覆盖从数据库加载的值
        isLoadingFromData.value = true;

        relationConfig.value = {
          target: options.relation.target || '',
          valueField: options.relation.valueField || 'code',
          labelField: options.relation.labelField || 'name',
        };

        // 解析 filter 为单个对象
        if (options.relation.filter) {
          console.log('Found filter:', options.relation.filter);
          const filterEntries = Object.entries(options.relation.filter);
          console.log('Filter entries:', filterEntries);
          if (filterEntries.length > 0) {
            const [field, value] = filterEntries[0]; // 只取第一个
            relationFilter.value = {
              field,
              value,
              options: []
            };
            console.log('Set relationFilter:', relationFilter.value);
          }
        }

        // 加载目标表字段后,再加载 filter 的选项,最后重置标志
        loadTargetTableFields(options.relation.target).then(() => {
          if (relationFilter.value.field) {
            loadFilterFieldOptions();
          }
          // 数据加载完成,重置标志
          isLoadingFromData.value = false;
        });
      }

      // 日期时间配置
      if (options.datetime) {
        datetimeConfig.value = {
          showTime: options.datetime.showTime || false,
          defaultToCurrentTime: options.datetime.defaultToCurrentTime || false,
        };
      }

      // 附件配置
      if (options.attachment) {
        console.log('Loading attachment config:', options.attachment);
        attachmentConfig.value = {
          multiple: options.attachment.multiple || false,
          acceptInput: options.attachment.accept ? options.attachment.accept.join(', ') : '',
          maxSizeMB: options.attachment.maxSize ? Math.round(options.attachment.maxSize / (1024 * 1024)) : 10,
        };
      }

      // 自动编码配置
      if (options.patterns && Array.isArray(options.patterns)) {
        autocodePatterns.value = options.patterns.map(p => ({
          type: p.type || 'string',
          options: { ...p.options }
        }));
      }
    }
  } else if (props.dataInfo) {
    // 创建模式:设置 tableCode 和默认排序值
    if (props.dataInfo.TableCode) {
      formData.value.tableCode = props.dataInfo.TableCode;
    }
    if (props.dataInfo.Sort !== undefined) {
      formData.value.sort = props.dataInfo.Sort;
    }
  }
};

// 加载表列表
const loadTables = async () => {
  try {
    const res = await findTableList({ page: 1, pageSize: -1 });
    console.log(res);
    tables.value = res.data || [];
  } catch (error) {
    console.error('Failed to load tables:', error);
  }
};

// 初始化
onMounted(async () => {
  loadFormData();
  await loadTables();

  // 如果是创建模式且 dataInfo 有 Sort 值,使用它作为默认排序
  if (!props.dataInfo?.ID && props.dataInfo?.Sort !== undefined) {
    formData.value.sort = props.dataInfo.Sort;
  }

  // 如果是创建模式,初始化默认关联配置
  if (!props.dataInfo?.ID && needsRelationConfig.value) {
    // 加载字典表的字段列表
    await loadTargetTableFields('dict_item');
    // 设置默认过滤条件为 dict
    relationFilter.value = {
      field: 'dict',
      value: '',
      options: []
    };
    // 加载 t_dict 的可选值
    loadFilterFieldOptions();
  }

  // 加载当前表的字段列表(用于自动编码字段选择)
  loadCurrentTableFields();
});

// 监听 dataInfo 变化
watch(() => props.dataInfo, (newDataInfo) => {
  if (newDataInfo && Object.keys(newDataInfo).length > 0) {
    loadFormData();

    // 如果是创建模式且有 Sort 值,设置默认排序
    if (!newDataInfo.ID && newDataInfo.Sort !== undefined) {
      formData.value.sort = newDataInfo.Sort;
    }
  }
}, { deep: true, immediate: true });



// 提交表单
const onSubmit = () => {
  errors.value = {};

  // 基本验证
  if (!formData.value.fieldType) {
    AppModal.alert({ title: '提示', content: '请选择字段类型' });
    return;
  }
  if (!formData.value.code) {
    errors.value.code = '字段编码必填';
    return;
  }
  if (!formData.value.name) {
    errors.value.name = '字段名称必填';
    return;
  }

  // 关联配置验证
  if (needsRelationConfig.value && !relationConfig.value.target) {
    AppModal.alert({ title: '提示', content: '请选择关联表' });
    return;
  }

  // 自动编码配置验证
  if (needsAutocodeConfig.value) {
    if (autocodePatterns.value.length === 0) {
      AppModal.alert({ title: '提示', content: '请至少添加一个编码模式' });
      return;
    }

    // 验证每个模式的必填项
    for (let i = 0; i < autocodePatterns.value.length; i++) {
      const pattern = autocodePatterns.value[i];

      if (pattern.type === 'string' && !pattern.options.value) {
        AppModal.alert({ title: '提示', content: `模式 ${i + 1}: 请填写固定字符串` });
        return;
      }

      if (pattern.type === 'date' && !pattern.options.format) {
        AppModal.alert({ title: '提示', content: `模式 ${i + 1}: 请填写日期格式` });
        return;
      }

      if (pattern.type === 'field' && !pattern.options.fieldCode) {
        AppModal.alert({ title: '提示', content: `模式 ${i + 1}: 请填写字段代码` });
        return;
      }

      if (pattern.type === 'integer') {
        if (!pattern.options.digits || pattern.options.digits < 1) {
          AppModal.alert({ title: '提示', content: `模式 ${i + 1}: 请填写序列号位数(1-10)` });
          return;
        }
        if (pattern.options.start === undefined || pattern.options.start === null || pattern.options.start === '') {
          AppModal.alert({ title: '提示', content: `模式 ${i + 1}: 请填写起始值` });
          return;
        }
      }
    }
  }

  // 组装 Options - 从原有数据中保留其他配置(如 ui)
  let options = {};
  if (props.dataInfo && props.dataInfo.Options) {
    const existingOptions = typeof props.dataInfo.Options === 'string'
      ? JSON.parse(props.dataInfo.Options)
      : props.dataInfo.Options;
    // 保留 ui 等其他配置
    if (existingOptions.ui) {
      options.ui = existingOptions.ui;
    }
  }

  // 关联配置
  if (needsRelationConfig.value && relationConfig.value.target) {
    options.relation = {
      target: relationConfig.value.target,
      valueField: relationConfig.value.valueField || 'code',
      labelField: relationConfig.value.labelField || 'name',
    };

    // 添加 filter (从单个对象转换)
    if (relationFilter.value.field && relationFilter.value.value) {
      options.relation.filter = {
        [relationFilter.value.field]: relationFilter.value.value
      };
    }
  }

  // 日期时间配置
  if (needsDateTimeConfig.value) {
    options.datetime = {
      showTime: datetimeConfig.value.showTime,
      defaultToCurrentTime: datetimeConfig.value.defaultToCurrentTime,
    };
  }

  // 附件配置
  if (needsAttachmentConfig.value) {
    options.attachment = {
      multiple: attachmentConfig.value.multiple,
      accept: attachmentConfig.value.acceptInput
        ? attachmentConfig.value.acceptInput.split(',').map(s => s.trim()).filter(s => s)
        : [],
      maxSize: attachmentConfig.value.maxSizeMB * 1024 * 1024, // 转换为字节
    };
  }

  // 自动编码配置
  if (needsAutocodeConfig.value && autocodePatterns.value.length > 0) {
    options.patterns = autocodePatterns.value.map(p => ({
      type: p.type,
      options: { ...p.options }
    }));
  }

  // 提交数据
  const submitData = {
    code: formData.value.code,
    table_code: formData.value.tableCode, // 后端需要 table_code
    name: formData.value.name,
    field_type: formData.value.fieldType, // 后端会根据此字段推断 Type 和 Length
    description: formData.value.description,
    sort: formData.value.sort,
    required: formData.value.required,
    is_show: formData.value.isShow,
    is_show: formData.value.isShow,
    is_filter: formData.value.isFilter,
    group_name: formData.value.groupName,
    options: Object.keys(options).length > 0 ? JSON.stringify(options) : null,
  };

  if (props.dataInfo && props.dataInfo.ID) {
    submitData.id = props.dataInfo.ID;
  }

  emit('submitForm', submitData);
};
</script>

<style scoped>
/* 表单样式 */
.required::after {
  content: ' *';
  color: #dc3545;
  font-weight: bold;
}

.form-actions {
  padding-top: 16px;
  border-top: 1px solid #e9ecef;
}

/* 配置区域标题 */
h6.text-muted {
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e9ecef;
}

/* v-select 样式调整 - 对齐修复 */
.v-select {
  min-height: calc(1.5em + 0.5rem + 2px);
  padding: 0;
}

.v-select .vs__dropdown-toggle {
  padding: 0 0 0 0.5rem;
  border: 1px solid #ced4da;
  border-radius: 0.25rem;
  min-height: calc(1.5em + 0.5rem + 2px);
}

.v-select .vs__selected-options {
  padding: 0;
}

.v-select .vs__search {
  margin: 0;
  padding: 0.25rem 0;
}

.v-select .vs__search::placeholder {
  color: #6c757d;
}

.v-select.vs--open .vs__dropdown-toggle {
  border-color: #86b7fe;
  box-shadow: 0 0 0 0.25rem rgba(13, 110, 253, 0.25);
}

.v-select .vs__dropdown-menu {
  min-width: 100%;
}

.v-select .vs__actions {
  padding: 0 0.5rem 0 0;
}
</style>

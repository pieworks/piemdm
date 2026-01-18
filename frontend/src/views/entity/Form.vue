<template>
  <form name="entityForm" :id="'entityForm-' + extTable" method="post" enctype="multipart/form-data">
    <div class="col-12 col-sm-12">
      <div class="row">
        <div class="row g-3 mb-3">
          <div class="form-group row">
            <legend class="col-form-label col-sm-2 required">{{ $t('Reason') }}:</legend>
            <div class="col-sm-10">
              <textarea class="form-control form-control-sm" :class="formErrors['reason'] ? 'is-invalid' : ''"
                v-model="dataInfo['reason']" id="Reason" :placeholder="$t('Reason')" maxlength="255"
                rows="3">{{ dataInfo.reason }}</textarea>
            </div>
          </div>
        </div>
        <template v-for="field in flatFields" :key="field.Code || field.name">
          <div class="col-12 mt-4" v-if="field.isHeader">
            <h6 class="text-secondary border-bottom pb-2">
              <i class="bi bi-bookmark me-2"></i>{{ field.name }}
            </h6>
          </div>
          <div class="col-sm-6" v-else>
            <div class="form-group row">
              <legend :for="field.Code" class="col-form-label col-sm-4"
                :class="field.Required == 'Yes' ? 'required' : ''">
                {{ field.Name }}
              </legend>

              <!-- 自动编码字段特殊处理 -->
              <div class="col-sm-8" v-if="isAutocode(field)">
                <!-- 创建模式:显示"自动生成"提示 -->
                <input v-if="!dataInfo.id" type="text" class="form-control form-control-sm" :value="'自动生成'" disabled />
                <!-- 更新模式:显示纯文本(只读) -->
                <input v-else type="text" class="form-control form-control-sm" :value="dataInfo[field.Code]" disabled />
              </div>


              <!-- 根据 field_type 或 widget 渲染不同的输入控件 -->

              <!-- Date 日期 -->
              <div class="col-sm-8" v-else-if="isDate(field)">
                <VueDatePicker v-model="dataInfo[field.Code]" model-type="yyyy-MM-dd" :format="'yyyy-MM-dd'"
                  :locale="currentLocale" :enable-time-picker="false" auto-apply :clearable="field.Required !== 'Yes'"
                  :class="formErrors[field.Code] ? 'is-invalid' : ''" :id="field.Code" />
              </div>

              <!-- Time 时间 -->
              <div class="col-sm-8" v-else-if="isTime(field)">
                <VueDatePicker v-model="dataInfo[field.Code]" time-picker model-type="HH:mm:ss" :format="'HH:mm:ss'"
                  :locale="currentLocale" auto-apply :clearable="field.Required !== 'Yes'"
                  :class="formErrors[field.Code] ? 'is-invalid' : ''" :id="field.Code" />
              </div>

              <!-- DateTime 日期时间 -->
              <div class="col-sm-8" v-else-if="isDateTime(field)">
                <VueDatePicker v-model="dataInfo[field.Code]" model-type="yyyy-MM-dd HH:mm:ss"
                  :format="'yyyy-MM-dd HH:mm:ss'" :locale="currentLocale" auto-apply
                  :clearable="field.Required !== 'Yes'" :class="formErrors[field.Code] ? 'is-invalid' : ''"
                  :id="field.Code" />
              </div>

              <!-- Upload 附件 -->
              <div class="col-sm-8" v-else-if="isUpload(field)">
                <Upload v-model="dataInfo[field.Code]" :multiple="field.Options?.attachment?.multiple"
                  :accept="field.Options?.attachment?.accept?.join(',')"
                  :max-size="field.Options?.attachment?.maxSize" />
              </div>

              <!-- Number / Integer -->
              <div class="col-sm-8" v-else-if="isNumber(field)">
                <input type="number" class="form-control form-control-sm"
                  :class="formErrors[field.Code] ? 'is-invalid' : ''" v-model.number="dataInfo[field.Code]"
                  :id="field.Code" :placeholder="field.Name" :step="getNumberStep(field)" />
              </div>

              <!-- Select 下拉框 -->
              <div class="col-sm-8 has-validation" v-else-if="isSelect(field)">
                <v-select :key="`select-${field.Code}`" v-model="dataInfo[field.Code]" :id="field.Code"
                  :multiple="field.FieldType === 'multiselect' || getWidget(field) === 'MultiSelect'"
                  :reduce="option => option?.[field.Options?.relation?.valueField || field.RelationCode || 'code']"
                  :options="dictionarys['dict-' + field.Code] || []"
                  :get-option-label="option => formatOptionLabel(option, field)" :filterable="true"
                  :class="formErrors[field.Code] ? 'is-invalid' : ''" :placeholder="$t('Please select')"
                  @open="() => { loadOptionsOnOpen(field); }" @search="
                    (search, loading) => {
                      fetchOptions(search, loading, field);
                    }
                  ">
                  <template v-slot:option="option">
                    {{ formatOptionLabel(option, field) }}
                  </template>
                </v-select>
              </div>

              <!-- Radio 单选框组 -->
              <div class="col-sm-8" v-else-if="isRadio(field)">
                <div class="form-check form-check-inline" v-for="option in dictionarys['dict-' + field.Code] || []"
                  :key="option[field.Options?.relation?.valueField || 'code']">
                  <input type="radio" class="form-check-input" :class="formErrors[field.Code] ? 'is-invalid' : ''"
                    v-model="dataInfo[field.Code]"
                    :id="field.Code + '-' + option[field.Options?.relation?.valueField || 'code']"
                    :value="option[field.Options?.relation?.valueField || 'code']" />
                  <label class="form-check-label"
                    :for="field.Code + '-' + option[field.Options?.relation?.valueField || 'code']">
                    {{ option[field.Options?.relation?.labelField || 'name'] }}
                  </label>
                </div>
              </div>

              <!-- CheckboxGroup 复选框组 -->
              <div class="col-sm-8" v-else-if="isCheckboxGroup(field)">
                <div class="form-check form-check-inline" v-for="option in dictionarys['dict-' + field.Code] || []"
                  :key="option[field.Options?.relation?.valueField || 'code']">
                  <input type="checkbox" class="form-check-input" :class="formErrors[field.Code] ? 'is-invalid' : ''"
                    v-model="dataInfo[field.Code]"
                    :id="field.Code + '-' + option[field.Options?.relation?.valueField || 'code']"
                    :value="option[field.Options?.relation?.valueField || 'code']" />
                  <label class="form-check-label"
                    :for="field.Code + '-' + option[field.Options?.relation?.valueField || 'code']">
                    {{ option[field.Options?.relation?.labelField || 'name'] }}
                  </label>
                </div>
              </div>

              <!-- Checkbox 勾选 -->
              <div class="col-sm-8" v-else-if="isCheckbox(field)">
                <div class="form-check">
                  <input type="checkbox" class="form-check-input" v-model="dataInfo[field.Code]" :id="field.Code" />
                </div>
              </div>

              <!-- TextArea 文本域 -->
              <div class="col-sm-8" v-else-if="isTextarea(field)">
                <textarea class="form-control form-control-sm" :class="formErrors[field.Code] ? 'is-invalid' : ''"
                  v-model="dataInfo[field.Code]" :id="field.Code" :placeholder="field.Name"
                  :maxlength="field.Length || 2000" rows="3"></textarea>
              </div>

              <!-- Text 输入框 -->
              <div class="col-sm-8" v-else-if="isTextInput(field)">
                <!-- 只读字段：使用 disabled input -->
                <input v-if="readonlyFields.includes(field.Code)" type="text" class="form-control form-control-sm"
                  :value="dataInfo[field.Code]" disabled />
                <!-- 可编辑字段 -->
                <input v-else type="text" class="form-control form-control-sm" v-model="dataInfo[field.Code]"
                  :id="field.Code" :placeholder="field.Name" :class="formErrors[field.Code] ? 'is-invalid' : ''"
                  :maxlength="field.Length || 255" />
              </div>

              <!-- 默认:文本输入 -->
              <div class="col-sm-8" v-else>
                <!-- 只读字段：使用 disabled input -->
                <input v-if="readonlyFields.includes(field.Code)" type="text" class="form-control form-control-sm"
                  :value="dataInfo[field.Code]" disabled />
                <!-- 可编辑字段 -->
                <input v-else type="text" class="form-control form-control-sm" v-model="dataInfo[field.Code]"
                  :id="field.Code" :placeholder="field.Name" :class="formErrors[field.Code] ? 'is-invalid' : ''" />
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </form>
</template>

<script setup>
import { watchEffect, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import vSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';
import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';
import Upload from '@/components/Upload.vue';

// i18n locale
const { locale } = useI18n();
const currentLocale = computed(() => {
  return locale.value === 'zh-CN' ? 'zh' : 'en';
});

const props = defineProps({
  dataInfo: {
    type: Object,
    default: {},
  },
  fields: {
    type: Object,
    default: {},
  },
  dictionarys: {
    type: Object,
    default: {},
  },
  baseInfo: {
    type: Object,
    default: {},
  },
  extTable: String,
  formErrors: {
    type: Object,
    default: {},
  },
  readonlyFields: {
    type: Array,
    default: () => [],
  },
});

// 按 Sort 字段升序排序
const sortedFields = computed(() => {
  if (!props.fields || typeof props.fields !== 'object') {
    return [];
  }

  // 将 fields 对象转换为数组
  const fieldsArray = Array.isArray(props.fields)
    ? props.fields
    : Object.values(props.fields);

  // 按 Sort/sort 字段升序排序（同时支持大写和小写）
  return [...fieldsArray].sort((a, b) => {
    const sortA = a.Sort ?? a.sort ?? 0;
    const sortB = b.Sort ?? b.sort ?? 0;
    return sortA - sortB;
  });
});

// 分组并扁平化字段列表
const flatFields = computed(() => {
  if (!sortedFields.value || sortedFields.value.length === 0) return [];

  const groups = {};
  const groupOrder = [];

  sortedFields.value.forEach(field => {
    const groupName = field.GroupName || '基本信息';
    if (!groups[groupName]) {
      groups[groupName] = [];
      groupOrder.push(groupName);
    }
    groups[groupName].push(field);
  });

  const result = [];
  groupOrder.forEach(name => {
    result.push({ isHeader: true, name });
    groups[name].forEach(f => result.push(f));
  });

  return result;
});

const emit = defineEmits(['update:dataInfo', 'fetchData']);

watchEffect(() => {
  emit('update:dataInfo', props.dataInfo);
});

function fetchOptions(search, loading, option) {
  if (search.length > 1) {
    emit('fetchData', search, loading, option);
  }
}

// 下拉框打开时加载选项
function loadOptionsOnOpen(field) {
  const options = props.dictionarys['dict-' + field.Code] || [];
  // 如果选项为空,触发加载
  if (options.length === 0) {
    emit('fetchData', '', null, field);
  }
}

// 格式化下拉选项显示标签: "code name" 格式
function formatOptionLabel(option, field) {
  if (!option) return '';

  const valueField = field.Options?.relation?.valueField || field.RelationCode || 'code';
  const labelField = field.Options?.relation?.labelField || field.RelationName || 'name';

  // 如果 option 是完整对象,直接格式化
  if (typeof option === 'object') {
    const code = option[valueField] || '';
    const name = option[labelField] || '';
    return `${code} ${name}`.trim();
  }

  // 如果 option 是 primitive 值(reduce 后的值),尝试从选项列表中查找匹配的对象
  const options = props.dictionarys['dict-' + field.Code] || [];
  const matchedOption = options.find(opt => opt[valueField] == option); // 使用 == 进行宽松比较

  if (matchedOption) {
    const code = matchedOption[valueField] || '';
    const name = matchedOption[labelField] || '';
    return `${code} ${name}`.trim();
  }

  // 兜底: 直接返回 primitive 值
  return String(option);
}

// 辅助方法:判断字段类型
function getWidget(field) {
  // 优先使用 Options.UI.Widget
  if (field.Options && field.Options.ui && field.Options.ui.widget) {
    return field.Options.ui.widget;
  }
  // 兼容旧的 Style 字段
  if (field.Style) {
    return field.Style;
  }
  // 根据 FieldType 推断
  const typeMap = {
    'text': 'Input',
    'textarea': 'Textarea',
    'phone': 'Input',
    'email': 'Input',
    'url': 'Input',
    'integer': 'InputNumber',
    'decimal': 'InputNumber',
    'percent': 'InputNumber',
    'date': 'DatePicker',
    'time': 'TimePicker',
    'datetime': 'DateTimePicker',
    'checkbox': 'Checkbox',
    'select': 'Select',
    'multiselect': 'MultiSelect',
    'radio': 'RadioGroup',
    'checkboxgroup': 'CheckboxGroup',
  };
  return typeMap[field.FieldType] || 'Input';
}

// 判断是否为自动编码字段
function isAutocode(field) {
  // 方法1: 检查 FieldType (推荐)
  if (field.FieldType === 'autocode' || field.field_type === 'autocode') {
    return true;
  }

  // 方法2: 检查 Options.patterns 是否存在(向后兼容)
  if (field.Options && field.Options.patterns && Array.isArray(field.Options.patterns) && field.Options.patterns.length > 0) {
    return true;
  }

  return false;
}

function isTextInput(field) {
  if (isAutocode(field)) {
    return false; // 排除自动编码字段
  }
  // 排除 Select 类型字段
  if (isSelect(field)) {
    return false;
  }
  // 排除 Radio 和 CheckboxGroup 字段
  if (isRadio(field) || isCheckboxGroup(field)) {
    return false;
  }
  // 排除日期、时间、日期时间字段
  if (isDate(field) || isTime(field) || isDateTime(field)) {
    return false;
  }
  // 排除附件字段
  if (isUpload(field)) {
    return false;
  }
  const widget = getWidget(field);
  return widget === 'Input' || field.Type === 'Text';
}

function isTextarea(field) {
  const widget = getWidget(field);
  return widget === 'Textarea';
}

function isDate(field) {
  // 先检查 FieldType
  if (field.FieldType === 'date') {
    return true;
  }
  // 再检查 Type
  if (field.Type === 'Date') {
    return true;
  }
  // 最后检查 widget
  const widget = getWidget(field);
  return widget === 'DatePicker';
}

function isTime(field) {
  // 先检查 FieldType
  if (field.FieldType === 'time') {
    return true;
  }
  // 检查 widget
  const widget = getWidget(field);
  return widget === 'TimePicker';
}

function isDateTime(field) {
  // 先检查 FieldType
  if (field.FieldType === 'datetime') {
    return true;
  }
  // 再检查 Type
  if (field.Type === 'DateTime') {
    return true;
  }
  // 最后检查 widget
  const widget = getWidget(field);
  return widget === 'DateTimePicker';
}

function isUpload(field) {
  // 先检查 FieldType
  if (field.FieldType === 'attachment') {
    return true;
  }
  // 检查 widget
  const widget = getWidget(field);
  return widget === 'Upload';
}

function isNumber(field) {
  const widget = getWidget(field);
  return widget === 'InputNumber' || field.Type === 'Number';
}

function isSelect(field) {
  // 优先检查 FieldType (排除 radio 和 checkboxgroup)
  if (field.FieldType === 'select' || field.FieldType === 'multiselect') {
    return true;
  }

  // 检查 widget
  const widget = getWidget(field);
  if (widget === 'Select' || widget === 'MultiSelect') {
    return true;
  }

  // 检查是否有关联配置(排除空字符串),但排除 radio 和 checkboxgroup
  if (field.FieldType === 'radio' || field.FieldType === 'checkboxgroup') {
    return false;
  }
  const hasRelation = field.Options?.relation?.target && field.Options.relation.target.trim() !== '';
  return hasRelation;
}

function isRadio(field) {
  // 检查 FieldType
  if (field.FieldType === 'radio') {
    return true;
  }
  // 检查 widget
  const widget = getWidget(field);
  return widget === 'RadioGroup';
}

function isCheckboxGroup(field) {
  // 检查 FieldType
  if (field.FieldType === 'checkboxgroup') {
    return true;
  }
  // 检查 widget
  const widget = getWidget(field);
  return widget === 'CheckboxGroup';
}

function isCheckbox(field) {
  const widget = getWidget(field);
  return widget === 'Checkbox';
}

function getNumberStep(field) {
  // 如果是 decimal 类型,返回小数步长
  if (field.FieldType === 'decimal' || field.FieldType === 'percent') {
    return 0.01;
  }
  return 1;
}
</script>

<style scoped>
legend.required:before {
  color: red;
  content: '*';
  vertical-align: middle;
  margin-right: 2px;
}

.v-select.is-invalid .vs__dropdown-toggle {
  border: 1px solid #dd0000;
}

/* VueDatePicker 样式调整 */
:deep(.dp__input) {
  padding: 0.15rem 0.25rem 0.15rem 2rem;
  font-size: 0.875rem;
  height: calc(1.5em + 0.5rem + 2px);
  /* 与 form-control-sm 一致 */
}

:deep(.dp__input_icon) {
  left: 0rem;
}

:deep(.dp__clear_icon) {
  right: 0.5rem;
}
</style>

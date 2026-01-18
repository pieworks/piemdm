<template>
  <div class="dynamic-field">
    <component
      :is="widgetComponent"
      v-model="fieldValue"
      :widgetProps="field.options?.ui?.widgetProps || {}"
      :options="fieldOptions"
      @update:modelValue="handleUpdate"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import { getWidget } from '@/config/widgetMap';
import { getTableOptions } from '@/api/table_field';

const props = defineProps({
  field: {
    type: Object,
    required: true
  },
  modelValue: {
    type: [String, Number, Boolean, Array, Object],
    default: null
  }
});

const emit = defineEmits(['update:modelValue']);

// 字段值
const fieldValue = ref(props.modelValue);

// 字段选项（用于下拉、单选、复选等）
const fieldOptions = ref([]);

// 选项缓存（避免重复请求）
const optionsCache = new Map();

// 获取 Widget 组件
const widgetComponent = computed(() => {
  const widgetName = props.field.options?.ui?.widget || 'Input';
  return getWidget(widgetName);
});

// 是否需要加载选项
const needsOptions = computed(() => {
  const widgetName = props.field.options?.ui?.widget || '';
  return ['Select', 'MultiSelect', 'RadioGroup', 'CheckboxGroup'].includes(widgetName);
});

// 加载字段选项
const loadOptions = async () => {
  console.log('loadOptions called', {
    needsOptions: needsOptions.value,
    field: props.field.code,
    widget: props.field.options?.ui?.widget
  });

  if (!needsOptions.value) {
    console.log('needsOptions is false, skipping');
    return;
  }

  // 检查是否有关联配置
  const relation = props.field.options?.relation;
  console.log('relation config:', relation);

  if (relation && relation.target) {
    // 构建缓存键（包含 filter）
    const filterKey = relation.filter ? JSON.stringify(relation.filter) : '';
    const cacheKey = `${relation.target}:${filterKey}`;

    if (optionsCache.has(cacheKey)) {
      fieldOptions.value = optionsCache.get(cacheKey);
      console.log('Using cached options:', fieldOptions.value);
      return;
    }

    try {
      // 调用 API 获取关联表选项,传递 filter 参数
      console.log('Calling getTableOptions:', { target: relation.target, filter: relation.filter });
      const response = await getTableOptions(relation.target, relation.filter);
      console.log('getTableOptions response:', response);

      // axios 直接返回 response.data
      const options = response.data || [];
      fieldOptions.value = options;
      console.log('Set fieldOptions:', options);

      // 缓存选项
      optionsCache.set(cacheKey, options);
    } catch (error) {
      console.error('Failed to load field options:', error);
      fieldOptions.value = [];
    }
  } else if (props.field.options?.datasource?.options) {
    // 使用静态选项
    console.log('Using static options:', props.field.options.datasource.options);
    fieldOptions.value = props.field.options.datasource.options;
  } else {
    console.log('No relation or datasource options found');
  }
};

// 监听字段值变化
watch(() => props.modelValue, (newValue) => {
  fieldValue.value = newValue;
});

// 处理值更新
const handleUpdate = (value) => {
  fieldValue.value = value;
  emit('update:modelValue', value);
};

// 初始化
onMounted(() => {
  loadOptions();
});
</script>

<style scoped>
.dynamic-field {
  width: 100%;
}
</style>

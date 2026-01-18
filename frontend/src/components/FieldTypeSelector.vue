<template>
  <div class="field-type-selector">
    <!-- 搜索框 -->
    <div class="search-box mb-3">
      <input
        type="text"
        class="form-control form-control-sm"
        v-model="searchQuery"
        placeholder="搜索字段类型..."
      />
    </div>

    <!-- 字段类型分组 -->
    <div v-for="group in filteredGroups" :key="group.name" class="type-group mb-4">
      <h6 class="group-title text-muted mb-2">
        <i class="bi bi-folder me-1"></i>
        {{ group.label }}
      </h6>
      <div class="type-grid">
        <div
          v-for="type in group.types"
          :key="type"
          class="type-card"
          :class="{ active: modelValue === type }"
          @click="selectType(type)"
        >
          <div class="type-icon">
            <i :class="fieldTypePresets[type].icon"></i>
          </div>
          <div class="type-label">{{ fieldTypePresets[type].label }}</div>
        </div>
      </div>
    </div>

    <!-- 无搜索结果提示 -->
    <div v-if="filteredGroups.length === 0" class="no-results text-center text-muted py-4">
      <i class="bi bi-search fs-1"></i>
      <p class="mt-2">未找到匹配的字段类型</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { fieldTypePresets, fieldTypeGroups } from '@/config/fieldTypePresets';

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['update:modelValue']);

const searchQuery = ref('');

// 过滤后的分组
const filteredGroups = computed(() => {
  if (!searchQuery.value) {
    return fieldTypeGroups;
  }

  const query = searchQuery.value.toLowerCase();
  return fieldTypeGroups.map(group => ({
    ...group,
    types: group.types.filter(type =>
      fieldTypePresets[type].label.toLowerCase().includes(query)
    )
  })).filter(group => group.types.length > 0);
});

// 选择字段类型
const selectType = (type) => {
  emit('update:modelValue', type);
};
</script>

<style scoped>
.field-type-selector {
  max-height: 500px;
  overflow-y: auto;
}

.search-box input {
  border-radius: 4px;
}

.group-title {
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.type-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
}

.type-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 16px 12px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  background-color: #fff;
}

.type-card:hover {
  border-color: #0d6efd;
  background-color: #f8f9fa;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.type-card.active {
  border-color: #0d6efd;
  background-color: #e7f1ff;
  box-shadow: 0 0 0 3px rgba(13, 110, 253, 0.1);
}

.type-icon {
  font-size: 1.5rem;
  color: #6c757d;
  margin-bottom: 8px;
}

.type-card.active .type-icon {
  color: #0d6efd;
}

.type-label {
  font-size: 0.875rem;
  color: #495057;
  text-align: center;
  font-weight: 500;
}

.type-card.active .type-label {
  color: #0d6efd;
  font-weight: 600;
}

.no-results {
  padding: 40px 20px;
}

.no-results i {
  opacity: 0.3;
}

/* 滚动条样式 */
.field-type-selector::-webkit-scrollbar {
  width: 6px;
}

.field-type-selector::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.field-type-selector::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 3px;
}

.field-type-selector::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>

<template>
  <ul class="list-group list-group-flush">
    <li v-for="node in data" :key="node.id" class="list-group-item border-0 p-1">
      <div class="form-check">
        <input class="form-check-input" type="checkbox" :value="node.id" :checked="modelValue.includes(node.id)"
          @change="toggle(node.id, $event.target.checked)" :id="'perm-' + node.id">
        <label class="form-check-label" :for="'perm-' + node.id">
          {{ node.name }} <span class="text-muted small">({{ node.code }})</span>
        </label>
      </div>
      <div v-if="node.children && node.children.length" class="ms-2 border-start ps-2">
        <PermissionTree :data="node.children" :modelValue="modelValue" @update:modelValue="updateModel" />
      </div>
    </li>
  </ul>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
  data: {
    type: Array,
    default: () => []
  },
  modelValue: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['update:modelValue']);

const toggle = (id, checked) => {
  let newValue = [...props.modelValue];
  if (checked) {
    if (!newValue.includes(id)) newValue.push(id);
  } else {
    newValue = newValue.filter(item => item !== id);
  }
  emit('update:modelValue', newValue);
};

const updateModel = (val) => {
  emit('update:modelValue', val);
};
</script>

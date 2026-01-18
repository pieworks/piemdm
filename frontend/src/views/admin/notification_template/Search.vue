<template>
  <div class="px-1">
    <form id="search">
      <div class="form-group row">
        <div class="col-auto">
          <label class="col-form-label-sm">{{ $t('Template Code') }}:</label>
        </div>
        <div class="col-auto">
          <input type="text" class="form-control form-control-sm" :placeholder="$t('Template Code')"
            v-model.lazy="formData.templateCode" maxlength="100" size="20" />
        </div>
        <div class="col-auto">
          <label class="col-form-label-sm">{{ $t('Template Name') }}:</label>
        </div>
        <div class="col-auto">
          <input type="text" class="form-control form-control-sm" :placeholder="$t('Template Name')"
            v-model.lazy="formData.templateName" maxlength="100" size="20" />
        </div>
        <div class="col-auto">
          <button type="button" class="btn btn-outline-primary btn-sm" @click="$emit('search', formData)">
            {{ $t('Go') }}
          </button>
          <button type="button" class="btn btn-outline-secondary btn-sm ms-1" @click="onReset">
            {{ $t('Reset') }}
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const emit = defineEmits(['search']);
const formData = ref({});

const onReset = () => {
  formData.value = {};
  emit('search', formData.value);
};
</script>

<style scoped>
/* 保持一致的样式 */
</style>

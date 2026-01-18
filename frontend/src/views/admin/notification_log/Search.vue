<template>
  <div class="px-1">
    <form id="search">
      <div class="form-group row">
        <!-- Recipient ID -->
        <div class="col-auto">
          <label class="col-form-label-sm">{{ $t('Recipient') }}:</label>
        </div>
        <div class="col-auto">
          <input type="text" class="form-control form-control-sm" :placeholder="$t('Recipient ID')"
            v-model.lazy="formData.recipientId" size="15" />
        </div>

        <!-- Time Range -->
        <div class="col-auto">
          <label class="col-form-label-sm">{{ $t('From') }}:</label>
        </div>
        <div class="col-auto">
          <VueDatePicker v-model="formData.startTime" placeholder="Start Time" :format="'yyyy-MM-dd HH:mm:ss'"
            locale="zh" auto-apply></VueDatePicker>
        </div>
        <div class="col-auto">
          <label class="col-form-label-sm">-</label>
        </div>
        <div class="col-auto">
          <VueDatePicker v-model="formData.endTime" placeholder="End Time" :format="'yyyy-MM-dd HH:mm:ss'" locale="zh"
            auto-apply></VueDatePicker>
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
import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';
import { ref } from 'vue';
import VSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';

const emit = defineEmits(['search']);
const formData = ref({});

const onReset = () => {
  formData.value = {};
  emit('search', formData.value);
};
</script>

<style scoped>
:deep(.dp__input) {
  padding: 0.15rem 0.25rem 0.15rem 1.8rem !important;
  font-size: 0.875rem !important;
  height: 31px;
  /* Match bootstrap sm input height */
}

/* 调整 VSelect 样式以匹配 bootstrap sm */
:deep(.vs__dropdown-toggle) {
  padding: 0 0 2px 0;
  height: 31px;
}

:deep(.vs__selected) {
  margin: 2px 2px 0;
  font-size: 0.875rem;
}

:deep(.vs__search) {
  margin: 2px 0 0;
  font-size: 0.875rem;
}
</style>

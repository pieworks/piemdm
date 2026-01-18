<template>
  <div class="px-1">
    <form id="search">
      <div class="form-group row">
        <div class="col-auto">
          <label>Show:</label>
        </div>
        <div class="col-auto">
          <input
            type="text"
            class="form-control form-control-sm"
            placeholder="Code"
            v-model.lazy="formData.code"
            maxlength="36"
            size="40"
          />
        </div>
        <div class="col-auto">
          <label>From:</label>
        </div>
        <div class="col-auto">
          <VueDatePicker
            v-model="formData.startDate"
            placeholder="Select Date"
            :format="'yyyy-MM-dd HH:mm:ss'"
            locale="zh"
          ></VueDatePicker>
        </div>
        <div class="col-auto">
          <label>To:</label>
        </div>
        <div class="col-auto">
          <VueDatePicker
            v-model="formData.endDate"
            placeholder="Select Date"
            :format="'yyyy-MM-dd HH:mm:ss'"
            locale="zh"
          ></VueDatePicker>
        </div>
        <div class="col-auto">
          <button
            type="button"
            class="btn btn-outline-primary btn-sm"
            @click="$emit('search', formData)"
          >
            {{ $t('Go') }}
          </button>
          <button
            type="button"
            class="btn btn-outline-secondary btn-sm ms-1"
            @click="onReset"
          >
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

  const formData = ref({});

  const onReset = () => {
    formData.value = {};
  };
</script>

<style scoped>
  :deep(.dp__input) {
    padding: 0.15rem 0.25rem 0.15rem 1.8rem;
    font-size: 0.875rem;
  }
</style>

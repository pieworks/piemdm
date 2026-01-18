<template>
  <div class="p-1">
    <form id="search">
      <div class="form-group row align-items-center">
        <div class="col-auto">
          <label>{{ $t('Show') }}:</label>
        </div>
        <div class="col-2">
          <v-select
            v-model="formData.showType"
            :reduce="option => option.value"
            :placeholder="$t('Please Select')"
            :options="approvalStatusOptions"
          ></v-select>
        </div>
        <div class="col-auto">
          <label>{{ $t('Time Range') }}:</label>
        </div>
        <div class="col-2">
          <v-select
            v-model="formData.timeRange"
            :reduce="option => option.value"
            :placeholder="$t('Please Select')"
            :options="timeRangeOptions"
          ></v-select>
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
  import { useFormOptions } from '@/composables/useFormOptions';
  import '@vuepic/vue-datepicker/dist/main.css';
  import { ref } from 'vue';
  import vSelect from 'vue-select';
  import 'vue-select/dist/vue-select.css';

  const { approvalStatusOptions, timeRangeOptions } = useFormOptions();

  const formData = ref({
    showType: 'Pending',
    timeRange: 'All',
  });

  const onReset = () => {
    formData.value = {
      showType: 'Pending',
      timeRange: 'All',
    };
  };
</script>

<style scoped></style>

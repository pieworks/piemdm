<template>
  <!-- data table -->
  <!-- margin-top: 115px;min-height: 658px; -->
  <div
    class="table-responsive text-nowrap p-1"
    style="min-height: 65vh; font-size: 0.9rem"
    v-if="total"
  >
    <table class="table table-sm table-bordered table-hover w-auto mb-0">
      <thead class="table-light">
        <tr>
          <th class="text-center">{{ $t('ID') }}</th>
          <th v-for="header in headers">{{ header }}</th>
          <th class="text-center">{{ $t('Actions') }}</th>
        </tr>
      </thead>
      <tbody id="tabletext">
        <tr v-for="(item, index) in list">
          <td class="text-center align-middle">
            {{ item.ID }}
          </td>
          <td>
            <a :href="'/approval/view?id=' + item.Code">{{ item.Code.substr(0, 8) }}</a>
          </td>
          <td>{{ item.Title }}</td>
          <td>{{ item.ApprovalCode }}</td>
          <td>{{ item.Status }}</td>
          <td>{{ item.OperateType }}</td>
          <td>{{ item.TaskNodeName }}</td>
          <td>{{ item.TaskUserName }}</td>
          <td>{{ formatDate(item.UpdatedAt) }}</td>
          <td>{{ formatDate(item.CreatedAt) }}</td>
          <td class="actions px-2">
            <a
              :href="'/approval/approval?id=' + item.ID"
              class="me-1"
            >
              <i class="bi bi-clipboard2-check"></i>
            </a>
            <a
              href="javascript:;"
              class="me-1"
              @click="$emit('getProcessInfo', item.Code)"
            >
              <i class="bi bi-grid-3x2-gap"></i>
            </a>
            <a :href="' /approval/view?id=' + item.ID">
              <i class="bi bi-file-text"></i>
            </a>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <div
    class="table-responsive text-nowrap p-1"
    style="min-height: 75vh"
    v-else
  >
    Your data is empty.
  </div>
</template>
<script setup>
  import { formatDate } from '@/utils/language.js';

  const props = defineProps({
    list: {
      type: Array,
      default: [],
    },
    total: {
      type: Number,
      default: 0,
    },
    headers: {
      type: Array,
      default: [],
    },
  });
</script>
<style></style>

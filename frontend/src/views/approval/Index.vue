<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">
      {{ $t('Approval List') }}
    </div>

    <AppSearch @search="onSearch" />

    <AppResult :page="page" :page-size="pageSize" :total="total"></AppResult>
    <!-- 56+47+38+12  57px 55px -->
    <div class="table-responsive text-nowrap p-1 d-none d-md-block"
      style="min-height: calc(100vh - 268px); overflow-y: auto; font-size: 0.9rem" v-if="total">
      <table class="table table-sm table-bordered table-hover w-100 mb-0">
        <thead class="table-light">
          <tr>
            <th>{{ $t('ID') }}</th>
            <th>{{ $t('Code') }}</th>
            <th>{{ $t('Title') }}</th>
            <th>{{ $t('ApprovalDefCode') }}</th>
            <th>{{ $t('EntityCode') }}</th>
            <th>{{ $t('SerialNumber') }}</th>
            <th>{{ $t('CurrentTaskID') }}</th>
            <th>{{ $t('CurrentTaskName') }}</th>
            <th>{{ $t('FormData') }}</th>
            <th>{{ $t('FormSchema') }}</th>
            <th>{{ $t('Priority') }}</th>
            <th>{{ $t('Urgency') }}</th>
            <th>{{ $t('Description') }}</th>
            <th>{{ $t('Status') }}</th>
            <th>{{ $t('CreatedBy') }}</th>
            <th>{{ $t('UpdatedBy') }}</th>
            <th>{{ $t('CreatedAt') }}</th>
            <th>{{ $t('UpdatedAt') }}</th>
            <th>{{ $t('DeletedAt') }}</th>
          </tr>
        </thead>
        <tbody id="tabletext">
          <tr v-for="item in tableData" :key="item.ID">
            <td class="text-center align-middle">
              {{ item.ID }}
            </td>
            <td>
              <a :href="'/approval/task?id=' + item.ID" :title="$t('Process')">
                {{ item.Code }}
              </a>
            </td>
            <td>{{ item.Title }}</td>
            <td>{{ item.ApprovalDefCode }}</td>
            <td>{{ item.EntityCode }}</td>
            <td>{{ item.SerialNumber }}</td>
            <td>{{ item.CurrentTaskID }}</td>
            <td>{{ item.CurrentTaskName }}</td>
            <td>{{ item.FormData }}</td>
            <td>{{ item.FormSchema }}</td>
            <td>{{ item.Priority }}</td>
            <td>{{ item.Urgency }}</td>
            <td>{{ item.Description }}</td>
            <td>
              <StatusBadge :status="item.Status" />
            </td>
            <td>{{ item.CreatedBy }}</td>
            <td>{{ item.UpdatedBy }}</td>
            <td>{{ formatDate(item.CreatedAt) }}</td>
            <td>{{ formatDate(item.UpdatedAt) }}</td>
            <td>{{ formatDate(item.DeletedAt) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="table-responsive text-nowrap p-1" style="min-height: 75vh" v-else>
      <div class="text-center py-5">
        <i class="bi bi-inbox display-1 text-muted"></i>
        <p class="text-muted mt-3">{{ $t('No data.') }}</p>
      </div>
    </div>

    <AppPagination :page="page" :page-size="pageSize" :total="total" @page-change="pageChange" />
  </div>
</template>

<script setup>
import AppPagination from '@/components/Pagination.vue';
import AppResult from '@/components/Result.vue';
import StatusBadge from '@/components/StatusBadge.vue';
import { useApprovalStore } from '@/pinia/modules/approval';
import { formatDate } from '@/utils/language.js';
import { onMounted, ref } from 'vue';
import AppSearch from './Search.vue';

const approvalStore = useApprovalStore();
const page = ref(1);
const pageSize = ref(15);
const total = ref(0);
const tableData = ref([]);
const formData = ref({});
const loading = ref(false);

onMounted(() => {
  getApprovalData();
});

// Frontend conditional search method
const onSearch = data => {
  formData.value = data;
  getApprovalData();
};

// Get approval data
const getApprovalData = async () => {
  loading.value = true;

  const params = {
    page: page.value,
    pageSize: pageSize.value,
    ...formData.value,
  };

  // Use store method to fetch data
  const result = await approvalStore.fetchApprovalList(params);

  tableData.value = result.list || [];
  total.value = result.total || 0;

  loading.value = false;
};

const pageChange = p => {
  page.value = p;
  getApprovalData();
};
</script>

<style scoped></style>

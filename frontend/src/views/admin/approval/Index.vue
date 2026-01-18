<template>
  <div class="mt-3">
    <!-- search criteria -->
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Approval') }}{{ $t('List') }}</div>

    <AppSearch @search="onSearch" />

    <AppResult :page="page" :page-size="pageSize" :total="total"></AppResult>

    <!-- operation list-->
    <div class="form-group row py-2 px-1">
      <div class="col-sm-10">
        <button type="button" class="btn btn-outline-danger btn-sm me-1" @click="handlerDelete">
          <i class="bi bi-trash3"></i>
          {{ selected.length ? '(' + selected.length + ')' : '' }}
        </button>
      </div>
    </div>

    <!-- data table -->
    <!-- margin-top: 115px;min-height: 658px; -->
    <div class="table-responsive text-nowrap p-1"
      style="min-height: calc(100vh - 310px); overflow-y: auto; font-size: 0.9rem" v-if="total">
      <table class="table table-sm table-bordered table-striped table-hover w-100 mb-0 sticky-table">
        <thead class="table-light">
          <tr>
            <th class="text-center align-middle sticky-col sticky-col-checkbox">
              <input type="checkbox" @click="selectAll" v-model="checked" />
            </th>
            <th>{{ $t('ID') }}</th>
            <th>{{ $t('Code') }}</th>
            <th>{{ $t('Title') }}</th>
            <th>{{ $t('Approval Define Code') }}</th>
            <th>{{ $t('Entity Code') }}</th>
            <th>{{ $t('Serial Number') }}</th>
            <th>{{ $t('Current Task ID') }}</th>
            <th>{{ $t('Current Task Name') }}</th>
            <th>{{ $t('Form Data') }}</th>
            <th>{{ $t('Form Schema') }}</th>
            <th>{{ $t('Priority') }}</th>
            <th>{{ $t('Urgency') }}</th>
            <th>{{ $t('Description') }}</th>
            <th>{{ $t('Status') }}</th>
            <th>{{ $t('Created By') }}</th>
            <th>{{ $t('Updated By') }}</th>
            <th>{{ $t('Created At') }}</th>
            <th>{{ $t('Updated At') }}</th>
            <th>{{ $t('Deleted At') }}</th>
          </tr>
        </thead>
        <tbody id="tabletext">
          <tr v-for="(item, index) in tableData">
            <td class="text-center align-middle sticky-col sticky-col-checkbox">
              <input type="checkbox" v-model="selected" :value="item.ID" number />
            </td>
            <td>{{ item.ID }}</td>
            <td>
              <a :href="'/admin/approval/view?id=' + item.ID">{{ item.Code }}</a>
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
            <td>{{ item.Status }}</td>
            <td>{{ item.CreatedBy }}</td>
            <td>{{ item.UpdatedBy }}</td>
            <td>{{ formatDate(item.UpdatedAt) }}</td>
            <td>{{ formatDate(item.CreatedAt) }}</td>
            <td>{{ formatDate(item.DeletedAt) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="table-responsive text-nowrap p-1" style="min-height: 75vh" v-else>
      {{ $t('Your approval is empty.') }}
    </div>

    <AppPagination :page="page" :page-size="pageSize" :total="total" @page-change="pageChange" />
  </div>
</template>

<script setup>
import { batchDeleteApproval, getApprovalList } from '@/api/approval';
import AppPagination from '@/components/Pagination.vue';
import AppResult from '@/components/Result.vue';
import { AppToast } from '@/components/toast.js';
import { formatDate } from '@/utils/language.js';
import httpLinkHeader from 'http-link-header';
import { onMounted, ref } from 'vue';
import AppSearch from './Search.vue';

const formData = ref({});
const page = ref(1);
const pageSize = ref(15);
const total = ref(0);
const tableData = ref([]);
const selected = ref([]);
const checked = ref(false);

onMounted(() => {
  getTableData();
});

// search approval
const onSearch = searchDate => {
  formData.value = searchDate;
  getTableData();
};

// get approval list
const getTableData = async () => {
  const res = await getApprovalList({
    page: page.value,
    pageSize: pageSize.value,
    ...formData.value,
  });
  if (res) {
    tableData.value = res.data;
    const links = httpLinkHeader.parse(res.headers.link).refs;
    links.forEach(link => {
      if (['last'].includes(link.rel)) {
        const url = new URL(link.uri);
        total.value = parseInt(url.searchParams.get('page')) || 1;
      }
    });
  } else {
  }
};

// select all/ unselect all
const selectAll = () => {
  var ids = [];
  if (!checked.value) {
    tableData.value.forEach(function (val) {
      ids.push(val.ID);
    });
    selected.value = ids;
  } else {
    selected.value = [];
  }
};

const handlerDelete = async row => {
  const res = await batchDeleteApproval({ ids: selected.value });
  if (res) {
    AppToast.show({
      message: '删除成功',
      color: 'success',
    });
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--;
    }
    selected.value = [];
    getTableData();
  }
};

const pageChange = p => {
  page.value = p;
  getTableData();
};
</script>

<style scoped></style>

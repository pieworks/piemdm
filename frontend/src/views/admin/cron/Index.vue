<template>
  <div class="mt-3">
    <!-- search criteria -->
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Cron') }}{{ $t('List') }}</div>

    <AppSearch @search="onSearch" />

    <AppResult :page="page" :page-size="pageSize" :total="total"></AppResult>

    <!-- operation list-->
    <!-- Normal 正常, Frozen 已冻结, Deleted 已删除 -->
    <!-- action: Unfreeze 解冻, Freeze 冻结, Delete 删除 -->
    <div class="form-group row py-2 px-1">
      <div class="col-sm-10">
        <a class="btn btn-outline-primary btn-sm me-1" href="/admin/cron/create" role="button">
          <i class="bi bi-file-earmark-plus"></i>
          {{ $t('New') }}
        </a>
        <button type="button" class="btn btn-outline-primary btn-sm me-1" @click="changeStatus('Frozen')">
          {{ $t('Freeze') }}
        </button>
        <button type="button" class="btn btn-outline-primary btn-sm me-1" @click="changeStatus('Normal')">
          {{ $t('Unfreeze') }}
        </button>
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
            <th class="col-1">{{ $t('Code') }}</th>
            <th class="col-2">{{ $t('Expression') }}</th>
            <th>{{ $t('Name') }}</th>
            <th>{{ $t('Entity Code') }}</th>
            <th>{{ $t('System') }}</th>
            <th>{{ $t('Status') }}</th>
            <th>{{ $t('Created By') }}</th>
            <th>{{ $t('Updated By') }}</th>
            <th>{{ $t('Created At') }}</th>
            <th>{{ $t('Updated At') }}</th>
            <th class="sticky-col sticky-col-actions text-center">{{ $t('Actions') }}</th>
          </tr>
        </thead>
        <tbody id="tabletext">
          <tr v-for="(item, index) in tableData">
            <td class="text-center align-middle sticky-col sticky-col-checkbox">
              <input type="checkbox" v-model="selected" :value="item.ID" number />
            </td>
            <td>
              <a :href="'/admin/cron/view?id=' + item.ID">{{ item.Code }}</a>
            </td>
            <td>{{ item.Expression }}</td>
            <td>{{ item.Name }}</td>
            <td>{{ item.EntityCode }}</td>
            <td>{{ item.System }}</td>
            <td>
              <StatusBadge :status="item.Status" />
            </td>
            <td>{{ item.CreatedBy }}</td>
            <td>{{ item.UpdatedBy }}</td>
            <td>{{ formatDate(item.UpdatedAt) }}</td>
            <td>{{ formatDate(item.CreatedAt) }}</td>
            <td class="sticky-col sticky-col-actions text-center">
              <a :href="'/admin/cron/update?id=' + item.ID">
                <i class="bi bi-pencil"></i>
              </a>
              <a :href="'/admin/cron/view?id=' + item.ID">
                <i class="bi bi-file-text"></i>
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="table-responsive text-nowrap p-1" style="min-height: 75vh" v-else>
      {{ $t('Your cron is empty.') }}
    </div>

    <AppPagination :page="page" :page-size="pageSize" :total="total" @page-change="pageChange" />
  </div>
</template>

<script setup>
import { batchDeleteCron, getCronList, updateCronStatus } from '@/api/cron';
import AppPagination from '@/components/Pagination.vue';
import AppResult from '@/components/Result.vue';
import StatusBadge from '@/components/StatusBadge.vue';
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

// search cron
const onSearch = searchDate => {
  formData.value = searchDate;
  getTableData();
};

// get cron list
const getTableData = async () => {
  const res = await getCronList({
    page: page.value,
    pageSize: pageSize.value,
    ...formData.value,
  });
  if (res.status === 200) {
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

const changeStatus = async status => {
  const res = await updateCronStatus({
    ids: selected.value,
    status: status,
  });
  if (res.status === 200) {
    selected.value = [];
    checked.value = false;
    AppToast.show({
      message: '更新成功',
      color: 'success',
    });
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--;
    }
    getTableData();
  }
};

const handlerDelete = async row => {
  const res = await batchDeleteCron({ ids: selected.value });
  if (res.status === 200) {
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

<template>
  <div class="mt-3">
    <!-- search criteria -->
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Cron Log') }} {{ $t('List') }}</div>

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
            <th class="col-1">{{ $t('ID') }}</th>
            <th class="col-1">{{ $t('Method') }}</th>
            <th class="col-2">{{ $t('Param') }}</th>
            <th>{{ $t('ErrMsg') }}</th>
            <th>{{ $t('Start Time') }}</th>
            <th>{{ $t('End Time') }}</th>
            <th>{{ $t('Exec Time') }}</th>
            <th>{{ $t('Created At') }}</th>
            <th>{{ $t('Updated At') }}</th>
          </tr>
        </thead>
        <tbody id="tabletext">
          <tr v-for="(item, index) in tableData">
            <td class="text-center align-middle sticky-col sticky-col-checkbox">
              <input type="checkbox" v-model="selected" :value="item.ID" number />
            </td>
            <td>{{ item.ID }}</td>
            <td>
              <a :href="'/admin/cron_log/view?id=' + item.ID">{{ item.Method }}</a>
            </td>
            <td>{{ item.Param }}</td>
            <td>{{ item.ErrMsg }}</td>
            <td>{{ formatDate(item.StartTime) }}</td>
            <td>{{ formatDate(item.EndTime) }}</td>
            <td>{{ item.Status }}</td>
            <td>{{ formatDate(item.UpdatedAt) }}</td>
            <td>{{ formatDate(item.CreatedAt) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="table-responsive text-nowrap p-1" style="min-height: 75vh" v-else>
      {{ $t('Your cron log is empty.') }}
    </div>

    <AppPagination :page="page" :page-size="pageSize" :total="total" @page-change="pageChange" />
  </div>
</template>

<script setup>
import { batchDeleteCronLog, getCronLogList } from '@/api/cron_log';
import AppPagination from '@/components/Pagination.vue';
import AppResult from '@/components/Result.vue';
import { AppToast } from '@/components/toast.js';
import { formatDate } from '@/utils/language.js';
import httpLinkHeader from 'http-link-header';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import AppSearch from './Search.vue';

const router = useRouter();
const formData = ref({});
const page = ref(1);
const pageSize = ref(15);
const total = ref(0);
const tableData = ref([]);
const selected = ref([]);
const checked = ref(false);
const params = ref({});

onMounted(() => {
  params.value = router.currentRoute.value.query;
  getCronLogData();
});

// search cron log
const onSearch = searchDate => {
  formData.value = searchDate;
  getCronLogData();
};

// get cron log list
const getCronLogData = async () => {
  const res = await getCronLogList({
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

const handlerDelete = async () => {
  const res = await batchDeleteCronLog({ ids: selected.value });
  if (res) {
    AppToast.show({
      message: '删除成功',
      color: 'success',
    });
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--;
    }
    selected.value = [];
    getCronLogData();
  }
};

const pageChange = p => {
  page.value = p;
  getCronLogData();
};
</script>

<style scoped></style>

<template>
  <div class="mt-3">
    <!-- search criteria -->
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Permission') }}{{ $t('List') }}</div>

    <AppSearch @search="onSearch" />

    <AppResult :page="page" :page-size="pageSize" :total="total"></AppResult>

    <!-- operation list-->
    <div class="form-group row py-2 px-1">
      <div class="col-sm-10">
        <a class="btn btn-outline-primary btn-sm me-1" href="/admin/permission/create" role="button">
          <i class="bi bi-file-earmark-plus"></i>
          {{ $t('New') }}
        </a>
        <button type="button" class="btn btn-outline-danger btn-sm me-1" @click="handlerDelete">
          <i class="bi bi-trash3"></i>
          {{ selected.length ? '(' + selected.length + ')' : '' }}
        </button>
      </div>
    </div>

    <!-- data table -->
    <div class="table-responsive text-nowrap p-1"
      style="min-height: calc(100vh - 310px); overflow-y: auto; font-size: 0.9rem" v-if="total">
      <table class="table table-sm table-bordered table-striped table-hover w-100 mb-0 sticky-table">
        <thead class="table-light">
          <tr>
            <th class="text-center align-middle sticky-col sticky-col-checkbox">
              <input type="checkbox" @click="selectAll" v-model="checked" />
            </th>
            <th>{{ $t('Code') }}</th>
            <th>{{ $t('Name') }}</th>
            <th>{{ $t('Resource') }}</th>
            <th>{{ $t('Action') }}</th>
            <th>{{ $t('ParentID') }}</th>
            <th>{{ $t('Description') }}</th>
            <th class="sticky-col sticky-col-actions text-center">{{ $t('Actions') }}</th>
          </tr>
        </thead>
        <tbody id="tabletext">
          <tr v-for="(item, index) in tableData" :key="item.ID">
            <td class="text-center align-middle sticky-col sticky-col-checkbox">
              <input type="checkbox" v-model="selected" :value="item.ID" number />
            </td>
            <td>
              <a :href="'/admin/permission/view?id=' + item.ID">{{ item.code }}</a>
            </td>
            <td>{{ item.name }}</td>
            <td>{{ item.resource }}</td>
            <td>{{ item.action }}</td>
            <td>{{ item.parent_id }}</td>
            <td>{{ item.description }}</td>
            <td class="sticky-col sticky-col-actions text-center">
              <a :href="'/admin/permission/update?id=' + item.ID">
                <i class="bi bi-pencil"></i>
              </a>
              <a :href="'/admin/permission/view?id=' + item.ID">
                <i class="bi bi-file-text"></i>
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="table-responsive text-nowrap p-1" style="min-height: 75vh" v-else>
      {{ $t('Your permission list is empty.') }}
    </div>

    <AppPagination :page="page" :page-size="pageSize" :total="total" @page-change="pageChange" />
  </div>
</template>

<script setup>
import { batchDeletePermission, getPermissionList } from '@/api/permission';
import AppPagination from '@/components/Pagination.vue';
import AppResult from '@/components/Result.vue';
import { AppToast } from '@/components/toast.js';
import httpLinkHeader from 'http-link-header';
import { onMounted, ref } from 'vue';
import AppSearch from './Search.vue';

const checked = ref(false);
const page = ref(1);
const pageSize = ref(15);
const selected = ref([]);
const total = ref(0);
const tableData = ref([]);
const searchInfo = ref({});

onMounted(() => {
  getData();
});

const onSearch = (data) => {
  searchInfo.value = data;
  page.value = 1;
  getData();
};

const getData = async () => {
  const res = await getPermissionList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value,
  });
  if (res) {
    tableData.value = res.data;
    if (res.headers.link) {
      const links = httpLinkHeader.parse(res.headers.link).refs;
      links.forEach(link => {
        if (['last'].includes(link.rel)) {
          const url = new URL(link.uri);
          total.value = parseInt(url.searchParams.get('page')) || 1;
        }
      });
    } else {
      // Fallback if no link header, though pagination component needs total pages usually
      // If API returns total count in body, use that. Assuming total pages for now based on header convention details.
      // If total is 0, list might be empty or 1 page.
      // Assuming API behavior is consistent.
      if (res.data.length > 0 && total.value === 0) total.value = 1;
    }
  }
};

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
  const res = await batchDeletePermission({ ids: selected.value });
  if (res) {
    AppToast.show({
      message: 'Delete success',
      color: 'success',
    });
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--;
    }
    selected.value = [];
    getData();
  }
};

const pageChange = p => {
  page.value = p;
  getData();
};
</script>

<style scoped></style>

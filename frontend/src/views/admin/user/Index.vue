<template>
  <div class="mt-3">
    <!-- search criteria -->
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('User') }}{{ $t('List') }}</div>

    <AppSearch @search="onSearch" />

    <AppResult :page="page" :page-size="pageSize" :total="total"></AppResult>

    <!-- operation list-->
    <div class="form-group row py-2 px-1">
      <div class="col-sm-10">
        <a class="btn btn-outline-primary btn-sm me-1" href="/admin/user/create" role="button">
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
            <th>{{ $t('ID') }}</th>
            <th class="col-1">{{ $t('EmployeeID') }}</th>
            <th class="col-2">{{ $t('Username') }}</th>
            <th>{{ $t('FirstName') }}</th>
            <th>{{ $t('LastName') }}</th>
            <th>{{ $t('DisplayName') }}</th>
            <th>{{ $t('Email') }}</th>
            <th>{{ $t('Phone') }}</th>
            <th>{{ $t('Language') }}</th>
            <th>{{ $t('Sex') }}</th>
            <th>{{ $t('Desc') }}</th>
            <th>{{ $t('Admin') }}</th>
            <th>{{ $t('Status') }}</th>
            <th>{{ $t('CreatedBy') }}</th>
            <th>{{ $t('UpdatedBy') }}</th>
            <th>{{ $t('CreatedAt') }}</th>
            <th>{{ $t('UpdatedAt') }}</th>
            <th class="sticky-col sticky-col-actions text-center">{{ $t('Actions') }}</th>
          </tr>
        </thead>
        <tbody id="tabletext">
          <tr v-for="(item, index) in tableData">
            <td class="text-center align-middle sticky-col sticky-col-checkbox">
              <input type="checkbox" v-model="selected" :value="item.ID" number />
            </td>
            <td>{{ item.ID }}</td>
            <td>
              <a :href="'/admin/user/view?id=' + item.ID">{{ item.EmployeeID }}</a>
            </td>
            <td>{{ item.Username }}</td>
            <td>{{ item.FirstName }}</td>
            <td>{{ item.LastName }}</td>
            <td>{{ item.DisplayName }}</td>
            <td>{{ item.Email }}</td>
            <td>{{ item.Phone }}</td>
            <td>{{ item.Language }}</td>
            <td>{{ item.Sex }}</td>
            <td>{{ item.Desc }}</td>
            <td>{{ item.Admin }}</td>
            <td>
              <StatusBadge :status="item.Status" />
            </td>
            <td>{{ item.CreatedBy }}</td>
            <td>{{ item.UpdatedBy }}</td>
            <td>{{ formatDate(item.UpdatedAt) }}</td>
            <td>{{ formatDate(item.CreatedAt) }}</td>
            <td class="sticky-col sticky-col-actions text-center">
              <a :href="'/admin/user/update?id=' + item.ID">
                <i class="bi bi-pencil"></i>
              </a>
              <a :href="'/admin/user/view?id=' + item.ID">
                <i class="bi bi-file-text"></i>
              </a>
              <a href="javascript:void(0)" @click="openUserRoles(item)" :title="$t('Roles')">
                <i class="bi bi-person-badge"></i>
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="table-responsive text-nowrap p-1" style="min-height: 75vh" v-else>
      {{ $t('Your data is empty.') }}
    </div>

    <UserRoleModal ref="userRoleModalRef" />

    <AppPagination :page="page" :page-size="pageSize" :total="total" @page-change="pageChange" />
  </div>
</template>

<script setup>
import { batchDeleteUser, getUserList, updateUserStatus } from '@/api/user';
import AppPagination from '@/components/Pagination.vue';
import AppResult from '@/components/Result.vue';
import StatusBadge from '@/components/StatusBadge.vue';
import { AppToast } from '@/components/toast.js';
import { formatDate } from '@/utils/language.js';
import httpLinkHeader from 'http-link-header';
import { onMounted, ref } from 'vue';
import AppSearch from './Search.vue';
import UserRoleModal from './UserRoleModal.vue';

const formData = ref({});
const page = ref(1);
const pageSize = ref(15);
const total = ref(0);
const tableData = ref([]);
const selected = ref([]);
const checked = ref(false);
const userRoleModalRef = ref(null);

onMounted(() => {
  getTableData();
});

// search user
const onSearch = searchDate => {
  formData.value = searchDate;
  getTableData();
};

// get user list
const getTableData = async () => {
  const res = await getUserList({
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

const changeStatus = async status => {
  const res = await updateUserStatus({
    ids: selected.value,
    status: status,
  });
  if (res) {
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
  const res = await batchDeleteUser({ ids: selected.value });
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

const openUserRoles = user => {
  userRoleModalRef.value.open(user);
};
</script>

<style scoped></style>

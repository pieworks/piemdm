<template>
  <div class="mt-3">
    <!-- search criteria -->
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Notification Template') }} {{ $t('List') }}</div>

    <AppSearch @search="onSearch" />

    <AppResult :page="page" :page-size="pageSize" :total="total"></AppResult>

    <!-- operation list-->
    <div class="form-group row py-2 px-1">
      <div class="col-sm-10">
        <a class="btn btn-outline-primary btn-sm me-1" href="/admin/notification_template/create" role="button">
          <i class="bi bi-file-earmark-plus"></i>
          {{ $t('New') }}
        </a>
        <button type="button" class="btn btn-outline-primary btn-sm me-1" @click="changeStatus('Frozen')"
          :disabled="!selected.length">
          {{ $t('Freeze') }}
        </button>
        <button type="button" class="btn btn-outline-primary btn-sm me-1" @click="changeStatus('Normal')"
          :disabled="!selected.length">
          {{ $t('Unfreeze') }}
        </button>
        <button type="button" class="btn btn-outline-danger btn-sm me-1" @click="handleDelete"
          :disabled="!selected.length">
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
            <th>{{ $t('Template Code') }}</th>
            <th>{{ $t('Template Name') }}</th>
            <th>Type</th>
            <th>Notification Type</th>

            <th>Status</th>
            <th>Updated By</th>
            <th>Updated At</th>
            <th class="sticky-col sticky-col-actions text-center">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in tableData" :key="item.ID">
            <td class="text-center align-middle sticky-col sticky-col-checkbox">
              <input type="checkbox" v-model="selected" :value="item.ID" />
            </td>
            <td>
              <a :href="'/admin/notification_template/view?id=' + item.ID">{{ item.TemplateCode }}</a>
            </td>
            <td>{{ item.TemplateName }}</td>
            <td>{{ item.TemplateType }}</td>
            <td>{{ item.NotificationType }}</td>

            <td>
              <StatusBadge :status="item.Status" />
            </td>
            <td>{{ item.UpdatedBy }}</td>
            <td>{{ formatDate(item.UpdatedAt) }}</td>
            <td class="sticky-col sticky-col-actions text-center">
              <a :href="'/admin/notification_template/update?id=' + item.ID">
                <i class="bi bi-pencil"></i>
              </a>
              <a :href="'/admin/notification_template/view?id=' + item.ID" class="ms-2">
                <i class="bi bi-file-text"></i>
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="table-responsive text-nowrap p-1" style="min-height: 75vh" v-else>
      {{ $t('No data found.') }}
    </div>

    <AppPagination :page="page" :page-size="pageSize" :total="total" @page-change="pageChange" />
  </div>
</template>

<script setup>
import { deleteNotificationTemplate, getNotificationTemplateList, updateNotificationTemplateStatus } from '@/api/notification_template';
import { AppModal } from '@/components/Modal/modal.js';
import AppPagination from '@/components/Pagination.vue';
import AppResult from '@/components/Result.vue';
import StatusBadge from '@/components/StatusBadge.vue';
import { AppToast } from '@/components/toast.js';
import { formatDate } from '@/utils/language.js';
import httpLinkHeader from 'http-link-header';
import { onMounted, ref } from 'vue';
import AppSearch from './Search.vue';

const page = ref(1);
const pageSize = ref(15);
const total = ref(0);
const tableData = ref([]);
const selected = ref([]);
const checked = ref(false);
const searchFormData = ref({});

onMounted(() => {
  getData();
});

const onSearch = (data) => {
  searchFormData.value = data;
  page.value = 1;
  getData();
};

const getData = async () => {
  const res = await getNotificationTemplateList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchFormData.value,
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
      total.value = 1; // 假设如果不分页，就是1页
    }
    // 如果后端不支持 link header 分页，可以使用总数
    if (res.data.total) { // 假设返回格式包含 total
      total.value = Math.ceil(res.data.total / pageSize.value);
    }
  }
};

const selectAll = () => {
  if (!checked.value) {
    selected.value = tableData.value.map(item => item.ID);
  } else {
    selected.value = [];
  }
};

const changeStatus = async (status) => {
  const res = await updateNotificationTemplateStatus({
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
    getData();
  }
};

const handleDelete = async () => {
  if (!selected.value.length) return;

  const confirmed = await AppModal.confirm({
    title: 'Delete Confirmation',
    content: 'Are you sure you want to delete these items?',
  });

  if (!confirmed) return;

  // 暂时只支持单个删除 (循环调用)
  for (const id of selected.value) {
    await deleteNotificationTemplate(id);
  }

  AppToast.show({
    message: '删除成功',
    color: 'success',
  });
  selected.value = [];
  checked.value = false;
  getData();
};

const pageChange = (p) => {
  page.value = p;
  getData();
};
</script>

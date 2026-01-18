<template>
  <div class="mt-3">
    <!-- search criteria -->
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Notification Log') }} {{ $t('List') }}</div>

    <AppSearch @search="onSearch" />

    <AppResult :page="page" :page-size="pageSize" :total="total"></AppResult>

    <!-- data table -->
    <div class="table-responsive text-nowrap p-1"
      style="min-height: calc(100vh - 266px); overflow-y: auto; font-size: 0.9rem" v-if="total">
      <table class="table table-sm table-bordered table-striped table-hover w-100 mb-0 sticky-table">
        <thead class="table-light">
          <tr>
            <th>ID</th>
            <th>{{ $t('Recipient') }}</th>
            <th>Type</th>
            <th>{{ $t('Template Code') }}</th>
            <th>Title</th>
            <th>Status</th>
            <th>Retry</th>
            <th>Send Time</th>
            <th>Created At</th>
            <th class="sticky-col sticky-col-actions text-center">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in tableData" :key="item.ID">
            <td>{{ item.ID }}</td>
            <td>
              <div>{{ item.RecipientID }}</div>
            </td>
            <td>{{ item.NotificationType }}</td>
            <td>{{ item.TemplateCode }}</td>
            <td :title="item.Title" class="text-truncate" style="max-width: 200px;">{{ item.Title }}</td>
            <td>
              <StatusBadge :status="item.Status" />
            </td>
            <td>{{ item.RetryCount }} / {{ item.MaxRetryCount }}</td>
            <td>{{ formatDate(item.SendTime) }}</td>
            <td>{{ formatDate(item.CreatedAt) }}</td>
            <td class="sticky-col sticky-col-actions text-center">
              <a :href="'/admin/notification_log/view?id=' + item.ID" class="">
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
import { getNotificationLogList } from '@/api/notification_log';
import AppPagination from '@/components/Pagination.vue';
import AppResult from '@/components/Result.vue';
import StatusBadge from '@/components/StatusBadge.vue';
import { formatDate } from '@/utils/language.js';
import httpLinkHeader from 'http-link-header';
import { onMounted, ref } from 'vue';
import AppSearch from './Search.vue';

const page = ref(1);
const pageSize = ref(15);
const total = ref(0);
const tableData = ref([]);
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
  const res = await getNotificationLogList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchFormData.value,
  });
  console.log(res);
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
      total.value = 1;
    }
    if (res.data.total) {
      total.value = Math.ceil(res.data.total / pageSize.value);
    }
  }
};

const pageChange = (p) => {
  page.value = p;
  getData();
};
</script>

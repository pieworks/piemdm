<template>
  <div class="mt-3">
    <!-- search criteria -->
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Table Field') }}{{ $t('List') }}</div>

    <AppSearch @search="onSearch" />

    <AppResult :page="page" :page-size="pageSize" :total="total"></AppResult>

    <!-- operation list-->
    <div class="form-group row py-2 px-1">
      <div class="col-sm-10">
        <a class="btn btn-outline-primary btn-sm me-1"
          :href="'/admin/table_field/create?table_code=' + params.table_code" role="button">
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
        <button type="button" class="btn btn-outline-primary btn-sm me-1" @click="handlerPublic">
          {{ $t('Public') }}
        </button>
      </div>
    </div>
    <!-- data table -->
    <!-- margin-top: 115px;min-height: 658px; -->
    <div class="table-responsive text-nowrap p-1"
      style="min-height: calc(100vh - 255px); overflow-y: auto; font-size: 0.9rem" v-if="total">
      <table class="table table-sm table-bordered table-striped table-hover w-100 mb-0 sticky-table">
        <thead class="table-light">
          <tr>
            <th class="text-center sticky-col sticky-col-checkbox">
              <input type="checkbox" @click="selectAll" v-model="checked" />
            </th>
            <th class="text-center">{{ $t('ID') }}</th>
            <th class="col-2">{{ $t('Code') }}</th>
            <th class="col-3">{{ $t('Name') }}</th>
            <th>{{ $t('Table Code') }}</th>
            <th>{{ $t('Type') }}</th>
            <th>{{ $t('Length') }}</th>
            <th>{{ $t('Required') }}</th>
            <th>{{ $t('Is Index') }}</th>
            <th>{{ $t('Is Unique') }}</th>
            <th>{{ $t('Index Name') }}</th>
            <th>{{ $t('Index Priority') }}</th>
            <th>{{ $t('Description') }}</th>
            <th>{{ $t('Is Filter') }}</th>
            <th>{{ $t('Is Show') }}</th>
            <th>{{ $t('Group Name') }}</th>
            <th>{{ $t('Sort') }}</th>
            <th>{{ $t('Status') }}</th>
            <th>{{ $t('Created By') }}</th>
            <th>{{ $t('UpdatedBy') }}</th>
            <th>{{ $t('Created At') }}</th>
            <th>{{ $t('UpdatedAt') }}</th>
            <th>{{ $t('DeletedAt') }}</th>
            <th class="sticky-col sticky-col-actions text-center">{{ $t('Actions') }}</th>
          </tr>
        </thead>
        <tbody id="tabletext">
          <tr v-for="(item, index) in tableData">
            <td class="text-center sticky-col sticky-col-checkbox">
              <input type="checkbox" v-model="selected" :value="item.ID" number />
            </td>
            <td class="text-center">
              {{ item.ID }}
            </td>
            <td>
              <a :href="'/admin/table_field/view?table_code=' + params.table_code + '&id=' + item.ID">
                {{ item.Code }}
              </a>
            </td>
            <td>{{ item.Name }}</td>
            <td>{{ item.TableCode }}</td>
            <td>{{ item.Type }}</td>
            <td>{{ item.Length }}</td>
            <td>{{ item.Required }}</td>
            <td>{{ item.IsIndex }}</td>
            <td>{{ item.IsUnique }}</td>
            <td>{{ item.IndexName }}</td>
            <td>{{ item.IndexPriority }}</td>
            <td>{{ item.Description }}</td>
            <td>{{ item.IsFilter }}</td>
            <td>{{ item.IsShow }}</td>
            <td>{{ item.GroupName }}</td>
            <td>{{ item.Sort }}</td>
            <td>
              <StatusBadge :status="item.Status" />
            </td>
            <td>{{ item.CreatedBy }}</td>
            <td>{{ item.UpdatedBy }}</td>
            <td class="text-center">
              {{ formatDate(item.UpdatedAt) }}
            </td>
            <td class="text-center">
              {{ formatDate(item.CreatedAt) }}
            </td>
            <td>{{ item.DeletedAt }}</td>
            <td class="sticky-col sticky-col-actions text-center">
              <a :href="'/admin/table_field/update?table_code=' + params.table_code + '&id=' + item.ID
                ">
                <i class="bi bi-pencil"></i>
              </a>
              <a :href="'/admin/table_field/view?table_code=' + params.table_code + '&id=' + item.ID">
                <i class="bi bi-file-text"></i>
              </a>
            </td>
          </tr>
        </tbody>
      </table>
      <div class="table-responsive text-nowrap p-1">-- end --</div>
    </div>
    <div class="table-responsive text-nowrap p-1" style="min-height: 75vh" v-else>
      {{ $t('Your field is empty.') }}
    </div>

    <!-- <AppPagination :page="page" :page-size="pageSize" :total="total" @page-change="pageChange" /> -->
  </div>
</template>

<script setup>
import {
  batchDeleteTableField,
  findTableFieldList,
  publicTable,
  updateTableFieldStatus,
} from '@/api/table_field';
import AppResult from '@/components/Result.vue';
import StatusBadge from '@/components/StatusBadge.vue';
import { AppToast } from '@/components/toast.js';
import { AppModal } from '@/components/Modal/modal.js';
import { formatDate } from '@/utils/language.js';
import httpLinkHeader from 'http-link-header';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import AppSearch from './Search.vue';

const router = useRouter();
const page = ref(1);
const pageSize = ref(150);
const total = ref(0);
const tableData = ref([]);
const selected = ref([]);
const checked = ref(false);
const formData = ref({});
const params = ref({});

onMounted(() => {
  params.value = router.currentRoute.value.query;
  getTableFieldData();
});

// search table field
const onSearch = searchDate => {
  formData.value = searchDate;
  getTableFieldData();
};

// get table field list
const getTableFieldData = async () => {
  const res = await findTableFieldList({
    table_code: params.value.table_code,
    page: page.value,
    pageSize: pageSize.value,
    ...formData.value,
  });
  if (res.data) {
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
  const res = await updateTableFieldStatus({
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
    getTableFieldData();
  }
};

const handlerDelete = async row => {
  const res = await batchDeleteTableField({ ids: selected.value });
  if (res) {
    AppToast.show({
      message: '删除成功',
      color: 'success',
    });
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--;
    }
    selected.value = [];
    getTableFieldData();
  }
};

const handlerPublic = async () => {
  const res = await publicTable({
    table_code: params.value.table_code,
  });
  if (res && res.data) {
    // 构建表列表的 HTML
    const tableList = res.data.tables
      ? res.data.tables.map(table => `<li>${table}</li>`).join('')
      : '';

    AppModal.alert({
      title: '表结构发布成功',
      bodyHtml: true,  // 启用 HTML 渲染
      bodyContent: `
          <p>以下表已更新:</p>
          <ul style="text-align: left; margin-left: 20px;">
            ${tableList}
          </ul>
        `,
    });

    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    selected.value = [];
    getTableFieldData();
  }
};

const pageChange = p => {
  page.value = p;
  getTableFieldData();
};
</script>

<style scoped></style>

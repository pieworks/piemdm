<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ entityName }} {{ $t('View') }}</div>
    <ul class="nav nav-tabs">
      <li class="nav-item">
        <button class="nav-link active" id="base-info-tab" data-bs-toggle="tab" data-bs-target="#base-info-tab-pane"
          type="button" @click="updateUrlTab('base')">
          {{ $t('Basic Info') }}
        </button>
      </li>
      <li class="nav-item" v-for="item in tableExts" :key="item.Code">
        <button class="nav-link" :id="item.Code + '-tab'" data-bs-toggle="tab"
          :data-bs-target="'#' + item.Code + '-tab-pane'" type="button" @click="updateUrlTab(item.Code)">
          {{ item.Name }}
        </button>
      </li>
    </ul>
    <div class="card-body overlay-wrapper mt-2">
      <div id="view_wrapper">
        <div class="tab-content" id="myTabContent">
          <div class="tab-pane fade show active" id="base-info-tab-pane" tabindex="0">
            <div class="col-12 col-sm-12">
              <div class="row">
                <div class="row mb-2">
                  <!-- 按分组显示字段 -->
                  <template v-for="field in flatFields" :key="field.Code || field.name">
                    <!-- 分组标题 -->
                    <div class="col-12 mt-3 mb-3" v-if="field.isHeader">
                      <h6 class="text-secondary border-bottom pb-1 mb-0 small fw-semibold">
                        <i class="bi bi-bookmark me-2"></i>{{ field.name }}
                      </h6>
                    </div>
                    <!-- 字段显示 -->
                    <div class="col-sm-6" v-else>
                      <div class="row mb-1">
                        <legend :for="field.Code" class="col-form-label col-sm-4 py-0 small text-muted">
                          {{ field.Name }}:
                        </legend>
                        <div class="col-sm-8 my-auto small">
                          <!-- status字段使用颜色标识 -->
                          <span v-if="field.Code === 'status'" :class="getStatusClass(dataInfo[field.Code])">
                            {{ dataInfo[field.Code] }}
                          </span>
                          <!-- 附件字段:显示图片预览 -->
                          <template v-else-if="field.FieldType === 'attachment'">
                            <div class="attachment-preview">
                              <template v-for="(url, index) in parseAttachmentUrls(dataInfo[field.Code])" :key="index">
                                <!-- 图片预览 -->
                                <img v-if="isImageUrl(url)" :src="getFullUrl(url)" :alt="getFilename(url)"
                                  class="attachment-thumbnail me-2 mb-2" :title="getFilename(url)" />
                                <!-- 非图片文件:显示文件名链接 -->
                                <a v-else :href="getFullUrl(url)" target="_blank"
                                  class="d-block text-decoration-none me-2 mb-1">
                                  <i class="bi bi-file-earmark"></i> {{ getFilename(url) }}
                                </a>
                              </template>
                            </div>
                          </template>
                          <template v-else>
                            <input type="text" class="form-control form-control-sm"
                              :value="dataInfo[`${field.Code}_display`] !== null && dataInfo[`${field.Code}_display`] !== undefined ? dataInfo[`${field.Code}_display`] : '-'"
                              disabled />
                          </template>
                        </div>
                      </div>
                    </div>
                  </template>
                </div>
                <!-- 返回按钮 -->
                <div class="form-group row mt-3 mb-3">
                  <div class="col-sm-auto">
                    <button class="btn btn-outline-secondary btn-sm" type="button" @click="goIndex">
                      {{ $t('Back') }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="tab-pane fade" :id="view.Code + '-tab-pane'" tabindex="0" v-for="view in tableExts"
            :key="view.Code" style="min-height: calc(100vh - 325px); overflow-y: auto; font-size: 0.9rem">
            <!-- operation list -->
            <div class="form-group py-2">
              <div class="col-sm-10">
                <button type="button" class="btn btn-outline-primary btn-sm me-1"
                  @click="handlerNew(view.Code, view.ItemField, view.EntityField)">
                  <i class="bi bi-file-earmark-plus"></i>
                  {{ $t('New') }}
                </button>
                <button type="button" class="btn btn-outline-primary btn-sm me-1"
                  @click="changeStatus(view.Code, 'BatchFreeze')">
                  {{ $t('Freeze') }}
                </button>
                <button type="button" class="btn btn-outline-primary btn-sm me-1"
                  @click="changeStatus(view.Code, 'BatchUnfreeze')">
                  {{ $t('Unfreeze') }}
                </button>
                <button type="button" class="btn btn-outline-danger btn-sm me-1"
                  @click="changeStatus(view.Code, 'BatchDelete')">
                  <i class="bi bi-trash3"></i>
                  {{ selected[view.Code]?.length ? '(' + selected[view.Code].length + ')' : '' }}
                </button>
              </div>
            </div>

            <div class="table-responsive text-nowrap p-1">
              <table class="table table-sm table-bordered table-striped table-hover mt-2 w-100 sticky-table">
                <thead class="table-light">
                  <tr>
                    <th class="text-center align-middle sticky-col sticky-col-checkbox">
                      <input type="checkbox" @click="selectAll(view.Code)" v-model="checked[view.Code]" />
                    </th>
                    <th class="sticky-col sticky-col-data">{{ $t('ID') }}</th>
                    <template v-if="extDatas[view.Code]?.tableFields">
                      <!-- 过滤掉 ID 字段防止重复显示 -->
                      <template v-for="field in extDatas[view.Code].tableFields" :key="field.Code">
                        <th v-if="field.Code.toLowerCase() !== 'id'">
                          {{ field.Name }}
                        </th>
                      </template>
                    </template>
                    <th class="text-center sticky-col sticky-col-actions">Actions</th>
                  </tr>
                </thead>
                <tbody v-if="extDatas[view.Code]?.list && extDatas[view.Code].list.length > 0">
                  <tr v-for="item in extDatas[view.Code].list" :key="item.id">
                    <td class="text-center align-middle sticky-col sticky-col-checkbox">
                      <input type="checkbox" v-model="selected[view.Code]" :value="item.id" number />
                    </td>
                    <td class="sticky-col sticky-col-data">{{ item.id }}</td>
                    <template v-if="extDatas[view.Code]?.tableFields">
                      <template v-for="field in extDatas[view.Code].tableFields" :key="field.Code">
                        <td v-if="field.Code.toLowerCase() !== 'id'">
                          <!-- status字段使用颜色标识 -->
                          <span v-if="field.Code.toLowerCase() === 'status'" :class="getStatusClass(item[field.Code])">
                            {{ item[field.Code] }}
                          </span>
                          <!-- 其他字段 -->
                          <template v-else>
                            {{ item[field.Code] }}
                          </template>
                        </td>
                      </template>
                    </template>
                    <td class="actions px-2 sticky-col sticky-col-actions text-center">
                      <a href="javascript:;" @click="getApprovalHistory(view.Code, item.id)" class="me-1"
                        :title="$t('Approval History')">
                        <i class="bi bi-clipboard2-check"></i>
                      </a>
                      <a href="#" @click.prevent="handlerUpdateItem(view.Code, item.id)" class="me-1">
                        <i class="bi bi-pencil"></i>
                      </a>
                    </td>
                  </tr>
                </tbody>
              </table>
              <div v-if="!extDatas[view.Code]?.list || extDatas[view.Code].list.length === 0"
                class="text-center p-3 text-muted">
                暂时没有数据，请扩展 {{ view.Name }} 视图。
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {
  findEntity, findEntityList,
  getEntityList,
  updateEntityStatus,
} from '@/api/entity';
import { getTableList } from '@/api/table';
import { getTableFields } from '@/api/table_field';
import { onMounted, ref, computed, reactive, h } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { formatFieldValue, preloadFieldDictionaries } from '@/utils/fieldFormatter';
import { AppToast } from '@/components/toast.js';
import { AppModal } from '@/components/Modal/modal.js';
import { Tab } from 'bootstrap';
import { nextTick } from 'vue';

const router = useRouter();
const { t } = useI18n();
const dataInfo = ref({});
const tableExts = ref([]);
const tableFields = ref([]);
const allFields = ref([]); // 所有字段(包括系统字段)
const params = ref({});
const extDatas = reactive({});
const entityName = ref('');
const selected = reactive({});
const checked = reactive({});
const reason = ref('');

onMounted(() => {
  params.value = router.currentRoute.value.query;
  getEntityName();
  getEntityInfo();
  getTableExts(params.value.table_code);
});

// 获取表名
const getEntityName = async () => {
  const where = {
    status: 'Normal',
  };
  if (params.value.table_code) {
    where.code = params.value.table_code;
  }
  const res = await getTableList(where);
  if (res && res.length > 0) {
    entityName.value = res.data[0].Name;
  }
};

const getEntityInfo = async () => {
  const res = await findEntity({
    table_code: params.value.table_code,
    id: params.value.id,
  });

  // 先保存原始数据和字段配置
  const rawData = res.data.info;
  tableFields.value = res.data.tableFields;

  // 获取所有字段(包括系统字段) - 用于格式化显示
  await getAllFields();

  // 1. 预加载所有需要的字典数据 (使用 tableFields, 保持与 Form.vue 一致)
  await preloadFieldDictionaries(tableFields.value);

  // 2. 预处理数据,添加 {field}_display 字段 (使用 tableFields 的 Code 字段)
  const processedData = { ...rawData };
  for (const field of tableFields.value) {
    const fieldCode = field.Code;
    const value = rawData[fieldCode];
    // 转换 field 格式以供 formatFieldValue 使用
    const fieldForFormat = {
      code: fieldCode,
      field_type: field.FieldType,
      options: field.Options,
      ...field,
    };
    const formattedValue = await formatFieldValue(value, fieldForFormat);
    // 添加 {field}_display 字段 (使用 PascalCase Code)
    processedData[`${fieldCode}_display`] = formattedValue ?? '';
  }
  dataInfo.value = processedData;

  // 捕获getTableExts错误,不影响字段显示
  try {
    await getTableExts(params.value.table_code); // 使用params而不是dataInfo
  } catch (error) {
    tableExts.value = [];
  }
};

// 获取所有字段(包括系统字段)
const getAllFields = async () => {
  try {
    const res = await getTableFields({
      table_code: params.value.table_code,
    });
    if (res && res.data) {
      // 显示所有字段(包括 is_show 为 false 的字段),并添加 Options 配置
      allFields.value = res.data.map(f => ({
        ...f,
        field_type: f.field_type,
        relation: f.options?.relation,  // 关联配置
      }));
    }
  } catch (error) {
    console.error('获取字段列表失败:', error);
    // 如果获取失败,回退到使用tableFields
    allFields.value = tableFields.value.map(f => ({
      code: f.Code,
      name: f.Name,
      type: f.Type,
      is_show: true,
      is_system: false,
    }));
  }
};

// 按 Sort 字段升序排序 (使用 tableFields 而不是 allFields, 保持与 Form.vue 一致)
const sortedFields = computed(() => {
  if (!tableFields.value || tableFields.value.length === 0) {
    return [];
  }
  // 按 Sort 字段升序排序
  return [...tableFields.value].sort((a, b) => {
    const sortA = a.Sort || 0;
    const sortB = b.Sort || 0;
    return sortA - sortB;
  });
});

// 分组并扁平化字段列表 (与 Form.vue 保持一致的逻辑)
const flatFields = computed(() => {
  if (!sortedFields.value || sortedFields.value.length === 0) return [];

  const groups = {};
  const groupOrder = [];

  sortedFields.value.forEach(field => {
    const groupName = field.GroupName || '基本信息';
    if (!groups[groupName]) {
      groups[groupName] = [];
      groupOrder.push(groupName);
    }
    groups[groupName].push(field);
  });

  const result = [];
  groupOrder.forEach(name => {
    result.push({ isHeader: true, name });
    groups[name].forEach(f => result.push(f));
  });

  return result;
});

// 注意: formatFieldValue 已从 @/utils/fieldFormatter.js 导入
// 数据预处理在 getEntityInfo 中完成,模板直接使用 {field}_display

// 获取status的样式类
const getStatusClass = (status) => {
  const statusMap = {
    'Normal': 'badge bg-light text-dark',
    'Deleted': 'badge bg-danger',
    'Frozen': 'badge bg-warning text-dark',
  };
  return statusMap[status] || 'badge bg-secondary';
};

// get Extension Entity list
const getExtTableData = async (item) => {
  const params = {
    table_code: item.Code,
    pageSize: 100,
    // status: 'Normal',
  };

  // 添加关联过滤
  // item.ItemField: 子表中的关联字段名（如 'supplier_id', 'entity_id'）
  // item.EntityField: 父表中的关联字段名（如 'id', 'code'）
  if (item.ItemField) {
    const parentField = item.EntityField || 'id'; // 默认使用 'id'
    const parentValue = dataInfo.value[parentField];

    if (parentValue) {
      params[item.ItemField] = parentValue;
    }
  }



  let res = { data: [] };
  try {
    res = await findEntityList(params);
  } catch (e) {
    console.warn(`Failed to fetch data for table ${item.Code}:`, e);
    // Suppress error as requested for non-existent tables
    res = { data: [] };
  }

  // 获取相关表的字段定义
  let fields = [];
  try {
    const fieldRes = await getTableFields({ table_code: item.Code });
    if (fieldRes && fieldRes.data) {
      fields = fieldRes.data.map(f => ({
        ...f,
        Name: f.Name || f.name,
        Code: f.Code || f.code
      }));
    }
  } catch (e) {
    console.error('getExtTableData fields error:', e);
  }

  // 构造模板期望的数据结构 { list: [], tableFields: [] }
  extDatas[item.Code] = {
    list: res.data || [],
    tableFields: fields
  };

  checked[item.Code] = false;
  selected[item.Code] = [];
};

// 获取 Extension 的表列表
const getTableExts = async table_code => {
  const res = await getTableList({
    parent_table: table_code,
    pageSize: 100,
    status: 'Normal',
  });

  if (res && res.data && Array.isArray(res.data)) {
    tableExts.value = res.data;
    for (const item of res.data) {
      getExtTableData(item);
    }
    // Deep linking: activate tab if specified in URL
    nextTick(() => {
      const tabParam = router.currentRoute.value.query.tab;
      if (tabParam) {
        const triggerEl = document.querySelector(`button[data-bs-target="#${tabParam}-tab-pane"]`);
        if (triggerEl) {
          const tab = new Tab(triggerEl);
          tab.show();
        }
      }
    });
  } else {
    tableExts.value = [];
  }
};

// Update URL query param when tab is clicked
const updateUrlTab = (tabName) => {
  const query = { ...router.currentRoute.value.query };
  if (tabName === 'base') {
    delete query.tab;
  } else {
    query.tab = tabName;
  }
  router.replace({ query });
};

// select all/ unselect all
const selectAll = table_code => {
  var ids = [];
  if (!checked[table_code]) {
    extDatas[table_code].list.forEach(function (val) {
      ids.push(val.id);
    });
    selected[table_code] = ids;
  } else {
    selected[table_code] = [];
  }
};

const handlerNew = (code, itemField, entityField) => {
  const query = {
    table_code: code,
  };
  // 传递关联字段默认值
  if (itemField) {
    const parentVal = dataInfo.value[entityField || 'id'];
    if (parentVal) {
      query[itemField] = parentVal;
    }
  }
  // Add redirect_url to query
  query.redirect_url = encodeURIComponent(router.currentRoute.value.fullPath);
  router.push({ path: '/entity/create', query: query });
};

const handlerUpdate = (code) => {
  if (!selected[code] || selected[code].length != 1) {
    AppToast.show({
      message: 'Please select one item.',
      color: 'danger',
    });
    return;
  }
  router.push({
    path: '/entity/update',
    query: {
      table_code: code,
      id: selected[code][0],
      redirect_url: encodeURIComponent(router.currentRoute.value.fullPath)
    }
  });
};

// 处理 Item 编辑（从操作列点击）
const handlerUpdateItem = (table_code, itemId) => {
  router.push({
    path: '/entity/update',
    query: {
      table_code: table_code,
      id: itemId,
      redirect_url: encodeURIComponent(router.currentRoute.value.fullPath)
    }
  });
};

// 查看审批历史
const getApprovalHistory = async (table_code, id) => {
  try {
    const res = await getEntityList({
      table_code: table_code,
      entity_id: id,
      is_draft: 'Yes',
    });

    let bodyContent = '<table class="table table-sm table-hover"><thead><tr><th>编码</th><th>流程编码</th><th>创建日期</th></tr></thead><tbody>';

    // 检查是否有数据
    if (res && res.data && res.data.length > 0) {
      res.data.forEach(item => {
        bodyContent += '<tr>' +
          '<td>' + (item.id || '') + '</td>' +
          '<td>' + (item.approval_code || '-') + '</td>' +
          '<td>' + (item.created_at || '') + '</td>' +
          '</tr>';
      });
    } else {
      // 没有数据时显示提示
      bodyContent += '<tr><td colspan="3" class="text-center text-muted">暂无审批历史</td></tr>';
    }

    bodyContent += '</tbody></table>';

    AppModal.alert({
      title: '审批历史',
      bodyHtml: true,
      bodyContent: bodyContent,
      size: 'lg',
    });
  } catch (error) {
    console.error('获取审批历史失败:', error);
    AppToast.show({
      message: '获取审批历史失败',
      color: 'danger',
    });
  }
};

const changeStatus = async (table_code, operation) => {
  // 检查是否选择了记录
  if (!selected[table_code] || selected[table_code].length === 0) {
    AppModal.alert({
      content: '请先选择要操作的记录',
      title: '提示',
    });
    return;
  }

  // 清空原因输入框
  reason.value = '';

  // 根据 operation 设置标题
  let title = '';
  let actionText = '';

  switch (operation) {
    case 'BatchFreeze':
      // 检查是否包含已删除的数据
      for (const id of selected[table_code]) {
        const item = extDatas[table_code]?.list.find(i => i.id === id);
        if (item && item.status === 'Deleted') {
          AppModal.alert({
            content: '包含已删除的数据，无法进行冻结操作',
            title: '提示',
          });
          return;
        }
      }
      title = '冻结确认';
      actionText = '冻结';
      break;
    case 'BatchUnfreeze':
      // 检查是否包含已删除的数据
      for (const id of selected[table_code]) {
        const item = extDatas[table_code]?.list.find(i => i.id === id);
        if (item && item.status === 'Deleted') {
          AppModal.alert({
            content: '包含已删除的数据，无法进行解冻操作',
            title: '提示',
          });
          return;
        }
      }
      title = '解冻确认';
      actionText = '解冻';
      break;
    case 'BatchDelete':
      title = '删除确认';
      actionText = '删除';
      break;
    default:
      title = '操作确认';
      actionText = '操作';
  }

  // 创建自定义内容函数，包含原因输入框
  const content = () => {
    return h('div', [
      h('div', { class: 'mb-3' }, [
        h('p', {}, `确定要${actionText}选中的 ${selected[table_code].length} 条记录吗？`),
      ]),
      h('div', { class: 'mb-3' }, [
        h('label', { class: 'form-label', htmlFor: 'reasonInput' }, '原因：'),
        h('textarea', {
          id: 'reasonInput',
          class: 'form-control',
          rows: 3,
          placeholder: '请输入操作原因...',
          value: reason.value || '',
          onInput: e => {
            reason.value = e.target.value;
          },
        }),
      ]),
    ]);
  };

  const result = await AppModal.confirm({
    title: title,
    content: content,
    size: 'md',
  });

  if (result) {
    // 用户点击了确认，调用 API
    const res = await updateEntityStatus({
      table_code: table_code,
      ids: selected[table_code],
      operation: operation,
      reason: reason.value,
    });
    if (res.status === 200) {
      selected[table_code] = [];
      checked[table_code] = false;
      AppToast.show({
        message: '更新成功',
        color: 'success',
      });
      // Refresh data
      const tableItem = tableExts.value.find(t => t.Code === table_code);
      if (tableItem) {
        getExtTableData(tableItem);
      }
    }
  }
};


// goto index page
function goIndex() {
  router.push('/entity/index?table_code=' + params.value.table_code);
}

// 解析附件 URL (支持单文件字符串和多文件 JSON 数组)
const parseAttachmentUrls = (value) => {
  if (!value) return [];

  // 如果是 JSON 数组字符串
  if (typeof value === 'string' && value.startsWith('[')) {
    try {
      return JSON.parse(value);
    } catch (e) {
      console.warn('Failed to parse attachment JSON:', value);
      return [value];
    }
  }

  // 如果是单个 URL 字符串
  if (typeof value === 'string') {
    return [value];
  }

  // 如果已经是数组
  if (Array.isArray(value)) {
    return value;
  }

  return [];
};

// 检测是否为图片 URL
const isImageUrl = (url) => {
  if (!url) return false;
  return /\.(jpg|jpeg|png|gif|webp|bmp|svg)$/i.test(url);
};

// 从 URL 中提取文件名
const getFilename = (url) => {
  if (!url) return '';
  const parts = url.split('/');
  return parts[parts.length - 1] || 'file';
};

// 获取完整 URL (添加 API 基础路径)
const getFullUrl = (url) => {
  if (!url) return '';
  // 如果已经是完整 URL,直接返回
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url;
  }
  // 添加 API 基础路径
  const baseUrl = import.meta.env.VITE_BASE_API;
  if (!baseUrl) {
    AppModal.alert({
      title: t('ConfigError'),
      bodyContent: `${t('BaseApiNotConfigured')}，${t('AttachmentDisplayDisabled')}。`,
    });
    return url;
  }
  return `${baseUrl}${url}`;
};
</script>

<style scoped>
.attachment-preview {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}

.attachment-thumbnail {
  max-width: 120px;
  max-height: 120px;
  object-fit: cover;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  cursor: pointer;
  transition: transform 0.2s;
}

.attachment-thumbnail:hover {
  transform: scale(1.05);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}
</style>

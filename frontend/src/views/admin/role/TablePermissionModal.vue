<template>
  <div class="modal fade" :id="modalId" tabindex="-1" aria-labelledby="tablePermissionModalLabel" aria-hidden="true">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="tablePermissionModalLabel">
            {{ $t('Data Permission') }} - {{ currentRole?.Name }}
          </h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <div v-if="loading" class="text-center py-5">
            <div class="spinner-border" role="status">
              <span class="visually-hidden">Loading...</span>
            </div>
            <div class="mt-2">{{ $t('Loading') }}...</div>
          </div>
          <div v-else>
            <div class="mb-3">
              <p class="text-muted">
                为角色 <strong>{{ currentRole?.Name }}</strong> 分配可访问的动态表
              </p>
            </div>
            <!-- 动态表列表 -->
            <div class="table-responsive" style="max-height: 500px; overflow-y: auto">
              <table class="table table-sm table-hover">
                <thead class="table-light sticky-top">
                  <tr>
                    <th class="text-center" style="width: 50px">
                      <input type="checkbox" @change="toggleAllTables" :checked="allTablesSelected" />
                    </th>
                    <th>{{ $t('Table Code') }}</th>
                    <th>{{ $t('Table Name') }}</th>
                    <th class="text-center">{{ $t('Status') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="table in allTables" :key="table.ID">
                    <td class="text-center">
                      <input type="checkbox" :checked="isTableSelected(table.Code)" @change="toggleTable(table.Code)" />
                    </td>
                    <td>{{ table.Code }}</td>
                    <td>{{ table.Name }}</td>
                    <td class="text-center">
                      <span v-if="isTableSelected(table.Code)" class="badge bg-success">
                        已授权
                      </span>
                      <span v-else class="badge bg-secondary">未授权</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-if="allTables.length === 0" class="text-center text-muted py-4">
              {{ $t('No tables available') }}
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-outline-secondary btn-sm" data-bs-dismiss="modal">
            {{ $t('Close') }}
          </button>
          <button type="button" class="btn btn-outline-primary btn-sm" @click="savePermissions"
            :disabled="saving || loading">
            <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
            {{ $t('Save') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getTableList } from '@/api/table';
import request from '@/utils/request';
import { AppToast } from '@/components/toast.js';
import { Modal } from 'bootstrap';
import { computed, onMounted, ref } from 'vue';

const modalId = 'tablePermissionModal';
let modalInstance = null;

const currentRole = ref(null);
const allTables = ref([]);
const selectedTables = ref(new Set()); // 选中的表代码集合
const originalPermissions = ref(new Map()); // Map<TableCode, {ID, Status}>
const loading = ref(false);
const saving = ref(false);

onMounted(() => {
  const modalElement = document.getElementById(modalId);
  if (modalElement) {
    modalInstance = new Modal(modalElement);
  }
});

const open = async role => {
  currentRole.value = role;
  selectedTables.value = new Set();
  originalPermissions.value = new Map();
  loading.value = true;

  if (modalInstance) {
    modalInstance.show();
  }

  try {
    await Promise.all([
      loadAllTables(),
      loadRolePermissions(role.ID)
    ]);
  } finally {
    loading.value = false;
  }
};

const loadAllTables = async () => {
  try {
    const res = await getTableList({ pageSize: -1, status: 'Normal' });
    if (res && res.data) {
      allTables.value = res.data;
    }
  } catch (error) {
    console.error('Failed to load tables:', error);
    AppToast.show({
      message: '加载表列表失败',
      color: 'danger',
    });
  }
};

const loadRolePermissions = async roleId => {
  try {
    const res = await request.get(`/admin/table_permissions`, {
      params: { role_id: roleId },
    });

    const data = Array.isArray(res.data) ? res.data : (res.data && res.data.data);

    if (data) {
      data.forEach(p => {
        const code = p.TableCode || p.table_code;
        originalPermissions.value.set(code, {
          ID: p.ID,
          Status: p.Status
        });

        if (p.Status === 'Normal') {
          selectedTables.value.add(code);
        }
      });
    }
  } catch (error) {
    console.error('Failed to load permissions:', error);
    AppToast.show({
      message: '加载权限失败',
      color: 'danger',
    });
  }
};

const isTableSelected = tableCode => {
  return selectedTables.value.has(tableCode);
};

const toggleTable = tableCode => {
  if (selectedTables.value.has(tableCode)) {
    selectedTables.value.delete(tableCode);
  } else {
    selectedTables.value.add(tableCode);
  }
  selectedTables.value = new Set(selectedTables.value);
};

const allTablesSelected = computed(() => {
  return allTables.value.length > 0 &&
    allTables.value.every(table => selectedTables.value.has(table.Code));
});

const toggleAllTables = event => {
  if (event.target.checked) {
    selectedTables.value = new Set(allTables.value.map(t => t.Code));
  } else {
    selectedTables.value = new Set();
  }
};

const savePermissions = async () => {
  if (!currentRole.value) return;

  saving.value = true;
  try {
    const toGrant = [];
    const toDeleteIds = [];

    // 1. 找出需要授权的 (被选中，且 (不存在于原权限 或 原权限状态不是Normal))
    for (const tableCode of selectedTables.value) {
      const original = originalPermissions.value.get(tableCode);
      if (!original || original.Status !== 'Normal') {
        toGrant.push(tableCode);
      }
    }

    // 2. 找出需要删除的 (原权限存在，但现在未选中)
    originalPermissions.value.forEach((val, tableCode) => {
      if (!selectedTables.value.has(tableCode)) {
        toDeleteIds.push(val.ID);
      }
    });

    // 执行授权
    for (const tableCode of toGrant) {
      try {
        await request.post('/admin/table_permissions', {
          role_id: currentRole.value.ID,
          table_code: tableCode,
        });
      } catch (error) {
        console.error(`Failed to grant ${tableCode}:`, error);
      }
    }

    // 执行删除
    if (toDeleteIds.length > 0) {
      try {
        // 假设后端支持批量删除: DELETE /batch { ids: [...] }
        await request.delete('/admin/table_permissions/batch', {
          data: { ids: toDeleteIds }
        });
      } catch (error) {
        console.error('Failed to delete permissions:', error);
      }
    }

    AppToast.show({
      message: '保存成功',
      color: 'success',
    });

    // 重新加载
    await loadRolePermissions(currentRole.value.ID);

  } catch (error) {
    console.error('Failed to save permissions:', error);
    AppToast.show({
      message: '保存过程中发生错误',
      color: 'danger',
    });
  } finally {
    saving.value = false;
  }
};

defineExpose({
  open,
});
</script>

<style scoped>
.sticky-top {
  position: sticky;
  top: 0;
  z-index: 1;
  background-color: #f8f9fa;
}
</style>

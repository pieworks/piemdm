<template>
  <div class="modal fade" :id="modalId" tabindex="-1" aria-labelledby="roleUserModalLabel" aria-hidden="true">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="roleUserModalLabel">
            {{ $t('User Management') }} - {{ currentRole?.Name }}
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
            <!-- 添加用户选择器 -->
            <div class="mb-3">
              <label class="form-label">{{ $t('Add User') }}</label>
              <v-select v-model="selectedUserId" :options="availableUsers" :reduce="user => user.ID" label="Username"
                :placeholder="$t('Select a user...')" @option:selected="addUser">
                <template #option="option">
                  <span>{{ option.Username }} <small class="text-muted" v-if="option.Name || option.DisplayName">({{
                    option.DisplayName || option.Name }})</small></span>
                </template>
                <template #selected-option="option">
                  <span>{{ option.Username }} <small class="text-muted" v-if="option.Name || option.DisplayName">({{
                    option.DisplayName || option.Name }})</small></span>
                </template>
              </v-select>
            </div>

            <!-- 当前用户列表 -->
            <div class="table-responsive" style="max-height: 400px; overflow-y: auto">
              <table class="table table-sm table-hover">
                <thead class="table-light sticky-top">
                  <tr>
                    <th>{{ $t('Username') }}</th>
                    <th>{{ $t('DisplayName') }}</th>
                    <th>{{ $t('Email') }}</th>
                    <th class="text-center">{{ $t('Actions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="user in roleUsers" :key="user.ID">
                    <td>{{ user.Username }}</td>
                    <td>{{ user.DisplayName || user.Name }}</td>
                    <td>{{ user.Email }}</td>
                    <td class="text-center">
                      <button type="button" class="btn btn-sm btn-outline-danger" @click="removeUser(user.ID)">
                        <i class="bi bi-trash"></i>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-if="roleUsers.length === 0" class="text-center text-muted py-4">
              {{ $t('No users in this role') }}
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-outline-secondary btn-sm" data-bs-dismiss="modal">
            {{ $t('Close') }}
          </button>
          <button type="button" class="btn btn-outline-primary btn-sm" @click="saveUsers" :disabled="saving || loading">
            <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
            {{ $t('Save') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getUserList } from '@/api/user';
import request from '@/utils/request';
import { AppToast } from '@/components/toast.js';
import { Modal } from 'bootstrap';
import vSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';
import { computed, onMounted, ref, nextTick } from 'vue';

const modalId = 'roleUserModal';
let modalInstance = null;

const currentRole = ref(null);
const allUsers = ref([]);
const roleUsers = ref([]);
const selectedUserId = ref(null);
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
  roleUsers.value = [];
  selectedUserId.value = '';
  loading.value = true;

  if (modalInstance) {
    modalInstance.show();
  }

  try {
    // 并行加载所有用户和角色当前用户
    await Promise.all([
      loadAllUsers(),
      loadRoleUsers(role.ID)
    ]);
  } finally {
    loading.value = false;
  }
};

const loadAllUsers = async () => {
  try {
    const res = await getUserList({ pageSize: -1, status: 'Normal' });
    if (res && res.data) {
      allUsers.value = res.data;
    }
  } catch (error) {
    console.error('Failed to load users:', error);
    AppToast.show({
      message: '加载用户列表失败',
      color: 'danger',
    });
  }
};

const loadRoleUsers = async roleId => {
  try {
    const res = await request.get(`/admin/roles/${roleId}/users`);
    // API直接返回数组,位于res.data中
    const data = Array.isArray(res.data) ? res.data : (res.data && res.data.data);
    if (data) {
      roleUsers.value = data;
    }
  } catch (error) {
    console.error('Failed to load role users:', error);
    // 如果获取失败,继续显示空列表
  }
};

// 可选择的用户(排除已在角色中的用户)
const availableUsers = computed(() => {
  const roleUserIds = new Set(roleUsers.value.map(u => u.ID));
  return allUsers.value.filter(u => !roleUserIds.has(u.ID));
});

const addUser = (val) => {
  // 支持 v-select 的 option:selected 事件直接传对象
  if (val && typeof val === 'object' && val.ID) {
    roleUsers.value.push(val);
    nextTick(() => {
      selectedUserId.value = null;
    });
    return;
  }

  if (!selectedUserId.value) return;

  const user = allUsers.value.find(u => u.ID == selectedUserId.value);
  if (user) {
    roleUsers.value.push(user);
    selectedUserId.value = null;
  }
};

const removeUser = userId => {
  roleUsers.value = roleUsers.value.filter(u => u.ID !== userId);
};

const saveUsers = async () => {
  saving.value = true;
  try {
    const userIds = roleUsers.value.map(u => u.ID);
    await request.put(`/admin/roles/${currentRole.value.ID}/users`, {
      user_ids: userIds,
    });

    AppToast.show({
      message: '保存成功',
      color: 'success',
    });

    if (modalInstance) {
      modalInstance.hide();
    }
  } catch (error) {
    console.error('Failed to save users:', error);
    AppToast.show({
      message: '保存失败: ' + (error.response?.data?.message || error.message),
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

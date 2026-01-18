<template>
  <div
    class="modal fade"
    :id="modalId"
    tabindex="-1"
    aria-labelledby="userRoleModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="userRoleModalLabel">
            {{ $t('Role Management') }} - {{ currentUser?.Username }}
          </h5>
          <button
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"
          ></button>
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
                为用户 <strong>{{ currentUser?.Username }}</strong> 选择角色
              </p>
            </div>
            <!-- 角色列表 -->
            <div class="table-responsive" style="max-height: 400px; overflow-y: auto">
              <table class="table table-sm table-hover">
                <thead class="table-light sticky-top">
                  <tr>
                    <th class="text-center" style="width: 50px">
                      <input
                        type="checkbox"
                        @change="toggleAllRoles"
                        :checked="allRolesSelected"
                      />
                    </th>
                    <th>{{ $t('Role Code') }}</th>
                    <th>{{ $t('Role Name') }}</th>
                    <th>{{ $t('Description') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="role in allRoles" :key="role.ID">
                    <td class="text-center">
                      <input
                        type="checkbox"
                        :checked="isRoleSelected(role.ID)"
                        @change="toggleRole(role.ID)"
                      />
                    </td>
                    <td>{{ role.Code }}</td>
                    <td>{{ role.Name }}</td>
                    <td>{{ role.Description }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-if="allRoles.length === 0" class="text-center text-muted py-4">
              {{ $t('No roles available') }}
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">
            {{ $t('Close') }}
          </button>
          <button
            type="button"
            class="btn btn-primary"
            @click="saveRoles"
            :disabled="saving || loading"
          >
            <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
            {{ $t('Save') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { getRoleList } from '@/api/role';
  import request from '@/utils/request';
  import { AppToast } from '@/components/toast.js';
  import { Modal } from 'bootstrap';
  import { computed, onMounted, ref } from 'vue';

  const modalId = 'userRoleModal';
  let modalInstance = null;

  const currentUser = ref(null);
  const allRoles = ref([]);
  const selectedRoleIds = ref(new Set());
  const loading = ref(false);
  const saving = ref(false);

  onMounted(() => {
    const modalElement = document.getElementById(modalId);
    if (modalElement) {
      modalInstance = new Modal(modalElement);
    }
  });

  const open = async user => {
    currentUser.value = user;
    selectedRoleIds.value = new Set();
    loading.value = true;

    if (modalInstance) {
      modalInstance.show();
    }

    try {
      // 并行加载所有角色和用户当前角色
      await Promise.all([
        loadAllRoles(),
        loadUserRoles(user.ID)
      ]);
    } finally {
      loading.value = false;
    }
  };

  const loadAllRoles = async () => {
    try {
      const res = await getRoleList({ pageSize: -1, status: 'Normal' });
      if (res && res.data) {
        allRoles.value = res.data;
      }
    } catch (error) {
      console.error('Failed to load roles:', error);
      AppToast.show({
        message: '加载角色列表失败',
        color: 'danger',
      });
    }
  };

  const loadUserRoles = async userId => {
    try {
      const res = await request.get(`/admin/users/${userId}/roles`);
      // API直接返回数组,位于res.data中
      const data = Array.isArray(res.data) ? res.data : (res.data && res.data.data);
      if (data) {
        const roleIds = data.map(role => role.ID);
        selectedRoleIds.value = new Set(roleIds);
      }
    } catch (error) {
      console.error('Failed to load user roles:', error);
      // 如果获取失败,继续显示空的选择状态
    }
  };

  const isRoleSelected = roleId => {
    return selectedRoleIds.value.has(roleId);
  };

  const toggleRole = roleId => {
    if (selectedRoleIds.value.has(roleId)) {
      selectedRoleIds.value.delete(roleId);
    } else {
      selectedRoleIds.value.add(roleId);
    }
    // 触发响应式更新
    selectedRoleIds.value = new Set(selectedRoleIds.value);
  };

  const allRolesSelected = computed(() => {
    return allRoles.value.length > 0 &&
           allRoles.value.every(role => selectedRoleIds.value.has(role.ID));
  });

  const toggleAllRoles = event => {
    if (event.target.checked) {
      // 全选
      selectedRoleIds.value = new Set(allRoles.value.map(r => r.ID));
    } else {
      // 全不选
      selectedRoleIds.value = new Set();
    }
  };

  const saveRoles = async () => {
    saving.value = true;
    try {
      await request.put(`/admin/users/${currentUser.value.ID}/roles`, {
        role_ids: Array.from(selectedRoleIds.value),
      });

      AppToast.show({
        message: '保存成功',
        color: 'success',
      });

      if (modalInstance) {
        modalInstance.hide();
      }
    } catch (error) {
      console.error('Failed to save roles:', error);
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

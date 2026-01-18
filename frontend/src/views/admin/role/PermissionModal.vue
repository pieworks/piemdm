<template>
  <div class="modal fade" id="permissionModal" tabindex="-1" aria-hidden="true" ref="modalRef">
    <div class="modal-dialog modal-lg modal-dialog-scrollable">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ $t('Assign Permissions') }} - {{ roleName }}</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body small">
          <PermissionTree :data="permissionTree" v-model="selectedPermissions" />
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-outline-secondary btn-sm" data-bs-dismiss="modal">{{ $t('Close')
            }}</button>
          <button type="button" class="btn btn-outline-primary btn-sm" @click="save">{{ $t('Save') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { Modal } from 'bootstrap';
import { getPermissionTree } from '@/api/permission';
import { getRolePermissions, assignPermissions } from '@/api/role';
import { AppToast } from '@/components/toast.js';
import PermissionTree from '@/components/PermissionTree.vue';

const modalRef = ref(null);
let modalInstance = null;
const permissionTree = ref([]);
const selectedPermissions = ref([]);
const roleId = ref(null);
const roleName = ref('');

onMounted(async () => {
  modalInstance = new Modal(modalRef.value);
  // Pre-load permission tree
  const res = await getPermissionTree();
  if (res && res.data) {
    permissionTree.value = res.data;
  }
});

const open = async (role) => {
  roleId.value = role.ID;
  roleName.value = role.Name;
  selectedPermissions.value = [];

  // Fetch current permissions
  try {
    const res = await getRolePermissions(role.ID);
    if (res && res.data) {
      // Map permission objects to IDs
      selectedPermissions.value = res.data.map(p => p.id);
    }
    modalInstance.show();
  } catch (e) {
    console.error("Failed to load permissions", e);
    AppToast.show({ message: 'Failed to load permissions', color: 'danger' });
  }
};

const save = async () => {
  try {
    const res = await assignPermissions(roleId.value, { permission_ids: selectedPermissions.value });
    if (res) {
      AppToast.show({ message: 'Permissions assigned successfully', color: 'success' });
      modalInstance.hide();
    }
  } catch (e) {
    console.error("Failed to assign permissions", e);
    AppToast.show({ message: 'Failed to assign permissions', color: 'danger' });
  }
};

defineExpose({ open });
</script>

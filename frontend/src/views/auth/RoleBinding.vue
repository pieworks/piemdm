<template>
  <div class="p-4">
    <div class="flex flex-row w-full space-x-[2rem]">
      <div class="input-group">
        <span class="input-group-text bg-white border-end-0">
          <i class="bi bi-search"></i>
        </span>
        <input
          type="text"
          class="form-control border-start-0 ps-0"
          v-model="search"
          placeholder="Type to search"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
  import request from '@/utils/request';
  import { deleteItem } from '@/utils';

  import { AppToast } from '@/components/toast.js';
  import { computed, onMounted, ref, unref } from 'vue';

  const props = defineProps({
    resource: { type: String },
    subject: { type: String },
  });

  const roleBindings = ref([]);
  const roles = ref([]);
  const subjects = ref([]);

  const showCreate = ref(false);
  const showDelete = ref(-1);

  const newRoleBinding = ref({});

  const createFormRef = ref();

  const search = ref('');
  const filterRoleBindings = computed(() =>
    roleBindings.value.filter(
      data => !search.value || data.name.toLowerCase().includes(search.value.toLowerCase())
    )
  );

  onMounted(() => {
    request.get(`/api/v1/${props.resource}`).then(response => {
      subjects.value = Array.from(response.data.data);
      for (let sub of Array.from(response.data.data)) {
        for (let role of Array.from(sub.roles)) {
          roleBindings.value.push({
            'subject': sub,
            'role': role,
          });
        }
      }
    });
    request.get(`/api/v1/roles`).then(response => {
      roles.value = Array.from(response.data.data);
    });
  });

  const createRoleBinding = () => {
    const form = unref(createFormRef);
    if (!form) {
      return;
    }
    form.validate((valid, err) => {
      if (valid) {
        request
          .post(
            `/api/v1/${props.resource}/${newRoleBinding.value.subject}/roles/${newRoleBinding.value.role}`
          )
          .then(response => {
             AppToast.show({ message: 'Create success', color: 'success' });
            let subject = roleBindings.value.find(
              item => item.subject.id == newRoleBinding.value.subject
            ).subject;
            let role = roles.value.find(item => item.id == newRoleBinding.value.role);
            roleBindings.value.push({
              'subject': subject,
              'role': role,
            });
            showCreate.value = false;
          });
      } else {
        AppToast.show({
          message: 'Input invalid, all fields required',
          color: 'danger',
        });
      }
    });
  };

  const deleteRoleBindings = row => {
    request.delete(`/api/v1/${props.resource}/${row.subject.id}/roles/${row.role.id}`).then(() => {
      AppToast.show({ message: 'Delete success', color: 'success' });
      deleteItem(roleBindings.value, row);
      showDelete.value = -1;
    });
  };
</script>

<style scoped></style>

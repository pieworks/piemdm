<template>
  <div class="w-full justify-center">
    <div class="flex flex-col h-full px-[4rem] py-[2rem] space-y-[1rem]">
      <div class="flex flex-col overflow-hidden rounded-md shadow-md border">
        <div class="flex w-full h-[5rem] items-center">
          <ListView
            class="ml-[1rem]"
            theme="filled"
            size="42"
            fill="#94A3B8"
          />
          <span class="m-[0.75rem] text-2xl font-600">RBAC Policies</span>
        </div>
        <div class="flex h-[3rem] items-center pl-[1rem] bg-slate-100">
          <ul class="nav nav-tabs">
            <li class="nav-item" v-for="(tab, index) in rbacTabs" :key="tab.name">
              <a
                class="nav-link"
                :class="{ active: currentTab === index }"
                href="#"
                @click.prevent="currentTab = index"
              >
                {{ tab.name }}
              </a>
            </li>
          </ul>
        </div>
      </div>

      <div class="card h-max w-full">
        <div class="card-body">
          <div class="flex flex-col w-full">
            <template v-for="(tab, index) in rbacTabs">
              <component
                v-if="currentTab == index"
                :is="tab.component"
                v-bind:key="index"
                :resource="tab.resource"
                :subject="tab.subject"
              />
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref } from 'vue';
  import Role from './role.vue';
  import RoleBinding from './RoleBinding.vue';

  const rbacTabs = [
    {
      name: 'Roles',
      component: Role,
    },
    {
      name: 'GroupRoleBinding',
      component: RoleBinding,
      resource: 'groups',
      subject: 'Group',
    },
    {
      name: 'UserRoleBinding',
      component: RoleBinding,
      resource: 'users',
      subject: 'User',
    },
  ];

  const currentTab = ref(0);

  const handlePolicy = tab => {
    currentTab.value = tab;
  };
</script>

<style scoped></style>

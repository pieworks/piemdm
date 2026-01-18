<template>
  <div class="mt-3" style="min-height: calc(100vh - 105px); font-size: 0.9rem;">
    <div class="row mt-3">
      <div class="col-sm-4 my-2">
        <div class="card">
          <!-- Change Requests -->
          <div class="card-body">
            <h6 class="card-title">{{ $t('Change Requests') }}</h6>
            <!-- <p class="card-text">{{ $t('Statistics') }}</p> -->
            <div class="row gx-2">
              <div class="col-6">
                <span class="text-secondary">{{ $t('Pending') }}:</span>
                {{ approvalStats.Pending || 0 }}
              </div>
              <div class="col-6">
                <span class="text-secondary">{{ $t('Canceled') }}:</span>
                {{ approvalStats.Canceled || 0 }}
              </div>
            </div>
            <div class="row gx-2">
              <div class="col-6">
                <span class="text-secondary">{{ $t('Approved') }}:</span>
                {{ approvalStats.Approved || 0 }}
              </div>
              <div class="col-6">
                <span class="text-secondary">{{ $t('Deleted') }}:</span>
                {{ approvalStats.Deleted || 0 }}
              </div>
            </div>
            <div class="row gx-2">
              <div class="col-6">
                <span class="text-secondary">{{ $t('Rejected') }}:</span>
                {{ approvalStats.Rejected || 0 }}
              </div>
              <div class="col-6">
                <span class="text-secondary">{{ $t('Expired') }}:</span>
                {{ approvalStats.Expired || 0 }}
              </div>
            </div>
            <div class="row gx-2">
              <div class="col-6">
                <span class="text-secondary">{{ $t('Total') }}:</span>
                {{ (approvalStats.Pending || 0)
                  + (approvalStats.Approved || 0)
                  + (approvalStats.Rejected || 0)
                  + (approvalStats.Canceled || 0)
                  + (approvalStats.Deleted || 0)
                  + (approvalStats.Expired || 0) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="row mt-3">
      <div v-for="entity in entityStatsList" :key="entity.code" class="col-sm-2 my-2">
        <div class="card">
          <div class="card-body">
            <h6 class="card-title">{{ entity.name }}</h6>
            <!-- <p class="card-text">{{ $t('Statistics') }}</p> -->
            <div>
              <span class="text-secondary">{{ $t('Normal') }}:</span>
              {{ entity.statistics.Normal || 0 }}
            </div>
            <div>
              <span class="text-secondary">{{ $t('Frozen') }}:</span>
              {{ entity.statistics.Frozen || 0 }}
            </div>
            <div>
              <span class="text-secondary">{{ $t('Deleted') }}:</span>
              {{ entity.statistics.Deleted || 0 }}
            </div>
            <div>
              <span class="text-secondary">{{ $t('Total') }}:</span>
              {{ (entity.statistics.Normal || 0)
                + (entity.statistics.Frozen || 0)
                + (entity.statistics.Deleted || 0) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getApprovalStats } from '@/api/approval';
import { getEntityStats } from '@/api/entity';
import { onMounted, ref } from 'vue';

const entityStatsList = ref([]);
const approvalStats = ref({});

onMounted(() => {
  getStats();
  getEntityStatistics();
});

const getStats = async () => {
  const res = await getApprovalStats();
  if (res && res.data) {
    approvalStats.value = res.data;
  }
};

const getEntityStatistics = async () => {
  const res = await getEntityStats();
  if (res && res.data) {
    entityStatsList.value = res.data;
  }
};
</script>

<style scoped></style>

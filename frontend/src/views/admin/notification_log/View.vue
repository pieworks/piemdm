<template>
  <div class="mt-3" style="min-height: calc(100vh - 105px); font-size: 0.9rem">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Notification Log') }}{{ $t('View') }}</div>
    <div class="card-body overlay-wrapper">
      <div id="create_wrapper">
        <div class="tab-content" id="myTabContent">
          <div class="tab-pane fade show active" id="baseinfo">
            <div id="create_wrapper container">
              <div class="row">
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="ID" class="col-form-label col-sm-2">ID:</legend>
                    <div class="col-sm-auto">{{ dataInfo.ID }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="ApprovalID" class="col-form-label col-sm-2">Approval ID:</legend>
                    <div class="col-sm-auto">{{ dataInfo.ApprovalID }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="TaskID" class="col-form-label col-sm-2">Task ID:</legend>
                    <div class="col-sm-auto">{{ dataInfo.TaskID }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="RecipientID" class="col-form-label col-sm-2">{{ $t('Recipient ID') }}:</legend>
                    <div class="col-sm-auto">{{ dataInfo.RecipientID }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="RecipientType" class="col-form-label col-sm-2">{{ $t('Recipient Type') }}:</legend>
                    <div class="col-sm-auto">{{ dataInfo.RecipientType }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="NotificationType" class="col-form-label col-sm-2">Notification Type:</legend>
                    <div class="col-sm-auto">{{ dataInfo.NotificationType }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="TemplateCode" class="col-form-label col-sm-2">{{ $t('Template Code') }}:</legend>
                    <div class="col-sm-auto">{{ dataInfo.TemplateCode }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="Title" class="col-form-label col-sm-2">Title:</legend>
                    <div class="col-sm-auto">{{ dataInfo.Title }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="Content" class="col-form-label col-sm-2">Content:</legend>
                    <div class="col-sm-auto">
                      <pre>{{ dataInfo.Content }}</pre>
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="Status" class="col-form-label col-sm-2">Status:</legend>
                    <div class="col-sm-auto">
                      <StatusBadge :status="dataInfo.Status" />
                    </div>
                  </div>
                </div>
                <div class="col-sm-12" v-if="dataInfo.ErrorMessage">
                  <div class="form-group row">
                    <legend for="ErrorMessage" class="col-form-label col-sm-2">Error Message:</legend>
                    <div class="col-sm-auto text-danger">{{ dataInfo.ErrorMessage }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="RetryCount" class="col-form-label col-sm-2">Retry Count:</legend>
                    <div class="col-sm-auto">{{ dataInfo.RetryCount }} / {{ dataInfo.MaxRetryCount }}</div>
                  </div>
                </div>
                <div class="col-sm-12" v-if="dataInfo.NextRetryTime">
                  <div class="form-group row">
                    <legend for="NextRetryTime" class="col-form-label col-sm-2">Next Retry Time:</legend>
                    <div class="col-sm-auto">{{ formatDate(dataInfo.NextRetryTime) }}</div>
                  </div>
                </div>
                <div class="col-sm-12" v-if="dataInfo.SendTime">
                  <div class="form-group row">
                    <legend for="SendTime" class="col-form-label col-sm-2">Send Time:</legend>
                    <div class="col-sm-auto">{{ formatDate(dataInfo.SendTime) }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="CreatedAt" class="col-form-label col-sm-2">Created At:</legend>
                    <div class="col-sm-auto">{{ formatDate(dataInfo.CreatedAt) }}</div>
                  </div>
                </div>
                <div class="col-sm-12" v-if="dataInfo.ExtraData">
                  <div class="form-group row">
                    <legend for="ExtraData" class="col-form-label col-sm-2">Extra Data:</legend>
                    <div class="col-sm-auto">
                      <pre>{{ dataInfo.ExtraData }}</pre>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-sm-12 my-3">
              <div class="form-group row">
                <div class="col-sm-10 col-form-label">
                  <button class="btn btn-outline-secondary btn-sm" type="button" @click="goIndex">
                    {{ $t('Back') }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { findNotificationLog } from '@/api/notification_log';
import StatusBadge from '@/components/StatusBadge.vue';
import { formatDate } from '@/utils/language.js';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const dataInfo = ref({});

onMounted(() => {
  getNotificationLogInfo();
});

const getNotificationLogInfo = async () => {
  const params = router.currentRoute.value.query;
  if (params.id) {
    const res = await findNotificationLog(params.id);
    if (res) {
      dataInfo.value = res.data;
    }
  }
};

function goIndex() {
  router.push('/admin/notification_log/index');
}
</script>

<style scoped></style>

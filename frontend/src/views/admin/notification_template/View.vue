<template>
  <div class="mt-3" style="min-height: calc(100vh - 105px); font-size: 0.9rem">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Notification Template') }}{{ $t('View') }}</div>
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
                    <legend for="TemplateCode" class="col-form-label col-sm-2">{{ $t('Template Code') }}:</legend>
                    <div class="col-sm-auto">{{ dataInfo.TemplateCode }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="TemplateName" class="col-form-label col-sm-2">{{ $t('Template Name') }}:</legend>
                    <div class="col-sm-auto">{{ dataInfo.TemplateName }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="TemplateType" class="col-form-label col-sm-2">Template Type:</legend>
                    <div class="col-sm-auto">{{ dataInfo.TemplateType }}</div>
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
                    <legend for="TitleTemplate" class="col-form-label col-sm-2">Title Template:</legend>
                    <div class="col-sm-auto">{{ dataInfo.TitleTemplate }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="ContentTemplate" class="col-form-label col-sm-2">Content Template:</legend>
                    <div class="col-sm-auto">
                      <pre>{{ dataInfo.ContentTemplate }}</pre>
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="Variables" class="col-form-label col-sm-2">Variables:</legend>
                    <div class="col-sm-auto">
                      <pre>{{ dataInfo.Variables }}</pre>
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
                <div class="col-sm-12" v-if="dataInfo.Description">
                  <div class="form-group row">
                    <legend for="Description" class="col-form-label col-sm-2">Description:</legend>
                    <div class="col-sm-auto">{{ dataInfo.Description }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="CreatedBy" class="col-form-label col-sm-2">Created By:</legend>
                    <div class="col-sm-auto">{{ dataInfo.CreatedBy }}</div>
                  </div>
                </div>
                <div class="col-sm-12" v-if="dataInfo.UpdatedBy">
                  <div class="form-group row">
                    <legend for="UpdatedBy" class="col-form-label col-sm-2">Updated By:</legend>
                    <div class="col-sm-auto">{{ dataInfo.UpdatedBy }}</div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend for="CreatedAt" class="col-form-label col-sm-2">Created At:</legend>
                    <div class="col-sm-auto">{{ formatDate(dataInfo.CreatedAt) }}</div>
                  </div>
                </div>
                <div class="col-sm-12" v-if="dataInfo.UpdatedAt">
                  <div class="form-group row">
                    <legend for="UpdatedAt" class="col-form-label col-sm-2">Updated At:</legend>
                    <div class="col-sm-auto">{{ formatDate(dataInfo.UpdatedAt) }}</div>
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
import { findNotificationTemplate } from '@/api/notification_template';
import StatusBadge from '@/components/StatusBadge.vue';
import { formatDate } from '@/utils/language.js';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const dataInfo = ref({});

onMounted(() => {
  getNotificationTemplateInfo();
});

const getNotificationTemplateInfo = async () => {
  const params = router.currentRoute.value.query;
  if (params.id) {
    const res = await findNotificationTemplate(params.id);
    if (res) {
      dataInfo.value = res.data;
    }
  }
};

function goIndex() {
  router.push('/admin/notification_template/index');
}
</script>

<style scoped></style>

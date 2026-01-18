<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">WebhookDelivery View</div>
    <div class="card-body overlay-wrapper">
      <div id="create_wrapper">
        <div
          class="tab-content"
          id="myTabContent"
        >
          <div
            class="tab-pane fade show active"
            id="baseinfo"
          >
            <div id="create_wrapper container">
              <div class="row">
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="ID"
                      class="col-form-label col-sm-2"
                    >
                      ID:
                    </legend>
                    <div class="col-sm-10 my-auto">
                      {{ dataInfo.ID }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="HookID"
                      class="col-form-label col-sm-2"
                    >
                      HookID:
                    </legend>
                    <div class="col-sm-10 my-auto">
                      {{ dataInfo.HookID }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="DeliveryCode"
                      class="col-form-label col-sm-2"
                    >
                      DeliveryCode:
                    </legend>
                    <div class="col-sm-10 my-auto">
                      {{ dataInfo.DeliveryCode }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Event"
                      class="col-form-label col-sm-2"
                    >
                      Event:
                    </legend>
                    <div class="col-sm-10 my-auto">
                      {{ dataInfo.Event }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="EntityID"
                      class="col-form-label col-sm-2"
                    >
                      EntityID:
                    </legend>
                    <div class="col-sm-10 my-auto">
                      {{ dataInfo.EntityID }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="RequestHeaders"
                      class="col-form-label col-sm-2"
                    >
                      RequestHeaders:
                    </legend>
                    <div
                      class="col-sm-10 my-auto"
                      style="word-break: break-all; word-wrap: break-all"
                    >
                      <div v-for="reqHeader in reqHeaders">{{ reqHeader }}</div>
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="RequestPayload"
                      class="col-form-label col-sm-2"
                    >
                      RequestPayload:
                    </legend>
                    <div
                      class="col-sm-10 my-auto"
                      style="word-break: break-all; word-wrap: break-all"
                    >
                      <pre v-text="reqPayload"></pre>
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="ResponseStatus"
                      class="col-form-label col-sm-2"
                    >
                      ResponseStatus:
                    </legend>
                    <div class="col-sm-10 my-auto">
                      {{ dataInfo.ResponseStatus }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="ResponseMessage"
                      class="col-form-label col-sm-2"
                    >
                      ResponseMessage:
                    </legend>
                    <div class="col-sm-10 my-auto">
                      {{ dataInfo.ResponseMessage }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="ResponseHeaders"
                      class="col-form-label col-sm-2"
                    >
                      ResponseHeaders:
                    </legend>
                    <div
                      class="col-sm-10 my-auto"
                      style="word-break: break-all; word-wrap: break-all"
                    >
                      <div v-for="repHeader in repHeaders">{{ repHeader }}</div>
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="ResponseBody"
                      class="col-form-label col-sm-2"
                    >
                      ResponseBody:
                    </legend>
                    <div
                      class="col-sm-10 my-auto"
                      style="word-break: break-all; word-wrap: break-all"
                    >
                      <pre v-text="repBody"></pre>
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="DeliveredAt"
                      class="col-form-label col-sm-2"
                    >
                      DeliveredAt:
                    </legend>
                    <div class="col-sm-10 my-auto">
                      {{ dataInfo.DeliveredAt }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="CompletedAt"
                      class="col-form-label col-sm-2"
                    >
                      CompletedAt:
                    </legend>
                    <div class="col-sm-10 my-auto">
                      {{ dataInfo.CompletedAt }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 my-3 mb-3">
                  <div class="form-group row">
                    <div class="col-sm-auto">
                      <button
                        class="btn btn-outline-secondary btn-sm"
                        type="button"
                        @click="goIndex"
                      >
                        Back
                      </button>
                    </div>
                  </div>
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
  import { getWebhookDelivery } from '@/api/webhook_delivery';
  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';

  const router = useRouter();

  const dataInfo = ref({});
  const reqHeaders = ref({});
  const reqPayload = ref({});
  const repHeaders = ref({});
  const repBody = ref({});

  onMounted(() => {
    getWebhookDeliveryInfo();
  });

  // get webhook delivery info
  const getWebhookDeliveryInfo = async () => {
    const params = router.currentRoute.value.query;
    const res = await getWebhookDelivery({
      id: params.id,
    });
    dataInfo.value = res.data;
    reqHeaders.value = dataInfo.value.RequestHeaders.split('\n');
    if (dataInfo.value.RequestPayload) {
      let reqPayloadObj = JSON.parse(dataInfo.value.RequestPayload);
      reqPayload.value = JSON.stringify(reqPayloadObj, null, 2);
    }
    repHeaders.value = dataInfo.value.ResponseHeaders.split('\n');
    if (dataInfo.value.ResponseBody) {
      let repBodyObj = JSON.parse(dataInfo.value.ResponseBody);
      repBody.value = JSON.stringify(repBodyObj, null, 2);
    }
  };

  // goto index page
  function goIndex() {
    router.push('/admin/webhook_delivery/index');
  }
</script>

<style scoped></style>

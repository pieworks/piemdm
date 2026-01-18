<template>
  <form
    name="entityForm"
    id="entityForm"
    method="post"
    enctype="multipart/form-data"
  >
    <div class="col-12 col-sm-12">
      <div class="row">
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.url }]">
              {{ $t('Url') }}:
            </legend>
            <div class="col-sm-10">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="url"
                v-bind="urlAttrs"
                name="url"
                placeholder="https://example.com/postreceive"
                maxlength="128"
                size="128"
              />
              <div
                v-if="errors.Url"
                class="text-danger small mt-1"
              >
                {{ errors.Url }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.tablecode }]">
              {{ $t('Table Code') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="tableCode"
                v-bind="tableCodeAttrs"
                name="tableCode"
                :placeholder="$t('Table Code')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.TableCode"
                class="text-danger small mt-1"
              >
                {{ errors.TableCode }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.username }]">
              {{ $t('Username') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="username"
                v-bind="usernameAttrs"
                name="username"
                :placeholder="$t('Username')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.Username"
                class="text-danger small mt-1"
              >
                {{ errors.Username }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend
              :class="['col-form-label', 'col-sm-2', { required: requiredFields.contenttype }]"
            >
              {{ $t('Content Type') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect
                v-model="contentType"
                v-bind="contentTypeAttrs"
                name="contentType"
                :reduce="option => option.value"
                :placeholder="$t('Please Select')"
                :options="headerContentTypeOptions"
              ></VSelect>
              <div
                v-if="errors.ContentType"
                class="text-danger small mt-1"
              >
                {{ errors.ContentType }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.secret }]">
              {{ $t('Secret') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="secret"
                v-bind="secretAttrs"
                name="secret"
                :placeholder="$t('Secret')"
                autocomplete="off"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.Secret"
                class="text-danger small mt-1"
              >
                {{ errors.Secret }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12 mb-1">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.events }]">
              {{ $t('Events') }}:
            </legend>
            <div class="col-sm-6">
              <textarea
                class="form-control form-control-sm"
                v-model="events"
                v-bind="eventsAttrs"
                name="events"
                :placeholder="$t('Events')"
                maxlength="255"
                rows="3"
              ></textarea>
              <div
                v-if="errors.Events"
                class="text-danger small mt-1"
              >
                {{ errors.Events }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.status }]">
              {{ $t('Status') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect
                v-model="status"
                v-bind="statusAttrs"
                name="status"
                :reduce="option => option.value"
                :placeholder="$t('Please Select')"
                :options="statusOptions"
              ></VSelect>
              <div
                v-if="errors.Status"
                class="text-danger small mt-1"
              >
                {{ errors.Status }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12 mb-1">
          <div class="form-group row">
            <legend
              :class="['col-form-label', 'col-sm-2', { required: requiredFields.description }]"
            >
              {{ $t('Description') }}:
            </legend>
            <div class="col-sm-6">
              <textarea
                class="form-control form-control-sm"
                v-model="description"
                v-bind="descriptionAttrs"
                name="description"
                :placeholder="$t('Description')"
                maxlength="255"
                rows="3"
              ></textarea>
              <div
                v-if="errors.Description"
                class="text-danger small mt-1"
              >
                {{ errors.Description }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="form-group row mt-3">
      <div class="col-sm-auto">
        <input
          v-if="dataInfo"
          type="hidden"
          name="id"
          value="{{ dataInfo.ID }}"
        />
        <button
          type="button"
          class="btn btn-outline-primary btn-sm me-2 "
          @click="onSubmit()"
        >
          {{ $t('Submit') }}
        </button>
        <button
          type="button"
          class="btn btn-outline-secondary btn-sm"
          @click="$emit('goIndex')"
        >
          {{ $t('Cancel') }}
        </button>
      </div>
    </div>
  </form>
</template>

<script setup>
  import { useFormOptions } from '@/composables/useFormOptions';
  import { yup } from '@/utils/yup-config';
  import { useForm } from 'vee-validate';
  import { computed, watch } from 'vue';
  import VSelect from 'vue-select';
  import 'vue-select/dist/vue-select.css';

  const { statusOptions, headerContentTypeOptions } = useFormOptions();

  // Define props to receive dataInfo from parent component
  const props = defineProps({
    dataInfo: {
      type: Object,
      default: () => ({}),
    },
  });

  // 使用全局配置的默认错误信息
  const validationSchema = yup.object({
    Url: yup
      .string()
      .required()
      .test('is-valid-url', '请输入有效的URL地址', function (value) {
        if (!value) return false;
        const url = /^http(s)?:\/\/[0-9a-zA-Z]+(\.[^\s]+)?(:[0-9]+)?$/;
        return url.test(value);
      }),
    TableCode: yup.string().required().max(64),
    Username: yup.string().max(64),
    ContentType: yup.string().required(),
    Secret: yup.string().max(64),
    Events: yup.string().max(255),
    Status: yup.string().required(),
    Description: yup.string().max(255),
  });

  // required字段映射
  const requiredFields = computed(() => {
    const requiredMap = {};
    Object.keys(validationSchema.fields).forEach(key => {
      const field = validationSchema.fields[key];
      requiredMap[key.toLowerCase()] = field.tests.some(test => test.OPTIONS?.name === 'required');
    });
    return requiredMap;
  });

  // 表单初始化
  const { values, errors, defineField, handleSubmit, setValues } = useForm({
    validationSchema,
    initialValues: props.dataInfo,
  });

  // 字段定义
  const [url, urlAttrs] = defineField('Url');
  const [tableCode, tableCodeAttrs] = defineField('TableCode');
  const [username, usernameAttrs] = defineField('Username');
  const [contentType, contentTypeAttrs] = defineField('ContentType');
  const [secret, secretAttrs] = defineField('Secret');
  const [events, eventsAttrs] = defineField('Events');
  const [status, statusAttrs] = defineField('Status');
  const [description, descriptionAttrs] = defineField('Description');

  const emit = defineEmits(['submitForm', 'goIndex']);

  const onSubmit = handleSubmit(values => {
    emit('submitForm', values);
  });

  // 监听 dataInfo 变化并更新表单值
  watch(
    () => props.dataInfo,
    newDataInfo => {
      if (newDataInfo && Object.keys(newDataInfo).length > 0) {
        setValues(newDataInfo);
      }
    },
    { immediate: true, deep: true }
  );
</script>

<style scoped>
  .required:after {
    content: ' *';
    color: #dc3545;
    font-weight: bold;
  }
</style>

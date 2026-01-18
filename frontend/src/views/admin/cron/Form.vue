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
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.code }]">
              {{ $t('Code') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="code"
                v-bind="codeAttrs"
                name="code"
                :placeholder="$t('Code')"
                maxlength="8"
                size="12"
              />
              <div
                v-if="errors.Code"
                class="text-danger small mt-1"
              >
                {{ errors.Code }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend
              :class="['col-form-label', 'col-sm-2', { required: requiredFields.expression }]"
            >
              {{ $t('Expression') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="expression"
                v-bind="expressionAttrs"
                name="expression"
                :placeholder="$t('Expression')"
                maxlength="32"
                size="32"
              />
              <div
                v-if="errors.Expression"
                class="text-danger small mt-1"
              >
                {{ errors.Expression }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.name }]">
              {{ $t('Name') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="name"
                v-bind="nameAttrs"
                name="name"
                :placeholder="$t('Name')"
                maxlength="128"
                size="64"
              />
              <div
                v-if="errors.Name"
                class="text-danger small mt-1"
              >
                {{ errors.Name }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend
              :class="['col-form-label', 'col-sm-2', { required: requiredFields.entitycode }]"
            >
              {{ $t('Entity Code') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="entityCode"
                v-bind="entityCodeAttrs"
                name="entityCode"
                :placeholder="$t('Entity Code')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.EntityCode"
                class="text-danger small mt-1"
              >
                {{ errors.EntityCode }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.system }]">
              {{ $t('System') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="system"
                v-bind="systemAttrs"
                name="system"
                :placeholder="$t('System')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.System"
                class="text-danger small mt-1"
              >
                {{ errors.System }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.url }]">
              {{ $t('Url') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="url"
                v-bind="urlAttrs"
                name="url"
                :placeholder="$t('Url')"
                maxlength="128"
                size="64"
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
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.protocol }]">
              {{ $t('Protocol') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect
                v-model="protocol"
                v-bind="protocolAttrs"
                name="protocol"
                :reduce="option => option.value"
                :placeholder="$t('Please Select')"
                :options="protocolOptions"
              ></VSelect>

              <div
                v-if="errors.Protocol"
                class="text-danger small mt-1"
              >
                {{ errors.Protocol }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.method }]">
              {{ $t('Method') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect
                v-model="method"
                v-bind="methodAttrs"
                name="method"
                :reduce="option => option.value"
                :placeholder="$t('Please Select')"
                :options="methodOptions"
              ></VSelect>
              <div
                v-if="errors.Method"
                class="text-danger small mt-1"
              >
                {{ errors.Method }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.appId }]">
              {{ $t('AppId') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="appId"
                v-bind="appIdAttrs"
                name="appId"
                :placeholder="$t('AppId')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.AppId"
                class="text-danger small mt-1"
              >
                {{ errors.AppId }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.appKey }]">
              {{ $t('AppKey') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="appKey"
                v-bind="appKeyAttrs"
                name="appKey"
                :placeholder="$t('AppKey')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.AppKey"
                class="text-danger small mt-1"
              >
                {{ errors.AppKey }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.signType }]">
              {{ $t('SignType') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="signType"
                v-bind="signTypeAttrs"
                name="signType"
                :placeholder="$t('SignType')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.SignType"
                class="text-danger small mt-1"
              >
                {{ errors.SignType }}
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
          v-if="props.dataInfo && props.dataInfo.ID"
          type="hidden"
          name="id"
          :value="props.dataInfo.ID"
        />
        <button
          type="button"
          class="btn btn-outline-primary btn-sm me-2"
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
  import { useForm } from 'vee-validate';
  import { computed, watch } from 'vue';
  import VSelect from 'vue-select';
  import 'vue-select/dist/vue-select.css';
  import * as yup from 'yup';

  const { statusOptions, protocolOptions, methodOptions } = useFormOptions();

  // Define props to receive dataInfo from parent component
  const props = defineProps({
    dataInfo: {
      type: Object,
      default: () => ({}),
    },
  });

  // 简化的验证模式 - 直接使用平面结构
  const validationSchema = yup.object({
    Code: yup.string().required().max(8),
    Expression: yup.string().required().max(32),
    Name: yup.string().required().max(128),
    EntityCode: yup.string().max(64),
    System: yup.string().max(64),
    Url: yup.string().max(255),
    Protocol: yup.string().max(32),
    Method: yup.string().max(16),
    AppId: yup.string().max(64),
    AppKey: yup.string().max(64),
    SignType: yup.string().max(64),
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
  const [code, codeAttrs] = defineField('Code');
  const [expression, expressionAttrs] = defineField('Expression');
  const [name, nameAttrs] = defineField('Name');
  const [entityCode, entityCodeAttrs] = defineField('EntityCode');
  const [system, systemAttrs] = defineField('System');
  const [url, urlAttrs] = defineField('Url');
  const [protocol, protocolAttrs] = defineField('Protocol');
  const [method, methodAttrs] = defineField('Method');
  const [appId, appIdAttrs] = defineField('AppId');
  const [appKey, appKeyAttrs] = defineField('AppKey');
  const [signType, signTypeAttrs] = defineField('SignType');
  const [description, descriptionAttrs] = defineField('Description');
  const [status, statusAttrs] = defineField('Status');

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

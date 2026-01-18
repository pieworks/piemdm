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
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.ip }]">
              {{ $t('IP') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="ip"
                v-bind="ipAttrs"
                name="ip"
                :placeholder="$t('IP')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.IP"
                class="text-danger small mt-1"
              >
                {{ errors.IP }}
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
    <div class="form-group row mt-3 mb-3">
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

  const { statusOptions } = useFormOptions();

  // Define props to receive dataInfo from parent component
  const props = defineProps({
    dataInfo: {
      type: Object,
      default: () => ({}),
    },
  });

  // 验证模式 - 基于后端 Application 模型
  const validationSchema = yup.object({
    Name: yup.string().required().max(128),
    IP: yup.string(),
    Status: yup.string().required().oneOf(['Normal', 'Frozen', 'Deleted']),
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
  const [name, nameAttrs] = defineField('Name');
  const [ip, ipAttrs] = defineField('IP');
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

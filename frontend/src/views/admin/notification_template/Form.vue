<template>
  <form name="templateForm" id="templateForm" method="post">
    <div class="col-12  col-sm-12">
      <div class="row">
        <!-- TemplateCode -->
        <div class="col-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: true }]">
              {{ $t('Template Code') }}:
            </legend>
            <div class="col-sm-10">
              <input type="text" class="form-control form-control-sm" v-model="templateCode" v-bind="templateCodeAttrs"
                name="templateCode" :placeholder="$t('Unique code for the template')" maxlength="100" />
              <div v-if="errors.TemplateCode" class="text-danger small mt-1">
                {{ errors.TemplateCode }}
              </div>
            </div>
          </div>
        </div>

        <!-- TemplateName -->
        <div class="col-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: true }]">
              {{ $t('Template Name') }}:
            </legend>
            <div class="col-sm-10">
              <input type="text" class="form-control form-control-sm" v-model="templateName" v-bind="templateNameAttrs"
                name="templateName" :placeholder="$t('Display name')" maxlength="200" />
              <div v-if="errors.TemplateName" class="text-danger small mt-1">
                {{ errors.TemplateName }}
              </div>
            </div>
          </div>
        </div>

        <!-- TemplateType -->
        <div class="col-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: true }]">
              {{ $t('Template Type') }}:
            </legend>
            <div class="col-sm-4">
              <VSelect v-model="templateType" v-bind="templateTypeAttrs" name="templateType"
                :options="templateTypeOptions" :reduce="option => option.code" label="label"
                :placeholder="$t('Select Type')"></VSelect>
              <div v-if="errors.TemplateType" class="text-danger small mt-1">
                {{ errors.TemplateType }}
              </div>
            </div>
          </div>
        </div>

        <!-- NotificationType -->
        <div class="col-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: true }]">
              {{ $t('Notification Type') }}:
            </legend>
            <div class="col-sm-4">
              <VSelect v-model="notificationType" v-bind="notificationTypeAttrs" name="notificationType"
                :options="notificationTypeOptions" :reduce="option => option.code" label="label"
                :placeholder="$t('Select Type')"></VSelect>
              <div v-if="errors.NotificationType" class="text-danger small mt-1">
                {{ errors.NotificationType }}
              </div>
            </div>
          </div>
        </div>

        <!-- TitleTemplate -->
        <div class="col-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: true }]">
              {{ $t('Title Template') }}:
            </legend>
            <div class="col-sm-10">
              <input type="text" class="form-control form-control-sm" v-model="titleTemplate"
                v-bind="titleTemplateAttrs" name="titleTemplate" :placeholder="$t('Title template with variables')"
                maxlength="500" />
              <div v-if="errors.TitleTemplate" class="text-danger small mt-1">
                {{ errors.TitleTemplate }}
              </div>
            </div>
          </div>
        </div>

        <!-- ContentTemplate -->
        <div class="col-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: true }]">
              {{ $t('Content Template') }}:
            </legend>
            <div class="col-sm-10">
              <textarea class="form-control form-control-sm" v-model="contentTemplate" v-bind="contentTemplateAttrs"
                name="contentTemplate" :placeholder="$t('Content template body...')" rows="6"></textarea>
              <div v-if="errors.ContentTemplate" class="text-danger small mt-1">
                {{ errors.ContentTemplate }}
              </div>
            </div>
          </div>
        </div>

        <!-- Variables -->
        <div class="col-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2']">
              {{ $t('Variables (JSON)') }}:
            </legend>
            <div class="col-sm-10">
              <textarea class="form-control form-control-sm" v-model="variables" v-bind="variablesAttrs"
                name="variables" :placeholder="$t('Default variables JSON')" rows="4"></textarea>
              <div v-if="errors.Variables" class="text-danger small mt-1">
                {{ errors.Variables }}
              </div>
            </div>
          </div>
        </div>



        <!-- Status -->
        <div class="col-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2']">
              {{ $t('Status') }}:
            </legend>
            <div class="col-sm-4">
              <VSelect v-model="status" v-bind="statusAttrs" name="status" :options="['Normal', 'Frozen']"
                :placeholder="$t('Select Status')"></VSelect>
              <div v-if="errors.Status" class="text-danger small mt-1">
                {{ errors.Status }}
              </div>
            </div>
          </div>
        </div>

        <!-- Description -->
        <div class="col-12 mb-1">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2']">
              {{ $t('Description') }}:
            </legend>
            <div class="col-sm-10">
              <textarea class="form-control form-control-sm" v-model="description" v-bind="descriptionAttrs"
                name="description" :placeholder="$t('Description')" maxlength="500" rows="3"></textarea>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="form-group row my-3">
      <div class="col-sm-auto">
        <button type="button" class="btn btn-outline-primary btn-sm me-2" @click="onSubmit()">
          {{ $t('Submit') }}
        </button>
        <button type="button" class="btn btn-outline-secondary btn-sm" @click="$emit('goIndex')">
          {{ $t('Cancel') }}
        </button>
      </div>
    </div>
  </form>
</template>

<script setup>
import { yup } from '@/utils/yup-config';
import { useForm } from 'vee-validate';
import { watch } from 'vue';
import VSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';

const props = defineProps({
  dataInfo: {
    type: Object,
    default: () => ({}),
  },
});

const emit = defineEmits(['submitForm', 'goIndex']);

// 验证规则
const validationSchema = yup.object({
  TemplateCode: yup.string().required().max(100),
  TemplateName: yup.string().required().max(200),
  TemplateType: yup.string().required(),
  NotificationType: yup.string().required(),
  TitleTemplate: yup.string().required().max(500),
  ContentTemplate: yup.string().required(),
  Variables: yup.string().test('is-json', 'Must be valid JSON Object', (value) => {
    if (!value) return true;
    try {
      const parsed = JSON.parse(value);
      // Backend expects map[string]any, so it must be a JSON object, not array or primitive
      return typeof parsed === 'object' && parsed !== null && !Array.isArray(parsed);
    } catch (e) {
      return false;
    }
  }),
});

// 表单初始化
const { values, errors, defineField, handleSubmit, setValues } = useForm({
  validationSchema,
  initialValues: {
    Variables: '{}',
    Status: 'Normal',
    ...props.dataInfo,
  },
});

// 字段定义
const [templateCode, templateCodeAttrs] = defineField('TemplateCode');
const [templateName, templateNameAttrs] = defineField('TemplateName');
const [templateType, templateTypeAttrs] = defineField('TemplateType');
const [notificationType, notificationTypeAttrs] = defineField('NotificationType');
const [titleTemplate, titleTemplateAttrs] = defineField('TitleTemplate');
const [contentTemplate, contentTemplateAttrs] = defineField('ContentTemplate');
const [variables, variablesAttrs] = defineField('Variables');

const [status, statusAttrs] = defineField('Status');
const [description, descriptionAttrs] = defineField('Description');

const templateTypeOptions = [
  { code: 'approval_start', label: 'approval_start (审批发起)' },
  { code: 'approval_pending', label: 'approval_pending (待审批)' },
  { code: 'approval_approved', label: 'approval_approved (审批通过)' },
  { code: 'approval_rejected', label: 'approval_rejected (审批拒绝)' },
  { code: 'approval_timeout', label: 'approval_timeout (审批超时)' },
  { code: 'approval_cancel', label: 'approval_cancel (审批取消)' },
  { code: 'task_assigned', label: 'task_assigned (任务分配)' },
  { code: 'task_transferred', label: 'task_transferred (任务转交)' },
  { code: 'task_reminder', label: 'task_reminder (任务催办)' },
];

const notificationTypeOptions = [
  { code: 'email', label: 'email (邮件)' },
  { code: 'sms', label: 'sms (短信)' },
  { code: 'internal', label: 'internal (站内信)' },
  { code: 'webhook', label: 'webhook (Webhook)' },
];

// 监听 dataInfo 变化
watch(
  () => props.dataInfo,
  newDataInfo => {
    if (newDataInfo && Object.keys(newDataInfo).length > 0) {
      const formData = { ...newDataInfo };
      // 如果 Variables 是对象，转换为 JSON 字符串格式以便编辑
      if (formData.Variables && typeof formData.Variables === 'object') {
        try {
          formData.Variables = JSON.stringify(formData.Variables, null, 2);
        } catch (e) {
          console.error('Failed to stringify Variables:', e);
        }
      } else if (!formData.Variables) {
        // Default to {} if empty
        formData.Variables = '{}';
      } else if (typeof formData.Variables === 'string') {
        // Try to pretty-print if it's already a string
        try {
          const parsed = JSON.parse(formData.Variables);
          formData.Variables = JSON.stringify(parsed, null, 2);
        } catch (e) {
          // Keep as is if parsing fails
        }
      }
      setValues(formData);
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

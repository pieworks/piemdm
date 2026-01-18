/**
 * Widget 组件映射表
 * 将 widget 名称映射到实际的 Vue 组件
 */

// 基础输入组件
const Input = {
  name: 'Input',
  template: `
    <input
      type="text"
      class="form-control form-control-sm"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      v-bind="widgetProps"
    />
  `,
  props: ['modelValue', 'widgetProps'],
  emits: ['update:modelValue']
};

const Textarea = {
  name: 'Textarea',
  template: `
    <textarea
      class="form-control form-control-sm"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      v-bind="widgetProps"
    ></textarea>
  `,
  props: ['modelValue', 'widgetProps'],
  emits: ['update:modelValue']
};

const InputNumber = {
  name: 'InputNumber',
  template: `
    <input
      type="number"
      class="form-control form-control-sm"
      :value="modelValue"
      @input="$emit('update:modelValue', parseFloat($event.target.value) || 0)"
      v-bind="widgetProps"
    />
  `,
  props: ['modelValue', 'widgetProps'],
  emits: ['update:modelValue']
};

// 选择组件
const Checkbox = {
  name: 'Checkbox',
  template: `
    <div class="form-check">
      <input
        type="checkbox"
        class="form-check-input"
        :checked="modelValue"
        @change="$emit('update:modelValue', $event.target.checked)"
        v-bind="widgetProps"
      />
    </div>
  `,
  props: ['modelValue', 'widgetProps'],
  emits: ['update:modelValue']
};

const Select = {
  name: 'Select',
  template: `
    <select
      class="form-select form-select-sm"
      :value="modelValue"
      @change="$emit('update:modelValue', $event.target.value)"
      v-bind="widgetProps"
    >
      <option value="">请选择</option>
      <option
        v-for="option in options"
        :key="option.code || option.value"
        :value="option.code || option.value"
      >
        {{ option.name || option.label }}
      </option>
    </select>
  `,
  props: ['modelValue', 'widgetProps', 'options'],
  emits: ['update:modelValue']
};

const MultiSelect = {
  name: 'MultiSelect',
  template: `
    <select
      multiple
      class="form-select form-select-sm"
      :value="modelValue"
      @change="handleChange"
      v-bind="widgetProps"
    >
      <option
        v-for="option in options"
        :key="option.code || option.value"
        :value="option.code || option.value"
      >
        {{ option.name || option.label }}
      </option>
    </select>
  `,
  props: ['modelValue', 'widgetProps', 'options'],
  emits: ['update:modelValue'],
  methods: {
    handleChange(event) {
      const selected = Array.from(event.target.selectedOptions).map(opt => opt.value);
      this.$emit('update:modelValue', selected);
    }
  }
};

const RadioGroup = {
  name: 'RadioGroup',
  template: `
    <div class="radio-group">
      <div
        v-for="option in options"
        :key="option.code || option.value"
        class="form-check"
      >
        <input
          type="radio"
          class="form-check-input"
          :value="option.code || option.value"
          :checked="modelValue === (option.code || option.value)"
          @change="$emit('update:modelValue', option.code || option.value)"
          :name="radioName"
        />
        <label class="form-check-label">
          {{ option.name || option.label }}
        </label>
      </div>
    </div>
  `,
  props: ['modelValue', 'widgetProps', 'options'],
  emits: ['update:modelValue'],
  computed: {
    radioName() {
      return 'radio_' + Math.random().toString(36).substr(2, 9);
    }
  }
};

const CheckboxGroup = {
  name: 'CheckboxGroup',
  template: `
    <div class="checkbox-group">
      <div
        v-for="option in options"
        :key="option.code || option.value"
        class="form-check"
      >
        <input
          type="checkbox"
          class="form-check-input"
          :value="option.code || option.value"
          :checked="isChecked(option.code || option.value)"
          @change="handleChange(option.code || option.value, $event.target.checked)"
        />
        <label class="form-check-label">
          {{ option.name || option.label }}
        </label>
      </div>
    </div>
  `,
  props: ['modelValue', 'widgetProps', 'options'],
  emits: ['update:modelValue'],
  methods: {
    isChecked(value) {
      return Array.isArray(this.modelValue) && this.modelValue.includes(value);
    },
    handleChange(value, checked) {
      let newValue = Array.isArray(this.modelValue) ? [...this.modelValue] : [];
      if (checked) {
        if (!newValue.includes(value)) {
          newValue.push(value);
        }
      } else {
        newValue = newValue.filter(v => v !== value);
      }
      this.$emit('update:modelValue', newValue);
    }
  }
};

// 日期时间组件
const DatePicker = {
  name: 'DatePicker',
  template: `
    <input
      type="date"
      class="form-control form-control-sm"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      v-bind="widgetProps"
    />
  `,
  props: ['modelValue', 'widgetProps'],
  emits: ['update:modelValue']
};

const TimePicker = {
  name: 'TimePicker',
  template: `
    <input
      type="time"
      class="form-control form-control-sm"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      v-bind="widgetProps"
    />
  `,
  props: ['modelValue', 'widgetProps'],
  emits: ['update:modelValue']
};

const DateTimePicker = {
  name: 'DateTimePicker',
  template: `
    <input
      type="datetime-local"
      class="form-control form-control-sm"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      v-bind="widgetProps"
    />
  `,
  props: ['modelValue', 'widgetProps'],
  emits: ['update:modelValue']
};

// 高级组件
const Upload = {
  name: 'Upload',
  template: `
    <input
      type="file"
      class="form-control form-control-sm"
      @change="handleFileChange"
      v-bind="widgetProps"
    />
  `,
  props: ['modelValue', 'widgetProps'],
  emits: ['update:modelValue'],
  methods: {
    handleFileChange(event) {
      const file = event.target.files[0];
      if (file) {
        this.$emit('update:modelValue', file.name);
      }
    }
  }
};

const RichTextEditor = {
  name: 'RichTextEditor',
  template: `
    <textarea
      class="form-control form-control-sm"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      rows="10"
      v-bind="widgetProps"
    ></textarea>
  `,
  props: ['modelValue', 'widgetProps'],
  emits: ['update:modelValue']
};

// Widget 映射表
export const widgetMap = {
  // 基础输入
  Input,
  Textarea,
  InputNumber,

  // 选择
  Checkbox,
  Select,
  MultiSelect,
  RadioGroup,
  CheckboxGroup,

  // 日期时间
  DatePicker,
  TimePicker,
  DateTimePicker,

  // 高级
  Upload,
  RichTextEditor,

  // 别名支持
  input: Input,
  textarea: Textarea,
  number: InputNumber,
  checkbox: Checkbox,
  select: Select,
  multiselect: MultiSelect,
  radio: RadioGroup,
  checkboxgroup: CheckboxGroup,
  date: DatePicker,
  time: TimePicker,
  datetime: DateTimePicker,
  upload: Upload,
  richtext: RichTextEditor,
};

/**
 * 获取 Widget 组件
 * @param {string} widgetName - Widget 名称
 * @returns {object} Vue 组件
 */
export function getWidget(widgetName) {
  return widgetMap[widgetName] || Input; // 默认返回 Input
}

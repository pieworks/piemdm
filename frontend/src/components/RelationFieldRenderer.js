/**
 * 关联字段渲染组件
 * 用于在列表中显示关联对象的 name 字段
 */

export const RelationFieldRenderer = {
  name: 'RelationFieldRenderer',
  props: {
    value: {
      type: [String, Object],
      default: null
    },
    field: {
      type: Object,
      required: true
    }
  },
  computed: {
    displayValue() {
      // 如果值是对象（嵌套对象），显示 name 字段
      if (this.value && typeof this.value === 'object') {
        return this.value.name || this.value.code || '-';
      }
      // 如果值是字符串（只有 code），直接显示
      return this.value || '-';
    }
  },
  template: `
    <span>{{ displayValue }}</span>
  `
};

/**
 * 字段类型预设配置
 * 与后端 field_type_presets.go 保持一致
 */

export const fieldTypePresets = {
  // ===== 基本类型 (Basic) =====
  text: {
    label: '单行文本',
    group: 'basic',
    icon: 'bi-input-cursor-text',
    dataType: 'Text',
    length: 255,
    ui: {
      widget: 'Input'
    },
    validation: {
      max: 255
    }
  },
  textarea: {
    label: '多行文本',
    group: 'basic',
    icon: 'bi-textarea-t',
    dataType: 'Text',
    length: 2000,
    ui: {
      widget: 'Textarea'
    },
    validation: {
      max: 2000
    }
  },
  phone: {
    label: '手机号码',
    group: 'basic',
    icon: 'bi-telephone',
    dataType: 'Text',
    length: 20,
    ui: {
      widget: 'Input',
      widgetProps: {
        type: 'tel'
      }
    },
    validation: {
      format: 'phone',
      pattern: '^1[3-9]\\d{9}$',
      max: 20
    }
  },
  email: {
    label: '电子邮箱',
    group: 'basic',
    icon: 'bi-envelope',
    dataType: 'Text',
    length: 100,
    ui: {
      widget: 'Input',
      widgetProps: {
        type: 'email'
      }
    },
    validation: {
      format: 'email',
      max: 100
    }
  },
  url: {
    label: 'URL',
    group: 'basic',
    icon: 'bi-link-45deg',
    dataType: 'Text',
    length: 255,
    ui: {
      widget: 'Input',
      widgetProps: {
        type: 'url'
      }
    },
    validation: {
      format: 'url',
      max: 255
    }
  },
  integer: {
    label: '整数',
    group: 'basic',
    icon: 'bi-123',
    dataType: 'Number',
    ui: {
      widget: 'InputNumber',
      widgetProps: {
        step: 1
      }
    },
    validation: {
      validator: 'integer'
    }
  },
  decimal: {
    label: '小数',
    group: 'basic',
    icon: 'bi-calculator',
    dataType: 'Number',
    precision: 10,
    scale: 2,
    ui: {
      widget: 'InputNumber',
      widgetProps: {
        step: 0.01
      }
    }
  },
  percent: {
    label: '百分比',
    group: 'basic',
    icon: 'bi-percent',
    dataType: 'Number',
    precision: 5,
    scale: 2,
    ui: {
      widget: 'InputNumber',
      widgetProps: {
        step: 0.01,
        min: 0,
        max: 100,
        suffix: '%'
      }
    },
    validation: {
      min: 0,
      max: 100
    }
  },
  password: {
    label: '密码',
    group: 'basic',
    icon: 'bi-key',
    dataType: 'Text',
    length: 128,
    ui: {
      widget: 'Input',
      widgetProps: {
        type: 'password'
      }
    },
    validation: {
      min: 6,
      max: 128
    }
  },

  // ===== 选择类型 (Choices) =====
  checkbox: {
    label: '勾选',
    group: 'choices',
    icon: 'bi-check-square',
    dataType: 'Boolean',
    ui: {
      widget: 'Checkbox'
    }
  },
  select: {
    label: '下拉单选',
    group: 'choices',
    icon: 'bi-list-ul',
    dataType: 'Text',
    length: 64,
    ui: {
      widget: 'Select'
    },
    requireDatasource: true,
    requireRelation: true
  },
  multiselect: {
    label: '下拉多选',
    group: 'choices',
    icon: 'bi-list-check',
    dataType: 'Text',
    length: 500,
    ui: {
      widget: 'MultiSelect'
    },
    requireDatasource: true,
    requireRelation: true
  },
  radio: {
    label: '单选框',
    group: 'choices',
    icon: 'bi-circle',
    dataType: 'Text',
    length: 64,
    ui: {
      widget: 'RadioGroup'
    },
    requireDatasource: true,
    requireRelation: true
  },
  checkboxgroup: {
    label: '复选框组',
    group: 'choices',
    icon: 'bi-check2-square',
    dataType: 'Text',
    length: 500,
    ui: {
      widget: 'CheckboxGroup'
    },
    requireDatasource: true,
    requireRelation: true
  },

  // ===== 日期时间 (DateTime) =====
  date: {
    label: '日期',
    group: 'datetime',
    icon: 'bi-calendar',
    dataType: 'Date',
    ui: {
      widget: 'DatePicker',
      widgetProps: {
        format: 'YYYY-MM-DD'
      }
    }
  },
  time: {
    label: '时间',
    group: 'datetime',
    icon: 'bi-clock',
    dataType: 'Text',
    length: 8,
    ui: {
      widget: 'TimePicker',
      widgetProps: {
        format: 'HH:mm:ss'
      }
    }
  },
  datetime: {
    label: '日期时间',
    group: 'datetime',
    icon: 'bi-calendar-event',
    dataType: 'Date',
    ui: {
      widget: 'DateTimePicker',
      widgetProps: {
        format: 'YYYY-MM-DD HH:mm:ss',
        showTime: true
      }
    }
  },

  // ===== 关系类型 (Relation) ===== [暂未实现]
  // belongsto: {
  //   label: '多对一',
  //   group: 'relation',
  //   icon: 'bi-arrow-right',
  //   dataType: 'Text',
  //   length: 64,
  //   ui: {
  //     widget: 'Select'
  //   },
  //   requireRelation: true
  // },
  // hasmany: {
  //   label: '一对多',
  //   group: 'relation',
  //   icon: 'bi-arrow-left-right',
  //   dataType: 'Text',
  //   length: 500,
  //   ui: {
  //     widget: 'SubTable'
  //   },
  //   requireRelation: true
  // },
  // manytomany: {
  //   label: '多对多',
  //   group: 'relation',
  //   icon: 'bi-diagram-3',
  //   dataType: 'Text',
  //   length: 500,
  //   ui: {
  //     widget: 'MultiSelect'
  //   },
  //   requireRelation: true
  // },

  // ===== 高级类型 (Advanced) =====
  autocode: {
    label: '自动编码',
    group: 'advanced',
    icon: 'bi-hash',
    dataType: 'Text',
    length: 64,
    ui: {
      widget: 'Input',
      widgetProps: {
        disabled: true
      }
    }
  },
  // formula: {  // 暂未实现
  //   label: '公式',
  //   group: 'advanced',
  //   icon: 'bi-calculator-fill',
  //   dataType: 'Text',
  //   length: 255,
  //   ui: {
  //     widget: 'Input',
  //     widgetProps: {
  //       disabled: true
  //     }
  //   }
  // },
  // json: {  // 暂时注释掉,以后实现
  //   label: 'JSON',
  //   group: 'advanced',
  //   icon: 'bi-braces',
  //   dataType: 'Text',
  //   length: 2000,
  //   ui: {
  //     widget: 'Textarea',
  //     widgetProps: {
  //       rows: 10
  //     }
  //   }
  // },
  attachment: {
    label: '附件',
    group: 'advanced',
    icon: 'bi-paperclip',
    dataType: 'Text',
    length: 500,
    ui: {
      widget: 'Upload'
    }
  }
  // richtext: {  // 暂未实现
  //   label: '富文本',
  //   group: 'advanced',
  //   icon: 'bi-file-richtext',
  //   dataType: 'Text',
  //   length: 10000,
  //   ui: {
  //     widget: 'RichTextEditor'
  //   }
  // },
  // sort: {  // 暂未实现
  //   label: '排序',
  //   group: 'advanced',
  //   icon: 'bi-sort-numeric-down',
  //   dataType: 'Number',
  //   ui: {
  //     widget: 'InputNumber'
  //   }
  // }
};

/**
 * 字段类型分组
 */
export const fieldTypeGroups = [
  {
    name: 'basic',
    label: '基本类型',
    types: ['text', 'textarea', 'phone', 'email', 'url', 'integer', 'decimal', 'percent', 'password']
  },
  {
    name: 'choices',
    label: '选择类型',
    types: ['checkbox', 'select', 'multiselect', 'radio', 'checkboxgroup']
  },
  {
    name: 'datetime',
    label: '日期时间',
    types: ['date', 'time', 'datetime']
  },
  // {  // 暂未实现
  //   name: 'relation',
  //   label: '关系类型',
  //   types: ['belongsto', 'hasmany', 'manytomany']
  // },
  {
    name: 'advanced',
    label: '高级类型',
    types: ['autocode', 'attachment']  // json, formula, richtext, sort 暂未实现
  }
];

/**
 * 根据字段类型获取预设配置
 * @param {string} fieldType 字段类型
 * @returns {object|null} 预设配置
 */
export function getFieldPreset(fieldType) {
  return fieldTypePresets[fieldType] || null;
}

/**
 * 获取所有字段类型分组
 * @returns {array} 分组列表
 */
export function getFieldTypeGroups() {
  return fieldTypeGroups;
}

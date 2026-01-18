module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: [
    '@vue/eslint-config-typescript',
    '@vue/eslint-config-prettier',
    'eslint:recommended',
    'plugin:vue/vue3-recommended',
  ],
  parserOptions: {
    ecmaVersion: 2021,
    sourceType: 'module',
  },
  rules: {
    'no-console': 'warn',
    'no-debugger': 'warn',
    'vue/multi-word-component-names': 'off',
  },
};

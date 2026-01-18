import App from '@/App.vue';
import permission from '@/directive/permission';
import { i18n } from '@/lang/index';
import { store } from '@/pinia';
import { useUserStore } from '@/pinia/modules/user';
import MyRouter from '@/router';
import { setupYupLocale } from '@/utils/yup-config';
import { createApp } from 'vue';

// bootstrap
import 'bootstrap';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/css/bootstrap.min.css';

// local workflow
import { WorkflowVue } from 'workflow-vue';
// import 'workflow-vue/style.css';

// local styles
import '@/assets/css/responsive.css';
import '@/assets/css/sticky-table.css';


if (import.meta.env.VITE_MOCK) {
  import('./mock');
}

const app = createApp(App);
app.use(WorkflowVue);
app.use(MyRouter);
app.use(store);
const userStore = useUserStore();
userStore.restoreSession(); // must be called at app startup
app.use(i18n);

// init yup locale
setupYupLocale();
// permission directive
app.directive('permission', permission);

app.mount('#app');

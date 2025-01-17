import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';

import vueDebounce from 'vue-debounce';

import Buefy from 'buefy';
import 'buefy/dist/buefy.css';

import 'video.js/dist/video-js.css';
import 'videojs-vr/dist/videojs-vr.css';
// import '@fortawesome/fontawesome-free/css/all.css';
import '@fortawesome/fontawesome-free/js/all';
import '@mdi/font/css/materialdesignicons.css';

Vue.config.productionTip = false;
Vue.use(Buefy);
Vue.use(vueDebounce);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app');

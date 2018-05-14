// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import FastClick from 'fastclick'
import VueRouter from 'vue-router'
import App from './App'
import store from './store/store'
import router from './router'
import axios from 'axios'
import types from './store/types'
import VueAxios from 'vue-axios'

// Use vue-axios to bind axios to this.axios
Vue.use(VueAxios, axios)

// Defaults of axios
axios.defaults.timeout = 5000
axios.defaults.baseURL = '/api'

// Set response interceptor of axios
// Redirect to login page if get 401 error.
axios.interceptors.response.use(
    response => {
	console.log("res interceptor")
        return response
    },
    error => {
	console.log("res interceptor err: " + error.message)
        if (error.response) {
	    console.log("error.response:" + error.response)

            switch (error.response.status) {
                case 401:
	            store.commit(types.LOGOUT)
	            router.replace({
		        path: '/login',
			query: {redirect: router.currentRoute.fullPath}
		    })
	    }
	}
	return Promise.reject(error)
    }
)

FastClick.attach(document.body)

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app-box')

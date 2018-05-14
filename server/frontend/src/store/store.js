import Vuex from 'vuex'
import Vue from 'vue'
import * as types from './types'

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        user: null,
    },
    mutations: {
        [types.LOGIN]: (state, data) => {
            localStorage.user = data;
            state.user = data;
        },
        [types.LOGOUT]: (state) => {
            localStorage.removeItem('user');
            state.user = null
        }
    }
})

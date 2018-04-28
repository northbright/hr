import Vue from 'vue'
import VueRouter from 'vue-router'
import store from './store/store'
import * as types from './store/types'
import Home from './components/Home.vue'
import Login from './components/Login.vue'
import NotFound from './components/NotFound.vue'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: '/',
        meta: {
            requireAuth: true
        },
        component: Home
    },
    {
        path: '/login',
        name: 'login',
        component: Login
    },
    {
        path: '/404',
	name: '404',
	component: NotFound,
	hidden: true
    },
    {
        path: '*',
        hidden: true,
	redirect: { path: '/404' }
    }
];

if (window.localStorage.getItem('token')) {
    //store.commit(types.LOGIN, window.localStorage.getItem('token'))
}

const router = new VueRouter({
    routes
});

router.beforeEach((to, from, next) => {
    if (to.matched.some(r => r.meta.requireAuth)) {
	console.log("beforeEach()")
        if (store.state.user) {
            next();
        }
        else {
            next({
                path: '/login',
                query: {redirect: to.fullPath}
            })
        }
    }
    else {
        next();
    }
})

export default router;

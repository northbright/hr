<template>
  <div>
    <group title="登录">
      <x-input name="username" placeholder="用户名／手机）" v-model="username"></x-input>
      <x-input name="password" placeholder="请输入密码" type="password" v-model="password"></x-input>
      <x-button type="primary" v-on:click.native="login">OK</x-button>
    </group>
  </div>
</template>

<script>
import { XInput, Group, XButton, Cell } from 'vux'
import * as types from '../store/types'

export default {
  components: {
    XInput,
    XButton,
    Group,
    Cell
  },
  data () {
    return {
	username: "",
	password: "",
        csrf_token: ""
    }
  },
  methods: {
    mounted() {
       console.log("mounted().")
    },
    login() {
       console.log("login().") 

       this.axios.get('/csrf_token').then( (response) => {
         console.log("csrf token: " + response.data.csrf_token)
         this.csrf_token = response.data.csrf_token

         this.axios.post('/login', {
           username: this.username,
           password: this.password,
           csrf_token: this.csrf_token
         })
         .then( (response) => {
           console.log("POST /login success. data: " + response.data)
           this.$store.commit(types.LOGIN, response.data.id)
           this.$router.push({
             path: '/',
           })
         })
         .catch( (error) =>  {
           console.log("POST /login error: " + error)
         })
       })
       .catch( (error) => {
           console.log("GET /login error:" + error)
       })
    }
  }
}
</script>

<style scoped>
.red {
  color: red;
}
.green {
  color: green;
}
</style>


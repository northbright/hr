<template>
  <div>
    <msg :title="title" :description="description" :buttons="buttons" :icon="icon"></msg>
  </div>
</template>

<script>
import { Msg, Divider, XButton } from 'vux'
import * as types from '../store/types'

export default {
  components: {
    Msg,
    Divider,
    XButton,
  },
  data () {
    return {
	title: 'Page not found',
	description: 'Current page is not found.',
        icon: 'warn',
        buttons: [{
          type: 'primary',
          text: 'Go to Home',
          link: '/'
        }]
    }
  },
  methods: {
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


import Vue from 'vue'
import Router from 'vue-router'
import store from './store'

import BackendError from '@/views/BackendError.vue'
import NotImplemented from '@/views/NotImplemented.vue'

import CameraList from '@/views/CameraList.vue'


Vue.use(Router)

var router = new Router({
  routes: [
    {
      path: '/',
      name: 'Cameras',
      component: CameraList,
    },
    {
      path: '/404',
      name: 'Not Found',
      component: NotImplemented,
    },
    {
      path: '/backend_error',
      name: 'Backend Error',
      component: BackendError,
    },
  ]
})

router.beforeEach((to, from, next) => {
  if (store.getters['config/isReady']) {
    next()
  } else {
    store.dispatch('config/getConfig').then(() => {
      next()
    }).catch(function(){
      router.push('/backend_error')
      next()
    })
  }
})


export default router
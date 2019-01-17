import Vue from 'vue'
import Vuex from 'vuex'

import config from './modules/config'
import cameras from './modules/cameras'

Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
  plugins: [ ],
  modules: {
    config,
    cameras,
  },
  strict: debug
})
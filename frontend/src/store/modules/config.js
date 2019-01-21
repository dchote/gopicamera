import Vue from 'vue'

var configURL = location.protocol + '//' + location.hostname + ':' + location.port + '/config'

// initial state
const state = {
  status: 'pending',
  APIURL: '',
  cameraURL: ''
}

// getters
const getters = {
  isReady: state => state.status == 'ready' || state.status == 'error'
}

// actions
const actions = {
  getConfig({commit}) {
    return new Promise((resolve, reject) => {
      Vue.axios.get(configURL)
      .then(response => {
        // set baseURL so we can just use relative url paths in the rest of the code
        Vue.axios.defaults.baseURL =  response.data.APIURL
        
        commit('load_config', response)
        
        resolve(response)
      }).catch(function (error) {
        commit('fetch_error')
        reject(error)
      })
    })
  }
}

// mutations
const mutations = {
  fetch_error(state) {
    state.status = 'error'
  },
  load_config(state, config) {
    state.status = 'ready'
    state.APIURL = config.data.APIURL
    state.cameraURL = config.data.cameraURL
    
    // eslint-disable-next-line
    console.log('config:', state)
  }
  
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
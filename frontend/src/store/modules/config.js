import Vue from 'vue'

const configURL = 'http://' + location.hostname + ':8088/config'


// initial state
const state = {
  status: 'pending',
  APIURL: '',
  cameraURL: ''
}

// getters
const getters = {
  isReady: state => state.status == 'ready'
}

// actions
const actions = {
  getConfig({commit}) {
    return new Promise((resolve, reject) => {
      Vue.axios.get(configURL)
      .then(response => {
        // set baseURL so we can just use relative url paths in the rest of the code
        Vue.axios.defaults.baseURL = response.data.APIURL
        
        commit('setConfig', response)
        
        resolve(response)
      }).catch(function (error) {
        reject(error)
      })
    })
  }
}

// mutations
const mutations = {
  setConfig(state, config) {
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
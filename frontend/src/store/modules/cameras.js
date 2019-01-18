import Vue from 'vue'

// initial state
const state = {
  status: '',
  cameras: []
}

// getters
const getters = {

}

// actions
const actions = {
  fetchCameras({commit}) {
    return new Promise((resolve, reject) => {
      commit('fetch_request')
      
      Vue.axios.get('/camera/list')
      .then(response => {
        commit('load_state', response.data)
        
        resolve(response)
      }).catch(function (error) {
        commit('fetch_error')
        
        reject(error)
      })
    })
  },
  
}

// mutations
const mutations = {
  fetch_request(state) {
    state.status = 'loading'
  },
  fetch_error(state) {
    state.status = 'error'
  },
  load_state(state, payload) {
    state.cameras = payload.cameras
    state.status = 'loaded'
    
    // eslint-disable-next-line
    console.log('cameras:', state)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
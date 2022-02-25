import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    id: '',
    name: ''
  },
  mutations: {
    getID (state, payload) {
      state.id = payload
    },
    getName (state, payload) {
      state.name = payload
    }
  },
  actions: {
  },
  modules: {
  }
})

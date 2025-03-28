import Vue from 'vue'
import Vuex from 'vuex'
import user from './modules/user'
import app from './modules/app'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    user,
    app
  },
  getters: {
    token: state => state.user.token,
    name: state => state.user.name,
    avatar: state => state.user.avatar,
    roles: state => state.user.roles,
    permissions: state => state.user.permissions
  }
}) 
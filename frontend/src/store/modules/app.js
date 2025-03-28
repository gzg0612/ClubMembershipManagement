const state = {
  sidebar: {
    opened: localStorage.getItem('sidebarOpened') ? !!+localStorage.getItem('sidebarOpened') : true
  }
}

const mutations = {
  TOGGLE_SIDEBAR: (state) => {
    state.sidebar.opened = !state.sidebar.opened
    localStorage.setItem('sidebarOpened', state.sidebar.opened ? '1' : '0')
  }
}

const actions = {
  toggleSidebar({ commit }) {
    commit('TOGGLE_SIDEBAR')
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
} 
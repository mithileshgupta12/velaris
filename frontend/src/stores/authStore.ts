import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export interface ILoggedInUser {
  id: number
  name: string
  email: string
  createdAt: string
  updatedAt: string
}

export const useAuthStore = defineStore('auth', () => {
  const initialized = ref<boolean>(false)
  const loggedInUser = ref<ILoggedInUser | null>(null)

  const isLoggedIn = computed(() => !!loggedInUser.value)

  const setLoggedInUser = (payload: ILoggedInUser) => {
    loggedInUser.value = payload
  }

  const checkAuth = async () => {
    if (!initialized.value) {
      console.log('First try')
      initialized.value = true
      return
    }

    console.log('Not first try')
  }

  return {
    loggedInUser,
    isLoggedIn,
    setLoggedInUser,
    checkAuth,
  }
})

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

  const setLoggedInUser = (payload: ILoggedInUser | null) => {
    loggedInUser.value = payload
  }

  const setInitialized = (payload: boolean) => {
    initialized.value = payload
  }

  return {
    loggedInUser,
    isLoggedIn,
    initialized,
    setLoggedInUser,
    setInitialized,
  }
})

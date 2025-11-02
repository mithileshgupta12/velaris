import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

interface IUser {
  id: number
  name: string
  email: string
}

export const useAuthStore = defineStore('auth', () => {
  const loggedInUser = ref<IUser | null>(null)

  const isLoggedIn = computed(() => !!loggedInUser.value)

  const setLoggedInUser = (payload: IUser) => {
    loggedInUser.value = payload
  }

  return {
    loggedInUser,
    isLoggedIn,
    setLoggedInUser,
  }
})

import axios from '@/utils/axios'
import { AxiosError } from 'axios'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

const useAuth = () => {
  const router = useRouter()

  const loading = ref(false)
  const error = ref<string | null>(null)

  const isLoading = computed(() => loading.value)
  const isError = computed(() => !!error.value)
  const errorMessage = computed(() => error.value)

  const login = async (data: { email: string; password: string }) => {
    loading.value = true
    error.value = null

    try {
      await axios.post('/auth/login', data)

      router.push('/dashboard')
    } catch (e) {
      if (e instanceof AxiosError) {
        error.value = e.response?.data.error.message || 'Internal server error'
      } else {
        error.value = 'Internal server error'
      }
    } finally {
      loading.value = false
    }
  }

  const register = async (data: {
    name: string
    email: string
    password: string
    passwordConfirmation: string
  }) => {
    loading.value = true
    error.value = null

    try {
      await axios.post('/auth/register', data)

      router.push('/auth/login')
    } catch (e) {
      if (e instanceof AxiosError) {
        error.value = e.response?.data.error.message || 'Internal server error'
      } else {
        error.value = 'Internal server error'
      }
    } finally {
      loading.value = false
    }
  }

  return {
    login,
    register,
    isLoading,
    isError,
    errorMessage,
  }
}

export default useAuth

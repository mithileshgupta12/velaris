import { useAuthStore } from '@/stores/authStore'
import axios from '@/utils/axios'
import { AxiosError } from 'axios'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

interface ISuccessResponse {
  success: boolean
  data: {
    id: number
    name: string
    email: string
    created_at: string
    updated_at: string
  }
}

const useAuth = () => {
  const router = useRouter()

  const { setLoggedInUser } = useAuthStore()

  const loading = ref(false)
  const error = ref<string | null>(null)

  const isLoading = computed(() => loading.value)
  const isError = computed(() => !!error.value)
  const errorMessage = computed(() => error.value)

  const login = async (data: { email: string; password: string }) => {
    loading.value = true
    error.value = null

    try {
      const response = await axios.post<ISuccessResponse>('/auth/login', data)

      setLoggedInUser({
        id: response.data.data.id,
        name: response.data.data.name,
        email: response.data.data.email,
        createdAt: response.data.data.created_at,
        updatedAt: response.data.data.updated_at,
      })

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

  const register = async ({
    name,
    email,
    password,
    passwordConfirmation,
  }: {
    name: string
    email: string
    password: string
    passwordConfirmation: string
  }) => {
    loading.value = true
    error.value = null

    try {
      await axios.post('/auth/register', {
        name,
        email,
        password,
        password_confirmation: passwordConfirmation,
      })

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

  const getLoggedInUser = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await axios.get<ISuccessResponse>('/auth/user')

      console.log(response.data)
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
    getLoggedInUser,
    isLoading,
    isError,
    errorMessage,
  }
}

export default useAuth

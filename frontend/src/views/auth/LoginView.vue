<template>
  <h1>Login</h1>

  <p v-if="isError">{{ error }}</p>

  <label for="email">Email</label>
  <input type="email" name="email" id="email" v-model="form.email" />
  <label for="password">Password</label>
  <input type="password" name="password" id="password" v-model="form.password" />
  <input
    type="submit"
    :value="loading ? 'Loading...' : 'Login'"
    :disabled="loading"
    @click="handleLogin"
  />
</template>

<script setup lang="ts">
import useAuth from '@/composables/useAuth'
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

const { login, loading, isError, error } = useAuth()
const router = useRouter()

const form = reactive({
  email: '',
  password: '',
})

const handleLogin = async () => {
  const response = await login({
    email: form.email,
    password: form.password,
  })

  if (response) {
    router.push({ name: 'dashboard' })
  }
}
</script>

<template>
  <h1>Login</h1>

  <p v-if="isError">{{ errorMessage }}</p>

  <label for="email">Email</label>
  <input type="email" name="email" id="email" v-model="form.email" />
  <label for="password">Password</label>
  <input type="password" name="password" id="password" v-model="form.password" />
  <input
    type="submit"
    :value="isLoading ? 'Loading...' : 'Login'"
    :disabled="isLoading"
    @click="handleLogin"
  />
</template>

<script setup lang="ts">
import useAuth from '@/composables/useAuth'
import { reactive } from 'vue'

const { login, isLoading, isError, errorMessage } = useAuth()

const form = reactive({
  email: '',
  password: '',
})

const handleLogin = async () => {
  await login({
    email: form.email,
    password: form.password,
  })
}
</script>

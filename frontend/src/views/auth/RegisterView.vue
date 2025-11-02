<template>
  <h1>Register</h1>

  <p v-if="isError">{{ errorMessage }}</p>

  <form>
    <label for="name">Name</label>
    <input type="text" name="name" id="name" v-model="form.name" />
    <label for="email">Email</label>
    <input type="email" name="email" id="email" v-model="form.email" />
    <label for="password">Password</label>
    <input type="password" name="password" id="password" v-model="form.password" />
    <label for="password_confirmation">Confirm Password</label>
    <input
      type="password"
      name="password_confirmation"
      id="password_confirmation"
      v-model="form.passwordConfirmation"
    />
    <input
      type="submit"
      :value="isLoading ? 'Loading...' : 'Register'"
      :disabled="isLoading"
      @click="handleRegister"
    />
  </form>
</template>

<script setup lang="ts">
import useAuth from '@/composables/useAuth'
import { reactive } from 'vue'

const { register, isLoading, isError, errorMessage } = useAuth()

const form = reactive({
  name: '',
  email: '',
  password: '',
  passwordConfirmation: '',
})

const handleRegister = async () => {
  await register({
    name: form.name,
    email: form.email,
    password: form.password,
    passwordConfirmation: form.passwordConfirmation,
  })
}
</script>

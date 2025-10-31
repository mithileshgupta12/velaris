<template>
  <h1>Login</h1>
  <input type="text" name="email" id="email" v-model="form.email" />
  <input type="password" name="password" id="password" v-model="form.password" />
  <input type="submit" value="Login" @click="submit" />

  <pre>{{
    JSON.stringify(
      {
        email: form.email,
        password: form.password,
      },
      null,
      4,
    )
  }}</pre>
</template>

<script setup lang="ts">
import { reactive } from 'vue'

const form = reactive({
  email: '',
  password: '',
})

const submit = () => {
  console.log('Logged in')
}
</script>

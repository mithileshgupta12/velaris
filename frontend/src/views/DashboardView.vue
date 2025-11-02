<template>
  <h1>Dashboard</h1>
  <pre>{{ JSON.stringify(loggedInUser, null, 4) }}</pre>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/authStore'

const { loggedInUser } = useAuthStore()
</script>

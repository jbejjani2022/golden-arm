<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  onMount(async () => {
    const response = await fetch('/api/admin/validate-session', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include'  // Ensures cookies are sent
    });

    const result = await response.json();

    if (!result.valid) {
      goto('/admin');
    }
  });
</script>
  
<main>
  <div class="admin-page">
    <slot></slot>
  </div>
</main>

<style>
  /* Override the global dark theme for admin pages */
  :global(body) {
    background-color: white;
    color: black;
    font-family: Arial, sans-serif;
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
  }

  :global(a) {
      color: #3498db;
      text-decoration: none;
  }

  :global(a:hover) {
      text-decoration: underline;
  }
</style>

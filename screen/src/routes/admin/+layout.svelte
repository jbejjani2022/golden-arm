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
  :global(a) {
      color: #3498db;
      text-decoration: none;
  }

  :global(a:hover) {
      text-decoration: underline;
  }
</style>

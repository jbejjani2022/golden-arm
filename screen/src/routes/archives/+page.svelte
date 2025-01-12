<script lang="ts">
  import { formatDate } from '$lib';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let archive: Array<any> = [];
  let error: string = '';

  // Fetch the movie archive using /api/movie/archive
  onMount(async () => {
    try {
      const response = await fetch('/api/movie/archive');
      const data = await response.json();

      if (data.success) {
        archive = data.data;
      } else {
        error = 'Failed to load the movie archive.';
      }
    } catch (err) {
      console.error(err);
      error = 'Something went wrong while fetching the movie archive.';
    }
  });
</script>

<main class="archive">
  <h1>The Golden Arm Archive</h1>

  {#if error}
    <p class="error">{error}</p>
  {/if}
</main>

<style>

</style>

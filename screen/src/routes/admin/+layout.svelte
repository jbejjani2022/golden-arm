<script lang="ts">
    // Add any global logic or imports here
</script>
  
<main>
  <div class="admin-page">
    <slot></slot>
  </div>
</main>

<style>
  /* Apply a white background to all admin pages */
  :global(body) {
    background-color: white;
    color: black;
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

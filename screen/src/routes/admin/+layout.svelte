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
  }

  .admin-page {
    padding: 20px;
  }
</style>

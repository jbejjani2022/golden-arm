<script lang="ts">
  // Add any global logic or imports here
  import { page } from '$app/state';

  // Create a derived store to check if the current path starts with '/admin'
  const isAdmin = page.url.pathname.startsWith('/admin');
</script>

<main>
  <slot></slot> <!-- This will render the content of the page -->

   <!-- Only show footer if it's not the admin page -->
  {#if !isAdmin}
    <footer class="global-footer">
      <p>Visit <a href="https://www.instagram.com/eliotgoldenarm/" target="_blank" rel="noopener noreferrer">@eliotgoldenarm</a> on Instagram.</p>
    </footer>
  {/if}
</main>

<style>
  /* Dark theme styling */
  :global(body) {
    margin: 0;
    font-family: Arial, sans-serif;
    background-color: #202020;
    color: #f0f0f0;
  }

  /* Global Footer Styling */
  .global-footer {
    position: fixed;
    bottom: 0;
    width: 100%;
    padding: 10px;
    text-align: center;
    font-size: 14px;
  }

  .global-footer a {
    color: #e50914;
    text-decoration: none;
    font-weight: bold;
  }

  .global-footer a:hover {
    text-decoration: underline;
  }

  /* Global Button Styling */
  :global(button) {
    padding: 10px 15px;
    background-color: #e50914; /* Button background color */
    color: white;
    border: none;
    border-radius: 5px;
    font-size: 16px;
    cursor: pointer;
    transition: background-color 0.3s ease;
  }

  :global(button:hover) {
    background-color: #b40710; /* Darker red on hover */
  }

  :global(button:disabled) {
    background-color: #666;
    cursor: not-allowed;
  }
</style>

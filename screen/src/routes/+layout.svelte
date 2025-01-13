<script lang="ts">
  // Add any global logic or imports here
  import { page } from '$app/state';

  // Create a derived store to check if the current path starts with '/admin'
  const isAdmin = page.url.pathname.startsWith('/admin');
</script>

<main>
  <div class="page-container">
    <slot></slot>
    <!-- Show footer if it's not the admin page -->
    {#if !isAdmin}
      <footer class="global-footer">
        <p>Visit <a href="https://www.instagram.com/eliotgoldenarm/" target="_blank" rel="noopener noreferrer">@eliotgoldenarm</a> on Instagram.</p>
      </footer>
    {/if}
  </div>
</main>

<style>
  :root {
    --gold: #d3ab0c;
    --dark-gold: #a28200;
    --dark: #202020;
  }

  /* Dark theme styling */
  :global(body) {
    margin: 0 auto;
    font-family: Arial, sans-serif;
    background-color: var(--dark);
    color: #f0f0f0;
    max-width: 800px;
    padding: 20px;
    text-align: center;
  }

  :global(.theater-info h1) {
    font-size: 36px;
    margin-bottom: 10px;
  }
  
  :global(.theater-info p) {
    font-size: 18px;
    margin-bottom: 30px;
    color: #bbb;
  }

  /* Make the page container fill the viewport */
  .page-container {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }

  /* Ensure the main content takes up all available space */
  .page-container > :global(main) {
    flex: 1;
  }

  /* Footer styling */
  .global-footer {
    background-color: var(--dark);
    color: #fff;
    text-align: center;
    padding: 1rem;
  }

  .global-footer a {
    color: var(--gold);
    text-decoration: none;
  }

  .global-footer a:hover {
    text-decoration: underline;
  }

  :global(.link) {
    color: var(--gold);
    text-decoration: none;
  }

  :global(.link:hover) {
    text-decoration: underline;
  }

  /* Global Button Styling */
  :global(button) {
    padding: 10px 15px;
    background-color: var(--gold);
    color: white;
    border: none;
    border-radius: 5px;
    font-size: 16px;
    cursor: pointer;
    transition: background-color 0.3s ease;
  }

  :global(button:hover) {
    background-color: var(--dark-gold);
  }

  :global(button:disabled) {
    background-color: #666;
    cursor: not-allowed;
  }

  :global(.modal) {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5); /* Semi-transparent background */
    display: flex;
    justify-content: center;  /* Horizontally center */
    align-items: center;      /* Vertically center */
    z-index: 1000;
}

  :global(.modal-content) {
      background-color: #fff;
      color: black;
      padding: 20px;
      border-radius: 8px;
      width: 400px;
      text-align: center; /* Center text */
      display: flex;
      flex-direction: column;  /* Arrange elements vertically */
      align-items: center;     /* Center items horizontally */
  }

  :global(.modal-content input) {
      margin: 5px 0; /* Space between input fields */
      padding: 8px;
      width: 100%;
      max-width: 250px; /* Limit input width */
  }

  :global(.modal-content button) {
      padding: 10px 15px;
      color: white;
      border: none;
      cursor: pointer;
      border-radius: 5px;
      width: 100%;
      max-width: 250px; /* Same as input field width */
      margin-top: 10px; /* Space between inputs and button */
  }

  :global(.cancel-button) {
    background-color: darkgrey;
  }

  :global(.cancel-button:hover) {
    background-color: grey;
  }

  :global(.form-group) {
    margin-bottom: 15px;
    display: flex;
    align-items: center;
  }

  :global(.form-group label) {
    margin-right: 10px;
    width: auto;
  }

  :global(.form-group input) {
    flex-grow: 1;
    padding: 5px;
    font-size: 14px;
  }
</style>

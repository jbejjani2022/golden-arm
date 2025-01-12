<script lang="ts">
    // Add any global logic or imports here
  </script>
  
  <main>
    <slot></slot> <!-- This will render the content of the page -->
  </main>
  
  <style>
    /* Dark theme styling */
    :global(body) {
      margin: 0;
      font-family: Arial, sans-serif;
      background-color: white;
      color: black;
    }
  </style>
  
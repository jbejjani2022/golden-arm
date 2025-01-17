<script lang="ts">
  // Add any global logic or imports here
  import { page } from '$app/state';
  import { onMount } from 'svelte';

  // Check if the current path starts with '/admin'
  const isAdmin = page.url.pathname.startsWith('/admin');

  // Navbar mobile
  let showMobileMenu = false;

  // Data for next movie
  let movie: any = null;
  let error: string = '';

  // Fetch the next movie using the /api/movie/next endpoint
  onMount(async () => {
    try {
      const response = await fetch('/api/movie/next');
      const data = await response.json();

      if (data.success) {
        movie = data.data;
      } else {
        error = 'Failed to load the next movie.';
      }
    } catch (err) {
      console.error(err);
      error = 'Something went wrong while fetching the movie data.';
    }
  });

  // Toggles the mobile menu
  function toggleMenu() {
    showMobileMenu = !showMobileMenu;
  }

  // Ensures the menu can be toggled via keyboard (Enter or Space)
  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Enter' || event.key === ' ') {
      toggleMenu();
    }
  }

</script>

<!-- Navbar -->
<nav class="navbar">
  <!-- Logo -->
  <a href="/" class="navbar-logo">
    <img src="/standinlogo.png" alt="Logo" class="logo" />
  </a>

  <!-- Navigation Links -->
  <ul class="navbar-links">
    <li>
      <a href="/archives" class:active={page.url.pathname === '/archives'}>Past Screenings</a>
    </li>
    <li>
      {#if movie?.ID}
      <a href={`/reservations/${movie.ID}`} class:active={page.url.pathname === `/reservations/${movie.ID}`}>Reserve a Seat</a>
    {/if}    </li>
    <li>
      <a href="/merch" class:active={page.url.pathname === '/merch'}>Merch</a>
    </li>
    <li>
      <a href="/filmfest" class:active={page.url.pathname === '/filmfest'}>Film Festival</a>
    </li>
  </ul>

    <!-- Hamburger Menu Button -->
    <button
    class="hamburger"
    on:click={toggleMenu}
    on:keydown={handleKeyDown}
    aria-label="Toggle navigation menu"
    aria-expanded={showMobileMenu ? 'true' : 'false'}
  >
    <span></span>
    <span></span>
    <span></span>
  </button>

 <!-- Mobile Menu -->
 {#if showMobileMenu}
 <div class="mobile-menu">
   <a href="/archives" on:click={() => (showMobileMenu = false)}>Past Screenings</a>
   <a href={`/reservations/${movie.ID}`} on:click={() => (showMobileMenu = false)}>Reserve a Seat</a>
   <a href="/merch" on:click={() => (showMobileMenu = false)}>Merch</a>
   <a href="/filmfest" on:click={() => (showMobileMenu = false)}>Film Festival</a>
 </div>
{/if}
</nav>

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
    max-width: 80%;
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

  main {
  margin-top: 70px; /* Matches the navbar height */
  padding: 1rem;
  box-sizing: border-box;
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

  /* Navbar Styles */
  .navbar {
    display: flex;
    align-items: center;
    /* justify-content: space-between; */
    padding: 0 1rem; /* Add some padding inside the navbar */
    background-color: #202020; 
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    position: fixed;
    top: 0;
    left: 0;
    width: 100%; /* Spans the full width */
    max-width: 100%; /* Ensures it doesn’t overflow the screen */
    z-index: 10;
    margin: 0;
    box-sizing: border-box; /* Includes padding and border in the element's total width/height */
    /* height: 5%; */ 
  }

 

  .navbar-logo img {
    height: 40px;
    padding: 10px 20px;
  }

  .navbar-links {
    list-style: none;
    display: flex;
    gap: 2rem;
    margin: 0;
    padding: 5px;
  }

  .navbar-links li {
    font-size: 1rem;
  }

  .navbar-links a {
    text-decoration: none;
    color: #ffffff;
    font-weight: 500;
    transition: color 0.3s ease;
  }

  .navbar-links a:hover {
    color: #f4c523; /* Highlight color on hover */
  }

  .navbar-links a.active {
    color: #f4c523; /* Active link color */
    /* font-weight: 600; */
  }

  /* for mobile */

  .hamburger {
    display: none;
    flex-direction: column;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    right: 0;
  }

  .hamburger span {
    background-color: white;
    height: 3px;
    margin: 3px 0;
    width: 25px;
  }

  .mobile-menu {
    display: none; /* Hidden by default */
    flex-direction: column;
    position: absolute;
    top: 100%;
    right: 0;
    left: 0;
    background-color: #202020;
    z-index: 10;
  }

  .mobile-menu a {
    color: #f4c523;
    text-decoration: none;
    padding: 1rem;
    text-align: center;
    border-bottom: 1px solid white;
  }

  @media screen and (max-width: 768px) {
    .navbar {
      justify-content: space-between; /* Space between logo and hamburger on mobile */
    }
    
    .navbar-links {
      display: none; /* Hide links for smaller screens */
    }

    .hamburger {
      display: flex;
      justify-content: flex-end; /* Align hamburger to the right side */

    }

    .mobile-menu {
      display: flex;
    }
  }
</style>

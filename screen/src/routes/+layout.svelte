<script lang="ts">
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

  let showModal = false;

  const confirmComment = () => {
    showModal = true;
  }

  const cancelComment = () => {
    showModal = false;
  }

  let name = '';
  let email = '';
  let comment = '';

  const handleComment = async () => {
    if (!name || !email || !comment) {
      alert("Please fill out the suggestion form.");
      return;
    }
    // Send suggestion to server
    try {
      const response = await fetch(`/api/comment`, {
        method: "POST",
        headers: { 
          "Content-Type": "application/json" 
        },
        body: JSON.stringify({
          name,
          email,
          comment,
        })
      });

      const result = await response.json();
      showModal = false;
      if (result.success) {
        alert("Thank you for your suggestion!");
      } else {
        alert("Failed to submit comment.");
      }
    } catch (err) {
      console.error(err);
      alert("Something went wrong while submitting the comment.");
    }
  };

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

  <!-- Instagram Icon -->
  <a 
    href="https://www.instagram.com/eliotgoldenarm/" 
    target="_blank" 
    rel="noopener noreferrer"
    class="instagram-link"
  >
    <img 
      src="/instagram.svg" 
      alt="Instagram" 
      class="instagram-icon"
    />
  </a>

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
 <div class="mobile-menu {showMobileMenu ? 'open' : ''}">
   <a href="/archives" on:click={() => (showMobileMenu = false)}>Past Screenings</a>
   <a href={`/reservations/${movie.ID}`} on:click={() => (showMobileMenu = false)}>Reserve a Seat</a>
   <a href="/merch" on:click={() => (showMobileMenu = false)}>Merch</a>
   <a href="/filmfest" on:click={() => (showMobileMenu = false)}>Film Festival</a>
 </div>
{/if}
</nav>

<div class="content-wrapper">
<main>
  <div class="page-container">
    {#if showModal}
    <div class="modal">
      <div class="modal-content">
          <div class="form-group">
            <label for="name">Name: </label>
            <input type="text" id="name" bind:value={name} placeholder="Enter your name" required />
          </div>
          <div class="form-group">
            <label for="email">Email: </label>
            <input type="email" id="email" bind:value={email} placeholder="Enter your email" required />
          </div>
          <div class="form-group">
            <label for="comment">What should we screen next?</label>
            <input type="text" id="comment" bind:value={comment} placeholder="Enter any movie suggestions!" required />
          </div>
          <button type="submit" on:click={handleComment}>Send</button>
          <button type="button" class="cancel-button" on:click={cancelComment}>Cancel</button>
      </div>
    </div>
    {/if}

    <slot></slot>
  </div>
</main>
</div>

 <!-- Show footer if it's not the admin page -->
 {#if !isAdmin}
 <footer class="global-footer">
   <div class="footer-content">
     <div class="footer-section">
       <a href="/" class="footer-logo-link">
         <img src="/standinlogo.png" alt="Logo" class="footer-logo-img" />
       </a>
       <div class="suggestions-section">
         <p>Have suggestions? We'd love to hear from you!</p>
         <button class="suggestions-button" on:click={confirmComment}>Comment</button>
       </div>
       <div class="social-links">
         <p>Stay updated on how The Golden Arm is shaping the cinema landscape of the Boston area.</p>
         <p>Follow us <a href="https://www.instagram.com/eliotgoldenarm/" target="_blank" rel="noopener noreferrer">@eliotgoldenarm</a>.</p>
       </div>
     </div>
   </div>
 </footer>
{/if}

<style>
  :root {
    --gold: #edbc0d;
    --dark-gold: #b08d00;
    --dark: #1f1f1f;
  }

  /* Dark theme styling */
  :global(body) {
    margin: 0 auto;
    padding: 0;
    font-family: Arial, sans-serif;
    background-color: var(--dark);
    color: #f0f0f0;
    min-height: 100vh;
    text-align: center;
  }
  
  .content-wrapper {
    width: 80%;
    margin: 0 auto;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
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

  :global(.link) {
    color: var(--gold);
    text-decoration: none;
  }

  :global(.link:hover) {
    text-decoration: underline;
  }

  /* Make the page container fill the viewport */
  .page-container {
    display: flex;
    flex-direction: column;
    flex: 1;
  }

  main {
    display: flex;
    flex-direction: column;
    flex: 1;
    margin-top: 10%; /* Matches the navbar height */
    box-sizing: border-box;
    width: 100%;
  }

  /* Footer styling */
  .global-footer {
    background-color: #0e0e0e;
    color: #fff;
    padding: 3rem 1rem;
    box-shadow: 0 -4px 6px rgba(0, 0, 0, 0.1);
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    box-sizing: border-box; /* Includes padding in width/height */
    margin-top: 2rem;
  }

  .footer-content {
    width: 80%;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 2rem; /* Spacing between footer sections */
  }

  .footer-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem; /* Spacing between elements in a section */
  }

  .footer-logo-link {
    display: block;
  }

  .footer-logo-img {
    height: 60px; /* Adjust logo size */
    width: auto;
  }

  .social-links {
    font-size: 0.9rem;
  }

  .social-links a {
    color: var(--gold);
    text-decoration: none;
    transition: color 0.3s ease;
  }

  .social-links a:hover {
    text-decoration: underline;
  }

  .suggestions-section {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .suggestions-button {
    padding: 8px 16px;
    font-size: 14px;
    border-radius: 5px;
  }

  /* Mobile styles */
  @media screen and (max-width: 768px) {
    .footer-content {
      padding: 0 1rem;
    }

    .footer-logo-img {
      height: 50px; /* Slightly smaller logo for mobile */
    }

    .social-links {
      font-size: 0.8rem;
    }
  }

  /* Global Button Styling */
  :global(button) {
    padding: 12px 24px;
    background-color: var(--gold);
    color: black;
    font-weight: bold;
    border: none;
    border-radius: 18px;
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
    background-color: #0e0e0e; 
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    position: fixed;
    top: 0;
    left: 0;
    width: 100%; /* Spans the full width */
    max-width: 100%; /* Ensures it doesn’t overflow the screen */
    z-index: 10;
    box-sizing: border-box; /* Includes padding and border in the element's total width/height */
    height: 10%;
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
    color: var(--gold);
  }

  .navbar-links a.active {
    color: var(--gold);
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

  .mobile-menu.open {
    display: flex;
  }

  .mobile-menu a {
    color: var(--gold);
    text-decoration: none;
    padding: 1rem;
    text-align: center;
    border-bottom: 1px solid white;
  }

  .instagram-link {
    margin-left: auto;
    padding: 0.5rem;
    display: flex;
    align-items: center;
    transition: filter 0.3s ease;
  }

  .instagram-icon {
    width: 24px;
    height: 24px;
    filter: invert(1);
  }

  .instagram-link:hover .instagram-icon {
    filter: invert(71%) sepia(97%) saturate(1061%) hue-rotate(359deg) brightness(103%) contrast(98%);
  }

  /* Update mobile styles */
  @media screen and (max-width: 768px) {
    .navbar {
      justify-content: space-between;
    }
    
    .instagram-link {
      margin-left: 0; /* Reset margin for mobile */
      order: -1; /* Places Instagram icon before hamburger menu */
    }

    .navbar-links {
      display: none;
    }

    .hamburger {
      display: flex;
    }

    .mobile-menu {
      display: flex;
    }
  }
</style>

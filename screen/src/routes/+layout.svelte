<script lang="ts">
  import { page } from '$app/state';
  import { onMount } from 'svelte';
  import { apiBaseUrl } from '$lib/api';

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
      const response = await fetch(`${apiBaseUrl}/movie/next`);
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
      const response = await fetch(`${apiBaseUrl}/comment`, {
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

<div class="layout">
  <!-- Navbar -->
<nav class="navbar">
  <!-- Logo -->
  <a href="/" class="navbar-logo">
    <img src="/logofr.png" alt="Logo" class="logo" />
  </a>

  <!-- Navigation Links -->
  <ul class="navbar-links">
    <li>
      <a href="/about" class:active={page.url.pathname === '/about'}>About</a>
    </li>
    <li>
      {#if movie?.ID}
      <a href={`/reservations/${movie.ID}`} class:active={page.url.pathname === `/reservations/${movie.ID}`}>Reserve a Seat</a>
    {/if}    
    </li>
    <li>
      <a href="/archives" class:active={page.url.pathname === '/archives'}>Past Screenings</a>
    </li>
    <!-- <li>
      <a href="/filmfest" class:active={page.url.pathname === '/filmfest'}>Film Festival</a>
    </li> -->
    <li>
      <a href="/merch" class:active={page.url.pathname === '/merch'}>Merch</a>
    </li>
    <li>
      <a href="https://roombook.fas.harvard.edu/EmsWebApp/RoomRequest.aspx?data=ity3Dem%2byxxGFZTQvNr97zkAqLedHLx6" target="_blank">Book the Theater</a>
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
   <a href="/about" on:click={() => (showMobileMenu = false)}>About</a>
   <a href={`/reservations/${movie.ID}`} on:click={() => (showMobileMenu = false)}>Reserve a Seat</a>
   <a href="/archives" on:click={() => (showMobileMenu = false)}>Past Screenings</a>
   <!-- <a href="/filmfest" on:click={() => (showMobileMenu = false)}>Film Festival</a> -->
   <a href="/merch" on:click={() => (showMobileMenu = false)}>Merch</a>
   <a href="https://roombook.fas.harvard.edu/EmsWebApp/RoomRequest.aspx?data=ity3Dem%2byxxGFZTQvNr97zkAqLedHLx6" on:click={() => (showMobileMenu = false)}>Book the Theater</a>
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
         <img src="/logofr.png" alt="Logo" class="footer-logo-img" />
       </a>
       <div class="footer-info">
         <p class="footer-elt">Have suggestions? We'd love to hear from you!</p>
          <p class="footer-elt">Email us at <a href="mailto:goldenarmtheater@gmail.com">goldenarmtheater@gmail.com</a> or drop a comment below.</p>
         <button class="suggestions-button" on:click={confirmComment}>Comment</button>
       </div>
       <div class="footer-info">
         <p>Stay updated on how The Golden Arm is shaping the cinema landscape of the Boston area <a href="https://www.instagram.com/eliotgoldenarm/" target="_blank" rel="noopener noreferrer">@eliotgoldenarm</a>.</p>
       </div>
     </div>
   </div>
 </footer>
{/if}

</div>

<style>
  .layout{
    display: flex;
  flex-direction: column;
  }

  :root {
    --gold: #edbc0d;
    --dark-gold: #b08d00;
    --dark: #1f1f1f;
  }

  @font-face {
    font-family: Telegraf;
    src: url('/fonts/Telegraf/PPTelegraf-Regular.otf') format('opentype');
    font-weight: normal;
    font-style: normal;
  }

  @font-face {
    font-family: Telegraf-Ultrabold;
    src: url('/fonts/Telegraf/PPTelegraf-Ultrabold.otf') format('opentype');
    font-weight: bold;
    font-style: normal;
  }

  /* Dark theme styling */
  :global(body) {
    margin: 0 auto;
    padding: 0;
    font-family: Ubuntu, Arial, sans-serif;
    background-color: var(--dark);
    color: #f0f0f0;
    min-height: 100vh;
    text-align: center;
  }
  
  .content-wrapper {
    width: 90%;
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
    margin-top: 10%;
    flex-grow: 1;
    overflow: auto;
  }
  
  main {
    display: flex;
    flex-direction: column;
    flex: 1;
    box-sizing: border-box;
    width: 100%;
  }

  /* Footer styling */
  .global-footer {
    background-color: #000;
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
    gap: 3rem; /* Spacing between footer sections */
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

  .footer-info {
    font-size: 1rem;
  }

  .footer-info a {
    color: var(--gold);
    text-decoration: none;
    transition: color 0.3s ease;
  }

  .footer-info a:hover {
    text-decoration: underline;
  }

  .footer-elt {
    margin-top: 1.5rem;
  }

  .suggestions-button {
    padding: 8px 16px;
    font-size: 14px;
    border-radius: 5px;
    margin-top: 0.5rem;

  }

  /* Mobile styles */
  @media screen and (max-width: 768px) {
    .footer-content {
      padding: 0 1rem;
    }

    .footer-logo-img {
      height: 50px; /* Slightly smaller logo for mobile */
    }

    .footer-info {
      font-size: 0.9rem;
    }

    :global(.modal-content) {
      width: 75%;
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
      padding-top: 3rem;
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
    margin-bottom: 20px;
  }

  :global(.cancel-button:hover) {
    background-color: grey;
  }

  :global(.form-group) {
    /* margin-top: 15px; */
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
  /* flex-direction: column; */
  align-items: center;
  padding: 0rem;
  background-color: #000;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 10;
  box-sizing: border-box;
}

/* Logo Row */
.navbar-logo img {
  left: 10px;
  height: 5rem;
  padding: 10px 20px;
}

/* Navigation Links */
.navbar-links {
  list-style: none;
  display: flex;
  /* justify-content: center; */
  align-items: center;
  gap: 2rem;
  margin: 0;
  padding: 0;
}

.navbar-links li {
  font-size: 1rem;
}

.navbar-links a {
  text-decoration: none;
  color: #ffffff;
  font-weight: bold;
  transition: color 0.3s ease;
}

.navbar-links a:hover {
  color: var(--gold);
}

.navbar-links a.active {
  color: var(--gold);
}

/* Instagram Icon */
.instagram-link {
  /* position: absolute; */
  top: 10px;
  right: 10px;
  /* align-items: center;
  display: flex; */
}

.instagram-icon {
  height: 30px;
  width: 30px;
  align-items: center;
}

  /* for mobile */
  .hamburger span {
    background-color: white;
    height: 3px;
    margin: 3px 0;
    width: 25px;
  }

  .hamburger {
    background: none; /* Ensure no background color */
    border: none; /* Remove any border */
    display: none;
    flex-direction: column;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    right: 0;
  }

  .mobile-menu {
    display: none; /* Hidden by default */
    flex-direction: column;
    position: absolute;
    top: 100%;
    right: 0;
    left: 0;
    background-color: #000;
    z-index: 10;
    font-family: Telegraf-Ultrabold, sans-serif;
  }

  .mobile-menu.open {
    display: flex;
  }

  .mobile-menu a {
    /* color: var(--gold); */
    color: white;
    text-decoration: none;
    padding: 1rem;
    text-align: center;
    border-bottom: 1px solid white;
  }

  .mobile-menu a:hover {
    color: var(--gold);
  }

  .instagram-link {
    margin-left: auto;
    padding: 1.5rem;
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

    :global(.modal-content) {

      width: 75%;

  }
  }
</style>

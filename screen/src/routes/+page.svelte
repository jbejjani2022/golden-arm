<script lang="ts">
	  import { formatDate } from '$lib';
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
  
    let movie: any = null;
    let error: string = '';
    let showModal = false;
  
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
  
  <!-- Home Page Layout -->
  <main class="homepage">
    <section class="theater-info">
      <h1>The Golden Arm</h1>
      <p>Eliot's student-run theater showcasing weekly films.</p>
    </section>
  
    {#if error}
      <p class="error">{error}</p>
    {/if}
  
    {#if movie}
      <section class="movie-info">
        <h2>Coming soon: "{movie.Title}"</h2>
        <h3>{formatDate(movie.Date)}</h3>
  
        <!-- Movie Poster -->
        <div class="poster">
          <img src={movie.PosterURL} alt="Movie poster for {movie.Title}" />
        </div>
  
        <!-- Movie Menu -->
        <div class="menu">
          <h3>Menu</h3>
          <img src={movie.MenuURL} alt="Menu for {movie.Title}" />
        </div>
        
        <button class="home-button" on:click={() => goto(`/reservations/${movie.ID}`)}>Reserve a seat</button>
      </section>
    {:else}
      <p>Loading movie information...</p>
    {/if}
    <br>
    <p>We want your suggestions!</p>
    <button class="home-button" on:click={confirmComment}>What should we screen next?</button>
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
            <label for="comment">What would you like us to screen next?</label>
            <input type="text" id="comment" bind:value={comment} placeholder="Enter any movie suggestions!" required />
          </div>
          <button type="submit" on:click={handleComment}>Send</button>
          <button type="button" class="cancel-button" on:click={cancelComment}>Cancel</button>
      </div>
    </div>
    {/if}
  </main>
  
  <style>
    .homepage {
      max-width: 800px;
      margin: 0 auto;
      padding: 20px;
      text-align: center;
    }
  
    .theater-info h1 {
      font-size: 36px;
      margin-bottom: 10px;
    }
  
    .theater-info p {
      font-size: 18px;
      margin-bottom: 30px;
      color: #bbb;
    }
  
    .movie-info h2 {
      font-size: 28px;
      margin: 20px 0;
    }
  
    .movie-info p {
      font-size: 16px;
      color: #ddd;
    }
  
    .poster img,
    .menu img {
      max-width: 100%;
      height: auto;
      border-radius: 8px;
      margin: 20px 0;
    }
  
    .home-button {
      display: inline-block;
      padding: 12px 24px;
      color: white;
      text-decoration: none;
      border-radius: 5px;
      margin-top: 20px;
    }
  
    .error {
      color: #ff5252;
      font-size: 16px;
    }

    .modal-content button {
      color: black;
    }
  </style>
  
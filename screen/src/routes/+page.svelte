<script lang="ts">
	import { formatDate } from '$lib';
    import { onMount } from 'svelte';
  
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
  
        <a class="reserve-button" href={`/reservations/${movie.ID}`}>Reserve a seat</a>
      </section>
    {:else}
      <p>Loading movie information...</p>
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
  
    .reserve-button {
      display: inline-block;
      padding: 12px 24px;
      background-color: #e50914;
      color: white;
      text-decoration: none;
      border-radius: 5px;
      font-weight: bold;
      margin-top: 20px;
    }
  
    .reserve-button:hover {
      background-color: #b20710;
    }
  
    .error {
      color: #ff5252;
      font-size: 16px;
    }
  </style>
  
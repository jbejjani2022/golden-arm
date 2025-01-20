<script lang="ts">
  import { formatDate } from '$lib';
  import { onMount } from 'svelte';

  let archive: Array<any> = [];
  let error: string = '';

  // Fetch the movie archive using /api/movie/archive
  onMount(async () => {
    try {
      const response = await fetch('/api/movie/archive');
      const data = await response.json();

      if (data.success) {
        archive = data.data;
      } else {
        error = 'Failed to load the movie archive.';
      }
    } catch (err) {
      console.error(err);
      error = 'Something went wrong while fetching the movie archive.';
    }
  });
</script>

<main class="archive">
  <section class="theater-info">
    <h1>The Golden Arm Archive</h1>
    <p>
      Discover our past screenings.
    </p>
  </section>

  {#if error}
    <p class="error">{error}</p>
  {/if}

  {#if archive.length > 0}
    <div class="movie-list">
      {#each archive as movie}
        <div class="movie-card">
          <h2>{movie.Title}</h2>
          <p>{formatDate(movie.Date)}</p>
          <div class="images">
            <img src={movie.PosterURL} alt="{movie.Title} poster" class="poster" />
            <img src={movie.MenuURL} alt="{movie.Title} menu" class="menu below" />
          </div>          
        </div>
      {/each}
    </div>
  {/if}
</main>

<style>
  main.archive {
    padding: 2rem;
    color: #f0f0f0;
    min-height: 100vh;
    margin-top: 20px;
  }

  h1 {
    text-align: center;
    font-size: 2rem;
    margin-bottom: 1.5rem;
  }

  .movie-list {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    overflow-y: auto;
  }

  .movie-card {
    background-color: #1a1a1a;
    padding: 2rem;
    border-radius: 25px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.5);
  }

  .movie-card h2 {
    margin: 0;
    font-size: 1.5rem;
  }

  .movie-card p {
    font-size: 1rem;
    color: #c0c0c0;
  }

  .movie-card:hover {
    transform: scale(1.02);
    transition: transform 0.2s ease;
  }

  .images {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .poster {
    width: 100%;
  }

  .menu.below {
    margin-top: 1rem;
    width: 80%;
  }

  .images img {
    width: 50%;
    height: auto;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  }

  .menu {
    flex: 1;
    max-width: 50%;
  }

  @media (max-width: 600px) {
    .images {
      flex-direction: column;
    }

    .poster,
    .menu {
      max-width: 100%;
    }
  }
</style>

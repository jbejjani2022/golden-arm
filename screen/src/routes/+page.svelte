<script lang="ts">
  import { formatDate, formatRuntime } from '$lib';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { apiBaseUrl } from '$lib/api';
  import { Splide, SplideSlide, SplideTrack } from '@splidejs/svelte-splide';
  import '@splidejs/svelte-splide/css';

  let movie: any = null;
  let calendar: any = null;
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

    try {
      const response = await fetch(`${apiBaseUrl}/calendar`);
      const data = await response.json();

      if (data.success) {
        calendar = data.data;
      } else {
        error = 'Failed to load the calendar.';
      }
    } catch (err) {
      console.error(err);
      error = 'Something went wrong while fetching the calendar.';
    }
  });

  // Past screening carousel
  type Movie = {
    PosterURL: string;
    Title: string;
  };

  let archive: Movie[] = [];

  onMount(async () => {
    try {
      const response = await fetch(`${apiBaseUrl}/movie/archive`);
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

  const options = {
    type: 'loop',
    drag: true,
    snap: true,
    perPage: 3,
    perMove: 1,
    focus: 0,
    autoplay: true,
    interval: 4000,
    speed: 2000,
    arrows: true,
    padding: '2rem',
    gap: '1.5rem',
    pagination: true,
    width: '100%', // use full container width
    breakpoints: {
      768: {
        perPage: 2,
        gap: '1rem',
        padding: '1rem'
      },
      480: {
        perPage: 1,
        gap: '0.5rem',
        padding: '0.5rem',
        arrows: false // hide arrows on mobile
      }
    }
  }

</script>
  
  <!-- Home Page Layout -->
  <main>
    <section class="top-text">
      ELIOT HOUSE'S STUDENT-RUN THEATER SHOWCASING WEEKLY FILMS.
    </section>
  
    {#if movie}
      <section class="movie-info">
        <!-- Left side: Movie Info -->
        <div class="movie-details">
          <h1 class="movie-title">{movie.Title}</h1>
          <div class="movie-screening">
            <div style="padding: 0.5rem">{formatRuntime(movie.Runtime)}</div>
            <div style="margin-bottom: 0.5rem">Screening {formatDate(movie.Date)}</div>
            <button class="reserve-button" on:click={() => goto(`/reservations/${movie.ID}`)}>Get Tickets</button>
          </div>
        </div>

        <!-- Right side: Movie Poster -->
        <div class="movie-poster">
          <img src={movie.PosterURL} alt="Movie poster for {movie.Title}" />
        </div>
      </section>
    {:else}
      <p>Loading movie information...</p>
    {/if}

  <div class="row-header">
    <h2 class="header-title">Past Screenings</h2>
    <a href="/archives" class="see-all-link">See All</a>
  </div>
    <!-- Separator line -->
  <div class="separator"></div>

  {#if archive.length > 0}
  <!-- <Splide aria-label="Past screening posters">
    {#each archive as movie}
    <SplideSlide>
      <img src={movie.PosterURL} alt={movie.Title}>
    </SplideSlide>
    {/each}
  </Splide> -->

  <Splide 
    options={ options } 
    hasTrack={ false } 
    aria-label="Past screening posters"
  >
  <div style="position: relative">
    <SplideTrack>
      {#each archive as movie}
        <SplideSlide>
          <img src={movie.PosterURL} alt={movie.Title}>
        </SplideSlide>
      {/each}
    </SplideTrack>
  </div>

    <div class="splide__progress">
      <div class="splide__progress__bar">
      </div>
    </div>
  </Splide>
  {:else}
    <p>Archive is empty...</p>
  {/if}

  <div class="row-header">
    <h2 class="header-title">Calendar</h2>
  </div>
  <div class="separator"></div>
  {#if calendar}
    <div class="calendar-image">
      <img src={calendar.ImageURL} alt="Calendar" />
    </div>
  {:else}
    <p>Calendar not found.</p>
  {/if}

</main>
  
<style>
  /* Update the Splide slide styles */
  :global(.splide__slide) {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  :global(.splide__slide img) {
    width: 100%;
    object-fit: cover;
    border-radius: 5px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    transition: transform 0.3s ease;
  }

  :global(.splide__slide:hover img) {
    transform: scale(1.05);
  }

  :global(.splide) {
    max-width: 1200px;
    margin: 0 auto;
    margin-bottom: 3rem;
  }

  .top-text {
    font-family: Telegraf-Ultrabold, sans-serif;
    font-size: 30px;
    font-weight: bold;
    margin-top: 20px;
    margin-bottom: 40px;
  }

  .movie-info {
    flex-wrap: wrap; /* Enables wrapping if content doesn't fit */
    display: flex; /* Use flexbox to align items side by side */
    justify-content: space-evenly; /* Space out the elements */
    /* min-width: 300px; Ensures a minimum width */
    justify-content: center;
    align-items: center; /* Vertically center the content */
    gap: 15%; /* Increase space between text and poster */
    padding: 20px;
  }

  .movie-details {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    margin-bottom: 1rem;
  }

  .movie-screening {
    margin-bottom: 10px;
    font-size: 1.2rem;

  }

  .movie-title {
    font-size: 2.5rem; /* Make title text larger */
    font-weight: bold; /* Make the title bold */
    margin-bottom: 5px;
    margin-top: 5px;

  }

  .reserve-button {
    cursor: pointer;
    font-size: 1rem;
    display: inline-block; /* Ensures the button only takes as much space as the text needs */
    margin-top: 10px;
    margin-bottom: 10px
  }

  .reserve-button:hover {
    background-color: var(--dark-gold);
  }

  .movie-poster {
    flex: 1; /* Take up the other half of the screen */
    max-width: 30%;
    min-width: 300px;
  }

  .movie-poster img {
    width: 100%;
    height: auto;
    border-radius: 5px;

  }

  /* .error {
    color: #ff5252;
    font-size: 16px;
  } */

  /* Separator line styling */
  .separator {
    width: 100%; /* Ensures it spans the full width of the container */
    height: 2px; /* Thickness of the line */
    background-color: #ddd; /* Color of the line */
    margin: 1rem 0; /* Optional: Adds spacing around the line */
    border: none; /* Removes any default borders */
  }

  /* Row header container */
  .row-header {
    display: flex;
    justify-content: space-between; /* Aligns the elements to opposite ends */
    align-items: center; /* Centers vertically */
    width: 100%;
    margin-bottom: 1rem; /* Adds space below the row */
    margin-top: 1.5rem;
  }

  /* Calendar styling */
  .calendar-image {
    display: flex;
    justify-content: center;
    align-items: center;
    margin: 2rem auto;
    max-width: 90%;  /* Prevents calendar from being too wide on large screens */
  }

  .calendar-image img {
    width: 100%;
    height: auto;  /* Maintains aspect ratio */
    max-width: 800px;  /* Adjust this value based on your needs */
    object-fit: contain;
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.25);
  }

  /* Title styling */
  .header-title {
    font-size: 1.5rem; /* Larger text */
    /* font-weight: bold; Bold font */
    margin: 0; /* Removes default margin */
    color: #fff; /* Dark text color */
  }

  /* See All link styling */
  .see-all-link {
    font-size: 1rem; /* Smaller text */
    color: #fff; /* Link color (blue) */
    text-decoration: none; /* Removes underline */
    font-weight: normal; /* Normal font weight */
    transition: color 0.3s; /* Smooth color transition on hover */
  }

  /* Hover effect for the link */
  .see-all-link:hover {
    color: #caac3e; /* Darker blue on hover */
    text-decoration: underline; /* Optional underline on hover */
  }


  @media screen and (max-width: 768px) {
    .top-text {
        margin-top: 100px; /* Adjust as needed to increase space */
    }

    .movie-screening {
      font-size: 13px;
    }

}

/* for da ipad kids */
@media screen and (min-width: 768px) and (max-width: 1024px) and (orientation: portrait), 
       screen and (min-width: 1024px) and (max-width: 1366px) and (orientation: landscape) {
    .top-text {
        margin-top: 80px; /* Adjust this value as needed */
        padding-top: 20px; /* Optional: Additional padding for spacing */
    }
}
</style>
  
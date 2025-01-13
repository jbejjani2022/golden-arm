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

  let currentSlide = 1; // Start at 1 to account for the first clone
  let isTransitioning = false;
  let interval: any;

  function nextSlide() {
    if (isTransitioning) return;
    isTransitioning = true;
    currentSlide++;
    updateCarousel();
  }

  function prevSlide() {
    if (isTransitioning) return;
    isTransitioning = true;
    currentSlide--;
    updateCarousel();
  }

  function updateCarousel() {
    const track = document.querySelector('.carousel-track') as HTMLElement;
    track.style.transition = 'transform 0.5s ease-in-out';
    track.style.transform = `translateX(-${currentSlide * 100}%)`;
  }

  function handleTransitionEnd() {
    const track = document.querySelector('.carousel-track') as HTMLElement;

    // Handle infinite looping
    if (currentSlide === archive.length + 1) {
      currentSlide = 1;
      track.style.transition = 'none';
      track.style.transform = `translateX(-${currentSlide * 100}%)`;
    }
    if (currentSlide === 0) {
      currentSlide = archive.length;
      track.style.transition = 'none';
      track.style.transform = `translateX(-${currentSlide * 100}%)`;
    }

    isTransitioning = false;
  }

  // Automatically rotate through the slides
  onMount(() => {
    interval = setInterval(nextSlide, 3000); // Change slide every 3 seconds
    return () => clearInterval(interval); // Clear interval on component unmount
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

  <div class="carousel">
    {#if archive.length > 0}
      <div
        class="carousel-track"
        style="--total:{archive.length}"
        on:transitionend={handleTransitionEnd}
      >
        <!-- Clone the first and last slides for seamless infinite looping -->
        <div class="carousel-slide">
          <img src={archive[archive.length - 1]?.PosterURL} alt="Last Clone Poster" />
        </div>
        {#each archive as movie}
          <div class="carousel-slide">
            <img src={movie.PosterURL} alt="{movie.Title} poster" />
          </div>
        {/each}
        <div class="carousel-slide">
          <img src={archive[0]?.PosterURL} alt="First Clone Poster" />
        </div>
      </div>
      <button class="carousel-btn prev" on:click={prevSlide}>&#8249;</button>
      <button class="carousel-btn next" on:click={nextSlide}>&#8250;</button>
    {/if}
  </div>

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
  }

  h1 {
    text-align: center;
    font-size: 2rem;
    margin-bottom: 1.5rem;
  }

  .carousel {
    position: relative;
    width: 100%;
    overflow: hidden;
    margin-bottom: 2rem;
  }

  .carousel-track {
    display: flex;
    transition: transform 0.5s ease-in-out;
  }

  .carousel-slide {
    min-width: 100%;
    display: flex;
    justify-content: center;
  }

  .carousel-slide img {
    width: 50%;
    height: auto;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  }

  .carousel-btn {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    color: #fff;
    border: none;
    font-size: 2rem;
    cursor: pointer;
    padding: 5px;
    z-index: 1;
  }

  .carousel-btn.prev {
    left: 0;
  }

  .carousel-btn.next {
    right: 0;
  }

  .carousel-btn:hover {
    background: rgba(255, 255, 255, 0.1);
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

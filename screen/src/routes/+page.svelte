<script lang="ts">
	  import { formatDate, formatRuntime } from '$lib';
    import { goto } from '$app/navigation';
    import { onMount, onDestroy } from 'svelte';

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

  // Past screening carousel
  type Movie = {
    PosterURL: string;
    Title: string;
  };

  let archive: Movie[] = [];
  let currentIndex: number = 0;
  let intervalId: any;
  
  // For infinite loop effect, duplicate the first 2 slides at the end
  const numberOfVisibleSlides = 3;
  // Create a duplicatedArchive of the same type
  let duplicatedArchive: Movie[] = [];

  // Automatically rotate every 3 seconds
  onMount(async () => {
    try {
      const response = await fetch('/api/movie/archive');
      const data = await response.json();

      if (data.success) {
        archive = data.data;
        duplicatedArchive = [...archive, ...archive]; // Duplicate for the seamless loop
      } else {
        error = 'Failed to load the movie archive.';
      }

      // Start the automatic carousel rotation
      intervalId = setInterval(() => {
        nextSlide();
      }, 3000); // 3 second interval

    } catch (err) {
      console.error(err);
      error = 'Something went wrong while fetching the movie archive.';
    }
  });
 
  function prevSlide() {
    // Move to the previous slide
    if (currentIndex === 0) {
      currentIndex = duplicatedArchive.length / 2 - numberOfVisibleSlides; // Set to last valid index
    } else {
      currentIndex--;
    }
  }

  function nextSlide() {
    // Move to the next slide
    if (currentIndex === duplicatedArchive.length / 2 - numberOfVisibleSlides) {
      currentIndex = 0; // Reset to the first position
    } else {
      currentIndex++;
    }
  }

  // Clear interval when component is destroyed
  onDestroy(() => {
    if (intervalId) clearInterval(intervalId);
  });
 
  </script>
  
  <!-- Home Page Layout -->
  <main>
    <section class="top-text">
      ELIOT HOUSE'S STUDENT-RUN THEATER SHOWCASING WEEKLY FILMS.
    </section>
  
    {#if error}
      <p class="error">{error}</p>
    {/if}
  
    {#if movie}
      <section class="movie-info">
        <!-- Left side: Movie Info -->
        <div class="movie-details">
          <h1 class="movie-title">{movie.Title}</h1>
          <div class="movie-screening">
            <h2>{formatRuntime(movie.Runtime)}</h2>
            <h2>Screening {formatDate(movie.Date)}</h2>
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

  {#if error}
  <p>{error}</p>
{:else}
  <div class="carousel-container">
    <div
      class="carousel"
      style="transform: translateX(-{(currentIndex % archive.length) * 33.33}%)"
    >
      {#each [...archive, ...duplicatedArchive] as movie, index}
        <div
          class="carousel-slide {index === currentIndex || index === (currentIndex + 1) % archive.length || index === (currentIndex - 1 + archive.length) % archive.length ? 'active' : ''}"
        >
          <img src={movie.PosterURL} alt="{movie.Title} poster" />
        </div>
      {/each}
    </div>

    <button class="carousel-arrow carousel-arrow-left" on:click={prevSlide}>&lt;</button>
    <button class="carousel-arrow carousel-arrow-right" on:click={nextSlide}>&gt;</button>
  </div>
{/if}
    
    <div class="separator"></div>
    <div class="suggestions">
      <p>We want your suggestions!</p>
      <button class="home-button" on:click={confirmComment}>What should we screen next?</button>
    </div>
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
  </main>
  
<style>
    .top-text {
      font-size: 30px;
      font-weight: bold;
      margin-bottom: 30px;
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
      /* flex: 1;  */
      display: flex;
      flex-direction: column;
      justify-content: center;
      /* gap: 1rem; Spacing between elements within the details section */
    }

    .movie-screening {
      margin-bottom: 10px;
    }

    .movie-title {
      font-size: 2.5rem; /* Make title text larger */
      font-weight: bold; /* Make the title bold */
      margin-bottom: 5px;
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
      border-radius: 8px;

    }

    .home-button {
      display: inline-block;
      font-weight: bold;
      text-decoration: none;
      margin-bottom: 20px;
    }
  
    .error {
      color: #ff5252;
      font-size: 16px;
    }

    .modal-content button {
      color: black;
    }

    /* Separator line styling */
    .separator {
      width: 100%; /* Ensures it spans the full width of the container */
      height: 2px; /* Thickness of the line */
      background-color: #ddd; /* Color of the line */
      margin: 1rem 0; /* Optional: Adds spacing around the line */
      border: none; /* Removes any default borders */
    }

    .suggestions {
      margin-top: 2rem;
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

  /* Title styling */
  .header-title {
    font-size: 1.5rem; /* Larger text */
    font-weight: bold; /* Bold font */
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

  /* carousel */
  .carousel-container {
      display: flex;
      justify-content: center;
      align-items: center;
      overflow: hidden;
      position: relative;
      width: 80%; /* Adjust the width as necessary */
      margin: 0 auto;
    }

  .carousel {
    display: flex;
    transition: transform 0.5s ease;
  }

  .carousel-slide {
    flex: 0 0 33.33%;
    display: flex;
    justify-content: center;
    align-items: center;
    opacity: 0.7;
  }

  .carousel-slide img {
    width: 80%;
    height: auto;
    transition: transform 0.5s ease;
  }

  .carousel-slide.active img {
    transform: scale(1.2);
    opacity: 1;
  }

  .carousel-slide img:not(.active) {
    transform: scale(0.8);
    opacity: 0.7;
  }

  .carousel-arrow {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    background-color: rgba(0, 0, 0, 0.5);
    color: white;
    padding: 10px;
    font-size: 24px;
    cursor: pointer;
    z-index: 10;
  }

  .carousel-arrow-left {
    left: 10px;
  }

  .carousel-arrow-right {
    right: 10px;}

</style>
  
<script lang="ts">
    import { formatDate } from '$lib';
    import { onMount } from 'svelte';
  
    // let archive: Array<any> = [];
    // let error: string = '';
  
    // // Fetch the movie archive using /api/movie/archive
    // onMount(async () => {
    //   try {
    //     const response = await fetch('/api/movie/archive');
    //     const data = await response.json();
  
    //     if (data.success) {
    //       archive = data.data;
    //     } else {
    //       error = 'Failed to load the movie archive.';
    //     }
    //   } catch (err) {
    //     console.error(err);
    //     error = 'Something went wrong while fetching the movie archive.';
    //   }
    // });
  
    // let currentSlide = 1; // Start at 1 to account for the first clone
    // let isTransitioning = false;
    // let interval: any;
  
    // function nextSlide() {
    //   if (isTransitioning) return;
    //   isTransitioning = true;
    //   currentSlide++;
    //   updateCarousel();
    // }
  
    // function prevSlide() {
    //   if (isTransitioning) return;
    //   isTransitioning = true;
    //   currentSlide--;
    //   updateCarousel();
    // }
  
    // function updateCarousel() {
    //   const track = document.querySelector('.carousel-track') as HTMLElement;
    //   track.style.transition = 'transform 0.5s ease-in-out';
    //   track.style.transform = `translateX(-${currentSlide * 100}%)`;
    // }
  
    // function handleTransitionEnd() {
    //   const track = document.querySelector('.carousel-track') as HTMLElement;
  
    //   // Handle infinite looping
    //   if (currentSlide === archive.length + 1) {
    //     currentSlide = 1;
    //     track.style.transition = 'none';
    //     track.style.transform = `translateX(-${currentSlide * 100}%)`;
    //   }
    //   if (currentSlide === 0) {
    //     currentSlide = archive.length;
    //     track.style.transition = 'none';
    //     track.style.transform = `translateX(-${currentSlide * 100}%)`;
    //   }
  
    //   isTransitioning = false;
    // }
  
    // // Automatically rotate through the slides
    // onMount(() => {
    //   interval = setInterval(nextSlide, 3000); // Change slide every 3 seconds
    //   return () => clearInterval(interval); // Clear interval on component unmount
    // });
  </script>
  
  <main class="archive">
    <div>
         <h1>Coming Soon ... 👀</h1>
    </div>
     
  
  
  
    
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
  
  </style>
  
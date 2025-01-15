<script lang="ts">
  import { page } from '$app/state';
  import { onMount } from 'svelte';
  import { formatDate } from '$lib';
  
  let movie: any = null;
  let error: string = '';

  // Fetch movie information
  onMount(async () => {
    try {
      const response = await fetch(`/api/movie/${page.params.movie_id}`);
      const data = await response.json();

      if (data.success) {
        movie = data.data;
      } else {
        error = 'Failed to load the movie data.';
      }
    } catch (err) {
      console.error(err);
      error = 'Something went wrong while fetching the movie data.';
    }
  });

  let reservedSeats: string[] = []; // Store reserved seats

  onMount(async () => {
    const movieId = page.params.movie_id;
    try {
      const response = await fetch(`/api/reserved/${movieId}`);
      const result = await response.json();
      if (result.success) {
        reservedSeats = result.data.reserved_seats; // Update reserved seats array
      } else {
        console.error('Failed to load reserved seats data');
      }
    } catch (err) {
      console.error('Error fetching reserved seats:', err);
    }
  });


  // Define the Seat interface
  interface Seat {
    num: string;
    selected: boolean;
  }

  // Create custom rows for the grid: 5 seats, 7 seats, 6 seats
  let seats: Seat[][] = [
    Array.from({ length: 5 }, (_, col) => ({ num: `C${col + 1}`, selected: false })),
    Array.from({ length: 7 }, (_, col) => ({ num: `B${col + 1}`, selected: false })),
    Array.from({ length: 6 }, (_, col) => ({ num: `A${col + 1}`, selected: false })),
  ];

  let selectedSeat: Seat | null = null;
  let showResModal = false;
  let showCommentModal = false;
  let name = '';
  let email = '';
  let comment = '';

  function toggleSeat(seat: Seat) {
    seats = seats.map(row =>
      row.map(s => {
        if (s.num === seat.num) {
          s.selected = !s.selected; // Toggle selection
          selectedSeat = s.selected ? s : null;
        } else {
          s.selected = false; // Deselect other seats
        }
        return s;
      })
    );
  }

  const confirmReservation = () => {
    showResModal = true;
  }

  const cancelReservation = () => {
    showResModal = false;
  }

  const handleReservation = async () => {
    if (!name || !email) {
      alert("Please enter your name and email.");
      return;
    }
    showResModal = false;
    // Send the reservation
    try {
      const response = await fetch(`/api/reserve`, {
        method: "POST",
        headers: { 
          "Content-Type": "application/json" 
        },
        body: JSON.stringify({
          movie_id: page.params.movie_id,
          seat_number: selectedSeat?.num,
          name,
          email,
        })
      });

      const result = await response.json();
      if (result.success) {
        reservedSeats.push(selectedSeat?.num || '');  // add new reserved seat
        toggleSeat(selectedSeat || { num: '', selected: false });
        alert("Reservation confirmed!");
        confirmComment();
      } else {
        alert("Failed to confirm reservation.");
      }
    } catch (err) {
      console.error(err);
      alert("Something went wrong while confirming the reservation.");
    }
  }

  const confirmComment = () => {
    showCommentModal = true;
  }

  const cancelComment = () => {
    showCommentModal = false;
  }

  const handleComment = async () => {
    if (!comment) {
      alert("Please tell us what you'd like to see next!");
      return;
    }
    showCommentModal = false;
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
      if (result.success) {
        alert("Thank you for your suggestion!");
      } else {
        alert("Failed to submit comment.");
      }
    } catch (err) {
      console.error(err);
      alert("Something went wrong while submitting the comment.");
    }
  }
</script>

<main class="reservation-page">
  {#if movie}
  <div class="movie-info">
    <div class="movie-icon">
      <img src={movie.PosterURL} alt="Movie poster for {movie.Title}" />
    </div>
    <div class="movie-details">
      <h1>{movie.Title}</h1>
      <p class="movie-date">{formatDate(movie.Date)}</p>
    </div>
  </div>
  {:else}
  <p>Loading movie information...</p>
  {/if}

<h3>Select Your Seat</h3>
<div class="grid">
  {#each seats as row, rowIndex}
    <div class="row">
      {#each row as seat, colIndex}
        {#if rowIndex === 2 && colIndex === 3}
          <div class="seat-space"></div> <!-- Add space in the middle of the 6 chairs -->
        {/if}
        <button
          class="seat {reservedSeats.includes(seat.num) ? 'reserved' : ''}"
          disabled={reservedSeats.includes(seat.num)} 
          on:click={() => toggleSeat(seat)}
        >
          <img
          src={reservedSeats.includes(seat.num) ? '/grey-chair.png' : (seat.selected ? "/yellow-chair.png" : "/white-chair.png")}
            alt="Seat {seat.num}"
          />
          <span class="seat-label">{seat.num}</span> <!-- Add seat number here -->
        </button>
      {/each}
    </div>
  {/each}
</div>

<div id="screen-container">
  <div id="screen">Screen</div>
</div>

<button class="reserve-button" on:click={confirmReservation} disabled={!selectedSeat}>Confirm</button>

{#if showResModal}
<div class="modal">
  <div class="modal-content">
      <h2>Seat {selectedSeat?.num}</h2>
      <div class="form-group">
        <label for="name">Name: </label>
        <input type="text" id="name" bind:value={name} placeholder="Enter your name" required />
      </div>
      <div class="form-group">
        <label for="email">Email: </label>
        <input type="email" id="email" bind:value={email} placeholder="Enter your email" required />
      </div>
      <button type="submit" on:click={handleReservation}>Reserve</button>
      <button type="button" class="cancel-button" on:click={cancelReservation}>Cancel</button>
  </div>
</div>
{/if}

{#if showCommentModal}
<div class="modal">
  <div class="modal-content">
      <h2>What would you like to see next at The Golden Arm?</h2>
      <div class="form-group">
        <input type="text" id="comment" bind:value={comment} placeholder="Enter movie suggestions!" required />
      </div>
      <button type="submit" on:click={handleComment}>Send</button>
      <button type="button" class="cancel-button" on:click={cancelComment}>Cancel</button>
  </div>
</div>
{/if}
</main>

<style>
.reservation-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

/* new */
.movie-info {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.movie-icon img {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
}

.movie-details h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: bold;
}

.movie-details .movie-date {
  margin: 4px 0 0;
  font-size: 0.9rem;
  color: gray;
}

h1 {
  text-align: center;
}

.row {
  display: flex;
  gap: 10px;
  margin-top: 20px;
  margin-bottom: 20px;
  justify-content: center; /* Centers the seats horizontally within the row */
}

.seat {
  width: 50px;
  height: 50px;
  background-color: transparent; /* Optional: Remove any background color */
  border: none; /* Remove any borders if unnecessary */
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: visible; /* Allow the image to show fully */
}

button {
  display: inline-block;
  padding: 10px;
  color: black;
  text-decoration: none;
  border-radius: 5px;
}

.reserve-button {
  display: inline-block;
  padding: 12px 24px;
  color: black;
  text-decoration: none;
  border-radius: 5px;
  font-weight: bold;
  margin-top: 20px;
}

/* seatstuff */
.seat {
  background: none;
  border: none;
  padding: 0;
  margin: 0;
  cursor: pointer;
  transition: transform 0.2s ease-in-out;
  font-size: 12px; /* Adjust font size for seat labels */
  text-align: center; /* Center-align the text */
}

.seat img {
  /* transform: rotate(180deg); */

  width: 50px;
  height: 50px;
  transition: transform 0.2s ease;
}

.seat:hover {
  transform: scale(1.1); /* Enlarges seat on hover */
}

.reserved {

  cursor: not-allowed; /* Change cursor to not-allowed */
}

.reserved:hover {
  transform: none; /* Disable enlarging on hover */
}

/* screen */
#screen-container {
  display: flex;
  justify-content: center; /* Center the rectangle horizontally */
  margin-bottom: 20px; /* Space between the rectangle and the Confirm button */
}

#screen {
  background-color: white; /* Rectangle background color */
  color: black; /* Text color */
  width: 300px; /* Width of the rectangle */
  height: 30px; /* Height of the rectangle */
  display: flex;
  justify-content: center; /* Center text horizontally */
  align-items: center; /* Center text vertically */
  font-size: 18px; /* Font size of the text */
  border-radius: 10px; /* Optional: Rounded corners */
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2); /* Optional: Add a subtle shadow */
}

/* seat space */
.seat-space {
  width: 20px; /* Adjust the gap width */
}

.seat-label {
  position: absolute; /* Position label absolutely to center it */
  z-index: 2; /* Ensure the label is above the image */
  color: black;
  font-size: 14px;
  font-weight: bold;
  pointer-events: none; /* Prevent interfering with button click */
  padding-top: 15px;
}
</style>

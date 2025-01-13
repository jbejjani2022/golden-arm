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

  // Define the Seat interface
  interface Seat {
    id: number;
    selected: boolean;
  }

  // Create 2 x 9 seats grid
  let seats: Seat[][] = Array.from({ length: 2 }, (_, row) =>
    Array.from({ length: 9 }, (_, col) => ({ id: row * 9 + col + 1, selected: false }))
  );

  let selectedSeat: Seat | null = null;
  let showResModal = false;
  let showCommentModal = false;
  let name = '';
  let email = '';
  let comment = '';

  function toggleSeat(seat: Seat) {
    seats = seats.map(row =>
      row.map(s => {
        if (s.id === seat.id) {
          s.selected = !s.selected;
          selectedSeat = s.selected ? s : null;
        } else {
          s.selected = false;
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
          seat_number: selectedSeat?.id,
          name,
          email,
        })
      });

      const result = await response.json();
      if (result.success) {
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
<h1>The Golden Arm's Screening of "{movie.Title}"</h1>
<h2>{formatDate(movie.Date)}</h2>
{:else}
<p>Loading movie information...</p>
{/if}
<h3>Select Your Seat</h3>
<div class="grid">
  {#each seats as row}
    <div class="row">
      {#each row as seat}
        <button
          class="seat"
          class:selected={seat.selected}
          on:click={() => toggleSeat(seat)}
        >
          {seat.id}
        </button>
      {/each}
    </div>
  {/each}
</div>
<button class="reserve-button" on:click={confirmReservation} disabled={!selectedSeat}>Confirm</button>

{#if showResModal}
<div class="modal">
  <div class="modal-content">
      <h2>Seat {selectedSeat?.id}</h2>
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
  background-color: #ccc;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.seat.selected {
  background-color: var(--gold);
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

.reserve-button:disabled {
  background-color: grey;
}
</style>

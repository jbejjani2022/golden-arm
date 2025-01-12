<script lang="ts">
  import { page } from '$app/state';

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
  let showModal = false;
  let name = '';
  let email = '';

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
    showModal = true;
  }

  const cancelReservation = () => {
    showModal = false;
  }

  const handleReservation = async () => {
    if (!name || !email) {
      alert("Both name and email are required.");
      return;
    }
    if (name && email) {
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
          showModal = false;
        } else {
          alert("Failed to confirm reservation.");
        }
      } catch (err) {
        console.error(err);
        alert("Something went wrong while confirming the reservation.");
      }
    }
  }
</script>

<main class="reservation-page">
<h1>Select Your Seat</h1>
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

{#if showModal}
<div class="modal">
  <div class="modal-content">
      <h2>Enter Your Information</h2>
      <div class="form-group">
        <label for="name">Name: </label>
        <input type="text" id="name" bind:value={name} placeholder="Enter your name" required />
      </div>

      <div class="form-group">
        <label for="email">Email: </label>
        <input type="email" id="email" bind:value={email} placeholder="Enter your email" required />
      </div>
      <button type="submit">Submit</button>
      <button type="button" class="cancel-button" on:click={cancelReservation}>Cancel</button>
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
  background-color: green;
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

.modal {
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

.modal-content {
    background-color: #fff;
    color: black;
    padding: 20px;
    border-radius: 8px;
    width: 300px;
    text-align: center; /* Center text */
    display: flex;
    flex-direction: column;  /* Arrange elements vertically */
    align-items: center;     /* Center items horizontally */
}

.modal-content input {
    margin: 10px 0; /* Space between input fields */
    padding: 8px;
    width: 100%;
    max-width: 250px; /* Limit input width */
}

.modal-content button {
    padding: 10px 15px;
    background-color: #4CAF50;
    color: white;
    border: none;
    cursor: pointer;
    border-radius: 5px;
    width: 100%;
    max-width: 250px; /* Same as input field width */
    margin-top: 10px; /* Space between inputs and button */
}

.modal-content button:hover {
    background-color: #45a049;
}

.modal-content .cancel-button {
  background-color: darkgrey;
}

.modal-content .cancel-button:hover {
  background-color: grey;
}

.form-group {
  margin-bottom: 15px;
  display: flex;
  align-items: center;
}

.form-group label {
  margin-right: 10px;
  width: auto;
}

.form-group input {
  flex-grow: 1;
  padding: 5px;
  font-size: 14px;
}
</style>

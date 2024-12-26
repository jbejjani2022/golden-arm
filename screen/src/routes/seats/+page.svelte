<script>
    let seats = Array.from({ length: 2 }, (_, row) =>
      Array.from({ length: 4 }, (_, col) => ({ id: row * 4 + col + 1, selected: false }))
    );
    let selectedSeat = null;
  
    function toggleSeat(seat) {
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
  
    function confirmReservation() {
      const name = prompt("Enter your name:");
      const email = prompt("Enter your email:");
      if (name && email) {
        fetch("/api/reserve", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ seat_number: selectedSeat.id, name, email })
        })
          .then(res => res.json())
          .then(data => alert(data.message))
          .catch(err => console.error(err));
      }
    }
</script>

<main>
  <h1>Select Your Seat</h1>
  <div class="grid">
    {#each seats as row}
      <div class="row">
        {#each row as seat}
          <button
            class:seat
            class:selected={seat.selected}
            on:click={() => toggleSeat(seat)}
          >
            {seat.id}
          </button>
        {/each}
      </div>
    {/each}
  </div>
  <button on:click={confirmReservation} disabled={!selectedSeat}>Confirm</button>
</main>

<style>
  .grid { display: grid; gap: 10px; }
  .row { display: flex; gap: 10px; }
  .seat { width: 50px; height: 50px; }
  .selected { background-color: green; }
</style>

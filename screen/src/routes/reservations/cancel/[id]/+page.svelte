<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';

  const reservationId = $page.params.id;

  const deleteReservation = async (reservationId: string) => {
      try {
        const response = await fetch(`/api/reservation/${reservationId}`, {
          method: 'DELETE',
        });

        if (response.ok) {
          alert('Reservation canceled!');
          goto('/');
        } else {
          console.error('Failed to delete reservation');
        }
      } catch (err) {
        console.error('Error during reservation deletion:', err);
      }
  };
</script>
  
<main class="confirm">
  <div>
    <h1>Are you sure you want to cancel your reservation?</h1>
    <div class="button-row">
      <button class="cancel" on:click={() => goto('/')}>Nah</button>
      <button class="confirm" on:click={() => deleteReservation(reservationId)}>Yeah</button>
    </div>
  </div>
</main>
  
<style>
  main.confirm {
    padding: 2rem;
    color: #ffffff;
    text-align: center;
  }

  h1 {
    font-size: 2rem;
    margin-bottom: 2rem;
  }

  .button-row {
    display: flex;
    justify-content: center;
    gap: 1rem;
  }

  .cancel {
    background-color: #555;
    color: white;
  }

  .cancel:hover {
    background-color: #777;
  }

  @media screen and (max-width: 768px) {
    main.confirm {
      margin-top: 4rem;
    }

    .button-row {
      flex-direction: column;
      gap: 0.75rem;
    }
  }
</style>

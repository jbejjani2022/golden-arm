<script lang="ts">
    import { page } from '$app/state';
    import { onMount } from 'svelte';
    import { formatDate } from '$lib';
  
    let movieTitle = '';
    let movieDate = '';
    let reservations: Array<any> = [];
    let error = '';
  
    onMount(async () => {
      const { movie_id } = page.params;
  
      try {
        const response = await fetch(`/api/reservations/${movie_id}`);
        const data = await response.json();
  
        if (data.success) {
          try {
            const movieResponse = await fetch(`/api/movie/${movie_id}`);
            const movieData = await movieResponse.json();
            if (movieData.success) {
              console.log('Movie data:', movieData.data);
              movieTitle = movieData.data.Title;
              movieDate = formatDate(movieData.data.Date);
            } else {
              error = 'Failed to load movie data.';
            }
          } catch (movieErr) {
            console.error('Error fetching movie data:', movieErr);
          }
          reservations = data.data;
        } else {
          error = 'Failed to load reservation data.';
        }
      } catch (err) {
        console.error(err);
        error = 'Something went wrong while fetching reservation data.';
      }
    });

    const deleteReservation = async (reservationId: string) => {
      const confirmation = confirm('Are you sure you want to delete this reservation?');
  
      if (confirmation) {
        try {
          const response = await fetch(`/api/reservation/${reservationId}`, {
            method: 'DELETE',
          });
  
          if (response.ok) {
            window.location.reload();
          } else {
            console.error('Failed to delete reservation');
          }
        } catch (err) {
          console.error('Error during reservation deletion:', err);
        }
      }
    };
</script>
  
<h1>{movieTitle} ({movieDate})</h1>

{#if error}
<p style="color: red;">{error}</p>
{/if}

<h2>Reservations</h2>
{#if reservations.length > 0}
<table>
    <thead>
    <tr>
        <th>ID</th>
        <th>Seat Number</th>
        <th>Name</th>
        <th>Email</th>
        <th>Date</th>
        <th>Actions</th>
    </tr>
    </thead>
    <tbody>
    {#each reservations as res (res.ID)}
        <tr>
        <td>{res.ID}</td>
        <td>{res.SeatNumber}</td>
        <td>{res.Name}</td>
        <td>{res.Email}</td>
        <td>{formatDate(res.Date)}</td>
        <td>
            <button on:click={() => deleteReservation(res.ID)} style="color: red; border: none; background: none; cursor: pointer;">X</button>
        </td>
        </tr>
    {/each}
    </tbody>
</table>
{:else}
<p>No reservations found for this movie.</p>
{/if}
  
<style>
    table {
        width: 100%;
        border-collapse: collapse;
        margin-top: 20px;
    }

    th, td {
        padding: 12px;
        text-align: left;
        border: 1px solid #ddd;
    }

    th {
        background-color: #f4f4f4;
    }
</style>
  
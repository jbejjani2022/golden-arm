<script lang="ts">
    import { page } from '$app/state';
    import { onMount } from 'svelte';
    import { formatDate } from '$lib';
    import { apiBaseUrl } from '$lib/api';
  
    let movieTitle = '';
    let movieDate = '';
    let reservations: Array<any> = [];
    let error = '';
  
    onMount(async () => {
      const { movie_id } = page.params;
  
      try {
        const response = await fetch(`${apiBaseUrl}/reservations/${movie_id}`);
        const data = await response.json();
  
        if (data.success) {
          try {
            const movieResponse = await fetch(`${apiBaseUrl}/movie/${movie_id}`);
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
          const response = await fetch(`${apiBaseUrl}/reservation/${reservationId}`, {
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

    const copyEmailList = async () => {
      // Extract unique emails from reservations
      const uniqueEmails = [...new Set(reservations.map(res => res.Email))];
      const emailList = uniqueEmails.join(', ');

      try {
        await navigator.clipboard.writeText(emailList);
        alert('Email list copied to clipboard!');
      } catch (err) {
        console.error('Failed to copy email list:', err);
        alert('Failed to copy email list. Please try again.');
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
        <td>{res.SeatNumber}</td>
        <td>{res.Name}</td>
        <td>{res.Email}</td>
        <td>{formatDate(res.Date)}</td>
        <td>
            <button class="link-button delete" on:click={() => deleteReservation(res.ID)}>Delete</button>
        </td>
        </tr>
    {/each}
    </tbody>
</table>
{:else}
<p>No reservations found for this movie.</p>
{/if}

<button on:click={copyEmailList} style="margin-top: 20px; padding: 10px 20px; cursor: pointer;">Get Movie Email List</button>
  
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

    .link-button {
        background: none;
        border: none;
        padding: 0;
        color: #0066cc;
        cursor: pointer;
        font: inherit;
        margin-right: 1rem;
    }

    .link-button.delete {
        color: red;
    }

    .link-button:hover {
        text-decoration: none;
    }
</style>
  
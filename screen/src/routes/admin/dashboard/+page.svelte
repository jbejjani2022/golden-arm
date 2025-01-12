<script lang="ts">
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';

    // Movie data storage
    let movies: Array<any> = [];

    // Error message for fetching movies
    let error: string = '';

    // Fetch movie data on page load
    onMount(async () => {
        try {
            const response = await fetch('/api/movie/all');
            const data = await response.json();

            if (data.success) {
                movies = data.data;
            } else {
                error = 'Failed to load movie data.';
            }
        } catch (err) {
            console.error(err);
            error = 'Something went wrong while fetching the movie data.';
        }
    });

    function formatDate(dateString: string): string {
      const date = new Date(dateString);
      const options: Intl.DateTimeFormatOptions = {
            month: '2-digit',
            day: '2-digit',
            year: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            hour12: true,
        };
        return date.toLocaleString('en-US', options);
    }
  
    const logout = async () => {
      try {
        const response = await fetch('/api/admin/logout', {
          method: 'POST',
        });

        if (response.ok) {
          // Redirect to the login page
          goto('/admin');
        } else {
          console.error('Logout failed');
        }
      } catch (err) {
        console.error('Error during logout:', err);
      }
  };

</script>
  
<h1>Hey Golden Arm Operator</h1>

<button on:click={logout}>Logout</button>

<!-- Display error message if data fetching fails -->
{#if error}
  <p style="color: red;">{error}</p>
{/if}

<!-- Movie Table -->
{#if movies.length > 0}
<table>
    <thead>
        <tr>
            <th>ID</th>
            <th>Title</th>
            <th>Screening Date</th>
            <th>Poster URL</th>
            <th>Menu URL</th>
        </tr>
    </thead>
    <tbody>
        {#each movies as movie (movie.ID)}
            <tr>
                <td>{movie.ID}</td>
                <td>{movie.Title}</td>
                <td>{formatDate(movie.Date)}</td>
                <td><a href={movie.PosterURL} target="_blank">View Poster</a></td>
                <td><a href={movie.MenuURL} target="_blank">View Menu</a></td>
            </tr>
        {/each}
    </tbody>
</table>
{:else}
  <p>No movies available at the moment.</p>
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

  a {
      color: #3498db;
      text-decoration: none;
  }

  a:hover {
      text-decoration: underline;
  }
</style>
  
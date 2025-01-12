<script lang="ts">
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';
    import { formatDate } from '$lib';

    let movies: Array<any> = [];
    let comments: Array<any> = [];
    let error: string = '';

    // Fetch movie and comments data on page load
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

        try {
            const response = await fetch('/api/comments');
            const data = await response.json();

            if (data.success) {
                comments = data.data;
            } else {
                error = 'Failed to load comment data.';
            }
        } catch (err) {
            console.error(err);
            error = 'Something went wrong while fetching the comment data.';
        }
    });
  
    // Add Movie form data and handling
    let showForm = false;
    let newMovie = {
      Title: '',
      Date: '',
      PosterUrl: '',
      MenuUrl: ''
    };

    const handleAddMovie = async () => {
      const formattedDate = new Date(newMovie.Date).toISOString();
      try {
        const response = await fetch('/api/movie', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            title: newMovie.Title,
            date: formattedDate,
            poster_url: newMovie.PosterUrl,
            menu_url: newMovie.MenuUrl
          }),
        });

        const result = await response.json();
        if (result.success) {
          window.location.reload();
        } else {
          error = 'Failed to add movie.';
        }
      } catch (err) {
        console.error(err);
        error = 'Something went wrong while adding the movie.';
      }
    };

    // Delete movie handler
    const deleteMovie = async (movieId: string) => {
      const confirmation = confirm('Are you sure you want to delete this movie?');

      if (confirmation) {
        try {
          const response = await fetch(`/api/movie/${movieId}`, {
            method: 'DELETE',
          });

          if (response.ok) {
            window.location.reload();
          } else {
            console.error('Failed to delete movie');
          }
        } catch (err) {
          console.error('Error during movie deletion:', err);
        }
      }
    };

    // Delete comment handler
    const deleteComment = async (commentId: string) => {
      const confirmation = confirm('Are you sure you want to delete this comment?');

      if (confirmation) {
        try {
          const response = await fetch(`/api/comment/${commentId}`, {
            method: 'DELETE',
          });

          if (response.ok) {
            window.location.reload();
          } else {
            console.error('Failed to delete comment');
          }
        } catch (err) {
          console.error('Error during comment deletion:', err);
        }
      }
    };

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
  
<h1>What's good, Golden Arm operator.</h1>

<button class="logout-button" on:click={logout}>Logout</button>

<!-- Display error message if data fetching fails -->
{#if error}
  <p style="color: red;">{error}</p>
{/if}

<!-- Movie Table -->
<h2>Movies</h2>
{#if movies.length > 0}
<table>
    <thead>
        <tr>
            <th>ID</th>
            <th>Title</th>
            <th>Screening Date</th>
            <th>Assets</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>
        {#each movies as movie (movie.ID)}
            <tr>
                <td>{movie.ID}</td>
                <td>{movie.Title}</td>
                <td>{formatDate(movie.Date)}</td>
                <td>
                  <a href={movie.PosterURL} target="_blank">Poster</a>,
                  <a href={movie.MenuURL} target="_blank">Menu</a>
                </td>
                <td>
                    <a href={`/admin/dashboard/reservations/${movie.ID}`} style="margin-left: 10px;">Reservations</a>
                    <button on:click={() => deleteMovie(movie.ID)} style="color: red; border: none; background: none; cursor: pointer;">X</button>
                </td>
            </tr>
        {/each}
    </tbody>
</table>
{:else}
  <p>No movies available at the moment. Add some!</p>
{/if}
<br>
<!-- Add Movie Button -->
<button on:click={() => showForm = !showForm}>Add Movie</button>

<!-- Add Movie Form Popup -->
{#if showForm}
  <div class="form-popup">
    <h2>Add Movie</h2>
    <form on:submit|preventDefault={handleAddMovie}>
      <label for="title">Title:</label>
      <input type="text" id="title" bind:value={newMovie.Title} required />

      <label for="date">Date:</label>
      <input type="datetime-local" id="date" bind:value={newMovie.Date} required />

      <label for="posterUrl">Poster URL:</label>
      <input type="url" id="posterUrl" bind:value={newMovie.PosterUrl} required />

      <label for="menuUrl">Menu URL:</label>
      <input type="url" id="menuUrl" bind:value={newMovie.MenuUrl} required />

      <button type="submit">Submit</button>
      <button type="button" on:click={() => showForm = false}>Cancel</button>
    </form>
  </div>
{/if}
<br><br>
<!-- Comments Table -->
<h2>Comments</h2>
{#if comments.length > 0}
<table>
    <thead>
        <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Email</th>
            <th>Comment</th>
            <th>Date</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>
        {#each comments as comment (comment.ID)}
            <tr>
                <td>{comment.ID}</td>
                <td>{comment.Name}</td>
                <td>{comment.Email}</td>
                <td>{comment.Comment}</td>
                <td>{formatDate(comment.Date)}</td>
                <td>
                    <button on:click={() => deleteComment(comment.ID)} style="color: red; border: none; background: none; cursor: pointer;">X</button>
                </td>
            </tr>
        {/each}
    </tbody>
</table>
{:else}
  <p>No comments found.</p>
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

  .form-popup {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background-color: white;
    padding: 20px;
    border: 1px solid #ddd;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    z-index: 100;
    width: 300px;
  }

  .form-popup input,
  .form-popup button {
    width: 100%;
    margin-bottom: 10px;
  }

  .form-popup button {
    cursor: pointer;
    background-color: #3498db;
    color: white;
    border: none;
    padding: 10px;
  }

  .logout-button {
    position: absolute;
    top: 20px;
    right: 20px;
    padding: 10px 20px;
  }

  button {
    padding: 10px 20px;
    font-size: 14px;
  }

</style>
  
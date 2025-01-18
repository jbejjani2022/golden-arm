<script lang="ts">
    import { onMount } from 'svelte';
    import { formatDate } from '$lib';

    let movies: Array<any> = [];
    let comments: Array<any> = [];
    let calendars: Array<any> = [];
    let emailList: Array<any> = [];
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

        try {
            const response = await fetch('/api/calendar/all');
            const data = await response.json();

            if (data.success) {
                calendars = data.data;
            } else {
                error = 'Failed to load calendar data.';
            }
        } catch (err) {
            console.error(err);
            error = 'Something went wrong while fetching the calendar data.';
        }

        try {
            const response = await fetch('/api/emails');
            const data = await response.json();

            if (data.success) {
                emailList = data.data;
            } else {
                error = 'Failed to load email data.';
            }
        } catch (err) {
            console.error(err);
            error = 'Something went wrong while fetching the email data.';
        }
    });
  
    // Add Movie form data and handling
    let showForm = false;
    let newMovie = {
      Title: '',
      Date: '',
      Runtime: 0,
      PosterFile: null as File | null,
      MenuFile: null as File | null
    };

    const handleAddMovie = async () => {
      const formData = new FormData();
      formData.append('title', newMovie.Title);
      formData.append('date', new Date(newMovie.Date).toISOString());
      formData.append('runtime', newMovie.Runtime.toString());
      if (newMovie.PosterFile) formData.append('poster', newMovie.PosterFile);
      if (newMovie.MenuFile) formData.append('menu', newMovie.MenuFile);

      try {
        const response = await fetch('/api/movie', {
          method: 'POST',
          body: formData,
        });

        const result = await response.json();
        if (result.success) {
          console.log('Movie added successfully.');
          window.location.reload();
        } else {
          error = 'Failed to add movie.';
        }
      } catch (err) {
        console.error(err);
        error = 'Something went wrong while adding the movie.';
      }
    };

    // Add Calendar form data and handling
    let showCalendarForm = false;
    let newCalendar = {
      StartDate: '',
      EndDate: '',
      CalendarFile: null as File | null
    };

    function validateDateRange(startDate: Date, endDate: Date): boolean {
      return startDate < endDate;
    }

    const handleAddCalendar = async () => {
      // Create dates in local timezone
      const start = new Date(newCalendar.StartDate + 'T00:00:00');  // sets time to start of day
      const end = new Date(newCalendar.EndDate + 'T23:59:59.999');  // sets time to end of day
      
      if (!validateDateRange(start, end)) {
        alert('End date must be after start date');
        return;
      }

      newCalendar.StartDate = start.toISOString();
      newCalendar.EndDate = end.toISOString();
      const formData = new FormData();
      formData.append('start_date', newCalendar.StartDate);
      formData.append('end_date', newCalendar.EndDate);
      if (newCalendar.CalendarFile) formData.append('image', newCalendar.CalendarFile);

      try {
        const response = await fetch('/api/calendar', {
          method: 'POST',
          body: formData,
        });

        const result = await response.json();
        if (result.success) {
          console.log('Calendar added successfully.');
          window.location.reload();
        } else {
          error = result.error;
          alert("Error: " + error);
        }
      } catch (err) {
        console.error(err);
        error = 'Something went wrong while adding the calendar.';
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

    // Delete calendar handler
    const deleteCalendar = async (calendarId: string) => {
      const confirmation = confirm('Are you sure you want to delete this calendar?');

      if (confirmation) {
        try {
          const response = await fetch(`/api/calendar/${calendarId}`, {
            method: 'DELETE',
          });

          if (response.ok) {
            window.location.reload();
          } else {
            console.error('Failed to delete calendar');
          }
        } catch (err) {
          console.error('Error during calendar deletion:', err);
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

    // Function to copy the email list to clipboard
    const copyEmailList = async () => {
      if (emailList.length === 0) {
        alert('No emails available to copy.');
        return;
      }
      try {
        await navigator.clipboard.writeText(emailList.join(', '));
        alert('Email list copied to clipboard!');
      } catch (err) {
        console.error('Failed to copy email list:', err);
        alert('Failed to copy email list. Please try again.');
      }
    };

</script>
  
<h1>What's good, Golden Arm operator.</h1>
<button on:click={copyEmailList} style="padding: 10px 20px; cursor: pointer;" title="Copy all unique emails from reservations and comments to clipboard.">Get Email List</button>

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
            <th>Runtime</th>
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
                <td>{movie.Runtime}</td>
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
  <p>No movies found. Add some!</p>
{/if}
<br>
<!-- Add Movie Button -->
<button on:click={() => showForm = !showForm}>Add Movie</button>

<!-- Add Movie Form Popup -->
{#if showForm}
<div class="modal">
  <div class="modal-content">
    <h2>Add Movie</h2>
    <form on:submit|preventDefault={handleAddMovie}>
      <div class="form-group">
        <label for="title">Title:</label>
        <input type="text" id="title" bind:value={newMovie.Title} required />
      </div>
      <div class="form-group">
        <label for="date">Date:</label>
        <input type="datetime-local" id="date" bind:value={newMovie.Date} required />
      </div>
      <div class="form-group">
        <label for="runtime">Runtime (minutes):</label>
        <input type="number" id="runtime" bind:value={newMovie.Runtime} required min="1"/>
      </div>
      <div class="form-group">
        <label for="posterFile">Poster Image:</label>
        <input type="file" id="posterFile" accept="image/jpg, image/png" on:change={(event) => newMovie.PosterFile = (event.target as HTMLInputElement).files?.[0] || null} required />
      </div>
      <div class="form-group">
        <label for="menuFile">Menu Image:</label>
        <input type="file" id="menuFile" accept="image/jpg, image/png" on:change={(event) => newMovie.MenuFile = (event.target as HTMLInputElement).files?.[0] || null} required />
      </div>      
      <button type="submit">Submit</button>
      <button type="button" class="cancel-button" on:click={() => showForm = false}>Cancel</button>
    </form>
  </div>
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
<br><br>
<!-- Calendars Table -->
<h2>Calendars</h2>
{#if calendars.length > 0}
<table>
    <thead>
        <tr>
            <th>ID</th>
            <th>Date Range</th>
            <th>Date Added</th>
            <th>Assets</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>
        {#each calendars as calendar (calendar.ID)}
            <tr>
                <td>{calendar.ID}</td>
                <td>{formatDate(calendar.StartDate)} - {formatDate(calendar.EndDate)}</td>
                <td>{formatDate(calendar.Date)}</td>
                <td>
                  <a href={calendar.ImageURL} target="_blank">Image</a>
                </td>
                <td>
                    <button on:click={() => deleteCalendar(calendar.ID)} style="color: red; border: none; background: none; cursor: pointer;">X</button>
                </td>
            </tr>
        {/each}
    </tbody>
</table>
{:else}
  <p>No calendars found. Add some.</p>
{/if}
<br>
<!-- Add Calendar Button -->
<button on:click={() => showCalendarForm = !showCalendarForm}>Add Calendar</button>

<!-- Add Calendar Form Popup -->
{#if showCalendarForm}
<div class="modal">
  <div class="modal-content">
    <h2>Add Calendar</h2>
    <form on:submit|preventDefault={handleAddCalendar}>
      <div class="form-group">
        <label for="dateRange">Date Range:</label>
        <div class="date-range-container">
          <input 
            type="date" 
            id="startDate"
            bind:value={newCalendar.StartDate} 
            required 
          />
          <span>to</span>
          <input 
            type="date" 
            id="endDate"
            bind:value={newCalendar.EndDate} 
            required 
          />
        </div>
      </div>
      <div class="form-group">
        <label for="File">Calendar Image:</label>
        <input type="file" id="calendarFile" accept="image/jpg, image/png" on:change={(event) => newCalendar.CalendarFile = (event.target as HTMLInputElement).files?.[0] || null} required />
      </div>   
      <button type="submit">Submit</button>
      <button type="button" class="cancel-button" on:click={() => showCalendarForm = false}>Cancel</button>
    </form>
  </div>
</div>
{/if}

<style>
  table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 20px;
      align-items: center;
  }

  th, td {
      padding: 12px;
      text-align: left;
      border: 1px solid #ddd;
  }

  th {
      background-color: #f4f4f4;
  }

  .date-range-container {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .date-range-container span {
    padding: 0 0.5rem;
  }
</style>
  
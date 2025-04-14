<script lang="ts">
    // Add any global logic or imports here
    import { goto } from '$app/navigation';
	  import { onMount } from 'svelte';

    let emailList: Array<any> = [];

    onMount(async () => {
      try {
            const response = await fetch('/api/emails');
            const data = await response.json();

            if (data.success) {
                emailList = data.data;
            } else {
              console.error("Failed to load email data.");
            }
        } catch (err) {
            console.error(err);
        }
    })

    const logout = async () => {
      try {
        const response = await fetch('/api/admin/logout', {
          method: 'POST',
        });

        if (response.ok) {
          console.log('Logout successful');
          goto('/admin');
        } else {
          console.error('Logout failed');
        }
      } catch (err) {
        console.error('Error during logout:', err);
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
  
<main>
  <header class="admin-header">
    <nav class="admin-nav">
      <button class="email-button" on:click={copyEmailList} style="padding: 10px 20px; cursor: pointer;" title="Copy all unique emails from reservations and comments to clipboard.">Get Full Email List</button>
      <div class="nav-wrapper">
        <a href="/admin/dashboard">Dashboard</a>
        <a href="/admin/dashboard/merch">Merch</a>
        <a href="/admin" on:click|preventDefault={logout}>Logout</a>
      </div>
    </nav>
  </header>
  <slot></slot>
</main>

<style>
  /* Styling the header */
  .admin-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px;
    margin-bottom: 2rem;
    border-bottom: 1px solid #ddd;
  }

  .email-button {
    padding: 10px 20px;
    cursor: pointer;
    text-align: center;
    border: none;
    border-radius: 0px;
    font-weight: bold;
  }

  .nav-wrapper {
    display: flex;
    gap: 20px;
    justify-content: flex-end;
  }

  .admin-nav {
    display: flex;
    width: 100%;
    justify-content: space-between;
    align-items: center;
  }

  .admin-header a {
    font-weight: bold;
    color: #007bff; /* Adjust color to match your theme */
    text-decoration: none;
  }
</style>

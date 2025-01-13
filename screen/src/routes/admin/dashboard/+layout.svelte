<script lang="ts">
    // Add any global logic or imports here
    import { goto } from '$app/navigation';

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
</script>
  
<main>
  <header class="admin-header">
    <nav>
      <a href="/admin/dashboard">Dashboard</a>
      <a href="/admin" on:click|preventDefault={logout}>Logout</a>
    </nav>
  </header>
  <slot></slot>
</main>

<style>
  /* Styling the header */
  .admin-header {
    display: flex;
    justify-content: flex-end;
    padding: 10px;
    border-bottom: 1px solid #ddd;
  }

  .admin-header nav {
    display: flex;
    gap: 20px;
  }

  .admin-header a {
    font-weight: bold;
  }
</style>

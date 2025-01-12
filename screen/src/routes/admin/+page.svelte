<script lang="ts">
    import { goto } from '$app/navigation';
  
    let passkey = '';
    let error = '';
  
    const handleLogin = async () => {
      error = ''; // Reset error message
      try {
        const response = await fetch('/api/admin/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ passkey })
        });
  
        const result = await response.json();
  
        if (result.success) {
          goto('/admin/dashboard');
        } else {
          error = 'Are you really a Golden Arm operator?';
        }
      } catch (err) {
        console.error(err);
        error = 'Something went wrong. Please try again.';
      }
    };
  </script>
  
  <h1>The Golden Arm Operator Room</h1>
  
  <form on:submit|preventDefault={handleLogin}>
    <input
      type="password"
      id="passkey"
      bind:value={passkey}
      placeholder="Passkey"
      required
    />
    <button type="submit">Login</button>
  </form>
  
  {#if error}
    <p style="color: red;">{error}</p>
  {/if}
  
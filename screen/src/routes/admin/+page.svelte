<script lang="ts">
    import { goto } from '$app/navigation';
    import { apiBaseUrl } from '$lib/api';
  
    let passkey = '';
    let error = '';
  
    const handleLogin = async () => {
      error = '';
      try {
        const response = await fetch(`${apiBaseUrl}/admin/login`, {
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
  
  <div class="container">
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
  </div>

  <style>
    .container {
      text-align: center;
      padding: 20px;
      border-radius: 8px;
    }
  
    input {
      padding: 10px;
      margin: 10px 0;
      width: 200px;
      font-size: 16px;
      border: 1px solid #ccc;
      border-radius: 4px;
    }
  </style>
  
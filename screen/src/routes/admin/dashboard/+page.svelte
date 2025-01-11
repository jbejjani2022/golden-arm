<script lang="ts">
    import { goto } from '$app/navigation';
  
    const logout = () => {
      document.cookie = 'isAdmin=; Max-Age=0'; // Clear the cookie
      goto('/admin');
    };
  </script>
  
  <h1>Hey Golden Arm Operator</h1>
  
  <button on:click={logout}>Logout</button>
  
  <style>
    button {
      padding: 8px 16px;
      font-size: 16px;
      cursor: pointer;
      background-color: #e74c3c;
      color: #fff;
      border: none;
      border-radius: 5px;
    }
  </style>
  
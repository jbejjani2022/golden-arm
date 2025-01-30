import { redirect } from '@sveltejs/kit';
import { apiBaseUrl } from '$lib/api';

// Ensures admin dashboard page is only accessible if a valid session exists
export const load = async ({ fetch }) => {
  const response = await fetch(`${apiBaseUrl}/admin/validate-session`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include'  // Ensures cookies are sent
  });

  const result = await response.json();

  if (!result.valid) {
    throw redirect(302, '/admin');
  }

  return {
    isAdmin: true
  };
};

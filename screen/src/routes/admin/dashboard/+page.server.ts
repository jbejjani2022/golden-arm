import { redirect } from '@sveltejs/kit';

// Ensures admin dashboard page is only accessible if a valid session exists
export const load = async ({ cookies, fetch }) => {
  const response = await fetch('/api/admin/validate-session', {
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

import { redirect } from '@sveltejs/kit';

// Ensures admin dashboard page is only accessible if a valid session exists
export const load = async ({ cookies, fetch }) => {
  const sessionToken = cookies.get('sessionToken');

  const response = await fetch('/api/admin/validate-session', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ sessionToken })
  });

  const result = await response.json();

  if (!result.valid) {
    throw redirect(302, '/admin');
  }

  return {
    isAdmin: true
  };
};

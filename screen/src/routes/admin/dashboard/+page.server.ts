import { redirect } from '@sveltejs/kit';

export const load = async ({ cookies }) => {
  const isAdmin = cookies.get('isAdmin');

  if (!isAdmin) {
    // Redirect to the admin login page if the session is not valid
    throw redirect(302, '/admin');
  }

  return {
    isAdmin: true
  };
};

const BASE = '/api/v1';

function getToken() {
  return localStorage.getItem('access_token');
}

async function req(method, path, body) {
  const headers = { 'Content-Type': 'application/json' };
  const token = getToken();
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(BASE + path, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Request failed');
  return data;
}

export const api = {
  // Auth
  register: (body) => req('POST', '/auth/register', body),
  login:    (body) => req('POST', '/auth/login', body),
  logout:   ()     => req('POST', '/auth/logout'),
  profile:  ()     => req('GET',  '/auth/profile'),
  updateProfile: (body) => req('PUT', '/auth/profile', body),

  // Movies
  listMovies:  (genre) => req('GET', `/movies${genre ? `?genre=${genre}` : ''}`),
  getMovie:    (id)    => req('GET', `/movies/${id}`),
  listShowtimes: (movieId, date) =>
    req('GET', `/movies/${movieId}/showtimes${date ? `?date=${date}` : ''}`),
  getSeats:    (showtimeId) => req('GET', `/showtimes/${showtimeId}/seats`),

  // Bookings
  createBooking: (body)      => req('POST', '/bookings', body),
  listBookings:  ()          => req('GET',  '/bookings'),
  getBooking:    (id)        => req('GET',  `/bookings/${id}`),
  cancelBooking: (id)        => req('POST', `/bookings/${id}/cancel`),

  // Payments
  listPayments: ()  => req('GET', '/payments'),
  getReceipt:   (id) => req('GET', `/payments/${id}/receipt`),
};

import { writable } from 'svelte/store';

export const email = writable( {
  to: '',
  subject: '',
  body: ''
} );


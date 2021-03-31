export const AUTH_API_URL =
  process.env.NEXT_PUBLIC_AUTH_API_URL ||
  'http://account.local.eitan-flash.com:4000';
export const GRAPHQL_API_URL =
  process.env.NEXT_PUBLIC_GRAPHQL_API_URL ||
  'http://api.local.eitan-flash.com:5000';

export const GOOGLE_LOGIN_URL = `${AUTH_API_URL}/auth/google`;

export const AUTH_API_URL =
  process.env.NEXT_PUBLIC_AUTH_API_URL ||
  'http://account.local.eitan-flash.com:4000';
export const GRAPHQL_API_URL =
  process.env.NEXT_PUBLIC_GRAPHQL_API_URL ||
  'http://api.local.eitan-flash.com:5000';

export const CREATE_EMAIL_CONFIRMATION_URL = `${AUTH_API_URL}/auth/email/confirmations`;
export const EMAIL_LOGIN_URL = `${AUTH_API_URL}/auth/email/login`;
export const EMAIL_SIGNUP_URL = `${AUTH_API_URL}/auth/email/sign_up`;
export const GOOGLE_LOGIN_URL = `${AUTH_API_URL}/auth/oauth/google`;
export const LOGOUT_URL = `${AUTH_API_URL}/auth/logout`;

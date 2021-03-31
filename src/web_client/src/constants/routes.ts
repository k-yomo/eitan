export const routes = {
  login: () => '/login',
  loginWithOriginalPath: (originalPath: string) =>
    `/login?originalUrl=${originalPath}`,
  signUp: () => '/sign_up',
  passwordReset: () => '/password_reset',
  accountSettings: () => 'account_settings',
};

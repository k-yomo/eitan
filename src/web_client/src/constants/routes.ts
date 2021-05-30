export const routes = {
  home: () => '/',
  login: () => '/login',
  loginWithOriginalPath: (originalPath: string) =>
    `/login?originalUrl=${originalPath}`,
  signUp: () => '/sign_up',
  passwordReset: () => '/password_reset',
  userSettings: () => 'user_settings',
  waitingRoom: () => 'waiting_room',
};

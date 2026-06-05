import { defineStore } from "pinia";

const TOKEN_KEY = "edu-admin-token";
const USER_NAME_KEY = "edu-admin-user-name";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    token: localStorage.getItem(TOKEN_KEY) ?? "",
    userName: localStorage.getItem(USER_NAME_KEY) ?? "",
  }),
  getters: {
    isLoggedIn: (state) => state.token.length > 0,
  },
  actions: {
    setToken(token: string) {
      this.token = token;
      localStorage.setItem(TOKEN_KEY, token);
    },
    setUserName(userName: string) {
      this.userName = userName;
      localStorage.setItem(USER_NAME_KEY, userName);
    },
    setSession(token: string, userName: string) {
      this.setToken(token);
      this.setUserName(userName);
    },
    clearSession() {
      this.token = "";
      this.userName = "";
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_NAME_KEY);
    },
  },
});

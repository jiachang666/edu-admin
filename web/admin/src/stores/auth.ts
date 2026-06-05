import { defineStore } from "pinia";
import {
  fetchMe,
  fetchMyPermissions,
  type CurrentUser,
  type PermissionGroup,
} from "../api/auth";

const TOKEN_KEY = "edu-admin-token";
const USER_NAME_KEY = "edu-admin-user-name";
const USER_ID_KEY = "edu-admin-user-id";
const USER_ROLES_KEY = "edu-admin-user-roles";
const USER_ROLE_NAMES_KEY = "edu-admin-user-role-names";
const USER_PERMISSIONS_KEY = "edu-admin-user-permissions";

function readStringArray(key: string) {
  const rawValue = localStorage.getItem(key);
  if (!rawValue) {
    return [] as string[];
  }

  try {
    const parsedValue = JSON.parse(rawValue);
    if (Array.isArray(parsedValue)) {
      return parsedValue.filter((item): item is string => typeof item === "string");
    }
  } catch (error) {
    console.warn(`failed to parse ${key}`, error);
  }

  return [] as string[];
}

export const useAuthStore = defineStore("auth", {
  state: () => ({
    token: localStorage.getItem(TOKEN_KEY) ?? "",
    userName: localStorage.getItem(USER_NAME_KEY) ?? "",
    userId: Number(localStorage.getItem(USER_ID_KEY) ?? 0),
    roles: readStringArray(USER_ROLES_KEY),
    roleNames: readStringArray(USER_ROLE_NAMES_KEY),
    permissions: readStringArray(USER_PERMISSIONS_KEY),
    permissionGroups: [] as PermissionGroup[],
  }),
  getters: {
    isLoggedIn: (state) => state.token.length > 0,
    primaryRole: (state) => state.roles[0] ?? "",
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
    setUserId(userId: number) {
      this.userId = userId;
      localStorage.setItem(USER_ID_KEY, String(userId));
    },
    setRoles(roles: string[]) {
      this.roles = roles;
      localStorage.setItem(USER_ROLES_KEY, JSON.stringify(roles));
    },
    setRoleNames(roleNames: string[]) {
      this.roleNames = roleNames;
      localStorage.setItem(USER_ROLE_NAMES_KEY, JSON.stringify(roleNames));
    },
    setPermissions(permissions: string[]) {
      this.permissions = permissions;
      localStorage.setItem(USER_PERMISSIONS_KEY, JSON.stringify(permissions));
    },
    setPermissionGroups(permissionGroups: PermissionGroup[]) {
      this.permissionGroups = permissionGroups;
    },
    hasPermission(permission: string) {
      return this.permissions.includes(permission);
    },
    setSession(
      token: string,
      user: { id: number; displayName: string },
      roles: string[],
      roleNames: string[],
      permissions: string[],
    ) {
      this.setToken(token);
      this.setUserId(user.id);
      this.setUserName(user.displayName);
      this.setRoles(roles);
      this.setRoleNames(roleNames);
      this.setPermissions(permissions);
    },
    setCurrentUser(user: CurrentUser) {
      this.setUserId(user.id);
      this.setUserName(user.displayName);
      this.setRoles(user.roles);
      this.setRoleNames(user.roleNames);
    },
    async hydrateSession() {
      if (!this.token) {
        return;
      }

      const [userResult, permissionResult] = await Promise.all([fetchMe(), fetchMyPermissions()]);
      this.setCurrentUser(userResult);
      this.setPermissions(permissionResult.permissions);
      this.setPermissionGroups(permissionResult.permissionGroups);
    },
    clearSession() {
      this.token = "";
      this.userName = "";
      this.userId = 0;
      this.roles = [];
      this.roleNames = [];
      this.permissions = [];
      this.permissionGroups = [];
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_NAME_KEY);
      localStorage.removeItem(USER_ID_KEY);
      localStorage.removeItem(USER_ROLES_KEY);
      localStorage.removeItem(USER_ROLE_NAMES_KEY);
      localStorage.removeItem(USER_PERMISSIONS_KEY);
    },
  },
});

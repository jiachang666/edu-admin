import http from "./http";

type ApiEnvelope<T> = {
  code: number;
  message: string;
  data: T;
  requestId: string;
};

export type PermissionDefinition = {
  code: string;
  label: string;
  description: string;
};

export type PermissionGroup = {
  key: string;
  label: string;
  description: string;
  permissions: PermissionDefinition[];
};

export type LoginResult = {
  accessToken: string;
  expiresIn: number;
  user: {
    id: number;
    username: string;
    displayName: string;
  };
  roles: string[];
  roleNames: string[];
  permissions: string[];
};

export type CurrentUser = {
  id: number;
  username: string;
  displayName: string;
  mobile: string;
  status: string;
  roles: string[];
  roleNames: string[];
  lastLoginAt: string;
};

export type PermissionResult = {
  permissions: string[];
  permissionGroups: PermissionGroup[];
};

async function unwrap<T>(request: Promise<{ data: ApiEnvelope<T> }>) {
  const response = await request;
  return response.data.data;
}

export function login(payload: { username: string; password: string }) {
  return unwrap<LoginResult>(http.post("/auth/login", payload));
}

export function fetchMe() {
  return unwrap<CurrentUser>(http.get("/auth/me"));
}

export function fetchMyPermissions() {
  return unwrap<PermissionResult>(http.get("/auth/me/permissions"));
}

export function logout() {
  return unwrap<{ loggedOut: boolean }>(http.post("/auth/logout"));
}

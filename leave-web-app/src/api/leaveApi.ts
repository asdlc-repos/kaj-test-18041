import { api } from './index';

const BASE = '/api/leave';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: 'employee' | 'manager';
  token: string;
}

export interface LeaveBalance {
  annual: number;
  sick: number;
  casual: number;
  [key: string]: number;
}

export interface LeaveRequest {
  id: string;
  employeeId: string;
  employeeName?: string;
  type: string;
  startDate: string;
  endDate: string;
  reason: string;
  status: 'pending' | 'approved' | 'rejected' | 'cancelled';
  createdAt: string;
  comment?: string;
}

export interface CreateLeaveRequest {
  type: string;
  startDate: string;
  endDate: string;
  reason: string;
}

export const leaveApi = {
  login: (data: LoginRequest) =>
    api.post<User>(`${BASE}/auth/login`, data),

  getBalance: () =>
    api.get<LeaveBalance>(`${BASE}/balance`),

  getMyRequests: () =>
    api.get<LeaveRequest[]>(`${BASE}/requests`),

  createRequest: (data: CreateLeaveRequest) =>
    api.post<LeaveRequest>(`${BASE}/requests`, data),

  cancelRequest: (id: string) =>
    api.patch<LeaveRequest>(`${BASE}/requests/${id}/cancel`, {}),
};

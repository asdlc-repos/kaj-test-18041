import { api } from './index';
import type { LeaveRequest } from './leaveApi';

const BASE = '/api/manager';

export interface ReviewRequest {
  comment?: string;
}

export const managerApi = {
  getPendingRequests: () =>
    api.get<LeaveRequest[]>(`${BASE}/requests`),

  approveRequest: (id: string, data: ReviewRequest) =>
    api.post<LeaveRequest>(`${BASE}/requests/${id}/approve`, data),

  rejectRequest: (id: string, data: ReviewRequest) =>
    api.post<LeaveRequest>(`${BASE}/requests/${id}/reject`, data),
};

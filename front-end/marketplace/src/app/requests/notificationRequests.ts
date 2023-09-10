import request from '../helpers/requestHelper';

export const getNotifications = () => request('/v1/users/me/notifications');
export const viewNotificationRequest = ({ notificationId }: {notificationId: string}) =>
  request(`/v1/users/me/notifications/${notificationId}/view`, 'patch');

export const deleteNotificationRequest = ({ notificationId }: {notificationId: string}) =>
  request(`/v1/users/me/notifications/${notificationId}`, 'delete');

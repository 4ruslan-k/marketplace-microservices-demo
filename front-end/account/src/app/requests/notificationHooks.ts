import { useMutation, useQuery } from '@tanstack/react-query';
import { getNotifications, viewNotificationRequest, deleteNotificationRequest } from './notificationRequests';

interface Notification {
  id: string;
  title: string;
  message: string;
  notificationType: string;
  viewedAt?: string;
  createdAt: string;
}

interface NotificationsResponse {
  notifications: Notification[];
}

export const useFetchNotifications = ({ user }) => {
  const {
    isLoading: isLoadingUser,
    data,
    error: fetchNotificationsError,
    refetch: refetchNotifications,
    remove: removeNotificationsData,
  } = useQuery<NotificationsResponse, unknown, NotificationsResponse>({
    queryKey: ['notifications', user],
    queryFn: () => {
      if (user) return getNotifications();
      return {};
    },
  });
  const notifications: Notification[] = data?.notifications;
  return {
    isLoadingUser,
    notifications,
    fetchNotificationsError,
    refetchNotifications,
    removeNotificationsData,
  };
};

export const useWatchNotification = ({
  onSuccess,
  onError,
}: {
  onSuccess?: (() => void) | undefined;
  onError: (error) => void;
}) => {
  const {
    isLoading: isLoadingWatchNotificataion,
    error: watchNotificationError,
    mutate: viewNotification,
  } = useMutation({
    mutationFn: ({ notificationId }: any) => viewNotificationRequest({ notificationId }),
    onError,
    onSuccess,
  });
  return {
    isLoadingWatchNotificataion,
    watchNotificationError,
    viewNotification,
  };
};

export const useDeleteNotification = ({ onSuccess, onError }: { onSuccess?: () => void; onError: (error) => void }) => {
  const {
    isLoading: isLoadingDeleteNotificataion,
    error: deleteNotificationError,
    mutate: deleteNotification,
  } = useMutation({
    mutationFn: ({ notificationId }: any) => deleteNotificationRequest({ notificationId }),
    onError,
    onSuccess,
  });
  return {
    isLoadingDeleteNotificataion,
    deleteNotificationError,
    deleteNotification,
  };
};

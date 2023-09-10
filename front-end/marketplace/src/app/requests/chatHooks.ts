import { useQuery, useMutation } from '@tanstack/react-query';
import {
  getUser,
} from './userRequests';
import { getMessages } from './chatRequests';

export const useFetchMessages = () => {
  const {
    isLoading: isLoadingUser,
    data,
    error: fetchMessagesError,
    refetch: refetchMessages,
  } = useQuery({
    queryKey: ['messages'],
    queryFn: getMessages,
  });
  const messages = data?.items;
  return {
    isLoadingUser,
    messages,
    fetchMessagesError,
    refetchMessages,
  };
};

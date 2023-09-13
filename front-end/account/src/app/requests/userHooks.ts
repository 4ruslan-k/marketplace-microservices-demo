import { useQuery, useMutation } from '@tanstack/react-query';
import {
  getUser,
  login,
  register,
  logout,
  sendMfaCode,
  generateTotpMfaSetupRequest,
  disableTotpMfaRequest,
  enableTotpMfaRequest,
} from './userRequests';

export const useFetchUser = () => {
  const {
    isLoading: isLoadingUser,
    data,
    error: fetchUserError,
    refetch: refetchUser,
    remove: removeUserData,
  } = useQuery({
    queryKey: ['user'],
    queryFn: getUser,
  });
  const user = data?.user;
  return {
    isLoadingUser,
    user,
    fetchUserError,
    refetchUser,
    removeUserData,
  };
};

export const useLoginUser = ({ onSuccess, onError }) => {
  const {
    isLoading: isLoadingLogin,
    error: loginError,
    mutate: loginUser,
  } = useMutation({
    mutationFn: ({ email, password }: any) => login(email, password),
    onError,
    onSuccess,
  });
  return {
    isLoadingLogin,
    loginError,
    loginUser,
  };
};

export const useSendMfaOtpCode = ({ onSuccess, onError }) => {
  const {
    isLoading: isLoadingSendingOtpCode,
    error: mfaSendCodeError,
    mutate: sendMfaOtpCode,
  } = useMutation({
    mutationFn: ({ passwordVerificationToken, otp }: any) => sendMfaCode({ passwordVerificationToken, otp }),
    onError,
    onSuccess,
  });
  return {
    isLoadingSendingOtpCode,
    mfaSendCodeError,
    sendMfaOtpCode,
  };
};

export const useGenerateTotpMfaSetup = ({ onSuccess, onError }) => {
  const {
    isLoading: isLoadingGenerateTotpMfaSetup,
    error: generateTotpMfaSetupError,
    mutate: generateTotpMfaSetup,
  } = useMutation({
    mutationFn: generateTotpMfaSetupRequest,
    onError,
    onSuccess,
  });
  return {
    isLoadingGenerateTotpMfaSetup,
    generateTotpMfaSetupError,
    generateTotpMfaSetup,
  };
};

export const useDisableTotpMfa = ({ onSuccess, onError }) => {
  const {
    isLoading: isLoadingDisableTotpMfa,
    error: disableTotpMfaError,
    mutate: disableTotpMfa,
  } = useMutation({
    mutationFn: disableTotpMfaRequest,
    onError,
    onSuccess,
  });
  return {
    isLoadingDisableTotpMfa,
    disableTotpMfaError,
    disableTotpMfa,
  };
};

export const useEnableTotpMfa = ({ onSuccess, onError }) => {
  const {
    isLoading: isLoadingEnableTotpMfa,
    error: enableTotpMfaError,
    mutate: enableTotpMfa,
  } = useMutation({
    mutationFn: enableTotpMfaRequest,
    onError,
    onSuccess,
  });
  return {
    isLoadingEnableTotpMfa,
    enableTotpMfaError,
    enableTotpMfa,
  };
};

export const useSignupUser = ({ onSuccess, onError }) => {
  const {
    isLoading: isLoadingSignup,
    error: signupError,
    mutate: signupUser,
  } = useMutation({
    mutationFn: ({ email, password, name }: any) => register(email, password, name),
    onError,
    onSuccess,
  });
  return {
    isLoadingSignup,
    signupError,
    signupUser,
  };
};

export const useLogoutUser = ({ onSuccess, onError }) => {
  const {
    isLoading: isLoadingLogout,
    error: logoutError,
    mutate: logoutUser,
  } = useMutation({
    mutationFn: logout,
    onError,
    onSuccess,
  });
  return {
    isLoadingLogout,
    logoutError,
    logoutUser,
  };
};

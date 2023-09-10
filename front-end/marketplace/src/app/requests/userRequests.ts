import request from '../helpers/requestHelper';

// eslint-disable-next-line import/prefer-default-export
export const login = (username: string, password: string) => request('/v1/auth/login', 'post', { username, password });
export const sendMfaCode = ({ otp, passwordVerificationToken }:{otp:string, passwordVerificationToken: string}) =>
  request('/v1/auth/login/mfa/totp', 'post', { tokenId: passwordVerificationToken, code: otp });
export const generateTotpMfaSetupRequest = () => request('/v1/auth/me/mfa/totp', 'put');
export const disableTotpMfaRequest = ({ code }:{code: string}) => request('/v1/auth/me/mfa/totp/disable', 'patch', { code });
export const enableTotpMfaRequest =  ({ code }:{code: string}) => request('/v1/auth/me/mfa/totp/enable', 'patch', { code });
export const logout = () => request('/v1/auth/logout');
export const getUser = () => request('/v1/users/me');
export const register = (email: string, password: string, name: string) => request('/v1/users', 'post', { email, password, name });

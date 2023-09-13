import request from '../helpers/requestHelper';

// eslint-disable-next-line import/prefer-default-export
export const getMessages = () => request('/v1/chat/messages');

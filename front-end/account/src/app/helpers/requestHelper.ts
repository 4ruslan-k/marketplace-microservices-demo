import axios from 'axios';
import config from '../config';

const request = async (url: string, method = 'GET', data = {}) => {
  const response = await axios({
    method,
    // give permission to include cookies on cross-origin requests
    withCredentials: true,
    url,
    data,
    baseURL: config.API_URL,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  });
  return response.data;
};

export default request;

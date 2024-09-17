import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

export const getBlockchain = async () => {
  const response = await axios.get(`${API_URL}/blockchain`);
  return response.data;
};

// ... (rest of the API functions)
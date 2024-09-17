import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

export const getBlockchain = async () => {
  const response = await axios.get(`${API_URL}/blockchain`);
  return response.data;
};

export const createWallet = async () => {
  const response = await axios.post(`${API_URL}/createWallet`);
  return response.data;
};

export const getBalance = async (address: string) => {
  const response = await axios.get(`${API_URL}/getBalance/${address}`);
  return response.data;
};

export const sendTransaction = async (from: string, to: string, amount: number) => {
  const response = await axios.post(`${API_URL}/sendTransaction`, { from, to, amount });
  return response.data;
};

export const forgeBlock = async () => {
  const response = await axios.post(`${API_URL}/forgeBlock`);
  return response.data;
};

export const stakeJUL = async (address: string, amount: number) => {
  const response = await axios.post(`${API_URL}/stakeJUL`, { address, amount });
  return response.data;
};

export const unstakeJUL = async (address: string, amount: number) => {
  const response = await axios.post(`${API_URL}/unstakeJUL`, { address, amount });
  return response.data;
};

export const getCommunityFund = async () => {
  const response = await axios.get(`${API_URL}/getCommunityFund`);
  return response.data;
};

export const getValidators = async () => {
  const response = await axios.get(`${API_URL}/validators`);
  return response.data;
};

export const purchaseJUL = async (address: string, usdAmount: number) => {
  const response = await axios.post(`${API_URL}/purchaseJUL`, { wallet: address, usdAmount });
  return response.data;
};
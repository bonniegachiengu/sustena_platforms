import React, { useEffect, useState } from 'react';
import './App.css';
import { getBlockchain, createWallet, getBalance, sendTransaction, forgeBlock, stakeJUL, unstakeJUL, getCommunityFund, getValidators, purchaseJUL } from './services/api';

interface Block {
  Index: number;
  Hash: string;
  Transactions: any[];
  Validator: string;
}

interface Validator {
  Address: string;
  Stake: number;
}

const App: React.FC = () => {
  const [blockchain, setBlockchain] = useState<Block[]>([]);
  const [walletAddress, setWalletAddress] = useState('');
  const [balance, setBalance] = useState(0);
  const [recipientAddress, setRecipientAddress] = useState('');
  const [amount, setAmount] = useState(0);
  const [stakeAmount, setStakeAmount] = useState(0);
  const [unstakeAmount, setUnstakeAmount] = useState(0);
  const [purchaseAmount, setPurchaseAmount] = useState(0);
  const [communityFund, setCommunityFund] = useState(0);
  const [validators, setValidators] = useState<Validator[]>([]);

  useEffect(() => {
    fetchBlockchain();
    fetchCommunityFund();
    fetchValidators();
  }, []);

  const fetchBlockchain = async () => {
    try {
      const data = await getBlockchain();
      setBlockchain(data);
    } catch (error) {
      console.error("Failed to fetch blockchain:", error);
    }
  };

  const fetchCommunityFund = async () => {
    try {
      const data = await getCommunityFund();
      setCommunityFund(data.balance);
    } catch (error) {
      console.error("Failed to fetch community fund:", error);
    }
  };

  const fetchValidators = async () => {
    try {
      const data = await getValidators();
      setValidators(data);
    } catch (error) {
      console.error("Failed to fetch validators:", error);
    }
  };

  const handleCreateWallet = async () => {
    try {
      const response = await createWallet();
      setWalletAddress(response.address);
      alert(`Wallet created with address: ${response.address}`);
    } catch (error) {
      console.error("Failed to create wallet:", error);
    }
  };

  const handleGetBalance = async () => {
    try {
      const response = await getBalance(walletAddress);
      setBalance(response.balance);
    } catch (error) {
      console.error("Failed to get balance:", error);
    }
  };

  const handleSendTransaction = async () => {
    try {
      await sendTransaction(walletAddress, recipientAddress, amount);
      alert("Transaction sent successfully");
      fetchBlockchain();
      handleGetBalance();
    } catch (error) {
      console.error("Failed to send transaction:", error);
    }
  };

  const handleForgeBlock = async () => {
    try {
      await forgeBlock();
      alert("New block forged successfully");
      fetchBlockchain();
    } catch (error) {
      console.error("Failed to forge block:", error);
    }
  };

  const handleStakeJUL = async () => {
    try {
      await stakeJUL(walletAddress, stakeAmount);
      alert(`Successfully staked ${stakeAmount} JUL`);
      handleGetBalance();
      fetchValidators();
    } catch (error) {
      console.error("Failed to stake JUL:", error);
    }
  };

  const handleUnstakeJUL = async () => {
    try {
      await unstakeJUL(walletAddress, unstakeAmount);
      alert(`Successfully unstaked ${unstakeAmount} JUL`);
      handleGetBalance();
      fetchValidators();
    } catch (error) {
      console.error("Failed to unstake JUL:", error);
    }
  };

  const handlePurchaseJUL = async () => {
    try {
      await purchaseJUL(walletAddress, purchaseAmount);
      alert(`Successfully purchased JUL for ${purchaseAmount} USD`);
      handleGetBalance();
    } catch (error) {
      console.error("Failed to purchase JUL:", error);
    }
  };

  return (
    <div className="App">
      <h1>Sustena Blockchain</h1>
      <div>
        <h2>Wallet</h2>
        <button onClick={handleCreateWallet}>Create Wallet</button>
        <p>Wallet Address: {walletAddress}</p>
        <button onClick={handleGetBalance}>Get Balance</button>
        <p>Balance: {balance} JUL</p>
      </div>
      <div>
        <h2>Purchase JUL</h2>
        <input 
          type="number" 
          placeholder="USD Amount" 
          value={purchaseAmount} 
          onChange={(e) => setPurchaseAmount(Number(e.target.value))} 
        />
        <button onClick={handlePurchaseJUL}>Purchase JUL</button>
      </div>
      <div>
        <h2>Send Transaction</h2>
        <input 
          type="text" 
          placeholder="Recipient Address" 
          value={recipientAddress} 
          onChange={(e) => setRecipientAddress(e.target.value)} 
        />
        <input 
          type="number" 
          placeholder="Amount" 
          value={amount} 
          onChange={(e) => setAmount(Number(e.target.value))} 
        />
        <button onClick={handleSendTransaction}>Send</button>
      </div>
      <div>
        <h2>Stake JUL</h2>
        <input 
          type="number" 
          placeholder="Stake Amount" 
          value={stakeAmount} 
          onChange={(e) => setStakeAmount(Number(e.target.value))} 
        />
        <button onClick={handleStakeJUL}>Stake</button>
      </div>
      <div>
        <h2>Unstake JUL</h2>
        <input 
          type="number" 
          placeholder="Unstake Amount" 
          value={unstakeAmount} 
          onChange={(e) => setUnstakeAmount(Number(e.target.value))} 
        />
        <button onClick={handleUnstakeJUL}>Unstake</button>
      </div>
      <div>
        <h2>Forge Block</h2>
        <button onClick={handleForgeBlock}>Forge New Block</button>
      </div>
      <div>
        <h2>Community Fund</h2>
        <p>Balance: {communityFund} JUL</p>
      </div>
      <div>
        <h2>Validators</h2>
        {validators.map((validator, index) => (
          <div key={index}>
            <p>Address: {validator.Address}</p>
            <p>Stake: {validator.Stake} JUL</p>
          </div>
        ))}
      </div>
      <div>
        <h2>Blockchain</h2>
        {blockchain.map((block: Block) => (
          <div key={block.Index}>
            <h3>Block {block.Index}</h3>
            <p>Hash: {block.Hash}</p>
            <p>Transactions: {block.Transactions.length}</p>
            <p>Validator: {block.Validator}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default App;
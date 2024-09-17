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
  console.log('App component rendering');
  const [blockchain, setBlockchain] = useState<Block[] | null>(null);
  const [walletAddress, setWalletAddress] = useState('');
  const [balance, setBalance] = useState(0);
  const [communityFund, setCommunityFund] = useState(0);
  const [validators, setValidators] = useState<Validator[]>([]);
  const [recipientAddress, setRecipientAddress] = useState('');
  const [amount, setAmount] = useState(0);
  const [stakeAmount, setStakeAmount] = useState(0);
  const [unstakeAmount, setUnstakeAmount] = useState(0);
  const [purchaseAmount, setPurchaseAmount] = useState(0);

  useEffect(() => {
    console.log('useEffect running');
    fetchBlockchain();
    fetchCommunityFund();
    fetchValidators();
  }, []);

  const fetchBlockchain = async () => {
    console.log('Fetching blockchain');
    try {
      const data = await getBlockchain();
      console.log('Blockchain data:', data);
      setBlockchain(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error("Failed to fetch blockchain:", error);
      setBlockchain([]);
    }
  };

  const fetchCommunityFund = async () => {
    try {
      const data = await getCommunityFund();
      setCommunityFund(data.balance || 0);
    } catch (error) {
      console.error("Failed to fetch community fund:", error);
    }
  };

  const fetchValidators = async () => {
    try {
      const data = await getValidators();
      setValidators(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error("Failed to fetch validators:", error);
      setValidators([]);
    }
  };

  const handleCreateWallet = async () => {
    try {
      const response = await createWallet();
      setWalletAddress(response.address || '');
      alert(`Wallet created with address: ${response.address}`);
    } catch (error) {
      console.error("Failed to create wallet:", error);
    }
  };

  const handleGetBalance = async () => {
    try {
      const response = await getBalance(walletAddress);
      setBalance(response.balance || 0);
    } catch (error) {
      console.error("Failed to get balance:", error);
    }
  };

  const handleSendTransaction = async () => {
    try {
      await sendTransaction(walletAddress, recipientAddress, amount);
      alert("Transaction sent successfully");
      handleGetBalance();
      fetchBlockchain();
    } catch (error) {
      console.error("Failed to send transaction:", error);
      alert("Failed to send transaction");
    }
  };

  const handleForgeBlock = async () => {
    try {
      await forgeBlock();
      alert("New block forged successfully");
      fetchBlockchain();
    } catch (error) {
      console.error("Failed to forge block:", error);
      alert("Failed to forge block");
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
      alert("Failed to stake JUL");
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
      alert("Failed to unstake JUL");
    }
  };

  const handlePurchaseJUL = async () => {
    try {
      await purchaseJUL(walletAddress, purchaseAmount);
      alert(`Successfully purchased JUL for ${purchaseAmount} USD`);
      handleGetBalance();
    } catch (error) {
      console.error("Failed to purchase JUL:", error);
      alert("Failed to purchase JUL");
    }
  };

  console.log('Rendering App component');
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
          value={purchaseAmount}
          onChange={(e) => setPurchaseAmount(Number(e.target.value))}
          placeholder="USD Amount"
        />
        <button onClick={handlePurchaseJUL}>Purchase JUL</button>
      </div>
      <div>
        <h2>Send Transaction</h2>
        <input 
          type="text" 
          value={recipientAddress}
          onChange={(e) => setRecipientAddress(e.target.value)}
          placeholder="Recipient Address"
        />
        <input 
          type="number" 
          value={amount}
          onChange={(e) => setAmount(Number(e.target.value))}
          placeholder="Amount"
        />
        <button onClick={handleSendTransaction}>Send</button>
      </div>
      <div>
        <h2>Stake JUL</h2>
        <input 
          type="number" 
          value={stakeAmount}
          onChange={(e) => setStakeAmount(Number(e.target.value))}
          placeholder="Stake Amount"
        />
        <button onClick={handleStakeJUL}>Stake</button>
      </div>
      <div>
        <h2>Unstake JUL</h2>
        <input 
          type="number" 
          value={unstakeAmount}
          onChange={(e) => setUnstakeAmount(Number(e.target.value))}
          placeholder="Unstake Amount"
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
        {validators.length === 0 ? (
          <p>No validators found</p>
        ) : (
          validators.map((validator, index) => (
            <div key={index}>
              <p>Address: {validator.Address}</p>
              <p>Stake: {validator.Stake} JUL</p>
            </div>
          ))
        )}
      </div>
      <div>
        <h2>Blockchain</h2>
        {blockchain === null ? (
          <p>Loading blockchain...</p>
        ) : blockchain.length === 0 ? (
          <p>No blocks in the blockchain</p>
        ) : (
          blockchain.map((block: Block) => (
            <div key={block.Index}>
              <h3>Block {block.Index}</h3>
              <p>Hash: {block.Hash}</p>
              <p>Transactions: {block.Transactions.length}</p>
              <p>Validator: {block.Validator}</p>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default App;

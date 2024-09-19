import React, { useState, useEffect } from 'react';
import { createWallet, getBalance, sendTransaction, purchaseJUL, stakeJUL, unstakeJUL } from '../services/api';

const Wallet: React.FC = () => {
  const [wallets, setWallets] = useState<string[]>([]);
  const [selectedWallet, setSelectedWallet] = useState('');
  const [balance, setBalance] = useState(0);
  const [recipientAddress, setRecipientAddress] = useState('');
  const [amount, setAmount] = useState(0);
  const [purchaseAmount, setPurchaseAmount] = useState(0);
  const [stakeAmount, setStakeAmount] = useState(0);
  const [unstakeAmount, setUnstakeAmount] = useState(0);
  const [senderWallet, setSenderWallet] = useState('');
  const [receiverWallet, setReceiverWallet] = useState('');

  useEffect(() => {
    // Load wallets from localStorage on component mount
    const savedWallets = localStorage.getItem('wallets');
    if (savedWallets) {
      setWallets(JSON.parse(savedWallets));
    }
  }, []);

  const handleCreateWallet = async () => {
    try {
      const response = await createWallet();
      const newAddress = response.address || '';
      setWallets(prevWallets => {
        const updatedWallets = [...prevWallets, newAddress];
        localStorage.setItem('wallets', JSON.stringify(updatedWallets));
        return updatedWallets;
      });
      setSelectedWallet(newAddress);
      alert(`Wallet created with address: ${newAddress}`);
    } catch (error) {
      console.error("Failed to create wallet:", error);
    }
  };

  const handleGetBalance = async () => {
    if (!selectedWallet) {
      alert("Please select a wallet first");
      return;
    }
    try {
      const response = await getBalance(selectedWallet);
      setBalance(response.balance || 0);
    } catch (error) {
      console.error("Failed to get balance:", error);
    }
  };

  const handleSendTransaction = async () => {
    if (!senderWallet) {
      alert("Please select a sender wallet");
      return;
    }
    try {
      await sendTransaction(senderWallet, receiverWallet, amount);
      alert("Transaction sent successfully");
      handleGetBalance();
    } catch (error) {
      console.error("Failed to send transaction:", error);
      alert("Failed to send transaction");
    }
  };

  const handlePurchaseJUL = async () => {
    if (!selectedWallet) {
      alert("Please select a wallet first");
      return;
    }
    try {
      const response = await purchaseJUL(selectedWallet, purchaseAmount);
      alert(`Successfully purchased ${response.julAmount} JUL for ${purchaseAmount} USD`);
      handleGetBalance();
    } catch (error) {
      console.error("Failed to purchase JUL:", error);
      alert("Failed to purchase JUL");
    }
  };

  const handleStakeJUL = async () => {
    if (!selectedWallet) {
      alert("Please select a wallet first");
      return;
    }
    try {
      await stakeJUL(selectedWallet, stakeAmount);
      alert(`Successfully staked ${stakeAmount} JUL`);
      handleGetBalance();
    } catch (error) {
      console.error("Failed to stake JUL:", error);
      alert("Failed to stake JUL");
    }
  };

  const handleUnstakeJUL = async () => {
    if (!selectedWallet) {
      alert("Please select a wallet first");
      return;
    }
    try {
      await unstakeJUL(selectedWallet, unstakeAmount);
      alert(`Successfully unstaked ${unstakeAmount} JUL`);
      handleGetBalance();
    } catch (error) {
      console.error("Failed to unstake JUL:", error);
      alert("Failed to unstake JUL");
    }
  };

  const createAndFundWallet = async () => {
    try {
      const response = await createWallet();
      const newAddress = response.address || '';
      setWallets(prevWallets => {
        const updatedWallets = [...prevWallets, newAddress];
        localStorage.setItem('wallets', JSON.stringify(updatedWallets));
        return updatedWallets;
      });
      setSelectedWallet(newAddress);
      
      // Automatically purchase some JUL for the new wallet
      const usdAmount = 100; // You can adjust this or make it user-input
      await purchaseJUL(newAddress, usdAmount);
      
      alert(`Wallet created with address: ${newAddress} and funded with ${usdAmount} USD worth of JUL`);
      handleGetBalance();
    } catch (error) {
      console.error("Failed to create and fund wallet:", error);
    }
  };

  return (
    <div>
      <h2>Wallet</h2>
      <button onClick={createAndFundWallet}>Create and Fund New Wallet</button>
      
      <h3>Select Wallet</h3>
      <select 
        value={selectedWallet} 
        onChange={(e) => setSelectedWallet(e.target.value)}
      >
        <option value="">Select a wallet</option>
        {wallets.map((wallet, index) => (
          <option key={index} value={wallet}>{wallet}</option>
        ))}
      </select>
      
      <p>Selected Wallet: {selectedWallet}</p>
      <button onClick={handleGetBalance}>Get Balance</button>
      <p>Balance: {balance} JUL</p>

      <h3>Purchase JUL</h3>
      <input 
        type="number" 
        value={purchaseAmount}
        onChange={(e) => setPurchaseAmount(Number(e.target.value))}
        placeholder="USD Amount"
      />
      <button onClick={handlePurchaseJUL}>Purchase JUL</button>

      <h3>Send Transaction</h3>
      <select 
        value={senderWallet} 
        onChange={(e) => setSenderWallet(e.target.value)}
      >
        <option value="">Select sender wallet</option>
        {wallets.map((wallet, index) => (
          <option key={index} value={wallet}>{wallet}</option>
        ))}
      </select>
      <select 
        value={receiverWallet} 
        onChange={(e) => setReceiverWallet(e.target.value)}
      >
        <option value="">Select receiver wallet</option>
        {wallets.map((wallet, index) => (
          <option key={index} value={wallet}>{wallet}</option>
        ))}
      </select>
      <input 
        type="number" 
        value={amount}
        onChange={(e) => setAmount(Number(e.target.value))}
        placeholder="Amount"
      />
      <button onClick={handleSendTransaction}>Send</button>

      <h3>Stake JUL</h3>
      <input 
        type="number" 
        value={stakeAmount}
        onChange={(e) => setStakeAmount(Number(e.target.value))}
        placeholder="Stake Amount"
      />
      <button onClick={handleStakeJUL}>Stake</button>

      <h3>Unstake JUL</h3>
      <input 
        type="number" 
        value={unstakeAmount}
        onChange={(e) => setUnstakeAmount(Number(e.target.value))}
        placeholder="Unstake Amount"
      />
      <button onClick={handleUnstakeJUL}>Unstake</button>
    </div>
  );
};

export default Wallet;
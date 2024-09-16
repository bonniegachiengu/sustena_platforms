import React, { useState, useEffect } from 'react';
import axios from 'axios';

function App() {
  const [accounts, setAccounts] = useState([]);
  const [chain, setChain] = useState([]);
  const [selectedAccount, setSelectedAccount] = useState('');
  const [recipientAddress, setRecipientAddress] = useState('');
  const [amount, setAmount] = useState('');
  const [fee, setFee] = useState('');
  const [newAccountName, setNewAccountName] = useState('');

  useEffect(() => {
    fetchChain();
    fetchAccounts();
  }, []);

  const fetchChain = async () => {
    try {
      const response = await axios.get('/get_chain');
      setChain(response.data);
    } catch (error) {
      console.error('Error fetching chain:', error);
    }
  };

  const fetchAccounts = async () => {
    try {
      const response = await axios.get('/get_accounts');
      setAccounts(response.data);
    } catch (error) {
      console.error('Error fetching accounts:', error);
    }
  };

  const createAccount = async () => {
    try {
      const response = await axios.post('/create_account', { name: newAccountName });
      fetchAccounts(); // Refresh the accounts list
      setNewAccountName('');
    } catch (error) {
      console.error('Error creating account:', error);
    }
  };

  const getBalance = async (address) => {
    try {
      const response = await axios.get(`/get_balance?address=${address}`);
      return response.data;
    } catch (error) {
      console.error('Error getting balance:', error);
    }
  };

  const transfer = async () => {
    try {
      await axios.post('/transfer', {
        from: selectedAccount,
        to: recipientAddress,
        amount: parseFloat(amount),
        fee: parseInt(fee)
      });
      fetchChain();
      fetchAccounts();
    } catch (error) {
      console.error('Error transferring:', error);
    }
  };

  return (
    <div>
      <h1>Blockchain UI</h1>
      
      <h2>Create Account</h2>
      <input 
        type="text" 
        value={newAccountName} 
        onChange={(e) => setNewAccountName(e.target.value)} 
        placeholder="Enter account name"
      />
      <button onClick={createAccount}>Create Account</button>

      <h2>Accounts</h2>
      <ul>
        {accounts.map(account => (
          <li key={account.Address}>
            Name: {account.Name}, Address: {account.Address}, Balance: {account.Balance}
          </li>
        ))}
      </ul>

      <h2>Transfer</h2>
      <div>
        From: 
        <select value={selectedAccount} onChange={(e) => setSelectedAccount(e.target.value)}>
          <option value="">Select account</option>
          {accounts.map(account => (
            <option key={account.Address} value={account.Address}>
              {account.Name} ({account.Address})
            </option>
          ))}
        </select>
      </div>
      <div>
        To: 
        <select value={recipientAddress} onChange={(e) => setRecipientAddress(e.target.value)}>
          <option value="">Select recipient</option>
          {accounts.map(account => (
            <option key={account.Address} value={account.Address}>
              {account.Name} ({account.Address})
            </option>
          ))}
        </select>
      </div>
      <div>
        Amount: <input type="number" value={amount} onChange={(e) => setAmount(e.target.value)} />
      </div>
      <div>
        Fee: <input type="number" value={fee} onChange={(e) => setFee(e.target.value)} />
      </div>
      <button onClick={transfer}>Transfer</button>

      <h2>Blockchain</h2>
      <ul>
        {chain.map(block => (
          <li key={block.Index}>
            Block {block.Index}: {block.Hash}
            <ul>
              {block.Transactions.map((tx, index) => (
                <li key={index}>
                  From: {tx.From}, To: {tx.To}, Amount: {tx.Amount}, Fee: {tx.Fee}
                </li>
              ))}
            </ul>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
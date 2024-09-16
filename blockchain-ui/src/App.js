import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [accounts, setAccounts] = useState([]);
  const [chain, setChain] = useState([]);
  const [selectedAccount, setSelectedAccount] = useState('');
  const [recipientAddress, setRecipientAddress] = useState('');
  const [amount, setAmount] = useState('');
  const [fee, setFee] = useState('');
  const [newAccountName, setNewAccountName] = useState('');
  const [qarAmount, setQarAmount] = useState('');
  const [mempool, setMempool] = useState([]);

  const NANO = 1000000000;

  useEffect(() => {
    fetchChain();
    fetchAccounts();
    fetchMempool();
    const interval = setInterval(fetchMempool, 5000); // Fetch mempool every 5 seconds
    return () => clearInterval(interval);
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

  const fetchMempool = async () => {
    try {
      const response = await axios.get('/get_mempool');
      setMempool(response.data);
    } catch (error) {
      console.error('Error fetching mempool:', error);
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
      await fetchChain();
      await fetchAccounts();
      await fetchMempool();
      setAmount('');
      setFee('');
      alert('Transfer successful!');
    } catch (error) {
      console.error('Error transferring:', error);
      alert(error.response?.data || 'Error during transfer');
    }
  };

  const buyJUL = async () => {
    try {
      await axios.post('/buy_jul', {
        address: selectedAccount,
        amount: parseFloat(qarAmount)
      });
      await fetchAccounts(); // Refresh accounts after buying JUL
      setQarAmount(''); // Clear the input field
      alert('JUL purchased successfully!');
    } catch (error) {
      console.error('Error buying JUL:', error);
      alert(error.response?.data || 'Error buying JUL');
    }
  };

  const renderBlockchain = () => {
    return (
      <div className="blockchain">
        {chain.map((block, index) => (
          <div key={block.Index} className="block">
            <h3>Block {block.Index}</h3>
            <p>Timestamp: {new Date(block.Timestamp * 1000).toLocaleString()}</p>
            <p>Validator: {block.Validator}</p>
            <p>Total Fee: {block.Transactions.reduce((sum, tx) => sum + tx.Fee, 0) / NANO} JUL</p>
            <h4>Transactions:</h4>
            <ul>
              {block.Transactions.map((tx, txIndex) => (
                <li key={txIndex}>
                  From: {tx.From} To: {tx.To} Amount: {tx.Amount / NANO} JUL Fee: {tx.Fee / NANO} JUL
                </li>
              ))}
            </ul>
            {index < chain.length - 1 && <div className="chain-link">↓</div>}
          </div>
        ))}
      </div>
    );
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
          {accounts.filter(account => account.Name !== "Exchange").map(account => (
            <option key={account.Address} value={account.Address}>
              {account.Name} ({account.Address}) - Balance: {account.Balance / NANO} JUL
            </option>
          ))}
        </select>
      </div>
      <div>
        To: 
        <select value={recipientAddress} onChange={(e) => setRecipientAddress(e.target.value)}>
          <option value="">Select recipient</option>
          {accounts.filter(account => account.Name !== "Exchange" && account.Address !== selectedAccount).map(account => (
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

      <h2>Buy JUL</h2>
      <div>
        Account: 
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
        QAR Amount: <input type="number" value={qarAmount} onChange={(e) => setQarAmount(e.target.value)} />
      </div>
      <button onClick={buyJUL}>Buy JUL</button>

      <h2>Mempool</h2>
      <ul>
        {mempool.map((tx, index) => (
          <li key={index}>
            From: {tx.From}, To: {tx.To}, Amount: {tx.Amount} JUL, Fee: {tx.Fee} JUL
          </li>
        ))}
      </ul>

      <h2>Blockchain</h2>
      {renderBlockchain()}
    </div>
  );
}

export default App;
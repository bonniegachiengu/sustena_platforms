import React, { useEffect, useState } from 'react';
import { getMempool } from '../services/api';

interface Transaction {
  ID: string;
  From: string;
  To: string;
  Amount: number;
  Fee: number;
}

const Mempool: React.FC = () => {
  const [mempool, setMempool] = useState<Transaction[]>([]);

  useEffect(() => {
    fetchMempool();
  }, []);

  const fetchMempool = async () => {
    try {
      const data = await getMempool();
      console.log("Mempool data:", data); // Add this line for debugging
      setMempool(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error("Failed to fetch mempool:", error);
      setMempool([]);
    }
  };

  return (
    <div>
      <h2>Mempool</h2>
      {mempool.length === 0 ? (
        <p>No transactions in the mempool</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>From</th>
              <th>To</th>
              <th>Amount</th>
              <th>Fee</th>
            </tr>
          </thead>
          <tbody>
            {mempool.map((tx) => (
              <tr key={tx.ID}>
                <td>{tx.ID}</td>
                <td>{tx.From}</td>
                <td>{tx.To}</td>
                <td>{tx.Amount} JUL</td>
                <td>{tx.Fee} JUL</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default Mempool;
import React, { useEffect, useState } from 'react';
import { getBlockchain, forgeBlock } from '../services/api';

interface Block {
  Index: number;
  Hash: string;
  Transactions: any[];
  Validator: string;
}

const Blockchain: React.FC = () => {
  const [blockchain, setBlockchain] = useState<Block[] | null>(null);

  useEffect(() => {
    fetchBlockchain();
  }, []);

  const fetchBlockchain = async () => {
    try {
      const data = await getBlockchain();
      console.log("Blockchain data:", data); // Add this line for debugging
      setBlockchain(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error("Failed to fetch blockchain:", error);
      setBlockchain([]);
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

  return (
    <div>
      <h2>Blockchain</h2>
      <button onClick={handleForgeBlock}>Forge New Block</button>
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
  );
};

export default Blockchain;
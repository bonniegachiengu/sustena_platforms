import React, { useEffect, useState } from 'react';
import { getValidators, getCommunityFund } from '../services/api';

interface Validator {
  Address: string;
  Stake: number;
}

const Validators: React.FC = () => {
  const [validators, setValidators] = useState<Validator[]>([]);
  const [communityFund, setCommunityFund] = useState(0);

  useEffect(() => {
    fetchValidators();
    fetchCommunityFund();
  }, []);

  const fetchValidators = async () => {
    try {
      const data = await getValidators();
      console.log("Validators data:", data); // Add this line for debugging
      setValidators(Array.isArray(data.validators) ? data.validators : []);
    } catch (error) {
      console.error("Failed to fetch validators:", error);
      setValidators([]);
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

  return (
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
      <h2>Community Fund</h2>
      <p>Balance: {communityFund} JUL</p>
    </div>
  );
};

export default Validators;
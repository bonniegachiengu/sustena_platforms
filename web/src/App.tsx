import React from 'react';
import { Route, Routes, Link } from 'react-router-dom';
import './App.css';
import Wallet from './components/Wallet';
import Blockchain from './components/Blockchain';
import Validators from './components/Validators';
import Mempool from './components/Mempool';

const App: React.FC = () => {
  return (
    <div className="App">
      <nav>
        <ul>
          <li><Link to="/">Home</Link></li>
          <li><Link to="/wallet">Wallet</Link></li>
          <li><Link to="/blockchain">Blockchain</Link></li>
          <li><Link to="/validators">Validators</Link></li>
          <li><Link to="/mempool">Mempool</Link></li>
        </ul>
      </nav>

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/wallet" element={<Wallet />} />
        <Route path="/blockchain" element={<Blockchain />} />
        <Route path="/validators" element={<Validators />} />
        <Route path="/mempool" element={<Mempool />} />
      </Routes>
    </div>
  );
};

const Home: React.FC = () => <h1>Welcome to Sustena Blockchain</h1>;

export default App;

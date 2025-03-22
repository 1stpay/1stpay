import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import InvoicePage from './components/InvoicePage';

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/payment/:payment_id" element={<InvoicePage />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
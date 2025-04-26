// src/pages/HomePage.tsx

import React from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../app/store';  // Assure-toi que RootState est bien exporté depuis store
import LinksList from '../features/links/LinksList';

const HomePage = () => {
  const token = useSelector((state: RootState) => state.auth.token);

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Welcome to Link Vault</h1>
      {token ? (<>
          <div className="text-right mb-4">
    <a href="/add" className="text-blue-500">Add a new link</a>
  </div>

        <LinksList /></>
      ) : (
        <p>Please <a href="/login" className="text-blue-500">log in</a> or <a href="/register" className="text-blue-500">register</a > to view your links.</p>
      )}
    </div>
  );
};

export default HomePage;


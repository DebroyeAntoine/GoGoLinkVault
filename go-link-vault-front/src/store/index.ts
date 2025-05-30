// src/app/store.ts

import { configureStore } from '@reduxjs/toolkit';
import authReducer from '../features/auth/authSlice';
import linksReducer from '../features/links/linksSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    links: linksReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;  // Export du type RootState

export default store;


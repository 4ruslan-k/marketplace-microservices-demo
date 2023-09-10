"use client"

import React, { createContext, useState } from 'react';

export const SnackbarContext = createContext({});

export interface SnackbarProviderContextType {
  snackbarState: {
    message: string;
    open: boolean;
    severity: string;
  };
  openSuccessSnackbar: (string: string) => void;
  openErrorSnackbar: (message: string) => void;
  closeSnackbar: () => void;
}

export function SnackbarProvider({ children }: any): any {
  const [snackbarState, setSnackbarState] = useState({
    message: '',
    open: false,
    severity: 'success',
  });

  const openSuccessSnackbar = (message: string): void => {
    setSnackbarState({ message, open: true, severity: 'success' });
  };
  const openErrorSnackbar = (message: string) => {
    setSnackbarState({ message, open: true, severity: 'error' });
  };
  const closeSnackbar = () => {
    setSnackbarState({ ...snackbarState, open: false, message: '' });
  };

  const value: SnackbarProviderProps = {
    snackbarState,
    openSuccessSnackbar,
    openErrorSnackbar,
    closeSnackbar,
  };

  return (
    <SnackbarContext.Provider
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      value={value}
    >
      {children}
    </SnackbarContext.Provider>
  );
}

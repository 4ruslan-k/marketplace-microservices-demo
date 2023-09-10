"use client";

import React from "react";
import { QueryClientProvider, QueryClient } from "@tanstack/react-query";
import { SnackbarProvider } from "../context/SnackbarContext";
import { ThemeProvider } from "@mui/material";
import CssBaseline from '@mui/material/CssBaseline';

import theme from "../theme";

function Providers({ children }: React.PropsWithChildren) {
  const queryClient = new QueryClient();


  return (
    <ThemeProvider theme={theme}>
      <CssBaseline>
        <QueryClientProvider client={queryClient}>
          <SnackbarProvider>
            {children}
          </SnackbarProvider>
        </QueryClientProvider>
      </CssBaseline>
    </ThemeProvider>
  );
}

export default Providers;

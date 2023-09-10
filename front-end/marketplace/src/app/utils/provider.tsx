"use client";

import React from "react";
import { QueryClientProvider, QueryClient } from "@tanstack/react-query";
import { SnackbarProvider } from "../context/SnackbarContext";
import { ThemeProvider } from "@mui/material";
import theme from "../theme";

function Providers({ children }: React.PropsWithChildren) {
  const queryClient = new QueryClient();


  return (
    <ThemeProvider theme={theme}>
    <QueryClientProvider client={queryClient}>
    <SnackbarProvider>
      {children}
      </SnackbarProvider>
      </QueryClientProvider>
      </ThemeProvider>
  );
}

export default Providers;

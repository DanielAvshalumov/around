// src/app/layout.tsx
'use client'

import { AppBar, Box, Button, CssBaseline, ThemeProvider, Toolbar, Typography, createTheme } from '@mui/material'
import type { ReactNode } from 'react'

import './globals.css'
import MatrixBackground from '@/components/matrix';
import { AuthProvider, useAuth } from '@/context/AuthContext';
import Navbar from '@/components/navbar';



const theme = createTheme({
  palette: {
    primary: { main: '#1976d2' }, // MUI blue
    background: { default: '#f9f9f9' },
  },
})

export default function RootLayout({ children }: { children: ReactNode }) {

  return (
  <html lang="en">
    <body style={{ margin: 0, padding: 0 }}>
      <AuthProvider>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          {/* Flex container with full viewport height */}
          <Box
            sx={{
              display: 'flex',
              flexDirection: 'column',
              minHeight: '100vh',
            }}
          >
            <Navbar />
            <MatrixBackground />
            
            {/* Main content area - grows to push footer down */}
            <Box component="main" sx={{ flexGrow: 1, position: 'relative', zIndex: 1 }}>
              {children}
            </Box>

            {/* Footer always at bottom */}
            <Box
              component="footer"
              sx={{
                backgroundColor: '#333',
                color: '#fff',
                padding: 3,
                textAlign: 'center',
                mt: 'auto', // Pushes footer to bottom when content is short
              }}
            >
              <Typography variant="body2">
                © 2023 Company Name. All rights reserved.
              </Typography>
            </Box>
          </Box>
        </ThemeProvider>
      </AuthProvider>
    </body>
  </html>
);
}

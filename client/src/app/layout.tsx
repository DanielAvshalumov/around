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
      <body>
        <AuthProvider>
          <ThemeProvider theme={theme}>
            <CssBaseline />
            <Navbar />
            <MatrixBackground />
            {children}
            <Box sx={{ backgroundColor: '#333', color: '#fff', padding: 3, textAlign: 'center' }}>
              <Typography variant="body2">© 2023 Company Name. All rights reserved.</Typography>
            </Box>
          </ThemeProvider>
        </AuthProvider>
      </body>
    </html>
  )
}

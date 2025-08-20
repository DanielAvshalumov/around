// src/app/layout.tsx
'use client'

import { AppBar, Box, Button, CssBaseline, ThemeProvider, Toolbar, Typography, createTheme } from '@mui/material'
import type { ReactNode } from 'react'
import Logo from "../../images/logo.png"
import Image from 'next/image';
import './globals.css'
import MatrixBackground from '@/components/matrix';

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
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <AppBar position="static" color="inherit">
            <Toolbar>
              <div style={{flexGrow: 1}}>
                <Image src={Logo} alt={''} priority/>
              </div>
              <Button variant='contained' color="inherit">Login</Button>
            </Toolbar>
          </AppBar>
          <MatrixBackground />
          {children}
          <Box sx={{ backgroundColor: '#333', color: '#fff', padding: 3, textAlign: 'center' }}>
            <Typography variant="body2">© 2023 Company Name. All rights reserved.</Typography>
          </Box>
        </ThemeProvider>
      </body>
    </html>
  )
}

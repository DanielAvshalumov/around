'use client'
import { Box, Typography, Button } from '@mui/material';
import React, { useEffect, useRef } from 'react';

const keywords = [
  "ads","SEO","PPC","click","leads","brand","reach","sales","growth","traffic",
  "content","funnel","impress","metrics","pixels","target","ranking","organic",
  "paid","search","Google","Bing","social","email","bounce","CTA","copy","retarget",
  "optimize","keyword","hashtag","audit","analytics","A/B test","domain","hosting",
  "mobile","video","shorts","buyer","persona","clickbait","caption","schedule",
  "campaign","headline","convo","story","subscribe","CTR","CPC","conversion","funnels",
  "nurture","scroll","segment","split test","inbound","outbound","influencer","reels",
  "UGC","brand lift","automation","landing page","lead magnet","open rate","bounce rate",
  "geo-target","lookalike","boost","sponsor","KPI","ROI","lifetime value","CAC","drip"
];

const Home = () => {
  const matrixRef = useRef(null);
  
  useEffect(() => {
    const container = matrixRef.current;
    if (!container) return;
    
    container.innerHTML = ""; // Clear previous runs if component reloads
    const numColumns = Math.floor(window.innerWidth / 100);
    
    for (let i = 0; i < numColumns; i++) {
      const word = document.createElement("div");
      word.className = "matrix-word";
      word.innerText = keywords[Math.floor(Math.random() * keywords.length)];
      word.style.left = `${i * 100}px`;
      word.style.animationDelay = `${Math.random() * 5}s`;
      word.style.animationDuration = `${3 + Math.random() * 3}s`;
      container.appendChild(word);
    }
    
    // Create new words periodically
    const interval = setInterval(() => {
      if (container) {
        const word = document.createElement("div");
        word.className = "matrix-word";
        word.innerText = keywords[Math.floor(Math.random() * keywords.length)];
        word.style.left = `${Math.floor(Math.random() * window.innerWidth)}px`;
        word.style.animationDuration = `${3 + Math.random() * 3}s`;
        container.appendChild(word);
        
        // Remove old words to prevent memory issues
        if (container.children.length > numColumns * 3) {
          container.removeChild(container.children[0]);
        }
      }
    }, 1000);
    
    return () => clearInterval(interval);
  }, []);

  return (
    <Box sx={{ position: 'relative', overflow: 'hidden' }}>
      {/* Matrix Rain Container */}
      <div 
        ref={matrixRef}
        style={{
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          zIndex: 0,
        }}
        className="matrix-container"
      />
      
      {/* Your original MUI component */}
      

        <Box sx={{backgroundColor: 'white'}}>
          <Typography variant='h1' fontSize='100px'>yall like disque</Typography>
          </Box>

      <Box 
        className='text-container' 
        sx={{ 
          padding: 5, 
          backgroundColor: 'rgba(135, 206, 235, 0.8)', 
          textAlign: 'center',
          position: 'relative',
          zIndex: 1
        }}
      >
        <Typography 
          variant="h3" 
          component="h2" 
          gutterBottom 
          align="center" 
          sx={{ mb: 2, fontWeight: 'bold', color: 'white' }}
        >
          Get Around With Free SEO Tools
        </Typography>
        <Typography 
          variant="h6" 
          align="center" 
          color="text.secondary" 
          sx={{ mb: 2, color: 'white' }}
        >
          Boost your online presence with our complimentary SEO analysis tools. Select a tool below to get started.
        </Typography>
        <Button variant="contained" color="primary" size="large">
          Run
        </Button>
      </Box>
      
      {/* Styles for matrix effect */}
      <style jsx>{`
        .matrix-word {
          position: absolute;
          top: -20px;
          color: rgba(255, 255, 255, 0.7);
          font-family: monospace;
          font-weight: bold;
          font-size: 16px;
          text-shadow: 0 0 5px rgba(255, 255, 255, 0.5);
          animation: fall linear infinite;
          opacity: 0.6;
        }
        
        @keyframes fall {
          0% {
            transform: translateY(-100px);
            opacity: 0;
          }
          10% {
            opacity: 0.6;
          }
          90% {
            opacity: 0.6;
          }
          100% {
            transform: translateY(100vh);
            opacity: 0;
          }
        }
      `}</style>
    </Box>
  );
};

export default Home;
'use client'
import React, { useEffect, useState } from 'react';
import { Container, Box, Typography, Button, AppBar, Toolbar, TextField } from '@mui/material';
import styles from "./page.module.css"
import BacklinkBuilderForm from '@/components/landing-page-content';
import SourceLinkTable from '@/components/serp/serp';
import CoolLoadingScreen from '@/components/serp/loading/serp-loading';
import getForumProductLinks from '../../lib/backlink';
import SEOAIPromo from '@/components/page-component-one';

// {source:"",backlink:"",dofollow:false}
const MainPage = () => {

    const [isCentered, setIsCentered] = useState(true);
    const [loading, setLoading] = useState(false)
    const [backlinks, setBacklinks] = useState<any[]>([])
    const [industry, setIndustry] = useState('');

    useEffect(() => {
      console.log('hello',industry)
    },[industry])

    const handleBacklink = async () => {
      console.log("work")
      setLoading(true)
      setIsCentered(false)
      try {
        const api = await getForumProductLinks({
          industry: industry,
          comp_domains: null,
          browser: ''
        })
        // await new Promise(res => setTimeout(res,1500))
        // const data = await api.data;
        // console.log(api)
        setBacklinks(api)
        console.log(api)
        if (backlinks.length == 0) {
          setIsCentered(true)
        }
      } catch(err) {
        console.log(err)
      } finally {
        setLoading(false)
      }
    }

  return (
    <Box display='flex' flexDirection='column' minHeight='84.3vh' gap={2}>
      <Box 
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
        <TextField onChange={(e) => setIndustry(e.target.value)}/>
        <Typography 
          variant="h6" 
          align="center" 
          color="text.secondary" 
          sx={{ mb: 2, color: 'white' }}
        >
          Boost your online presence with our complimentary SEO analysis tools. Select a tool below to get started.
        </Typography>
        <div style={{display:'flex',gap:'10px',justifyContent:'center'}}>
          <Button variant="contained" color="primary" size="large">
            Backlink Builder
          </Button>
          <Button variant="contained" color="primary" size="large">
            Ad Copy Generator
          </Button>
          <Button variant="contained" color="primary" size="large" onClick={handleBacklink}>
            Page Rank Finder
          </Button>
        </div>
      </Box>
      <div className={styles.layoutContainer}>
        <div className={`${styles.originalComponent} ${isCentered ? styles.centered : styles.leftAligned}`}>
          <BacklinkBuilderForm handleBacklink={handleBacklink} industry={industry} setIndustry={setIndustry}/>
        </div>
        <div className={`${styles.newComponent} ${isCentered ? styles.hidden : styles.visible}`}>

          {loading ? <CoolLoadingScreen isCentered={false} /> : <SourceLinkTable backlinks={backlinks} isCentered={false}/>}
        </div>
      </div>
      <SEOAIPromo />
    </Box>
  );
};

export default MainPage;
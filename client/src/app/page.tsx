'use client'
import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import { Container, Box, Typography, Button, AppBar, Toolbar, TextField } from '@mui/material';
import styles from "./page.module.css"
import BacklinkBuilderForm from '@/components/landing-page-content';
import SourceLinkTable from '@/components/serp/serp';
import CoolLoadingScreen from '@/components/serp/loading/serp-loading';
import getForumProductLinks from '../../lib/backlink';
import SEOAIPromo from '@/components/page-component-one';
import getBacklinkDomains from '../../lib/hugHelper';
import { AuthService } from '../../lib/auth';
import BacklinkScraperBanner from '@/components/front-page-banner';

// {source:"",backlink:"",dofollow:false}
const MainPage = () => {

    const [isCentered, setIsCentered] = useState(true);
    const [loading, setLoading] = useState(false);
    const [domainLoad, setDomainLoad] = useState(false);
    const [backlinks, setBacklinks] = useState<any[]>([]);
    const [industry, setIndustry] = useState('');
    const [competitorDomains, setCompetitorDomains] = useState(['']);

    
    

    const handleDomainsNames = async () => {
      console.log('more work')
      setDomainLoad(true)
      try {
        const api = await getBacklinkDomains(industry)
        setCompetitorDomains(api)
      } catch(err) {
        console.log(err)
      } finally {
        setDomainLoad(false)
      }
    }

    const handleBacklink = async () => {
      console.log("work")
      setLoading(true)
      setIsCentered(false)
      try {
        const api = await getForumProductLinks({
          industry: industry,
          Comp_domains: competitorDomains,
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

    const backlinkForm = 
    <div className={`${styles.originalComponent} ${isCentered ? styles.centered : styles.leftAligned}`}>
      <BacklinkBuilderForm handleBacklink={handleBacklink} industry={industry} setIndustry={setIndustry} competitorDomains={competitorDomains} setCompetitorDomains={setCompetitorDomains}/>
    </div>

  return (
    <Box display='flex' flexDirection='column' minHeight='84.3vh' gap={2}>
      <SEOAIPromo />

      {/* <Box 
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
          <Button variant="contained" color="primary" size="large" onClick={handleDomainsNames}>
            Ad Copy Generator
          </Button>
          <Button variant="contained" color="primary" size="large" onClick={handleBacklink}>
            Page Rank Finder
          </Button>
        </div>
      </Box> */}

      <BacklinkScraperBanner />
      <div className={styles.layoutContainer}>
        <div className={`${styles.originalComponent} ${isCentered ? styles.centered : styles.leftAligned}`}>
          <BacklinkBuilderForm handleBacklink={handleBacklink} industry={industry} setIndustry={setIndustry} competitorDomains={competitorDomains} setCompetitorDomains={setCompetitorDomains}/>
        </div>
        <div className={`${styles.newComponent} ${isCentered ? styles.hidden : styles.visible}`}>
          {loading ? <CoolLoadingScreen isCentered={false} /> : <SourceLinkTable backlinks={backlinks} isCentered={false}/>}
        </div>
      </div>
    </Box>
  );
};

export default MainPage;
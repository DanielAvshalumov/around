// BacklinkResponseBuilder.tsx
'use client'
import React, { useState, useEffect } from 'react';
import styles from './page.module.css';
import { getBacklink, getGenAIReply, publishBacklink, testGenAiReply } from '../../../../lib/backlink';
import { useParams } from 'next/navigation';

export default function BacklinkResponseBuilder() {
  const [response, setResponse] = useState("");
  const [productUrl, setProductUrl] = useState('');
  const [copied, setCopied] = useState(false);
  const [isRegenerating, setIsRegenerating] = useState(false);

  const { id } = useParams()
  const [backlink,setBacklink] = useState({
    Source: "",
    Id: 0,
  })

  const opportunity = {
    platform: "Reddit",
    subreddit: "r/digitalmarketing",
    title: "What are your thoughts on sustainable marketing practices?",
    url: "reddit.com/r/digitalmarketing/...",
  };

  const handleCopy = async () => {
    await navigator.clipboard.writeText(response);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };
  
  const handleRegenerate = () => {
    setIsRegenerating(true);
    setTimeout(() => {
      setIsRegenerating(false);
    }, 1000);
  };


  const handlePublish = async () => {
    const res = await publishBacklink(backlink.Id, response)
    console.log(res)    
  }

  const handleBacklink = async (id: string) => {
    const data = await getBacklink(id)
    setBacklink(data)
  }

  const generateReply = async () => {
    console.log(backlink)
    const res = await getGenAIReply(backlink?.Source)
    // const res = await testGenAiReply()
    // const data = await res.data
    setResponse(res);
    console.log(res)
    return res
  }

  useEffect(() => {
    const bId = id as string
    handleBacklink(bId)
  },[])

  useEffect(() => {
    console.log(backlink)
  },[backlink])

  useEffect(() => {
    console.log(response)
  },[response])

  // const charCount = response.length || 0;
  // const wordCount = response.split(/\s+/).filter(Boolean).length || 0;

  return (
    <div className={styles.container}>
      <div className={styles.bgOrb1}></div>
      <div className={styles.bgOrb2}></div>

      <div className={styles.content}>
        {/* Header */}
        <div className={styles.header}>
          <div className={styles.iconWrapper}>
            <svg className={styles.headerIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"></path>
            </svg>
          </div>
          <h1 className={styles.title}>Craft Your Response</h1>
          <p className={styles.subtitle}>AI-powered engagement for {opportunity.platform}</p>
        </div>

        {/* Opportunity Badge */}
        <div className={styles.opportunityBadge}>
          <div className={styles.emoji}>💬</div>
          <div className={styles.opportunityContent}>
            <div className={styles.subreddit}>{backlink.Source}</div>
            <div className={styles.opportunityTitle}>{backlink?.Title}</div>
          </div>
          <a 
            href={`${backlink.Link}`}
            target="_blank"
            rel="noopener noreferrer"
            className={styles.viewBtn}
          >
            View
            <svg className={styles.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path>
            </svg>
          </a>
        </div>

        {/* Product URL */}
        <div className={styles.urlSection}>
          <div className={styles.glowWrapper}>
            <div className={styles.urlCard}>
              <div className={styles.urlLabel}>
                <svg className={styles.linkIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path>
                </svg>
                <span className={styles.urlLabelText}>Your Link</span>
              </div>
              <input
                type="url"
                value={productUrl}
                onChange={(e) => setProductUrl(e.target.value)}
                placeholder="https://yoursite.com/resource"
                className={styles.urlInput}
              />
            </div>
          </div>
        </div>

        {/* Response Editor */}
        <div className={styles.editorWrapper}>
          <div className={styles.editorCard}>
            <div className={styles.editorHeader}>
              <div className={styles.editorTitle}>
                <div className={styles.statusDot}></div>
                <span>Response</span>
              </div>
              <button
                onClick={generateReply}
                className={`${styles.regenerateBtn} ${isRegenerating ? styles.spinning : ''}`}
              >
                <svg className={styles.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                </svg>
                Regenerate
              </button>
            </div>
            
            <div className={styles.editorBody}>
              <textarea
                value={response}
                onChange={(e) => setResponse(e.target.value)}
                className={styles.textarea}
                placeholder="Your response appears here..."
              />
            </div>
          </div>
        </div>

        {/* Footer Actions */}
        <div className={styles.footerActions}>
          {/* <div className={styles.stats}>
            <div className={styles.stat}>
              <div className={styles.statNumber}>{charCount}</div>
              <div className={styles.statLabel}>Characters</div>
            </div>
            <div className={styles.stat}>
              <div className={`${styles.statNumber} ${styles.words}`}>{wordCount}</div>
              <div className={styles.statLabel}>Words</div>
            </div>
          </div> */}

          <div className={styles.actions}>
            <button onClick={handleCopy} className={styles.copyBtn}>
              {copied ? (
                <>
                  <svg className={styles.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7"></path>
                  </svg>
                  Copied
                </>
              ) : (
                <>
                  <svg className={styles.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
                  </svg>
                  Copy
                </>
              )}
            </button>
            <button onClick={handlePublish} className={styles.publishBtn}>Publish Response</button>
          </div>
        </div>
      </div>
    </div>
  );
}
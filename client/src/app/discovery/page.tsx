'use client';

import { useState, KeyboardEvent, FormEvent } from 'react';
import styles from './page.module.css';

interface BacklinkStrategy {
  title: string;
  description: string;
}

export default function BacklinkBuilder() {
  const [features, setFeatures] = useState<string[]>([]);
  const [tagInput, setTagInput] = useState('');
  const [productType, setProductType] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [results, setResults] = useState<BacklinkStrategy[] | null>(null);

  const handleTagKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      const value = tagInput.trim();
      if (value && !features.includes(value)) {
        setFeatures([...features, value]);
        setTagInput('');
      }
    }
  };

  const removeTag = (tagToRemove: string) => {
    setFeatures(features.filter(tag => tag !== tagToRemove));
  };

  const generateStrategies = (type: string, goals: string[]): BacklinkStrategy[] => {
    const goalList = goals.join(', ');
    
    return [
      {
        title: `Guest Posting Campaign for ${type}`,
        description: `Launch a strategic guest posting initiative targeting high-authority sites in your niche. Focus on ${goals[0] || 'industry-relevant topics'} to build contextual backlinks. Leverage ${goalList} to create compelling, data-driven content that editors want to publish. Expected DA range: 40-70.`
      },
      {
        title: `Digital PR & Journalist Outreach`,
        description: `Build relationships with journalists and bloggers covering ${type.toLowerCase()}. Create newsworthy content around ${goalList}, pitch exclusive insights and data studies. This approach generates high-quality editorial backlinks from news sites and industry publications.`
      },
      {
        title: `Resource Page Link Building`,
        description: `Identify resource pages and curated lists in the ${type.toLowerCase()} space. Your content addressing ${goals[0] || 'key pain points'} can be pitched as a valuable addition. Target pages linking to competitors, offering superior resources around ${goalList}.`
      },
      {
        title: `Broken Link Replacement Strategy`,
        description: `Find broken links on authoritative sites related to ${type.toLowerCase()}. Create replacement content covering ${goalList} and reach out to webmasters offering your working link as a solution. This win-win approach typically has a 15-25% success rate.`
      },
      {
        title: `Industry Authority Content Hub`,
        description: `Develop comprehensive guides and tools focused on ${goalList} that become go-to resources for ${type.toLowerCase()}. Natural backlinks accumulate as others reference your authoritative content. Promote through communities, forums, and social channels for amplification.`
      }
    ];
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    
    if (!productType.trim()) {
      alert('Please enter your target website or niche');
      return;
    }

    if (features.length === 0) {
      alert('Please add at least one link building goal');
      return;
    }

    setIsLoading(true);
    setResults(null);

    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 2500));

    const strategies = generateStrategies(productType, features);
    setResults(strategies);
    setIsLoading(false);
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1 className={styles.title}>Backlink Strategy Builder</h1>
        <p className={styles.subtitle}>AI-powered link building intelligence for enterprise SEO teams</p>
        <span className={styles.techBadge}>Enterprise Analytics Platform</span>
      </header>

      <div className={styles.layout}>
        <div className={styles.formSection}>
          <form onSubmit={handleSubmit}>
            <div className={styles.formGroup}>
              <label htmlFor="productType" className={styles.label}>
                Target Website/Niche
              </label>
              <input 
                type="text" 
                id="productType"
                className={styles.input}
                placeholder="e.g., Tech Blog, E-commerce, SaaS Platform"
                value={productType}
                onChange={(e) => setProductType(e.target.value)}
                required
              />
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="tagInput" className={styles.label}>
                Link Building Goals (press Enter after each)
              </label>
              <div className={styles.tagInputContainer}>
                {features.map((feature, index) => (
                  <div key={index} className={styles.tag}>
                    {feature}
                    <button 
                      type="button" 
                      onClick={() => removeTag(feature)}
                      className={styles.tagButton}
                      aria-label={`Remove ${feature}`}
                    >
                      ×
                    </button>
                  </div>
                ))}
                <input 
                  type="text" 
                  id="tagInput"
                  className={styles.tagInput}
                  placeholder="Add a goal and press Enter"
                  value={tagInput}
                  onChange={(e) => setTagInput(e.target.value)}
                  onKeyDown={handleTagKeyDown}
                />
              </div>
            </div>

            <button type="submit" className={styles.submitBtn}>
              Generate Strategy
            </button>
          </form>
        </div>

        <div className={styles.resultsSection}>
          <div className={styles.resultsContainer}>
            {isLoading ? (
              <div className={styles.loadingContainer}>
                <div className={styles.loadingSpinner}></div>
                <p className={styles.loadingText}>Analyzing backlink opportunities...</p>
              </div>
            ) : results ? (
              <div className={styles.resultsContent}>
                <div className={styles.resultsHeader}>
                  <h2 className={styles.resultsTitle}>Link Strategies</h2>
                  <span className={styles.resultsCount}>{results.length} Strategies</span>
                </div>
                {results.map((strategy, index) => (
                  <div key={index} className={styles.resultCard}>
                    <div className={styles.resultNumber}>{index + 1}</div>
                    <h4>{strategy.title}</h4>
                    <p>{strategy.description}</p>
                  </div>
                ))}
              </div>
            ) : (
              <div className={styles.emptyState}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
                </svg>
                <h3>Ready to Build Links?</h3>
                <p>Configure your strategy to generate backlink opportunities</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
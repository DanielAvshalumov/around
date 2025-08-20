'use client'
import React, { useState } from 'react';
import { Plus, X, Link, Building, Rocket, Trash2 } from 'lucide-react';
import styles from './page.module.css';


export default function BacklinkBuilderForm({ handleBacklink } : any) {
  const [industry, setIndustry] = useState('');
  const [competitorDomains, setCompetitorDomains] = useState(['']);

  const addDomain = () => {
    setCompetitorDomains([...competitorDomains, '']);
  };

  const removeDomain = (index) => {
    if (competitorDomains.length > 1) {
      setCompetitorDomains(competitorDomains.filter((_, i) => i !== index));
    }
  };

  const updateDomain = (index, value) => {
    const newDomains = [...competitorDomains];
    newDomains[index] = value;
    setCompetitorDomains(newDomains);
  };

  const handleKeyPress = (e, index) => {
    if (e.key === 'Enter' && competitorDomains[index].trim()) {
      e.preventDefault();
      if (index === competitorDomains.length - 1) {
        addDomain();
      }
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const validDomains = competitorDomains.filter(domain => domain.trim());
    console.log({
      industry: industry.trim(),
      competitorDomains: validDomains
    });
    // Handle form submission here
  };

  const isFormValid = industry.trim() && competitorDomains.some(domain => domain.trim());
  const validDomains = competitorDomains.filter(d => d.trim());

  return (
    <div>
      {/* Background Elements */}
      <div className={styles.backgroundElements}>
        <div className={styles.backgroundBlob1}></div>
        <div className={styles.backgroundBlob2}></div>
        <div className={styles.backgroundBlob3}></div>
      </div>

      {/* Main Form Container */}
      <div className={styles.formWrapper}>
        <div className={styles.formContainer}>
          {/* Animated Border Container */}
          <div className={styles.cardContainer}>
            {/* Animated Border */}
            <div className={styles.animatedBorder}></div>
            <div className={styles.cardBackground}></div>
            
            {/* Content */}
            <div className={styles.cardContent}>
              {/* Header */}
              <div className={styles.header}>
                <div className={styles.iconContainer}>
                  <Link className={styles.headerIcon} />
                </div>
                <h1 className={styles.title}>
                  Backlink Builder
                </h1>
                <p className={styles.subtitle}>
                  Analyze your competition and discover powerful backlink opportunities
                </p>
              </div>

              <form onSubmit={handleSubmit} className={styles.form}>
                {/* Industry Input */}
                <div className={styles.inputGroup}>
                  <label className={styles.label}>
                    Industry
                  </label>
                  <div className={styles.inputWrapper}>
                    <Building className={styles.inputIcon} />
                    <input
                      type="text"
                      value={industry}
                      onChange={(e) => setIndustry(e.target.value)}
                      placeholder="e.g., Digital Marketing, SaaS, E-commerce"
                      className={styles.input}
                    />
                  </div>
                </div>

                {/* Competitor Domains */}
                <div className={styles.domainsSection}>
                  <label className={styles.label}>
                    Competitor Domains
                  </label>
                  
                  <div className={styles.domainsList}>
                    {competitorDomains.map((domain, index) => (
                      <div key={index} className={styles.domainRow}>
                        <div className={styles.domainInputWrapper}>
                          <Link className={styles.inputIcon} />
                          <input
                            type="text"
                            value={domain}
                            onChange={(e) => updateDomain(index, e.target.value)}
                            onKeyPress={(e) => handleKeyPress(e, index)}
                            placeholder="example.com"
                            className={styles.input}
                          />
                        </div>
                        {competitorDomains.length > 1 && (
                          <button
                            type="button"
                            onClick={() => removeDomain(index)}
                            className={styles.removeButton}
                          >
                            <Trash2 className={styles.removeIcon} />
                          </button>
                        )}
                      </div>
                    ))}
                  </div>

                  {/* Add Domain Button */}
                  <button
                    type="button"
                    onClick={addDomain}
                    className={styles.addButton}
                  >
                    <Plus className={styles.addIcon} />
                    <span>Add Another Domain</span>
                  </button>
                </div>

                {/* Domain Preview */}
                {validDomains.length > 0 && (
                  <div className={styles.preview}>
                    <p className={styles.previewLabel}>Domains to analyze:</p>
                    <div className={styles.previewTags}>
                      {validDomains.map((domain, index) => (
                        <span key={index} className={styles.previewTag}>
                          {domain}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Submit Button */}
                <button
                  type="submit"
                  disabled={!isFormValid}
                  className={`${styles.submitButton} ${isFormValid ? styles.submitButtonEnabled : styles.submitButtonDisabled}`}
                   onClick={handleBacklink}
                >
                  <Rocket className={`${styles.submitIcon} ${isFormValid ? styles.submitIconAnimated : ''}`} />
                  <span>Analyze Backlinks</span>
                </button>
              </form>

              {/* Footer */}
              <div className={styles.footer}>
                <p className={styles.footerText}>
                  Discover high-quality backlink opportunities from your top competitors
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
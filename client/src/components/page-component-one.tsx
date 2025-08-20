import React, { useState, useEffect } from 'react';
import { Brain, Search, TrendingUp, ArrowRight } from 'lucide-react';

const SEOAIPromo = () => {
  const [mounted, setMounted] = useState(false);
  const [hovered, setHovered] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const styles = {
    container: {
      minHeight: '100vh',
      position: 'relative',
      overflow: 'hidden',
      fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif'
    },
    content: {
      maxWidth: '1400px',
      margin: '0 auto',
      padding: '0 2rem',
      position: 'relative',
      zIndex: 1
    },
    mainSection: {
      display: 'grid',
      gridTemplateColumns: '1fr 1fr',
      gap: '4rem',
      alignItems: 'center',
      minHeight: '100vh',
      paddingTop: '2rem'
    },
    leftColumn: {
      paddingLeft: '2rem',
      opacity: mounted ? 1 : 0,
      transform: mounted ? 'translateX(0)' : 'translateX(-50px)',
      transition: 'all 1s ease'
    },
    rightColumn: {
      display: 'flex',
      flexDirection: 'column',
      gap: '2rem',
      paddingRight: '2rem',
      opacity: mounted ? 1 : 0,
      transform: mounted ? 'translateX(0)' : 'translateX(50px)',
      transition: 'all 1s ease 0.2s'
    },
    eyebrow: {
      fontSize: '0.9rem',
      color: '#64748b',
      fontWeight: '500',
      marginBottom: '1rem',
      letterSpacing: '0.5px'
    },
    title: {
      fontSize: 'clamp(2.5rem, 6vw, 4rem)',
      fontWeight: '700',
      color: '#1e293b',
      lineHeight: 1.1,
      marginBottom: '1.5rem'
    },
    highlight: {
      background: 'linear-gradient(135deg, #3b82f6, #1d4ed8)',
      backgroundClip: 'text',
      WebkitBackgroundClip: 'text',
      WebkitTextFillColor: 'transparent'
    },
    description: {
      fontSize: '1.2rem',
      color: '#475569',
      lineHeight: 1.6,
      marginBottom: '2.5rem',
      maxWidth: '500px'
    },
    ctaButton: {
      display: 'inline-flex',
      alignItems: 'center',
      gap: '0.75rem',
      background: 'linear-gradient(135deg, #3b82f6, #1d4ed8)',
      color: 'white',
      padding: '1rem 2rem',
      borderRadius: '12px',
      fontSize: '1.1rem',
      fontWeight: '600',
      border: 'none',
      cursor: 'pointer',
      transition: 'all 0.3s ease',
      boxShadow: '0 4px 20px rgba(59, 130, 246, 0.3)',
      textDecoration: 'none'
    },
    featureCard: {
      background: 'rgba(255, 255, 255, 0.7)',
      backdropFilter: 'blur(20px)',
      border: '1px solid rgba(255, 255, 255, 0.2)',
      borderRadius: '20px',
      padding: '2rem',
      transition: 'all 0.4s ease',
      cursor: 'pointer'
    },
    featureIcon: {
      width: '50px',
      height: '50px',
      borderRadius: '12px',
      background: 'linear-gradient(135deg, #3b82f6, #1d4ed8)',
      color: 'white',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      marginBottom: '1.5rem'
    },
    featureTitle: {
      fontSize: '1.25rem',
      fontWeight: '600',
      color: '#1e293b',
      marginBottom: '0.75rem'
    },
    featureDescription: {
      color: '#64748b',
      lineHeight: 1.5,
      fontSize: '0.95rem'
    },
    floatingElement: {
      position: 'absolute',
      background: 'rgba(255, 255, 255, 0.1)',
      borderRadius: '50%',
      pointerEvents: 'none'
    },
    mobileStyles: {
      '@media (max-width: 768px)': {
        gridTemplateColumns: '1fr',
        gap: '2rem',
        padding: '2rem 1rem'
      }
    }
  };

  return (
    <div style={styles.container}>
      {/* Floating background elements */}
      <div style={{
        ...styles.floatingElement,
        width: '300px',
        height: '300px',
        top: '10%',
        right: '10%',
        background: 'radial-gradient(circle, rgba(59, 130, 246, 0.1) 0%, transparent 70%)'
      }}></div>
      <div style={{
        ...styles.floatingElement,
        width: '200px',
        height: '200px',
        bottom: '20%',
        left: '5%',
        background: 'radial-gradient(circle, rgba(29, 78, 216, 0.08) 0%, transparent 70%)'
      }}></div>

      <div style={styles.content}>
        <div style={{
          ...styles.mainSection,
          '@media (maxWidth: 768px)': {
            gridTemplateColumns: '1fr',
            gap: '2rem'
          }
        }}>
          {/* Left Column */}
          <div style={styles.leftColumn}>
            <div style={styles.eyebrow}>
              AI-POWERED SEO TOOLS
            </div>
            <h1 style={styles.title}>
              Free SEO analysis that{' '}
              <span style={styles.highlight}>trains your AI</span>
            </h1>
            <p style={styles.description}>
              Get comprehensive SEO insights while our AI learns your domain's unique patterns, 
              creating personalized recommendations that improve over time.
            </p>
            <button
              style={styles.ctaButton}
              onMouseEnter={(e) => {
                e.target.style.transform = 'translateY(-2px)';
                e.target.style.boxShadow = '0 8px 30px rgba(59, 130, 246, 0.4)';
              }}
              onMouseLeave={(e) => {
                e.target.style.transform = 'translateY(0)';
                e.target.style.boxShadow = '0 4px 20px rgba(59, 130, 246, 0.3)';
              }}
            >
              Start Free Analysis
              <ArrowRight size={20} />
            </button>
          </div>

          {/* Right Column */}
          <div style={styles.rightColumn}>
            <div 
              style={{
                ...styles.featureCard,
                transform: 'translateX(2rem)',
                marginBottom: '1rem'
              }}
              onMouseEnter={(e) => {
                e.target.style.transform = 'translateX(2rem) translateY(-5px)';
                e.target.style.boxShadow = '0 10px 40px rgba(0, 0, 0, 0.1)';
              }}
              onMouseLeave={(e) => {
                e.target.style.transform = 'translateX(2rem)';
                e.target.style.boxShadow = 'none';
              }}
            >
              <div style={styles.featureIcon}>
                <Search size={24} />
              </div>
              <h3 style={styles.featureTitle}>Instant SEO Audit</h3>
              <p style={styles.featureDescription}>
                Complete website analysis with actionable recommendations in seconds
              </p>
            </div>

            <div 
              style={{
                ...styles.featureCard,
                transform: 'translateX(-1rem)'
              }}
              onMouseEnter={(e) => {
                e.target.style.transform = 'translateX(-1rem) translateY(-5px)';
                e.target.style.boxShadow = '0 10px 40px rgba(0, 0, 0, 0.1)';
              }}
              onMouseLeave={(e) => {
                e.target.style.transform = 'translateX(-1rem)';
                e.target.style.boxShadow = 'none';
              }}
            >
              <div style={styles.featureIcon}>
                <Brain size={24} />
              </div>
              <h3 style={styles.featureTitle}>AI Domain Training</h3>
              <p style={styles.featureDescription}>
                Personalized AI model learns your industry and content patterns
              </p>
            </div>

            <div 
              style={{
                ...styles.featureCard,
                transform: 'translateX(1.5rem)',
                marginTop: '1rem'
              }}
              onMouseEnter={(e) => {
                e.target.style.transform = 'translateX(1.5rem) translateY(-5px)';
                e.target.style.boxShadow = '0 10px 40px rgba(0, 0, 0, 0.1)';
              }}
              onMouseLeave={(e) => {
                e.target.style.transform = 'translateX(1.5rem)';
                e.target.style.boxShadow = 'none';
              }}
            >
              <div style={styles.featureIcon}>
                <TrendingUp size={24} />
              </div>
              <h3 style={styles.featureTitle}>Smart Insights</h3>
              <p style={styles.featureDescription}>
                Domain-specific recommendations that evolve with your content
              </p>
            </div>
          </div>
        </div>
      </div>

      <style jsx>{`
        @media (max-width: 768px) {
          .main-section {
            grid-template-columns: 1fr !important;
            gap: 2rem !important;
          }
          .left-column, .right-column {
            padding-left: 0 !important;
            padding-right: 0 !important;
          }
          .feature-card {
            transform: none !important;
            margin: 0 !important;
          }
        }
      `}</style>
    </div>
  );
};

export default SEOAIPromo;
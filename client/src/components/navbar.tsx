'use client';

import Link from 'next/link';
import { useState, useEffect } from 'react';
import styles from './page.module.css';
import { useAuth } from '@/context/AuthContext';
import { AuthService } from '../../lib/auth';

export default function Navbar() {
  const [scrolled, setScrolled] = useState(false);
  const [activeSection, setActiveSection] = useState('');
  const { isAuthenticated, isLoading } = useAuth()

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 20);
    };
    
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const navItems = [
    { label: 'Features', href: '/features', icon: '✦' },
    { label: 'Solutions', href: '/solutions', icon: '◆' },
    { label: 'Resources', href: '/resources', icon: '▲' },
    { label: 'Pricing', href: '/pricing', icon: '●' },
  ];

  return (
    <nav className={`${styles.navbar} ${scrolled ? styles.scrolled : ''}`}>
      <div className={styles.gradientOverlay} />
      <div className={styles.noiseTexture} />
      
      <div className={styles.navContent}>
        <Link href="/" className={styles.logoContainer}>
          <div className={styles.logoIcon}>
            <div className={styles.logoIconInner} />
          </div>
          <span className={styles.logoText}>
            GettAround
            <span className={styles.logoBeta}>beta</span>
          </span>
        </Link>
        
        <div className={styles.navCenter}>
          {navItems.map((item, idx) => (
            <Link 
              key={item.href} 
              href={item.href} 
              className={styles.navLink}
              style={{ '--delay': `${idx * 0.1}s` } as React.CSSProperties}
            >
              <span className={styles.navIcon}>{item.icon}</span>
              <span className={styles.navLabel}>{item.label}</span>
              <div className={styles.navLinkGlow} />
            </Link>
          ))}
        </div>

        <div className={styles.navActions}>
         { !isLoading && <button onClick={AuthService.login} className={styles.secondaryButton}>
            {isAuthenticated ? "Sign Out" : "Log In"}
          </button>}
          
            <button className={styles.primaryButton}>
                <span className={styles.buttonShine} />
                <span className={styles.buttonText}>Start Free Trial</span>
                <svg className={styles.buttonArrow} width="16" height="16" viewBox="0 0 16 16" fill="none">
                <path d="M6 12L10 8L6 4" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                </svg>
            </button>
          
        </div>
      </div>

      <div className={styles.navBorder} />
    </nav>
  );
}
'use client';

import { Search, Link2, Sparkles, Globe } from 'lucide-react';
import styles from './front-page-banner.page.module.css';

export default function BacklinkScraperBanner() {
  return (
    <div className={styles.banner}>
      {/* Subtle background pattern */}
      <div className={styles.backgroundPattern}>
        <div className={styles.bgCircleTop}></div>
        <div className={styles.bgCircleBottom}></div>
      </div>

      <div className={styles.container}>
        <div className={styles.grid}>
          {/* Left content */}
          <div>
            <div className={styles.badge}>
              <Sparkles className={styles.badgeIcon} />
              <span className={styles.badgeText}>AI-Powered Analysis</span>
            </div>

            <h1 className={styles.heading}>
              Discover Forum Backlinks
            </h1>

            <p className={styles.description}>
              Track your brand mentions across forums, discussion boards, and communities worldwide.
            </p>

            <div className={styles.features}>
              <div className={styles.feature}>
                <div className={styles.featureIcon}>
                  <Search className={styles.icon} />
                </div>
                <span className={styles.featureText}>Deep Scan</span>
              </div>

              <div className={styles.feature}>
                <div className={styles.featureIcon}>
                  <Link2 className={styles.icon} />
                </div>
                <span className={styles.featureText}>Link Quality</span>
              </div>

              <div className={styles.feature}>
                <div className={styles.featureIcon}>
                  <Globe className={styles.icon} />
                </div>
                <span className={styles.featureText}>Real-time</span>
              </div>
            </div>
          </div>

          {/* Right visual */}
          <div className={styles.visualContainer}>
            <div className={styles.visualWrapper}>
              {/* Animated circles */}
              <div className={styles.circlesContainer}>
                <div className={styles.circleLarge}></div>
                <div className={styles.circleMedium}></div>
                <div className={styles.circleSmall}></div>
              </div>

              {/* Floating icons */}
              <div className={`${styles.floatingIcon} ${styles.floatingIcon1}`}>
                <Link2 className={styles.floatingIconSvg} />
              </div>

              <div className={`${styles.floatingIcon} ${styles.floatingIcon2}`}>
                <Search className={styles.floatingIconSvg} />
              </div>

              <div className={`${styles.floatingIcon} ${styles.floatingIcon3}`}>
                <Globe className={styles.floatingIconSvg} />
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Bottom accent line */}
      <div className={styles.accentLine}></div>
    </div>
  );
}
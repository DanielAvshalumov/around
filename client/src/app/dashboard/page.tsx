'use client'
import React, { useState } from 'react';
import { BarChart3, ExternalLink, MousePointerClick, Eye, TrendingUp, Calendar } from 'lucide-react';
import styles from './page.module.css';

export default function Dashboard() {
  const [dashboardData] = useState({
    urls: [
      {
        id: 1,
        url: 'https://example.com/product-launch',
        shortUrl: 'bit.ly/prod2024',
        clicks: 12847,
        impressions: 45231,
        ctr: 28.4,
        trend: '+12.5%'
      },
      {
        id: 2,
        url: 'https://example.com/marketing-campaign',
        shortUrl: 'bit.ly/mktg24',
        clicks: 8392,
        impressions: 31024,
        ctr: 27.1,
        trend: '+8.3%'
      }
    ],
    totalClicks: 21239,
    totalImpressions: 76255,
    avgCtr: 27.8
  });

  const handleViewAnalytics = () => {
    alert('Analytics view opened');
  };

  const formatNumber = (num) => {
    return num.toLocaleString();
  };

  return (
    <div className={styles.container}>
      <div className={styles.maxWidth}>
        {/* Header */}
        <div className={styles.header}>
          <div className={styles.headerContent}>
            <div className={styles.headerText}>
              <h1>Analytics Dashboard</h1>
              <p>
                <Calendar className={styles.trendIcon} />
                Last 30 days
              </p>
            </div>
            <button onClick={handleViewAnalytics} className={styles.button}>
              View Full Analytics
            </button>
          </div>
        </div>

        {/* Main Content Grid */}
        <div className={styles.mainGrid}>
          {/* Stats Sidebar */}
          <div className={styles.sidebar}>
            <div className={styles.statCard}>
              <div className={styles.statCardHeader}>
                <div className={`${styles.iconWrapper} ${styles.iconWrapperBlue}`}>
                  <MousePointerClick className={styles.iconBlue} />
                </div>
                <span className={styles.trend}>
                  <TrendingUp className={styles.trendIcon} />
                  +10.2%
                </span>
              </div>
              <h3 className={styles.statLabel}>Total Clicks</h3>
              <p className={styles.statValue}>{formatNumber(dashboardData.totalClicks)}</p>
            </div>

            <div className={styles.statCard}>
              <div className={styles.statCardHeader}>
                <div className={`${styles.iconWrapper} ${styles.iconWrapperPurple}`}>
                  <Eye className={styles.iconPurple} />
                </div>
                <span className={styles.trend}>
                  <TrendingUp className={styles.trendIcon} />
                  +15.8%
                </span>
              </div>
              <h3 className={styles.statLabel}>Total Impressions</h3>
              <p className={styles.statValue}>{formatNumber(dashboardData.totalImpressions)}</p>
            </div>

            <div className={styles.statCard}>
              <div className={styles.statCardHeader}>
                <div className={`${styles.iconWrapper} ${styles.iconWrapperEmerald}`}>
                  <BarChart3 className={styles.iconEmerald} />
                </div>
                <span className={styles.trend}>
                  <TrendingUp className={styles.trendIcon} />
                  +5.1%
                </span>
              </div>
              <h3 className={styles.statLabel}>Average CTR</h3>
              <p className={styles.statValue}>{dashboardData.avgCtr}%</p>
            </div>
          </div>

          {/* URL Performance Table */}
          <div className={styles.tableContainer}>
            <div className={styles.tableHeader}>
              <h2>
                <ExternalLink className={styles.tableHeaderIcon} />
                URL Performance
              </h2>
            </div>
            
            <div className={styles.tableWrapper}>
              <table className={styles.table}>
                <thead>
                  <tr>
                    <th className={styles.alignLeft}>URL</th>
                    <th className={styles.alignLeft}>Short Link</th>
                    <th className={styles.alignRight}>Clicks</th>
                    <th className={styles.alignRight}>Impressions</th>
                    <th className={styles.alignRight}>CTR</th>
                    <th className={styles.alignRight}>Trend</th>
                  </tr>
                </thead>
                <tbody>
                  {dashboardData.urls.map((item) => (
                    <tr key={item.id}>
                      <td>
                        <div className={styles.urlCell}>
                          <span className={styles.urlText}>{item.url}</span>
                          <span className={styles.urlId}>ID: {item.id}</span>
                        </div>
                      </td>
                      <td>
                        <a 
                          href={item.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className={styles.shortLink}
                        >
                          {item.shortUrl}
                          <ExternalLink className={styles.linkIcon} />
                        </a>
                      </td>
                      <td>
                        <div className={styles.metricCell}>
                          <span className={styles.metricValue}>{formatNumber(item.clicks)}</span>
                          <MousePointerClick className={styles.metricIcon} />
                        </div>
                      </td>
                      <td>
                        <div className={styles.metricCell}>
                          <span className={styles.metricValue}>{formatNumber(item.impressions)}</span>
                          <Eye className={styles.metricIcon} />
                        </div>
                      </td>
                      <td>
                        <div className={styles.metricCell}>
                          <div className={styles.ctrBadge}>
                            <span className={styles.ctrValue}>{item.ctr}%</span>
                          </div>
                        </div>
                      </td>
                      <td>
                        <span className={styles.trendCell}>
                          <TrendingUp className={styles.trendCellIcon} />
                          {item.trend}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className={styles.footer}>
          <p>Data refreshes every 5 minutes • Last updated: Just now</p>
        </div>
      </div>
    </div>
  );
}
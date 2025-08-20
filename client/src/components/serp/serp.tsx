// SourceLinkTable.jsx
import React from 'react';
import { ExternalLink } from 'lucide-react';
import styles from './page.module.css';

const SourceLinkTable = ({ backlinks, isCentered }: {backlinks: any[], isCentered: boolean}) => {
  // Sample data - replace with your actual data
  const sourceData = [
    {
      sourceLink: "https://example.com/article-1",
      domainAuthority: 85,
      bankLinks: "chase.com, bankofamerica.com"
    },
    {
      sourceLink: "https://another-site.com/news",
      domainAuthority: 72,
      bankLinks: "wellsfargo.com, citibank.com"
    },
    {
      sourceLink: "https://finance-blog.com/post",
      domainAuthority: 58,
      bankLinks: "usbank.com, pnc.com"
    },
    {
      sourceLink: "https://news-outlet.com/report",
      domainAuthority: 91,
      bankLinks: "jpmorgan.com, goldmansachs.com"
    }
  ];

  const handleViewSource = (url) => {
    window.open(url, '_blank');
  };

  const getDomainAuthorityClass = (da) => {
    if (da >= 80) return styles.daHigh;
    if (da >= 60) return styles.daMedium;
    if (da >= 40) return styles.daLow;
    return styles.daVeryLow;
  };

  return (
    <div className={`${styles.tableContainer} ${isCentered ? styles.hidden : styles.visible}`}>
      {!isCentered && (
        <div className={styles.tableWrapper}>
          <table className={styles.table}>
            <thead className={styles.tableHead}>
              <tr>
                <th className={styles.tableHeader}>Source Link</th>
                <th className={`${styles.tableHeader} ${styles.textCenter}`}>Domain Authority</th>
                <th className={styles.tableHeader}>Bank Links</th>
                <th className={`${styles.tableHeader} ${styles.textCenter}`}>Action</th>
              </tr>
            </thead>
            <tbody className={styles.tableBody}>
              {backlinks?.map((item, index) => (
                <tr key={index} className={styles.tableRow}>
                  <td className={styles.tableCell}>
                    <div className={styles.linkContainer}>
                      <ExternalLink className={styles.linkIcon} />
                      <span 
                        className={styles.sourceLink} 
                        onClick={() => handleViewSource(item.Source)}
                        title={item.Source}
                      >
                        {item.Source}
                      </span>
                    </div>
                  </td>
                  {/* <td className={`${styles.tableCell} ${styles.textCenter}`}>
                    <span className={`${styles.daBadge} ${getDomainAuthorityClass(item.domainAuthority)}`}>
                      {item.domainAuthority}
                    </span>
                  </td> */}
                  <td className={styles.tableCell}>
                    <div className={styles.bankLinks} title={item.Backlink}>
                      {item.Backlink}
                    </div>
                  </td>
                  <td className={`${styles.tableCell} ${styles.textCenter}`}>
                    <button 
                      onClick={() => handleViewSource(item.Source)}
                      className={styles.viewButton}
                    >
                      View
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default SourceLinkTable;
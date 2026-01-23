// SourceLinkTable.jsx
import React from 'react';
import { ExternalLink } from 'lucide-react';
import styles from './page.module.css';

const SourceLinkTable = ({ backlinks, isCentered }: {backlinks: any[], isCentered: boolean}) => {
  

  const handleViewSource = (url: string) => {
    window.open(url, '_blank');
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
                <th className={`${styles.tableHeader} ${styles.textCenter}`}>Add to Training Data</th>
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
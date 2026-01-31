// SourceLinkTable.jsx
import React, { useEffect } from 'react';
import { ExternalLink } from 'lucide-react';
import styles from './page.module.css';
import Link from 'next/link';

const SourceLinkTable = ({ backlinks, isCentered }: {backlinks: any[], isCentered: boolean}) => {
  
  useEffect(() => {
    console.log(backlinks);
  },[backlinks])

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
                <th className={`${styles.tableHeader} ${styles.textCenter}`}>Generate Response</th>
              </tr>
            </thead>
            <tbody className={styles.tableBody}>
              {backlinks?.map((item, index) => (
                <tr key={index} id={item.Id} className={styles.tableRow}>
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
                  <td className={`${styles.tableCell} ${styles.textCenter}`}>
                    <Link
                      href={`/backlink/${item.Id}`}
                      className={styles.viewButton}
                    >
                      Add
                    </Link>
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
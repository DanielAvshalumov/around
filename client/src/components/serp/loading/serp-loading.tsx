import React from 'react';

const LoadingTable = ({ isCentered }) => {
  const styles = {
    tableContainer: {
      transition: 'all 0.7s ease-in-out',
      overflow: 'hidden',
      borderRadius: '0.75rem',
      padding: '2px',
      background: 'linear-gradient(90deg, #3b82f6, #9333ea, #ec4899)',
      width: isCentered ? '0' : 'flex: 1',
      opacity: isCentered ? '0' : '1'
    },
    tableWrapper: {
      backgroundColor: 'white',
      borderRadius: '0.5rem',
      overflow: 'hidden',
      height: '100%'
    },
    table: {
      width: '100%',
      borderCollapse: 'collapse'
    },
    tableHead: {
      background: 'linear-gradient(90deg, #eff6ff, #faf5ff)'
    },
    tableHeader: {
      padding: '1rem 1.5rem',
      textAlign: 'left',
      fontSize: '0.875rem',
      fontWeight: '600',
      color: '#374151'
    },
    tableHeaderCenter: {
      textAlign: 'center'
    },
    tableRow: {
      borderBottom: '1px solid #e5e7eb'
    },
    tableCell: {
      padding: '1rem 1.5rem',
      verticalAlign: 'middle'
    },
    tableCellCenter: {
      textAlign: 'center'
    },
    skeleton: {
      background: 'linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%)',
      backgroundSize: '200% 100%',
      animation: 'loading 1.5s infinite',
      borderRadius: '4px'
    },
    skeletonText: {
      height: '16px',
      marginBottom: '4px'
    },
    skeletonBadge: {
      height: '24px',
      width: '48px',
      borderRadius: '12px'
    },
    skeletonButton: {
      height: '32px',
      width: '64px',
      borderRadius: '6px'
    }
  };

  const keyframes = `
    @keyframes loading {
      0% {
        background-position: 200% 0;
      }
      100% {
        background-position: -200% 0;
      }
    }
  `;

  return (
    <>
      <style>{keyframes}</style>
      <div style={styles.tableContainer}>
        {!isCentered && (
          <div style={styles.tableWrapper}>
            <table style={styles.table}>
              <thead style={styles.tableHead}>
                <tr>
                  <th style={styles.tableHeader}>Source Link</th>
                  <th style={{...styles.tableHeader, ...styles.tableHeaderCenter}}>Domain Authority</th>
                  <th style={styles.tableHeader}>Bank Links</th>
                  <th style={{...styles.tableHeader, ...styles.tableHeaderCenter}}>Action</th>
                </tr>
              </thead>
              <tbody>
                {[...Array(4)].map((_, index) => (
                  <tr key={index} style={styles.tableRow}>
                    <td style={styles.tableCell}>
                      <div style={{...styles.skeleton, ...styles.skeletonText, width: '200px'}}></div>
                    </td>
                    <td style={{...styles.tableCell, ...styles.tableCellCenter}}>
                      <div style={{...styles.skeleton, ...styles.skeletonBadge, margin: '0 auto'}}></div>
                    </td>
                    <td style={styles.tableCell}>
                      <div style={{...styles.skeleton, ...styles.skeletonText, width: '150px'}}></div>
                    </td>
                    <td style={{...styles.tableCell, ...styles.tableCellCenter}}>
                      <div style={{...styles.skeleton, ...styles.skeletonButton, margin: '0 auto'}}></div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </>
  );
};

export default LoadingTable;
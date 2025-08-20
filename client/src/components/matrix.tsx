'use client'
import { useEffect, useRef } from 'react';

const keywords = [
  "ads","SEO","PPC","click","leads","brand","reach","sales","growth","traffic",
  "content","funnel","impress","metrics","pixels","target","ranking","organic",
  "paid","search","Google","Bing","social","email","bounce","CTA","copy","retarget",
  "optimize","keyword","hashtag","audit","analytics","A/B test","domain","hosting",
  "mobile","video","shorts","buyer","persona","clickbait","caption","schedule",
  "campaign","headline","convo","story","subscribe","CTR","CPC","conversion","funnels",
  "nurture","scroll","segment","split test","inbound","outbound","influencer","reels",
  "UGC","brand lift","automation","landing page","lead magnet","open rate","bounce rate",
  "geo-target","lookalike","boost","sponsor","KPI","ROI","lifetime value","CAC","drip"
];

const MatrixBackground = () => {
  const matrixRef = useRef(null);
  
  useEffect(() => {
    const container = matrixRef.current;
    if (!container) return;
    
    container.innerHTML = "";
    const numColumns = Math.floor(window.innerWidth / 100);
    
    for (let i = 0; i < numColumns; i++) {
      const word = document.createElement("div");
      word.className = "matrix-word";
      word.innerText = keywords[Math.floor(Math.random() * keywords.length)];
      word.style.left = `${i * 100}px`;
      word.style.animationDelay = `${Math.random() * 5}s`;
      word.style.animationDuration = `${3 + Math.random() * 3}s`;
      container.appendChild(word);
    }
    
    const interval = setInterval(() => {
      if (container) {
        const word = document.createElement("div");
        word.className = "matrix-word";
        word.innerText = keywords[Math.floor(Math.random() * keywords.length)];
        word.style.left = `${Math.floor(Math.random() * window.innerWidth)}px`;
        word.style.animationDuration = `${3 + Math.random() * 3}s`;
        container.appendChild(word);
        
        if (container.children.length > numColumns * 3) {
          container.removeChild(container.children[0]);
        }
      }
    }, 1000);
    
    return () => clearInterval(interval);
  }, []);

  return (
    <>
      <div 
        ref={matrixRef}
        className="matrix-container"
        style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          zIndex: -1,
          pointerEvents: 'none',
          backgroundColor: 'white'
        }}
      />
      
      <style jsx>{`
        .matrix-word {
          position: absolute;
          top: -20px;
          color: rgba(0, 0, 0, 0.7);
          font-family: monospace;
          font-weight: bold;
          font-size: 16px;
          text-shadow: 0 0 5px rgba(0, 0, 0, 0.2);
          animation: fall linear infinite;
          opacity: 0.6;
          pointer-events: none;
        }
        
        @keyframes fall {
          0% {
            transform: translateY(-100px);
            opacity: 0;
          }
          10% {
            opacity: 0.6;
          }
          90% {
            opacity: 0.6;
          }
          100% {
            transform: translateY(100vh);
            opacity: 0;
          }
        }
      `}</style>
    </>
  );
};

export default MatrixBackground;
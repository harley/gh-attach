import React from 'react';
import {loadFont as loadSora} from '@remotion/google-fonts/Sora';
import {AbsoluteFill} from 'remotion';

const {fontFamily} = loadSora('normal', {
  weights: ['400', '600', '700'],
  subsets: ['latin'],
});

export type StillProps = {
  eyebrow: string;
  headline: string;
  subheadline: string;
  command: string;
};

const panelStyle: React.CSSProperties = {
  position: 'absolute',
  borderRadius: 36,
  boxShadow: '0 40px 120px rgba(15, 23, 42, 0.18)',
};

export const GhAttachStill: React.FC<StillProps> = ({eyebrow, headline, subheadline, command}) => {
  return (
    <AbsoluteFill
      style={{
        background: 'radial-gradient(circle at top left, #fff0d7 0%, #fffaf4 30%, #e0f2fe 72%, #f8fafc 100%)',
        fontFamily,
        overflow: 'hidden',
      }}
    >
      <div
        style={{
          position: 'absolute',
          top: 72,
          left: 88,
          width: 780,
          display: 'flex',
          flexDirection: 'column',
          gap: 20,
        }}
      >
        <div
          style={{
            fontSize: 22,
            fontWeight: 700,
            textTransform: 'uppercase',
            letterSpacing: '0.14em',
            color: '#c2410c',
          }}
        >
          {eyebrow}
        </div>
        <div
          style={{
            fontSize: 94,
            lineHeight: 0.92,
            letterSpacing: '-0.07em',
            fontWeight: 700,
            color: '#0f172a',
          }}
        >
          {headline}
        </div>
        <div
          style={{
            fontSize: 34,
            lineHeight: 1.25,
            letterSpacing: '-0.04em',
            color: '#334155',
            maxWidth: 720,
          }}
        >
          {subheadline}
        </div>
      </div>

      <div
        style={{
          ...panelStyle,
          top: 112,
          right: 98,
          width: 560,
          height: 620,
          background: 'rgba(8,12,23,0.94)',
          border: '1px solid rgba(255,255,255,0.08)',
          padding: 34,
          display: 'flex',
          flexDirection: 'column',
          gap: 24,
        }}
      >
        <div style={{display: 'flex', gap: 10}}>
          {['#fb7185', '#f59e0b', '#22c55e'].map((color) => (
            <div
              key={color}
              style={{
                width: 16,
                height: 16,
                borderRadius: 999,
                background: color,
              }}
            />
          ))}
        </div>
        <div style={{fontSize: 18, color: '#94a3b8', fontWeight: 600}}>Create a PR with attachments</div>
        <div
          style={{
            fontSize: 36,
            lineHeight: 1.28,
            letterSpacing: '-0.05em',
            color: '#f8fafc',
            fontWeight: 600,
            whiteSpace: 'pre-wrap',
          }}
        >
          <span style={{color: '#fb923c'}}>$ </span>
          {command}
        </div>
        <div style={{display: 'flex', gap: 14, flexWrap: 'wrap'}}>
          {[
            ['launch.png', '#f97316'],
            ['walkthrough.mp4', '#0ea5e9'],
            ['private repo safe', '#22c55e'],
          ].map(([label, accent]) => (
            <div
              key={label}
              style={{
                padding: '12px 18px',
                background: '#fff',
                borderRadius: 999,
                fontSize: 20,
                fontWeight: 700,
                color: '#0f172a',
                display: 'inline-flex',
                alignItems: 'center',
                gap: 12,
              }}
            >
              <span
                style={{
                  width: 12,
                  height: 12,
                  borderRadius: 999,
                  background: accent,
                  display: 'inline-block',
                }}
              />
              {label}
            </div>
          ))}
        </div>
      </div>

      <div
        style={{
          ...panelStyle,
          left: 96,
          bottom: 88,
          width: 1080,
          height: 250,
          background: 'rgba(255,255,255,0.9)',
          border: '1px solid rgba(15,23,42,0.05)',
          padding: 32,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          backdropFilter: 'blur(20px)',
        }}
      >
        <div style={{display: 'flex', flexDirection: 'column', gap: 14}}>
          <div style={{fontSize: 28, fontWeight: 700, color: '#0f172a', letterSpacing: '-0.04em'}}>
            Native GitHub attachments inside your PR body
          </div>
          <div style={{fontSize: 22, lineHeight: 1.35, color: '#475569', maxWidth: 720}}>
            Built for the real workflow: upload first, compose body second, and keep media privacy aligned with the repo.
          </div>
        </div>
        <div
          style={{
            width: 220,
            height: 220,
            borderRadius: 28,
            background:
              'linear-gradient(150deg, rgba(249,115,22,0.18), rgba(14,165,233,0.24), rgba(251,113,133,0.16))',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: 28,
            fontWeight: 700,
            color: '#0f172a',
            letterSpacing: '-0.05em',
          }}
        >
          gh attach
        </div>
      </div>
    </AbsoluteFill>
  );
};

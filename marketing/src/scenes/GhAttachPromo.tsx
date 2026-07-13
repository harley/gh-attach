import React from 'react';
import {loadFont as loadSpaceGrotesk} from '@remotion/google-fonts/SpaceGrotesk';
import {AbsoluteFill, Easing, interpolate, spring, useCurrentFrame, useVideoConfig} from 'remotion';

const {fontFamily} = loadSpaceGrotesk('normal', {
  weights: ['400', '500', '700'],
  subsets: ['latin'],
});

export type PromoProps = {
  title: string;
  subtitle: string;
  command: string;
};

const shellPanel: React.CSSProperties = {
  width: 740,
  height: 250,
  borderRadius: 28,
  background: 'rgba(8, 12, 23, 0.92)',
  border: '1px solid rgba(255,255,255,0.08)',
  boxShadow: '0 32px 80px rgba(2, 6, 23, 0.45)',
  padding: 28,
  display: 'flex',
  flexDirection: 'column',
  gap: 20,
};

const githubCard: React.CSSProperties = {
  width: 420,
  height: 430,
  borderRadius: 32,
  background: 'rgba(248, 250, 252, 0.95)',
  boxShadow: '0 36px 90px rgba(15, 23, 42, 0.22)',
  padding: 28,
  display: 'flex',
  flexDirection: 'column',
  gap: 20,
};

const commandWindow: React.CSSProperties = {
  fontFamily,
  fontSize: 28,
  lineHeight: 1.35,
  color: '#f8fafc',
  fontWeight: 500,
  whiteSpace: 'pre-wrap',
};

const bgOrb = (top: number, left: number, color: string, size: number): React.CSSProperties => ({
  position: 'absolute',
  top,
  left,
  width: size,
  height: size,
  borderRadius: size / 2,
  background: color,
  filter: 'blur(90px)',
  opacity: 0.55,
});

const AttachmentChip: React.FC<{
  frame: number;
  delay: number;
  label: string;
  accent: string;
}> = ({frame, delay, label, accent}) => {
  const {fps} = useVideoConfig();
  const reveal = spring({
    frame: frame - delay,
    fps,
    config: {damping: 200},
    durationInFrames: 24,
  });

  const translateY = interpolate(reveal, [0, 1], [24, 0]);
  const opacity = interpolate(reveal, [0, 1], [0, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
  });

  return (
    <div
      style={{
        transform: `translateY(${translateY}px)`,
        opacity,
        borderRadius: 999,
        background: '#ffffff',
        color: '#0f172a',
        padding: '12px 18px',
        fontSize: 20,
        fontWeight: 700,
        letterSpacing: '-0.03em',
        display: 'inline-flex',
        alignItems: 'center',
        gap: 12,
        boxShadow: '0 18px 30px rgba(15, 23, 42, 0.12)',
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
  );
};

export const GhAttachPromo: React.FC<PromoProps> = ({title, subtitle, command}) => {
  const frame = useCurrentFrame();
  const {fps} = useVideoConfig();

  const heroReveal = spring({
    frame,
    fps,
    config: {damping: 200},
    durationInFrames: 36,
  });

  const shellProgress = spring({
    frame: frame - 18,
    fps,
    config: {damping: 200},
    durationInFrames: 36,
  });

  const cardProgress = spring({
    frame: frame - 34,
    fps,
    config: {damping: 200},
    durationInFrames: 36,
  });

  const pulse = interpolate(frame, [0, 90, 180], [0.96, 1, 0.98], {
    easing: Easing.inOut(Easing.sin),
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
  });

  const commandChars = Math.floor(
    interpolate(frame, [12, 72], [0, command.length], {
      extrapolateLeft: 'clamp',
      extrapolateRight: 'clamp',
    }),
  );

  const prOpacity = interpolate(frame, [72, 102], [0, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
  });

  const titleOpacity = interpolate(heroReveal, [0, 1], [0, 1]);
  const titleY = interpolate(heroReveal, [0, 1], [36, 0]);
  const shellX = interpolate(shellProgress, [0, 1], [-90, 0]);
  const cardX = interpolate(cardProgress, [0, 1], [90, 0]);

  return (
    <AbsoluteFill
      style={{
        background: 'linear-gradient(140deg, #fff8ef 0%, #ffe5d7 28%, #d7f1ff 62%, #eff6ff 100%)',
        fontFamily,
        overflow: 'hidden',
      }}
    >
      <div style={bgOrb(60, 820, '#f97316', 280)} />
      <div style={bgOrb(440, 180, '#0ea5e9', 340)} />
      <div style={bgOrb(130, 270, '#fb7185', 220)} />

      <div
        style={{
          position: 'absolute',
          inset: 0,
          display: 'flex',
          flexDirection: 'column',
          padding: '56px 64px',
          transform: `scale(${pulse})`,
        }}
      >
        <div
          style={{
            transform: `translateY(${titleY}px)`,
            opacity: titleOpacity,
            display: 'flex',
            flexDirection: 'column',
            gap: 16,
            width: 760,
          }}
        >
          <div
            style={{
              fontSize: 18,
              fontWeight: 700,
              letterSpacing: '0.12em',
              textTransform: 'uppercase',
              color: '#c2410c',
            }}
          >
            gh extension for media-rich PRs
          </div>
          <div
            style={{
              fontSize: 72,
              lineHeight: 0.95,
              fontWeight: 700,
              letterSpacing: '-0.06em',
              color: '#111827',
            }}
          >
            {title}
          </div>
          <div
            style={{
              fontSize: 30,
              lineHeight: 1.25,
              color: '#334155',
              letterSpacing: '-0.04em',
              maxWidth: 720,
            }}
          >
            {subtitle}
          </div>
        </div>

        <div
          style={{
            marginTop: 42,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-start',
            gap: 32,
          }}
        >
          <div
            style={{
              ...shellPanel,
              transform: `translateX(${shellX}px)`,
            }}
          >
            <div style={{display: 'flex', gap: 10}}>
              {['#fb7185', '#f59e0b', '#22c55e'].map((color) => (
                <div
                  key={color}
                  style={{
                    width: 14,
                    height: 14,
                    borderRadius: 999,
                    background: color,
                  }}
                />
              ))}
            </div>
            <div style={{fontSize: 18, color: '#94a3b8', fontWeight: 500}}>Terminal</div>
            <div style={commandWindow}>
              <span style={{color: '#f97316', marginRight: 12}}>$</span>
              {command.slice(0, commandChars)}
              <span
                style={{
                  opacity: frame % 18 < 9 ? 1 : 0,
                  color: '#fdba74',
                }}
              >
                |
              </span>
            </div>
            <div style={{display: 'flex', gap: 12}}>
              <AttachmentChip frame={frame} delay={74} label="walkthrough.mp4" accent="#2563eb" />
              <AttachmentChip frame={frame} delay={80} label="launch.png" accent="#f97316" />
            </div>
          </div>

          <div
            style={{
              ...githubCard,
              transform: `translateX(${cardX}px)`,
              opacity: prOpacity,
            }}
          >
            <div style={{display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}>
              <div style={{fontSize: 26, fontWeight: 700, color: '#0f172a'}}>Pull request</div>
              <div
                style={{
                  background: '#dcfce7',
                  color: '#166534',
                  borderRadius: 999,
                  padding: '8px 14px',
                  fontSize: 16,
                  fontWeight: 700,
                }}
              >
                Open
              </div>
            </div>

            <div
              style={{
                background: '#f8fafc',
                borderRadius: 24,
                padding: 22,
                border: '1px solid #e2e8f0',
                display: 'flex',
                flexDirection: 'column',
                gap: 16,
              }}
            >
              <div style={{fontSize: 22, fontWeight: 700, color: '#0f172a'}}>
                Demo: native screenshot and video attachments
              </div>
              <div style={{fontSize: 18, lineHeight: 1.4, color: '#475569'}}>
                This PR was created with <span style={{fontWeight: 700}}>gh-attach</span> and includes native GitHub media in the body.
              </div>
              <div
                style={{
                  height: 140,
                  borderRadius: 20,
                  background:
                    'linear-gradient(135deg, rgba(14,165,233,0.18), rgba(249,115,22,0.24), rgba(251,113,133,0.16))',
                  border: '1px solid rgba(15,23,42,0.06)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: '#0f172a',
                  fontSize: 22,
                  fontWeight: 700,
                  letterSpacing: '-0.04em',
                }}
              >
                inline screenshot + demo video
              </div>
            </div>

            <div style={{display: 'flex', gap: 12, flexWrap: 'wrap'}}>
              <AttachmentChip frame={frame} delay={92} label="private repo safe" accent="#16a34a" />
              <AttachmentChip frame={frame} delay={98} label="native media URLs" accent="#ea580c" />
            </div>
          </div>
        </div>
      </div>
    </AbsoluteFill>
  );
};

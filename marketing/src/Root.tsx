import React from 'react';
import {Composition, Still} from 'remotion';
import {GhAttachPromo, type PromoProps} from './scenes/GhAttachPromo';
import {GhAttachStill, type StillProps} from './scenes/GhAttachStill';

const promoProps = {
  title: 'Native GitHub attachments for PRs',
  subtitle: 'Upload screenshots and video. Keep private repos private.',
  command: 'gh attach pr create --attach demo.mp4 --attach screenshot.png',
} satisfies PromoProps;

const stillProps = {
  eyebrow: 'gh extension',
  headline: 'Create PRs with native attachments',
  subheadline: 'No release-asset hack. No public image host. Real GitHub attachments.',
  command: 'gh attach pr create --attach launch.png --attach walkthrough.mp4',
} satisfies StillProps;

export const RemotionRoot: React.FC = () => {
  return (
    <>
      <Composition
        id="GhAttachPromo"
        component={GhAttachPromo}
        durationInFrames={180}
        fps={30}
        width={1280}
        height={720}
        defaultProps={promoProps}
      />
      <Still
        id="GhAttachStill"
        component={GhAttachStill}
        width={1600}
        height={900}
        defaultProps={stillProps}
      />
    </>
  );
};

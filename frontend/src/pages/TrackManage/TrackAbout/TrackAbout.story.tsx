import { mockTrack } from '@/api/fixtures/track';
import { StoryPreview } from '@/components/StoryBook/StoryPreview';
import { TrackAbout } from './TrackAbout';

export default {
	title: 'TrackAbout',
};

const save = () => Promise.resolve();

/** Как трек читается: тексты и кнопка, открывающая их правку на месте. */
export function Default() {
	return (
		<StoryPreview isPage>
			<TrackAbout track={mockTrack.Track} onSave={save} />
		</StoryPreview>
	);
}
